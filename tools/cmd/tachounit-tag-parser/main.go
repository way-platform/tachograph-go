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

// TachounitData is the root structure for the output JSON file.
type TachounitData struct {
	DownloadTypes   []Tag `json:"downloadTypes"`
	DataIdentifiers []Tag `json:"dataIdentifiers"`
}

func main() {
	var app7File, app8File, outputDir string
	flag.StringVar(&app7File, "app7", "docs/regulation/15-appendix-7-data-downloading-protocols.html", "Input Appendix 7 HTML file")
	flag.StringVar(&app8File, "app8", "docs/regulation/16-appendix-8-calibration-protocol.html", "Input Appendix 8 HTML file")
	flag.StringVar(&outputDir, "o", "tachounit", "Output directory for the JSON file")
	flag.Parse()

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Parse the two documents
	downloadTypes := parseAppendix7(app7File)
	dataIdentifiers := parseAppendix8(app8File)

	if len(downloadTypes) == 0 && len(dataIdentifiers) == 0 {
		log.Fatal("No tags were parsed from any appendix. Check HTML structure and selectors.")
	}

	// Marshal the data to JSON
	outputData := TachounitData{DownloadTypes: downloadTypes, DataIdentifiers: dataIdentifiers}
	jsonData, err := json.MarshalIndent(outputData, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Write the JSON file
	outputFile := filepath.Join(outputDir, "tags.json")
	if err := os.WriteFile(outputFile, jsonData, 0o644); err != nil {
		log.Fatalf("Failed to write output file %s: %v", outputFile, err)
	}

	fmt.Printf("Successfully generated %s with %d download types and %d data identifiers.\n", outputFile, len(downloadTypes), len(dataIdentifiers))
}

func parseAppendix7(inputFile string) []Tag {
	file, err := os.Open(inputFile)
	if err != nil {
		log.Printf("Warning: could not open %s: %v", inputFile, err)
		return nil
	}
	defer file.Close()

	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		log.Printf("Warning: could not parse %s: %v", inputFile, err)
		return nil
	}

	tags := make(map[uint64]Tag)
	nameCount := make(map[string]int) // Track name usage to handle duplicates

	// Find the table for TRTP values, likely after DDP_011
	doc.Find("p:contains('DDP_011')").Parent().NextAllFiltered("div.centered").First().Find("table tbody tr").Each(func(i int, row *goquery.Selection) {
		if i == 0 {
			return
		} // Skip header row
		cols := row.Find("td")
		if cols.Length() < 4 {
			return
		} // Need at least 4 columns

		descStr := strings.TrimSpace(cols.Eq(0).Find("p").Text())
		if descStr == "" {
			return
		}

		// Process TRTP values from columns 1, 2, 3 (generation 1, gen 2 v1, gen 2 v2)
		for j := 1; j < 4; j++ {
			idStr := strings.TrimSpace(cols.Eq(j).Find("p").Text())

			// Skip "Not used" entries
			if idStr == "Not used" || idStr == "" {
				continue
			}

			// Try to parse as hex (some might be decimal like "00")
			var id uint64
			var err error
			if len(idStr) <= 2 {
				// Try decimal first for short values like "00"
				if id, err = strconv.ParseUint(idStr, 10, 8); err != nil {
					// If decimal fails, try hex
					id, err = strconv.ParseUint(idStr, 16, 8)
				}
			} else {
				// For longer values, try hex first
				if id, err = strconv.ParseUint(idStr, 16, 8); err != nil {
					id, err = strconv.ParseUint(idStr, 10, 8)
				}
			}

			if err == nil {
				if _, exists := tags[id]; !exists {
					baseName := generateGoName(descStr, "TRTP")
					finalName := baseName

					// Check if this name already exists, if so append the hex ID
					nameCount[baseName]++
					if nameCount[baseName] > 1 {
						finalName = fmt.Sprintf("%s_0x%02X", baseName, id)
					}

					tags[id] = Tag{
						Name:        finalName,
						ID:          id,
						Description: descStr,
						Source:      "Appendix 7",
					}
				}
			}
		}
	})

	return sortAndSlice(tags)
}

func parseAppendix8(inputFile string) []Tag {
	file, err := os.Open(inputFile)
	if err != nil {
		log.Printf("Warning: could not open %s: %v", inputFile, err)
		return nil
	}
	defer file.Close()

	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		log.Printf("Warning: could not parse %s: %v", inputFile, err)
		return nil
	}

	tags := make(map[uint64]Tag)

	// Find the table for RDI values, likely after CPR_053
	doc.Find("p:contains('CPR_053')").Parent().NextAllFiltered("div.centered").First().Find("table tbody tr").Each(func(i int, row *goquery.Selection) {
		if i == 0 {
			return
		} // Skip header row
		cols := row.Find("td")
		if cols.Length() < 5 {
			return
		} // Need at least 5 columns

		idStr := strings.TrimSpace(cols.Eq(0).Find("p").Text())
		descStr := strings.TrimSpace(cols.Eq(2).Find("p").Text())

		if idStr == "" || descStr == "" {
			return
		}

		if id, err := strconv.ParseUint(idStr, 16, 16); err == nil {
			if _, exists := tags[id]; !exists {
				tags[id] = Tag{
					Name:        generateGoName(descStr, "RDI"),
					ID:          id,
					Description: descStr,
					Source:      "Appendix 8",
				}
			}
		}
	})

	return sortAndSlice(tags)
}

func sortAndSlice(tags map[uint64]Tag) []Tag {
	var sortedTags []Tag
	for _, tag := range tags {
		sortedTags = append(sortedTags, tag)
	}
	sort.Slice(sortedTags, func(i, j int) bool {
		return sortedTags[i].ID < sortedTags[j].ID
	})
	return sortedTags
}

func generateGoName(s, prefix string) string {
	s = strings.ReplaceAll(s, "/", " ")
	s = strings.ReplaceAll(s, "(", "")
	s = strings.ReplaceAll(s, ")", "")

	reg := regexp.MustCompile(`[^a-zA-Z0-9_ ]+`)
	s = reg.ReplaceAllString(s, "")

	parts := strings.Fields(s)
	for i, part := range parts {
		parts[i] = strings.Title(strings.ToLower(part))
	}

	name := strings.Join(parts, "")
	if prefix != "" {
		return prefix + "_" + name
	}
	return name
}
