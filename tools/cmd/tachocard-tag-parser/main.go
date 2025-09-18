package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Tag represents a single parsed tag definition.
type Tag struct {
	Name        string `json:"name"`
	ID          uint64 `json:"id"`
	Description string `json:"description"`
	Source      string `json:"source"`
}

// TachocardData is the root structure for the output JSON file.
type TachocardData struct {
	Files []Tag `json:"files"`
}

func main() {
	var inputFile string
	var outputDir string
	flag.StringVar(&inputFile, "i", "docs/regulation/10-appendix-2-tachograph-cards-specification.html", "Input HTML file to parse")
	flag.StringVar(&outputDir, "o", "tachocard", "Output directory for the JSON file")
	flag.Parse()

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatalf("Failed to open input file %s: %v", inputFile, err)
	}
	defer file.Close()

	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		log.Fatalf("Failed to parse HTML: %v", err)
	}

	// Use a map to avoid duplicate file IDs.
	tags := make(map[uint64]Tag)

	// Look for tables that contain file identifiers in the format we found
	doc.Find("table").Each(func(i int, table *goquery.Selection) {
		// Find rows that contain hex file identifiers
		table.Find("tr").Each(func(j int, row *goquery.Selection) {
			cols := row.Find("td")
			if cols.Length() < 2 {
				return
			}

			// Look for hex values in p.tbl-norm elements
			cols.Each(func(k int, col *goquery.Selection) {
				idStr := strings.TrimSpace(col.Find("p.tbl-norm").Text())

				// Check if this looks like a 4-character hex file identifier
				if len(idStr) == 4 {
					if id, err := strconv.ParseUint(idStr, 16, 16); err == nil {
						// Try to find a description in adjacent columns
						var descStr string

						// Look in the next column
						if k+1 < cols.Length() {
							descStr = strings.TrimSpace(cols.Eq(k + 1).Find("p").Text())
						}

						// Look in the previous column if next was empty
						if descStr == "" && k > 0 {
							descStr = strings.TrimSpace(cols.Eq(k - 1).Find("p").Text())
						}

						// Skip if no meaningful description found
						if descStr == "" || len(descStr) < 3 {
							return
						}

						// Skip common non-file-identifier values
						if descStr == "File ID" || descStr == "SFID" || strings.Contains(descStr, "Access") {
							return
						}

						if _, exists := tags[id]; !exists {
							tags[id] = Tag{
								Name:        generateGoName(descStr),
								ID:          id,
								Description: descStr,
								Source:      "Appendix 2",
							}
						}
					}
				}
			})
		})
	})

	if len(tags) == 0 {
		log.Fatal("No tags were parsed. Check the HTML structure and selectors.")
	}

	var sortedTags []Tag
	for _, tag := range tags {
		sortedTags = append(sortedTags, tag)
	}
	sort.Slice(sortedTags, func(i, j int) bool {
		return sortedTags[i].ID < sortedTags[j].ID
	})

	outputData := TachocardData{Files: sortedTags}
	jsonData, err := json.MarshalIndent(outputData, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	outputFile := filepath.Join(outputDir, "tags.json")
	if err := os.WriteFile(outputFile, jsonData, 0o644); err != nil {
		log.Fatalf("Failed to write output file %s: %v", outputFile, err)
	}

	fmt.Printf("Successfully generated %s with %d tags.\n", outputFile, len(sortedTags))
}

func generateGoName(s string) string {
	s = strings.ReplaceAll(s, "/", " ")
	s = strings.ReplaceAll(s, "(", "")
	s = strings.ReplaceAll(s, ")", "")

	reg := regexp.MustCompile(`[^a-zA-Z0-9_ ]+`)
	s = reg.ReplaceAllString(s, "")

	parts := strings.Fields(s)
	var resultParts []string
	for _, part := range parts {
		if strings.ToUpper(part) == "EF" || strings.ToUpper(part) == "DF" {
			resultParts = append(resultParts, strings.ToUpper(part))
		} else {
			resultParts = append(resultParts, strings.Title(strings.ToLower(part)))
		}
	}

	name := strings.Join(resultParts, "")

	if len(resultParts) > 1 && (resultParts[0] == "EF" || resultParts[0] == "DF") {
		name = resultParts[0] + "_" + strings.Join(resultParts[1:], "")
	}

	return name
}
