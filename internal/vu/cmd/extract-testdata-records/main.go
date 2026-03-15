// extract-testdata-records is a one-time tool for extracting anonymized VU transfer
// records from real DDD files into hexdump fixtures for use in unit tests.
//
// Usage:
//
//	go run ./internal/vu/cmd/extract-testdata-records/ [-o dir] [-max-activities N] [-max-speed N] <file.DDD> [...]
//
// Each input file (index N) produces a subdirectory named "NNN-<basename>/" under -o.
// Within it, each extracted transfer is written as "NNN-<TRANSFER_TYPE>.hexdump".
//
// Records are processed individually so mixed-generation files (Gen2 V1 + V2 in the
// same file) are handled correctly.  ACTIVITIES_GEN2_V2 and DETAILED_SPEED_GEN2 are
// capped at -max-activities / -max-speed to keep fixture size manageable.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/way-platform/tachograph-go/internal/hexdump"
	"github.com/way-platform/tachograph-go/internal/vu"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

var (
	outputDir        = flag.String("o", "internal/vu/testdata/records", "Output directory for hexdump files")
	maxActivities    = flag.Int("max-activities", 2, "Max ACTIVITIES_GEN2_V2 records per input file")
	maxDetailedSpeed = flag.Int("max-speed", 2, "Max DETAILED_SPEED_GEN2 records per input file")
)

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		log.Fatal("Usage: extract-testdata-records [-o dir] <file.DDD> [...]")
	}

	for i, dddPath := range args {
		if err := processFile(dddPath, i); err != nil {
			log.Fatalf("Failed to process %s: %v", dddPath, err)
		}
	}
}

func processFile(filePath string, fileIndex int) error {
	log.Printf("Processing [%03d]: %s", fileIndex, filePath)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	// Unmarshal all raw records (Strict: false to tolerate unknown tags).
	rawFile, err := vu.UnmarshalOptions{Strict: false}.UnmarshalRawVehicleUnitFile(data)
	if err != nil {
		return fmt.Errorf("unmarshal raw VU file: %w", err)
	}

	baseName := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
	outDir := filepath.Join(*outputDir, fmt.Sprintf("%03d-%s", fileIndex, baseName))

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	// Track how many records of each limited type we have written.
	typeCounts := make(map[vuv1.TransferType]int)
	recordIndex := 0

	for _, record := range rawFile.GetRecords() {
		tt := record.GetType()

		// Apply per-type limits; skip Gen1, CARD_DOWNLOAD, and unknown types.
		switch tt {
		case vuv1.TransferType_OVERVIEW_GEN2_V1,
			vuv1.TransferType_EVENTS_AND_FAULTS_GEN2_V1,
			vuv1.TransferType_TECHNICAL_DATA_GEN2_V1,
			vuv1.TransferType_OVERVIEW_GEN2_V2,
			vuv1.TransferType_EVENTS_AND_FAULTS_GEN2_V2,
			vuv1.TransferType_TECHNICAL_DATA_GEN2_V2:
			// No per-type limit.
		case vuv1.TransferType_ACTIVITIES_GEN2_V2:
			if typeCounts[tt] >= *maxActivities {
				continue
			}
		case vuv1.TransferType_DETAILED_SPEED_GEN2:
			if typeCounts[tt] >= *maxDetailedSpeed {
				continue
			}
		default:
			continue
		}

		anonValue, err := processRecord(record)
		if err != nil {
			log.Printf("  [%v] skip (error): %v", tt, err)
			continue
		}

		hexdumpData, err := hexdump.Marshal(anonValue)
		if err != nil {
			return fmt.Errorf("hexdump marshal for %v: %w", tt, err)
		}

		filename := fmt.Sprintf("%03d-%s.hexdump", recordIndex, tt.String())
		outPath := filepath.Join(outDir, filename)

		if err := os.WriteFile(outPath, hexdumpData, 0o644); err != nil {
			return fmt.Errorf("write hexdump %s: %w", outPath, err)
		}

		log.Printf("  Wrote %s (%d bytes)", filename, len(anonValue))
		typeCounts[tt]++
		recordIndex++
	}

	log.Printf("  Done: %d records written to %s", recordIndex, outDir)
	return nil
}

// processRecord parses, anonymizes, and marshals a single raw VU transfer record,
// returning the anonymized binary transfer value (without the 2-byte TV tag).
//
// Each record is processed in isolation via a mini single-record RawVehicleUnitFile so
// that mixed-generation files (e.g. OVERVIEW_GEN2_V1 alongside OVERVIEW_GEN2_V2) do not
// cause a parse error.
func processRecord(record *vuv1.RawVehicleUnitFile_Record) ([]byte, error) {
	miniRaw := &vuv1.RawVehicleUnitFile{}
	miniRaw.SetRecords([]*vuv1.RawVehicleUnitFile_Record{record})

	parsed, err := vu.ParseOptions{PreserveRawData: false}.ParseRawVehicleUnitFile(miniRaw)
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}

	anon, err := vu.AnonymizeOptions{}.AnonymizeVehicleUnitFile(parsed)
	if err != nil {
		return nil, fmt.Errorf("anonymize: %w", err)
	}

	anonRaw, err := vu.UnparseVehicleUnitFile(anon)
	if err != nil {
		return nil, fmt.Errorf("unparse: %w", err)
	}
	if len(anonRaw.GetRecords()) != 1 {
		return nil, fmt.Errorf("expected 1 unparsed record, got %d", len(anonRaw.GetRecords()))
	}

	return anonRaw.GetRecords()[0].GetValue(), nil
}
