package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// Represents a single field within a data type definition.
type Field struct {
	Name        string `json:"name"`
	ASN1Tag     string `json:"asn1_tag"`
	Description string `json:"description"`
	Size        string `json:"size"`
}

// Represents a complete data type definition.
type DataType struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Fields      []Field `json:"fields"`
}

func main() {
	inputFile := flag.String("input", "docs/regulation/09-appendix-1-data-dictionary.html", "Path to the input HTML file")
	outputFile := flag.String("output", "internal/gen/tachomodel.json", "Path to the output JSON file")
	flag.Parse()

	if err := os.MkdirAll(filepath.Dir(*outputFile), 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	content, err := os.ReadFile(*inputFile)
	if err != nil {
		log.Fatalf("Failed to read input file: %v", err)
	}

	doc, err := html.Parse(strings.NewReader(string(content)))
	if err != nil {
		log.Fatalf("Failed to parse HTML: %v", err)
	}

	dataTypes := parseDataTypes(doc)

	jsonData, err := json.MarshalIndent(dataTypes, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	if err := os.WriteFile(*outputFile, jsonData, 0644); err != nil {
		log.Fatalf("Failed to write output file: %v", err)
	}

	fmt.Printf("Successfully parsed %d data types and wrote to %s\n", len(dataTypes), *outputFile)
}

// Traverses the HTML document to find and parse data type definitions.
func parseDataTypes(n *html.Node) []DataType {
	var dataTypes []DataType
	var currentDataType *DataType

	// Regex to find the section titles like "2.1. ActivityChangeInfo"
	titleRegex := regexp.MustCompile(`^\s*(\d+\.\d+[a-z]?\.)\s*([a-zA-Z0-9]+)`)

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "p" {
			pClass := getAttribute(n, "class")
			if pClass == "title-gr-seq-level-3" {
				titleText := strings.TrimSpace(getText(n))
				matches := titleRegex.FindStringSubmatch(titleText)
				if len(matches) == 3 {
					if currentDataType != nil {
						dataTypes = append(dataTypes, *currentDataType)
					}
					currentDataType = &DataType{
						ID:   strings.TrimSuffix(matches[1], "."),
						Name: matches[2],
					}
				}
			} else if pClass == "norm" && currentDataType != nil && len(currentDataType.Fields) == 0 {
				// Assume the first <p class="norm"> after a title is the description
				if currentDataType.Description == "" {
					currentDataType.Description = strings.TrimSpace(getText(n))
				}
			}
		}

		if n.Type == html.ElementNode && n.Data == "table" && currentDataType != nil {
			// Check if this is a data definition table
			isDefTable := false
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.ElementNode && c.Data == "tbody" {
					// A simple heuristic: check for a header row with "Data element"
					headerText := strings.ToLower(getText(c.FirstChild))
					if strings.Contains(headerText, "data element") && strings.Contains(headerText, "length in bytes") {
						isDefTable = true
						break
					}
				}
			}

			if isDefTable {
				currentDataType.Fields = parseFieldsTable(n)
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

		f(n)

		// Append the last parsed data type
		if currentDataType != nil {
			dataTypes = append(dataTypes, *currentDataType)
		}

		return dataTypes
}

// Parses an HTML table node to extract field definitions.
func parseFieldsTable(n *html.Node) []Field {
	var fields []Field
	var tbody *html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "tbody" {
			tbody = c
			break
		}
	}

	if tbody == nil {
		return nil
	}

	// Skip the header row
	for tr := tbody.FirstChild.NextSibling; tr != nil; tr = tr.NextSibling {
		if tr.Type == html.ElementNode && tr.Data == "tr" {
			var cells []string
			for td := tr.FirstChild; td != nil; td = td.NextSibling {
				if td.Type == html.ElementNode && td.Data == "td" {
					cells = append(cells, strings.TrimSpace(getText(td)))
				}
			}
			if len(cells) >= 4 {
				field := Field{
					Name:        cells[0],
					ASN1Tag:     cells[1],
					Description: cells[2],
					Size:        cells[3],
				}
			fields = append(fields, field)
			}
		}
	}
	return fields
}

// Extracts all text from a node and its children.
func getText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	if n.Type != html.ElementNode {
		return ""
	}
	var buf strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		buf.WriteString(getText(c))
	}
	return buf.String()
}

// Helper to get an attribute value from a node.
func getAttribute(n *html.Node, key string) string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}
