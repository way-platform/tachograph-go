package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/way-platform/tachograph-go/internal/card"
	"github.com/way-platform/tachograph-go/internal/hexdump"
	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

var (
	inputDir   = flag.String("i", "", "Input directory containing .DDD files")
	outputDir  = flag.String("o", "internal/card/testdata/records", "Output directory for extracted hexdump files")
	startIndex = flag.Int("start", 0, "Starting index for output directory numbering")
)

func main() {
	flag.Parse()

	// Support both -i directory and positional file args
	var inputFiles []string
	if *inputDir != "" {
		if info, err := os.Stat(*inputDir); err != nil || !info.IsDir() {
			log.Fatalf("Input directory does not exist or is not a directory: %s", *inputDir)
		}
		err := filepath.WalkDir(*inputDir, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() && strings.HasSuffix(strings.ToUpper(d.Name()), ".DDD") {
				inputFiles = append(inputFiles, path)
			}
			return nil
		})
		if err != nil {
			log.Fatalf("Error walking directory: %v", err)
		}
	} else {
		for _, p := range flag.Args() {
			if strings.HasSuffix(strings.ToUpper(p), ".DDD") {
				inputFiles = append(inputFiles, p)
			} else {
				log.Printf("Skipping non-.DDD file: %s", p)
			}
		}
	}

	if len(inputFiles) == 0 {
		log.Fatal("No .DDD files found: provide -i <dir> or positional file args")
	}

	for i, path := range inputFiles {
		if err := processCardFile(path, *startIndex+i); err != nil {
			log.Printf("Warning: failed to process %s: %v", path, err)
		}
	}
}

func processCardFile(filePath string, fileIndex int) error {
	log.Printf("Processing [%03d]: %s", fileIndex, filePath)

	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Unmarshal to RawCardFile
	unmarshalOpts := card.UnmarshalOptions{}
	rawFile, err := unmarshalOpts.UnmarshalRawCardFile(data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal card file: %w", err)
	}

	// Calculate output directory paths
	// Get just the filename without extension
	baseName := filepath.Base(filePath)
	baseNameWithoutExt := strings.TrimSuffix(baseName, filepath.Ext(baseName))

	// Original directory: NNN-<filename>
	originalDir := filepath.Join(*outputDir, fmt.Sprintf("%03d-%s", fileIndex, baseNameWithoutExt))

	// Infer card type to determine if we should anonymize
	cardType := card.InferFileType(rawFile)
	isDriverCard := cardType == cardv1.CardType_DRIVER_CARD

	// Write original hexdumps
	log.Printf("  Writing original records to: %s", originalDir)
	if err := writeRecordsToDir(originalDir, rawFile.GetRecords()); err != nil {
		return fmt.Errorf("failed to write original records: %w", err)
	}

	// For driver cards, create anonymized version
	if isDriverCard {
		anonymizedDir := filepath.Join(*outputDir, fmt.Sprintf("%03d-anonymized", fileIndex))
		log.Printf("  Creating anonymized version...")

		// Parse RawCardFile → DriverCardFile
		parseOpts := card.ParseOptions{PreserveRawData: true}
		driverFile, err := parseOpts.ParseRawDriverCardFile(rawFile)
		if err != nil {
			return fmt.Errorf("failed to parse driver card file: %w", err)
		}

		// Anonymize
		anonOpts := card.AnonymizeOptions{}
		anonFile, err := anonOpts.AnonymizeDriverCardFile(driverFile)
		if err != nil {
			return fmt.Errorf("failed to anonymize driver card file: %w", err)
		}

		// Unparse back to RawCardFile
		anonRawFile, err := card.UnparseDriverCardFile(anonFile)
		if err != nil {
			return fmt.Errorf("failed to unparse anonymized driver card file: %w", err)
		}

		// Write anonymized hexdumps
		log.Printf("  Writing anonymized records to: %s", anonymizedDir)
		if err := writeRecordsToDir(anonymizedDir, anonRawFile.GetRecords()); err != nil {
			return fmt.Errorf("failed to write anonymized records: %w", err)
		}
	}

	return nil
}

func writeRecordsToDir(dirPath string, records []*cardv1.RawCardFile_Record) error {
	// Create directory
	if err := os.MkdirAll(dirPath, 0o755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write each record as hexdump
	for i, record := range records {
		// Get enum string representations
		fileType := record.GetFile().String()
		generation := record.GetGeneration().String()
		contentType := record.GetContentType().String()

		// Format filename: NNN-<FILE>-<GENERATION>-<CONTENT_TYPE>.hexdump
		filename := fmt.Sprintf("%03d-%s-%s-%s.hexdump", i, fileType, generation, contentType)
		outputPath := filepath.Join(dirPath, filename)

		// Marshal to hexdump format
		hexdumpData, err := hexdump.Marshal(record.GetValue())
		if err != nil {
			return fmt.Errorf("failed to marshal record %d to hexdump: %w", i, err)
		}

		// Write hexdump to file
		if err := os.WriteFile(outputPath, hexdumpData, 0o644); err != nil {
			return fmt.Errorf("failed to write hexdump file %s: %w", outputPath, err)
		}
	}

	return nil
}
