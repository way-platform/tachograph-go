package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	var inputFile string
	var outputDir string
	var enableCleanup bool
	var chunkLevel int

	flag.StringVar(&inputFile, "i", "", "Input HTML file to parse")
	flag.StringVar(&outputDir, "d", "", "Output directory for parsed files")
	flag.BoolVar(&enableCleanup, "cleanup", true, "Enable HTML cleanup for improved information density")
	flag.IntVar(&chunkLevel, "chunk-level", 1, "Chunking level: 1=main sections only, 2=include level-2 sections, 3=include level-3 sections")
	flag.Parse()

	if inputFile == "" {
		log.Fatal("Input file (-i) is required")
	}
	if outputDir == "" {
		log.Fatal("Output directory (-d) is required")
	}
	if chunkLevel < 1 || chunkLevel > 3 {
		log.Fatal("Chunk level must be 1, 2, or 3")
	}

	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Process the HTML file
	if err := processHTML(inputFile, outputDir, enableCleanup, chunkLevel); err != nil {
		log.Fatalf("Failed to process HTML file: %v", err)
	}

	fmt.Printf("Successfully processed %s and created section files in %s/\n", inputFile, outputDir)
}

func processHTML(inputFile, outputDir string, enableCleanup bool, chunkLevel int) error {
	// Open input file
	file, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer file.Close()

	fmt.Println("Reading HTML file...")

	// Read the entire file (we need to do this for goquery to work properly)
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file content: %w", err)
	}

	fmt.Println("Parsing HTML with goquery...")

	// Parse HTML with goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(fileContent)))
	if err != nil {
		return fmt.Errorf("failed to parse HTML: %w", err)
	}

	fmt.Println("Removing image tags...")

	// Remove all img tags
	doc.Find("img").Remove()

	// Also remove any base64 data URIs that might be in other attributes
	// This is a more aggressive cleanup to reduce file size
	doc.Find("*").Each(func(i int, s *goquery.Selection) {
		// Check all attributes for data URIs
		for _, attr := range []string{"src", "href", "data", "style"} {
			if val, exists := s.Attr(attr); exists {
				if strings.Contains(val, "data:image/") {
					s.RemoveAttr(attr)
				}
			}
		}
	})

	// Always split the document first, then apply cleanup to individual sections
	return splitDocument(doc, outputDir, enableCleanup, chunkLevel)
}

func splitDocument(doc *goquery.Document, outputDir string, enableCleanup bool, chunkLevel int) error {
	fmt.Printf("Splitting document into sections (chunk level %d)...\n", chunkLevel)

	// Build selector based on chunk level
	selectors := []string{"p.title-article-norm", "p.title-annex-1", "p.title-gr-seq-level-1"}

	if chunkLevel >= 2 {
		selectors = append(selectors, "p.title-gr-seq-level-2")
	}
	if chunkLevel >= 3 {
		selectors = append(selectors, "p.title-gr-seq-level-3")
	}

	selectorString := strings.Join(selectors, ", ")
	fmt.Printf("Using selector: %s\n", selectorString)

	// Find all section boundaries
	sections := []struct {
		Title     string
		FileName  string
		Selection *goquery.Selection
		Level     int
	}{}

	doc.Find(selectorString).Each(func(i int, s *goquery.Selection) {
		// Determine the level of this section
		level := getSectionLevel(s)
		// Skip title-gr-seq-level-1 elements that immediately follow title-annex-1 elements
		// (those are subtitles, not separate sections)
		if s.HasClass("title-gr-seq-level-1") {
			prev := s.Prev()
			if prev.HasClass("title-annex-1") {
				// This is a subtitle for the previous annex, skip it
				return
			}
		}

		title := strings.TrimSpace(s.Text())

		// Extract the descriptive subtitle
		subtitle := extractSectionSubtitle(s)

		fileName := generateFileName(len(sections)+2, title, subtitle, level) // Start at 02 for sections (01 is index)
		displayTitle := title
		if subtitle != "" {
			displayTitle = title + " - " + subtitle
		}

		sections = append(sections, struct {
			Title     string
			FileName  string
			Selection *goquery.Selection
			Level     int
		}{
			Title:     displayTitle,
			FileName:  fileName,
			Selection: s,
			Level:     level,
		})
	})

	if len(sections) == 0 {
		return fmt.Errorf("no sections found to split")
	}

	fmt.Printf("Found %d sections to split\n", len(sections))

	// Extract document head for reuse (after cleanup)
	headHTML, err := doc.Find("head").Html()
	if err != nil {
		return fmt.Errorf("failed to extract head HTML: %w", err)
	}

	// Create each section file
	for i, section := range sections {
		fmt.Printf("Creating section %d: %s -> %s (level %d)\n", i+1, section.Title, section.FileName, section.Level)

		if err := createSectionFile(section.Selection, sections, i, headHTML, outputDir, section.FileName, section.Title, enableCleanup); err != nil {
			return fmt.Errorf("failed to create section %s: %w", section.Title, err)
		}
	}

	// Create the sections index
	if err := createSectionsIndex(sections, outputDir); err != nil {
		return fmt.Errorf("failed to create sections index: %w", err)
	}

	fmt.Printf("Successfully created %d section files + index\n", len(sections))
	return nil
}

func getSectionLevel(s *goquery.Selection) int {
	if s.HasClass("title-article-norm") || s.HasClass("title-annex-1") || s.HasClass("title-gr-seq-level-1") {
		return 1
	}
	if s.HasClass("title-gr-seq-level-2") {
		return 2
	}
	if s.HasClass("title-gr-seq-level-3") {
		return 3
	}
	return 1 // default
}

func createSectionFile(sectionStart *goquery.Selection, allSections []struct {
	Title     string
	FileName  string
	Selection *goquery.Selection
	Level     int
}, sectionIndex int, headHTML, outputDir, fileName, title string, enableCleanup bool,
) error {
	// Create a temporary document for this section
	sectionHTML := "<!DOCTYPE html>\n<html>\n<head>\n" + headHTML + "\n</head>\n<body>\n"

	// Find the next section boundary (or end of document)
	var nextSectionStart *goquery.Selection
	if sectionIndex+1 < len(allSections) {
		nextSectionStart = allSections[sectionIndex+1].Selection
		fmt.Printf("    Next section: %s (ID: %s)\n", allSections[sectionIndex+1].Title, nextSectionStart.AttrOr("id", ""))
	} else {
		fmt.Printf("    This is the last section\n")
	}

	// Use a completely different approach: collect content by traversing siblings
	current := sectionStart

	// Add the section header itself
	if headerHTML, err := goquery.OuterHtml(current); err == nil {
		sectionHTML += headerHTML + "\n"
	}

	// Collect all following siblings until we hit the next section boundary
	for current = current.Next(); current.Length() > 0; current = current.Next() {
		// Check if this element is the start of the next section
		if nextSectionStart != nil && current.Get(0) == nextSectionStart.Get(0) {
			fmt.Printf("    Found section boundary, stopping collection\n")
			break
		}

		// Check if this element contains a section boundary marker
		if current.Find("p.title-article-norm, p.title-annex-1, p.title-gr-seq-level-1, p.title-gr-seq-level-2, p.title-gr-seq-level-3").Length() > 0 {
			// This element contains a section marker, check if it's the next section
			if nextSectionStart != nil {
				nextID := nextSectionStart.AttrOr("id", "")
				if current.Find(fmt.Sprintf("#%s", nextID)).Length() > 0 {
					fmt.Printf("    Found section boundary in container, stopping collection\n")
					break
				}
			}
		}

		// Add this element to the section
		if elementHTML, err := goquery.OuterHtml(current); err == nil {
			sectionHTML += elementHTML + "\n"
		}
	}

	fmt.Printf("    Collected content for section %s\n", title)

	sectionHTML += "</body>\n</html>"

	// Apply cleanup if enabled
	// if enableCleanup {
	// 	// Parse the section HTML and apply cleanup
	// 	sectionDoc, err := goquery.NewDocumentFromReader(strings.NewReader(sectionHTML))
	// 	if err == nil {
	// 		fmt.Printf("  Applying cleanup to %s...\n", fileName)
	// 		cleanupHTML(sectionDoc)

	// 		// Get the cleaned HTML without DOCTYPE (goquery.Html() returns just the <html> element)
	// 		cleanedHTML, err := sectionDoc.Html()
	// 		if err == nil {
	// 			// Remove any duplicate DOCTYPE that might be in the content
	// 			cleanedHTML = regexp.MustCompile(`<!DOCTYPE[^>]*>\s*`).ReplaceAllString(cleanedHTML, "")
	// 			// Remove HTML comments
	// 			cleanedHTML = regexp.MustCompile(`<!--.*?-->`).ReplaceAllString(cleanedHTML, "")
	// 			// Ensure single DOCTYPE and clean HTML structure
	// 			sectionHTML = "<!DOCTYPE html>\n" + cleanedHTML
	// 		}
	// 	}
	// }

	// Format the HTML using gohtml for better readability
	formattedHTML := sectionHTML

	// Write the final HTML to file
	outputFile := filepath.Join(outputDir, fileName)
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	writer.WriteString(formattedHTML)
	return nil
}

func createSectionsIndex(sections []struct {
	Title     string
	FileName  string
	Selection *goquery.Selection
	Level     int
}, outputDir string,
) error {
	indexFile := filepath.Join(outputDir, "01-sections-index.html")
	file, err := os.Create(indexFile)
	if err != nil {
		return fmt.Errorf("failed to create index file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	// Write HTML header
	writer.WriteString(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Document Sections Index</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .section { margin: 20px 0; padding: 10px; border-left: 4px solid #007acc; }
        .section.level-2 { margin-left: 20px; border-left-color: #4CAF50; }
        .section.level-3 { margin-left: 40px; border-left-color: #FF9800; }
        .section-title { font-weight: bold; color: #007acc; }
        .section.level-2 .section-title { color: #4CAF50; }
        .section.level-3 .section-title { color: #FF9800; }
        .section-info { color: #666; font-size: 0.9em; }
        a { text-decoration: none; color: inherit; }
        a:hover .section { background-color: #f0f8ff; }
    </style>
</head>
<body>
    <h1>Document Sections</h1>
`)

	// Write section links
	for _, section := range sections {
		levelClass := ""
		if section.Level > 1 {
			levelClass = fmt.Sprintf(" level-%d", section.Level)
		}
		writer.WriteString(fmt.Sprintf(`    <a href="%s">
        <div class="section%s">
            <div class="section-title">%s</div>
            <div class="section-info">File: %s (Level %d)</div>
        </div>
    </a>
`, section.FileName, levelClass, section.Title, section.FileName, section.Level))
	}

	writer.WriteString(`</body>
</html>`)

	// Format the entire index HTML
	writer.Flush()
	file.Close()

	// Read the file back, format it, and rewrite
	indexContent, err := os.ReadFile(indexFile)
	if err != nil {
		return fmt.Errorf("failed to read index file for formatting: %w", err)
	}

	formattedContent := string(indexContent)

	err = os.WriteFile(indexFile, []byte(formattedContent), 0o644)
	if err != nil {
		return fmt.Errorf("failed to write formatted index file: %w", err)
	}

	return nil
}

func extractSectionSubtitle(sectionHeader *goquery.Selection) string {
	// For articles, look for the next element with stitle-article-norm class
	if sectionHeader.HasClass("title-article-norm") {
		subtitle := sectionHeader.Next().Find("p.stitle-article-norm").Text()
		return strings.TrimSpace(subtitle)
	}

	// For appendices/annexes, look for the next element with title-gr-seq-level-1 class
	if sectionHeader.HasClass("title-annex-1") {
		nextElement := sectionHeader.Next()
		if nextElement.HasClass("title-gr-seq-level-1") {
			// Extract text, preferring boldface content if available
			boldText := nextElement.Find("span.boldface").Text()
			if boldText != "" {
				return strings.TrimSpace(boldText)
			}
			return strings.TrimSpace(nextElement.Text())
		}
	}

	// For title-gr-seq-level-1 elements (appendices), extract boldface content
	if sectionHeader.HasClass("title-gr-seq-level-1") {
		// Extract text, preferring boldface content if available
		boldText := sectionHeader.Find("span.boldface").Text()
		if boldText != "" {
			return strings.TrimSpace(boldText)
		}
		return strings.TrimSpace(sectionHeader.Text())
	}

	// For level-2 and level-3 sections, extract boldface content
	if sectionHeader.HasClass("title-gr-seq-level-2") || sectionHeader.HasClass("title-gr-seq-level-3") {
		boldText := sectionHeader.Find("span.boldface").Text()
		if boldText != "" {
			return strings.TrimSpace(boldText)
		}
		return strings.TrimSpace(sectionHeader.Text())
	}

	return ""
}

func generateFileName(index int, title string, subtitle string, level int) string {
	// Clean the main title
	cleanTitle := strings.TrimSpace(title)

	// Handle special cases for better readability
	cleanTitle = strings.ReplaceAll(cleanTitle, "ANNEX I C", "annex-1c")
	cleanTitle = strings.ReplaceAll(cleanTitle, "ANNEX II", "annex-2")
	cleanTitle = strings.ReplaceAll(cleanTitle, "Article ", "")
	cleanTitle = strings.ReplaceAll(cleanTitle, "Appendix ", "")
	cleanTitle = strings.ReplaceAll(cleanTitle, "Addendum", "addendum")

	// Handle appendix titles from title-gr-seq-level-1 elements
	if subtitle != "" {
		// These are appendices, so prefix with "appendix" if not already present
		if !strings.Contains(strings.ToLower(cleanTitle), "appendix") && !strings.Contains(strings.ToLower(cleanTitle), "annex") {
			// Determine appendix number based on common patterns
			if strings.Contains(subtitle, "DATA DICTIONARY") {
				cleanTitle = "appendix-1"
			} else if strings.Contains(subtitle, "TACHOGRAPH CARDS") {
				cleanTitle = "appendix-2"
			} else if strings.Contains(subtitle, "PICTOGRAMS") {
				cleanTitle = "appendix-3"
			} else if strings.Contains(subtitle, "PRINTOUTS") {
				cleanTitle = "appendix-4"
			} else if strings.Contains(subtitle, "DISPLAY") {
				cleanTitle = "appendix-5"
			} else if strings.Contains(subtitle, "FRONT CONNECTOR") {
				cleanTitle = "appendix-6"
			} else if strings.Contains(subtitle, "DATA DOWNLOADING") {
				cleanTitle = "appendix-7"
			} else if strings.Contains(subtitle, "CALIBRATION PROTOCOL") {
				cleanTitle = "appendix-8"
			} else if strings.Contains(subtitle, "TYPE APPROVAL") {
				cleanTitle = "appendix-9"
			} else if strings.Contains(subtitle, "SECURITY REQUIREMENTS") {
				cleanTitle = "appendix-10"
			} else if strings.Contains(subtitle, "COMMON SECURITY") {
				cleanTitle = "appendix-11"
			} else if strings.Contains(subtitle, "POSITIONING") || strings.Contains(subtitle, "GNSS") {
				cleanTitle = "appendix-12"
			} else if strings.Contains(subtitle, "APPENDIX 13") {
				cleanTitle = "appendix-13"
			} else if strings.Contains(subtitle, "REMOTE COMMUNICATION") {
				cleanTitle = "appendix-14"
			} else if strings.Contains(subtitle, "MIGRATION") {
				cleanTitle = "appendix-15"
			} else if strings.Contains(subtitle, "ADAPTOR") {
				cleanTitle = "appendix-16"
			} else if strings.Contains(subtitle, "APPENDIX 17") {
				cleanTitle = "appendix-17"
			} else {
				// Fallback: use a generic appendix name
				cleanTitle = "appendix-" + strings.ToLower(strings.Fields(subtitle)[0])
			}
		}
	}

	// Clean the subtitle if provided
	cleanSubtitle := ""
	if subtitle != "" {
		cleanSubtitle = strings.TrimSpace(subtitle)
		// Remove special characters and convert to kebab-case
		re := regexp.MustCompile(`[^a-zA-Z0-9\s-]`)
		cleanSubtitle = re.ReplaceAllString(cleanSubtitle, "")
		cleanSubtitle = strings.ToLower(cleanSubtitle)
		cleanSubtitle = regexp.MustCompile(`\s+`).ReplaceAllString(cleanSubtitle, "-")
		cleanSubtitle = strings.Trim(cleanSubtitle, "-")
	}

	// Remove any remaining special characters and convert to lowercase
	re := regexp.MustCompile(`[^a-zA-Z0-9\s-]`)
	cleanTitle = re.ReplaceAllString(cleanTitle, "")
	cleanTitle = strings.ToLower(cleanTitle)

	// Replace multiple spaces with single hyphens
	cleanTitle = regexp.MustCompile(`\s+`).ReplaceAllString(cleanTitle, "-")

	// Remove leading/trailing hyphens
	cleanTitle = strings.Trim(cleanTitle, "-")

	// Combine title and subtitle
	fullName := cleanTitle
	if cleanSubtitle != "" {
		fullName = fmt.Sprintf("%s-%s", cleanTitle, cleanSubtitle)
	}

	// No level prefixes needed - directory structure provides organization
	levelPrefix := ""

	// Truncate if too long to avoid filesystem limits (255 chars total, minus prefix and extension)
	maxNameLength := 240 - len(levelPrefix) - len(".html") - 3 // 3 for index digits
	if len(fullName) > maxNameLength {
		fullName = fullName[:maxNameLength]
		// Try to end at a word boundary
		if lastDash := strings.LastIndex(fullName, "-"); lastDash > maxNameLength-20 {
			fullName = fullName[:lastDash]
		}
	}

	// Clean up redundant repetitions and unnecessary words
	fullName = cleanupRedundantWords(fullName)

	// Generate filename with zero-padded index and descriptive title
	return fmt.Sprintf("%02d-%s%s.html", index, levelPrefix, fullName)
}

// cleanupRedundantWords removes redundant word repetitions and unnecessary prefixes
func cleanupRedundantWords(name string) string {
	// Remove common redundant prefixes
	name = strings.TrimPrefix(name, "appendix-")

	// Split into parts and remove duplicates
	parts := strings.Split(name, "-")
	seen := make(map[string]bool)
	var cleaned []string

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" && !seen[part] {
			// Skip very short parts that are likely noise
			if len(part) > 1 {
				cleaned = append(cleaned, part)
				seen[part] = true
			}
		}
	}

	result := strings.Join(cleaned, "-")

	// Handle specific patterns
	result = strings.ReplaceAll(result, "data-data", "data")
	result = strings.ReplaceAll(result, "introduction-introduction", "introduction")
	result = strings.ReplaceAll(result, "definitions-definitions", "definitions")
	result = strings.ReplaceAll(result, "appendix-appendix", "appendix")

	// Clean up any remaining double hyphens
	result = strings.ReplaceAll(result, "--", "-")
	result = strings.Trim(result, "-")

	return result
}

// cleanupHTML applies ultra-minimal HTML cleanup for maximum information density
func cleanupHTML(doc *goquery.Document) {
	// 1. Remove ALL styling - rely on browser defaults
	removeAllStyling(doc)

	// 2. Fix duplicate DOCTYPE and HTML structure issues
	fixHTMLStructure(doc)

	// 5. Remove redundant wrappers and empty elements
	removeRedundantElements(doc)

	// 6. Remove unnecessary spans and convert semantic ones
	removeUnnecessarySpans(doc)

	// 7. Remove redundant containers
	removeRedundantContainers(doc)

	// 9. Consolidate and optimize content
	consolidateContent(doc)
}

// removeAllStyling removes all CSS and styling to rely on browser defaults
func removeAllStyling(doc *goquery.Document) {
	// Remove external CSS links
	doc.Find("link[rel='stylesheet']").Remove()

	// Remove any existing style elements
	doc.Find("style").Remove()

	// Remove all style attributes from all elements
	doc.Find("*").Each(func(i int, elem *goquery.Selection) {
		elem.RemoveAttr("style")
	})
}

// fixHTMLStructure fixes duplicate DOCTYPE and other structural issues
func fixHTMLStructure(doc *goquery.Document) {
	// Remove any existing DOCTYPE declarations from the body content
	// (they shouldn't be there, but just in case)
	doc.Find("body").Contents().Each(func(i int, s *goquery.Selection) {
		if s.Is("html") {
			// This shouldn't happen, but if it does, unwrap it
			s.Children().Each(func(j int, child *goquery.Selection) {
				childHTML, _ := goquery.OuterHtml(child)
				s.BeforeHtml(childHTML)
			})
			s.Remove()
		}
	})
}

// convertToNativeSemantics converts custom structures to native HTML semantics
func convertToNativeSemantics(doc *goquery.Document) {
	// Convert custom definition lists to native <dl><dt><dd>
	convertToNativeDefinitionLists(doc)

	// Convert titles to proper headings
	convertTitlesToHeadings(doc)

	// Convert table of contents to native navigation
	convertTOCToNativeNav(doc)
}

// convertToNativeDefinitionLists converts custom div-based definitions to native <dl><dt><dd>
func convertToNativeDefinitionLists(doc *goquery.Document) {
	// Find all grid-container elements that represent definition lists
	doc.Find("div.grid-container.grid-list").Each(func(i int, container *goquery.Selection) {
		// Check if this looks like a definition item (has numbered content)
		numberCol := container.Find(".grid-list-column-1")
		contentCol := container.Find(".grid-list-column-2")

		if numberCol.Length() > 0 && contentCol.Length() > 0 {
			// Extract the number/marker and content
			number := strings.TrimSpace(numberCol.Text())

			// Preserve the HTML structure but clean it up
			contentHTML, _ := contentCol.Html()

			// Clean up the content HTML - remove redundant paragraph wrappers for simple content
			// but preserve structure for complex content
			if contentCol.Find("div.grid-container").Length() > 0 || contentCol.Find("p").Length() > 1 {
				// Complex content - preserve structure but clean it up
				contentCol.Find("p.norm").Each(func(j int, p *goquery.Selection) {
					p.RemoveAttr("class")
				})
				contentCol.Find("p.list").Each(func(j int, p *goquery.Selection) {
					p.RemoveAttr("class")
				})
				contentHTML, _ = contentCol.Html()
			} else {
				// Simple content - use text only
				contentHTML = strings.TrimSpace(contentCol.Text())
			}

			// Create native definition list item
			dtHTML := fmt.Sprintf(`<dt>%s</dt>`, number)
			ddHTML := fmt.Sprintf(`<dd>%s</dd>`, contentHTML)

			// Replace with native elements
			container.ReplaceWithHtml(dtHTML + ddHTML)
		}
	})

	// Group consecutive dt/dd pairs into dl elements
	doc.Find("dt").Each(func(i int, dt *goquery.Selection) {
		// Check if this dt is not already in a dl
		if !dt.Parent().Is("dl") {
			// Create a new dl and move consecutive dt/dd pairs into it
			dl := `<dl></dl>`
			dt.BeforeHtml(dl)
			dlElement := dt.Prev()

			current := dt
			for current.Length() > 0 && (current.Is("dt") || current.Is("dd")) {
				next := current.Next()
				currentHTML, _ := goquery.OuterHtml(current)
				dlElement.AppendHtml(currentHTML)
				current.Remove()
				current = next
			}
		}
	})
}

// convertTitlesToHeadings converts title elements to proper headings
func convertTitlesToHeadings(doc *goquery.Document) {
	// Convert title elements to proper headings
	doc.Find("p.title-article-norm").Each(func(i int, title *goquery.Selection) {
		titleText := strings.TrimSpace(title.Text())

		// Look for subtitle
		subtitle := ""
		next := title.Next()
		if next.HasClass("eli-title") {
			subtitleElem := next.Find("p.stitle-article-norm")
			if subtitleElem.Length() > 0 {
				subtitle = strings.TrimSpace(subtitleElem.Text())
				next.Remove() // Remove the subtitle container
			}
		}

		// Create combined heading
		fullTitle := titleText
		if subtitle != "" {
			fullTitle = titleText + ": " + subtitle
		}

		anchorID := generateAnchorID(fullTitle)
		headingHTML := fmt.Sprintf(`<h1 id="%s">%s</h1>`, anchorID, fullTitle)
		title.ReplaceWithHtml(headingHTML)
	})

	// Convert annex and appendix titles
	doc.Find("p.title-annex-1, p.title-gr-seq-level-1").Each(func(i int, title *goquery.Selection) {
		titleText := strings.TrimSpace(title.Text())

		// Look for subtitle
		subtitle := ""
		if title.HasClass("title-annex-1") {
			next := title.Next()
			if next.HasClass("title-gr-seq-level-1") {
				boldText := next.Find("span.boldface").Text()
				if boldText != "" {
					subtitle = strings.TrimSpace(boldText)
				} else {
					subtitle = strings.TrimSpace(next.Text())
				}
				next.Remove() // Remove the subtitle container
			}
		} else if title.HasClass("title-gr-seq-level-1") {
			// For title-gr-seq-level-1, the subtitle is within the element itself
			boldText := title.Find("span.boldface").Text()
			if boldText != "" {
				subtitle = strings.TrimSpace(boldText)
			}
		}

		// Create combined heading
		fullTitle := titleText
		if subtitle != "" {
			fullTitle = titleText + ": " + subtitle
		}

		anchorID := generateAnchorID(fullTitle)
		headingHTML := fmt.Sprintf(`<h1 id="%s">%s</h1>`, anchorID, fullTitle)
		title.ReplaceWithHtml(headingHTML)
	})
}

// convertTOCToNativeNav converts table-based TOCs to native navigation
func convertTOCToNativeNav(doc *goquery.Document) {
	doc.Find("table").Each(func(i int, table *goquery.Selection) {
		// Check if this is a table of contents
		titleCell := table.Find("p.title-toc")
		if titleCell.Length() > 0 {
			// This is a TOC table, convert it to native nav
			navHTML := `<nav><h2>Table of Contents</h2><ol>`

			// Process each row
			table.Find("tr").Each(func(j int, row *goquery.Selection) {
				cells := row.Find("td")
				if cells.Length() == 2 {
					contentText := strings.TrimSpace(cells.Last().Text())

					// Skip the header row
					if contentText != "TABLE OF CONTENT" && contentText != "" {
						// Generate anchor ID
						anchorID := generateAnchorID(contentText)
						navHTML += fmt.Sprintf(`<li><a href="#%s">%s</a></li>`, anchorID, contentText)
					}
				}
			})

			navHTML += `</ol></nav>`

			// Replace the table with the new nav
			table.ReplaceWithHtml(navHTML)
		}
	})
}

// simplifyTables removes all table attributes and simplifies structure
func simplifyTables(doc *goquery.Document) {
	// Remove all table styling attributes
	doc.Find("table").Each(func(i int, table *goquery.Selection) {
		// Remove all styling attributes
		table.RemoveAttr("width")
		table.RemoveAttr("border")
		table.RemoveAttr("cellspacing")
		table.RemoveAttr("cellpadding")
	})

	// Remove colgroup and col elements
	doc.Find("colgroup, col").Remove()

	// Clean up table cells
	doc.Find("td, th").Each(func(i int, cell *goquery.Selection) {
		// Remove ALL table formatting attributes
		cell.RemoveAttr("width")
		cell.RemoveAttr("valign")
		cell.RemoveAttr("align")
		cell.RemoveAttr("style")
		cell.RemoveAttr("colspan")
		cell.RemoveAttr("rowspan")
		// Additional table-specific formatting attributes
		cell.RemoveAttr("bgcolor")
		cell.RemoveAttr("height")
		cell.RemoveAttr("nowrap")

		// Remove paragraph wrappers in cells if they contain simple text
		paragraphs := cell.Find("p")
		if paragraphs.Length() == 1 {
			pText := strings.TrimSpace(paragraphs.Text())
			cellText := strings.TrimSpace(cell.Text())
			// If the paragraph contains all the cell's text, unwrap it
			if pText == cellText {
				cell.SetText(pText)
			}
		}
	})
}

// removeRedundantElements removes empty elements and redundant wrappers
func removeRedundantElements(doc *goquery.Document) {
	// Remove empty paragraphs and divs
	doc.Find("p:empty, div:empty").Remove()

	// Remove <br> tags
	doc.Find("br").Remove()

	// Remove wrapper divs that don't add semantic value
	doc.Find("div.centered, div[style]").Each(func(i int, div *goquery.Selection) {
		// Move children up and remove the wrapper
		div.Children().Each(func(j int, child *goquery.Selection) {
			childHTML, _ := goquery.OuterHtml(child)
			div.BeforeHtml(childHTML)
		})
		// Also preserve any direct text content
		text := strings.TrimSpace(div.Text())
		if text != "" && div.Children().Length() == 0 {
			div.BeforeHtml(text)
		}
		div.Remove()
	})

	// Remove unnecessary div wrappers
	doc.Find("div.eli-title").Each(func(i int, div *goquery.Selection) {
		// Move children up and remove the wrapper
		div.Children().Each(func(j int, child *goquery.Selection) {
			childHTML, _ := goquery.OuterHtml(child)
			div.BeforeHtml(childHTML)
		})
		div.Remove()
	})
}

// stripStylingAttributes removes all styling-related attributes and classes
func stripStylingAttributes(doc *goquery.Document) {
	// Remove all class attributes (since we're not using any styling)
	doc.Find("*").Each(func(i int, elem *goquery.Selection) {
		elem.RemoveAttr("class")
	})

	// Clean up IDs - keep only meaningful ones
	doc.Find("*[id]").Each(func(i int, elem *goquery.Selection) {
		id := elem.AttrOr("id", "")
		// Remove UUID-style IDs that aren't meaningful
		if strings.Contains(id, "-") && len(id) > 20 {
			elem.RemoveAttr("id")
		}
	})
}

// consolidateContent consolidates and optimizes content structure
func consolidateContent(doc *goquery.Document) {
	// Convert modification references to simple text
	doc.Find("p.modref").Each(func(i int, modref *goquery.Selection) {
		link := modref.Find("a")
		if link.Length() > 0 {
			title := link.AttrOr("title", "")
			text := strings.TrimSpace(link.Text())

			// Extract amendment info from title and text
			amendmentInfo := title
			if amendmentInfo == "" {
				amendmentInfo = text
			}

			// Create simple text note
			noteText := fmt.Sprintf("Amendment: %s", amendmentInfo)
			modref.ReplaceWithHtml(fmt.Sprintf("<p><em>%s</em></p>", noteText))
		} else {
			modref.Remove()
		}
	})

	// Remove any remaining external links
	doc.Find("a[href^='http']").Each(func(i int, link *goquery.Selection) {
		// Keep the text content but remove the link
		text := link.Text()
		if text != "" {
			link.ReplaceWithHtml(text)
		} else {
			link.Remove()
		}
	})

	// Clean up whitespace in text content
	doc.Find("*").Each(func(i int, elem *goquery.Selection) {
		// Only process text nodes, not elements with children
		if elem.Children().Length() == 0 {
			text := elem.Text()
			if text != "" {
				// Normalize whitespace
				cleanText := regexp.MustCompile(`\s+`).ReplaceAllString(strings.TrimSpace(text), " ")
				if cleanText != text {
					elem.SetText(cleanText)
				}
			}
		}
	})
}

// removeUnnecessarySpans removes span tags that don't add semantic value
func removeUnnecessarySpans(doc *goquery.Document) {
	// Convert semantic spans to proper HTML elements
	doc.Find("span.superscript").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		s.ReplaceWithHtml(fmt.Sprintf("<sup>%s</sup>", text))
	})

	doc.Find("span.italics").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		s.ReplaceWithHtml(fmt.Sprintf("<em>%s</em>", text))
	})

	doc.Find("span.boldface").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		s.ReplaceWithHtml(fmt.Sprintf("<strong>%s</strong>", text))
	})

	// Remove spans that only contain simple text (no semantic value)
	// Do this after class-based spans are processed, since stripStylingAttributes removes classes
	doc.Find("span").Each(func(i int, s *goquery.Selection) {
		// Check if span has no meaningful attributes
		attrs := s.Get(0).Attr
		hasClass := false
		for _, attr := range attrs {
			if attr.Key == "class" && attr.Val != "" {
				hasClass = true
				break
			}
		}

		// If span has no class or only empty attributes, unwrap it
		if !hasClass {
			text := s.Text()
			if text != "" {
				s.ReplaceWithHtml(text)
			} else {
				s.Remove()
			}
		}
	})
}

// removeRedundantContainers removes div wrappers that don't add semantic value
func removeRedundantContainers(doc *goquery.Document) {
	// Remove divs that only wrap a single paragraph
	doc.Find("div").Each(func(i int, div *goquery.Selection) {
		children := div.Children()
		if children.Length() == 1 && children.First().Is("p") {
			// Move the paragraph up and remove the div wrapper
			p := children.First()
			pHTML, _ := goquery.OuterHtml(p)
			div.ReplaceWithHtml(pHTML)
		}
	})

	// Remove nested divs with same or similar classes (this will run before class removal)
	doc.Find("div.norm div.norm, div.inline-element div.inline-element").Each(func(i int, innerDiv *goquery.Selection) {
		// Unwrap the inner div by moving its children up
		innerDiv.Children().Each(func(j int, child *goquery.Selection) {
			childHTML, _ := goquery.OuterHtml(child)
			innerDiv.BeforeHtml(childHTML)
		})
		// Also preserve any direct text content
		text := strings.TrimSpace(innerDiv.Text())
		if text != "" && innerDiv.Children().Length() == 0 {
			innerDiv.BeforeHtml(text)
		}
		innerDiv.Remove()
	})
}

// generateAnchorID creates a clean anchor ID from text
func generateAnchorID(text string) string {
	// Convert to lowercase and replace spaces/special chars with hyphens
	id := strings.ToLower(text)
	id = regexp.MustCompile(`[^a-z0-9\s-]`).ReplaceAllString(id, "")
	id = regexp.MustCompile(`\s+`).ReplaceAllString(id, "-")
	id = strings.Trim(id, "-")

	// Limit length
	if len(id) > 50 {
		id = id[:50]
		id = strings.Trim(id, "-")
	}

	return id
}
