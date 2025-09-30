//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"os"

	"github.com/way-platform/tachograph-go/internal/card"
	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

func main() {
	data, err := os.ReadFile("../../testdata/card/driver/proprietary-Nuutti_Nestori_Sahala_2025-09-11_08-38-24.DDD")
	if err != nil {
		fmt.Printf("Failed to read file: %v\n", err)
		os.Exit(1)
	}

	// Step 1: Parse as raw card
	rawCard, err := card.UnmarshalRawCardFile(data)
	if err != nil {
		fmt.Printf("Failed to unmarshal raw card: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Raw card has %d records\n\n", len(rawCard.GetRecords()))

	// Now set the File fields to see which one is EF_IDENTIFICATION
	for i, record := range rawCard.GetRecords() {
		fid := uint16(record.GetTag() >> 8)
		fileType := card.MapFidToElementaryFileType(fid)
		record.SetFile(fileType)

		fmt.Printf("Record %d:\n", i)
		fmt.Printf("  Tag: 0x%06X (FID: 0x%04X)\n", record.GetTag(), fid)
		fmt.Printf("  File: %s\n", record.GetFile())
		fmt.Printf("  ContentType: %s\n", record.GetContentType())
		fmt.Printf("  Length: %d\n", record.GetLength())
		fmt.Printf("  Value len: %d\n", len(record.GetValue()))

		if record.GetFile() == cardv1.ElementaryFileType_EF_IDENTIFICATION {
			fmt.Printf("\n>>> Found EF_IDENTIFICATION!\n")
			fmt.Printf("    Length field says: %d bytes\n", record.GetLength())
			fmt.Printf("    Value actually has: %d bytes\n", len(record.GetValue()))
			fmt.Printf("    Need at least: 143 bytes\n")
		}
		fmt.Println()
	}
}
