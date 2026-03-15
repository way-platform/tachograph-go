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
	outputDir  = flag.String("o", "internal/vu/testdata/records", "Output directory for extracted hexdump files")
	startIndex = flag.Int("start", 0, "Starting index for output directory numbering")
)

func main() {
	flag.Parse()

	inputFiles := flag.Args()
	if len(inputFiles) == 0 {
		log.Fatal("At least one input .DDD file is required as a positional argument")
	}

	for i, path := range inputFiles {
		if !strings.HasSuffix(strings.ToUpper(path), ".DDD") {
			log.Printf("Skipping non-.DDD file: %s", path)
			continue
		}
		if err := processVUFile(path, *startIndex+i); err != nil {
			log.Printf("Warning: failed to process %s: %v", path, err)
		}
	}
}

func processVUFile(filePath string, fileIndex int) error {
	log.Printf("Processing [%03d]: %s", fileIndex, filePath)

	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Unmarshal to RawVehicleUnitFile
	unmarshalOpts := vu.UnmarshalOptions{}
	rawFile, err := unmarshalOpts.UnmarshalRawVehicleUnitFile(data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal VU file: %w", err)
	}

	// Calculate output directory paths
	// Get just the filename without extension
	baseName := filepath.Base(filePath)
	baseNameWithoutExt := strings.TrimSuffix(baseName, filepath.Ext(baseName))

	// Original directory: NNN-<filename>
	originalDir := filepath.Join(*outputDir, fmt.Sprintf("%03d-%s", fileIndex, baseNameWithoutExt))

	// Write original hexdumps
	log.Printf("  Writing original records to: %s", originalDir)
	if err := writeRecordsToDir(originalDir, rawFile.GetRecords()); err != nil {
		return fmt.Errorf("failed to write original records: %w", err)
	}

	// Create anonymized version
	anonymizedDir := filepath.Join(*outputDir, fmt.Sprintf("%03d-anonymized", fileIndex))
	log.Printf("  Creating anonymized version...")

	// Parse RawVehicleUnitFile → VehicleUnitFile
	parseOpts := vu.ParseOptions{PreserveRawData: true}
	vuFile, err := parseOpts.ParseRawVehicleUnitFile(rawFile)
	if err != nil {
		return fmt.Errorf("failed to parse VU file: %w", err)
	}

	// Anonymize
	anonOpts := vu.AnonymizeOptions{}
	anonFile, err := anonOpts.AnonymizeVehicleUnitFile(vuFile)
	if err != nil {
		return fmt.Errorf("failed to anonymize VU file: %w", err)
	}

	// Unparse back to RawVehicleUnitFile
	anonRawFile, err := vu.UnparseVehicleUnitFile(anonFile)
	if err != nil {
		return fmt.Errorf("failed to unparse anonymized VU file: %w", err)
	}

	// Write anonymized hexdumps
	log.Printf("  Writing anonymized records to: %s", anonymizedDir)
	if err := writeRecordsToDir(anonymizedDir, anonRawFile.GetRecords()); err != nil {
		return fmt.Errorf("failed to write anonymized records: %w", err)
	}

	return nil
}

func writeRecordsToDir(dirPath string, records []*vuv1.RawVehicleUnitFile_Record) error {
	// Create directory
	if err := os.MkdirAll(dirPath, 0o755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write each record as a complete transfer (data + signature) hexdump
	for i, record := range records {
		// Get transfer type enum string representation
		transferType := record.GetType().String()

		// Get complete transfer value (already includes signature)
		transferValue := record.GetValue()

		// Write complete transfer: NNN-<TRANSFER_TYPE>.hexdump
		filename := fmt.Sprintf("%03d-%s.hexdump", i, transferType)
		outputPath := filepath.Join(dirPath, filename)

		hexdumpData, err := hexdump.Marshal(transferValue)
		if err != nil {
			return fmt.Errorf("failed to marshal transfer %d to hexdump: %w", i, err)
		}

		if err := os.WriteFile(outputPath, hexdumpData, 0o644); err != nil {
			return fmt.Errorf("failed to write hexdump file %s: %w", outputPath, err)
		}
	}

	return nil
}
