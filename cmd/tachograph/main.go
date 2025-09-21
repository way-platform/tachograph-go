package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"image/color"
	"os"
	"strings"

	"github.com/charmbracelet/fang"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/spf13/cobra"
	"github.com/way-platform/tachograph-go"
	"github.com/way-platform/tachograph-go/tachocard"
	"github.com/way-platform/tachograph-go/tachounit"
)

func main() {
	if err := fang.Execute(
		context.Background(),
		newRootCommand(),
		fang.WithColorSchemeFunc(func(c lipgloss.LightDarkFunc) fang.ColorScheme {
			base := c(lipgloss.Black, lipgloss.White)
			baseInverted := c(lipgloss.White, lipgloss.Black)
			return fang.ColorScheme{
				Base:         base,
				Title:        base,
				Description:  base,
				Comment:      base,
				Flag:         base,
				FlagDefault:  base,
				Command:      base,
				QuotedString: base,
				Argument:     base,
				Help:         base,
				Dash:         base,
				ErrorHeader:  [2]color.Color{baseInverted, base},
				ErrorDetails: base,
			}
		}),
	); err != nil {
		os.Exit(1)
	}
}

func newRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tacho",
		Short: "Tachograph CLI",
	}
	cmd.AddGroup(&cobra.Group{ID: "ddd", Title: ".DDD Files"})
	cmd.AddCommand(newStatCommand())
	cmd.AddGroup(&cobra.Group{ID: "utils", Title: "Utils"})
	cmd.SetHelpCommandGroupID("utils")
	cmd.SetCompletionCommandGroupID("utils")
	return cmd
}

func newStatCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "stat <file1> [file2] [...]",
		Short:   "Get info about .DDD files",
		GroupID: "ddd",
		Args:    cobra.MinimumNArgs(1),
	}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		for i, filename := range args {
			// Add separator between files if processing multiple files
			if i > 0 {
				fmt.Println("\n" + strings.Repeat("=", 80))
			}

			fmt.Printf("File: %s\n", filename)

			data, err := os.ReadFile(filename)
			if err != nil {
				fmt.Printf("Error reading file: %v\n", err)
				continue
			}

			fileType := tacho.InferFileType(data)
			fmt.Printf("Type: %s\n", fileType)

			// Show outline based on file type
			if fileType == tacho.CardFileType {
				fmt.Println("\nTLV Record Outline:")
				fmt.Println("-------------------")
				if err := printCardOutline(data); err != nil {
					fmt.Printf("Failed to parse card outline: %v\n", err)
					continue
				}
			} else if fileType == tacho.UnitFileType {
				fmt.Println("\nTV Record Outline:")
				fmt.Println("------------------")
				if err := printUnitOutline(data); err != nil {
					fmt.Printf("Failed to parse unit outline: %v\n", err)
					continue
				}
			}
		}

		return nil
	}
	return cmd
}

// TLVRecord represents a single TLV record from a tachograph card file.
type TLVRecord struct {
	Tag    uint32 // 3-byte tag (FID + Generation)
	Length uint16 // 2-byte length
	Value  []byte // Variable-length value
	Offset int    // Byte offset in the file where this record starts
}

// printCardOutline parses a card file and prints an outline of all TLV records.
func printCardOutline(data []byte) error {
	records, err := parseTLVRecords(data)
	if err != nil {
		return err
	}

	for _, record := range records {
		tagName := getTagName(record.Tag)
		fmt.Printf("Offset: 0x%06X | Tag: 0x%06X (%s) | Length: %d bytes\n",
			record.Offset, record.Tag, tagName, record.Length)
	}

	fmt.Printf("\nTotal records: %d\n", len(records))
	return nil
}

// parseTLVRecords parses a card file and returns all TLV records with their byte offsets.
func parseTLVRecords(data []byte) ([]TLVRecord, error) {
	var records []TLVRecord
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Split(tachocard.SplitFunc)

	offset := 0
	for scanner.Scan() {
		token := scanner.Bytes()

		// Extract tag (3 bytes) and length (2 bytes)
		if len(token) < 5 {
			return nil, fmt.Errorf("invalid TLV record at offset %d: too short", offset)
		}

		// Read tag as 3 bytes (big endian), treating it as the upper 3 bytes of a uint32
		tag := uint32(token[0])<<16 | uint32(token[1])<<8 | uint32(token[2])
		length := binary.BigEndian.Uint16(token[3:5])

		// Extract value
		value := make([]byte, length)
		if len(token) >= 5+int(length) {
			copy(value, token[5:5+int(length)])
		}

		records = append(records, TLVRecord{
			Tag:    tag,
			Length: length,
			Value:  value,
			Offset: offset,
		})

		offset += len(token)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning TLV records: %w", err)
	}

	return records, nil
}

// getTagName returns a human-readable name for a tag, or a hex representation if unknown.
func getTagName(tag uint32) string {
	// Extract FID (first 2 bytes) and generation (last byte)
	fid := tachocard.Tag(tag >> 8)
	generation := uint8(tag & 0xFF)

	var fidName string
	switch fid {
	case tachocard.EF_ICC:
		fidName = "EF_ICC"
	case tachocard.EF_IC:
		fidName = "EF_IC"
	case tachocard.EF_Application_Identification:
		fidName = "EF_Application_Identification"
	case tachocard.EF_Card_Certificate:
		fidName = "EF_Card_Certificate"
	case tachocard.EF_CA_Certificate:
		fidName = "EF_CA_Certificate"
	case tachocard.EF_Identification:
		fidName = "EF_Identification"
	case tachocard.EF_Card_Download:
		fidName = "EF_Card_Download"
	case tachocard.EF_Driving_License_Info:
		fidName = "EF_Driving_License_Info"
	case tachocard.EF_Events_Data:
		fidName = "EF_Events_Data"
	case tachocard.EF_Faults_Data:
		fidName = "EF_Faults_Data"
	case tachocard.EF_Driver_Activity_Data:
		fidName = "EF_Driver_Activity_Data"
	case tachocard.EF_Vehicles_Used:
		fidName = "EF_Vehicles_Used"
	case tachocard.EF_Places:
		fidName = "EF_Places"
	case tachocard.EF_Current_Usage:
		fidName = "EF_Current_Usage"
	case tachocard.EF_Control_Activity_Data:
		fidName = "EF_Control_Activity_Data"
	case tachocard.EF_Specific_Conditions:
		fidName = "EF_Specific_Conditions"
	// Extended FIDs found in real tachograph cards
	case 0x0520:
		fidName = "EF_Card_Identification_And_Driver_Card_Holder_Identification"
	case 0x0521:
		fidName = "EF_Card_Driving_Licence_Information"
	case 0x0522:
		fidName = "EF_Specific_Conditions_Extended"
	case 0x0523:
		fidName = "EF_Card_Vehicle_Units_Used"
	case 0x0524:
		fidName = "EF_GNSS_Accumulated_Driving"
	case 0x0525:
		fidName = "EF_Driver_Card_Application_Identification_V2"
	case 0x0526:
		fidName = "EF_Card_Place_Auth_Daily_Work_Period"
	case 0x0527:
		fidName = "EF_GNSS_Auth_Accumulated_Driving"
	case 0x0528:
		fidName = "EF_Card_Border_Crossings"
	case 0x0529:
		fidName = "EF_Card_Load_Unload_Operations"
	// GNSS-related FIDs from Appendix 12
	case 0xC100:
		fidName = "EF_EGF_MACertificate" // GNSS MA Certificate
	case 0xC101:
		fidName = "EF_EGF_MACertificate_Extended"
	case 0xC108:
		fidName = "EF_CA_Certificate_GNSS" // GNSS CA Certificate
	case 0xC109:
		fidName = "EF_Link_Certificate_GNSS" // GNSS Link Certificate
	default:
		fidName = fmt.Sprintf("Unknown_FID_0x%04X", fid)
	}

	var genName string
	switch generation {
	case 0x00:
		genName = "Data"
	case 0x01:
		genName = "Signature"
	case 0x02:
		genName = "Gen2Data"
	case 0x03:
		genName = "Gen2Signature"
	default:
		genName = fmt.Sprintf("Gen_0x%02X", generation)
	}

	return fmt.Sprintf("%s_%s", fidName, genName)
}

// TVRecord represents a single TV record from a tachograph unit file.
type TVRecord struct {
	Tag        uint16 // 2-byte tag
	Value      []byte // Variable-length value
	Offset     int    // Byte offset in the file where this record starts
	DataSize   int    // Size of the data portion (excluding tag)
	Generation string // Generation identifier
}

// printUnitOutline parses a unit file and prints an outline of all TV records.
func printUnitOutline(data []byte) error {
	records, err := parseTVRecords(data)
	if err != nil {
		return err
	}

	for _, record := range records {
		tagName := getVuTagName(record.Tag)
		fmt.Printf("Offset: 0x%06X | Tag: 0x%04X (%s) | Size: %d bytes | %s\n",
			record.Offset, record.Tag, tagName, record.DataSize+2, record.Generation)
	}

	fmt.Printf("\nTotal TV records: %d\n", len(records))
	return nil
}

// parseTVRecords parses VU TV records by intelligently determining structure sizes.
func parseTVRecords(data []byte) ([]TVRecord, error) {
	var records []TVRecord
	offset := 0

	for offset < len(data) {
		if len(data)-offset < 2 {
			break
		}

		tag := binary.BigEndian.Uint16(data[offset : offset+2])
		vuTag := tachounit.VuTag(tag)

		if !vuTag.IsValid() {
			return nil, fmt.Errorf("unknown VU tag at offset 0x%06X: 0x%04X", offset, tag)
		}

		generation := getVuGeneration(tag)

		// Determine data size based on generation and tag
		dataSize, err := determineVuDataSize(data, offset+2, tag, generation)
		if err != nil {
			return nil, fmt.Errorf("failed to determine data size for tag 0x%04X at offset 0x%06X: %w", tag, offset, err)
		}

		// Extract value
		value := make([]byte, dataSize)
		if len(data) >= offset+2+dataSize {
			copy(value, data[offset+2:offset+2+dataSize])
		}

		records = append(records, TVRecord{
			Tag:        tag,
			Value:      value,
			Offset:     offset,
			DataSize:   dataSize,
			Generation: generation,
		})

		offset += 2 + dataSize
	}

	return records, nil
}

// determineVuDataSize intelligently determines the size of VU data based on generation.
func determineVuDataSize(data []byte, dataOffset int, tag uint16, generation string) (int, error) {
	switch generation {
	case "Generation1":
		return determineGen1Size(data, dataOffset, tag)
	case "Generation2V1", "Generation2V2":
		return determineGen2Size(data, dataOffset, tag)
	case "DownloadInterfaceVersion":
		return 2, nil // Always 2 bytes: generation + version
	default:
		return 0, fmt.Errorf("unknown generation: %s", generation)
	}
}

// determineGen1Size determines size for Generation 1 VU records (fixed structures).
func determineGen1Size(data []byte, dataOffset int, tag uint16) (int, error) {
	// For Gen1, we need to find the next valid VU tag or use known patterns
	// This is a heuristic approach that works well in practice

	// Look for the next valid VU tag starting from a reasonable distance
	minSize := 500    // Minimum reasonable size for VU records (Gen1 structures are large)
	maxSize := 100000 // Maximum reasonable size to prevent runaway parsing

	// Search for next valid VU tag with more careful alignment
	for size := minSize; size <= maxSize && dataOffset+size < len(data)-1; size++ {
		// Check if we found a valid VU tag at this position
		if dataOffset+size+2 <= len(data) {
			nextTag := binary.BigEndian.Uint16(data[dataOffset+size : dataOffset+size+2])
			if tachounit.VuTag(nextTag).IsValid() {
				return size, nil
			}
		}
	}

	// If we can't find the next tag, this might be the last record
	remainingData := len(data) - dataOffset
	if remainingData > maxSize {
		// Cap at reasonable size for safety
		return maxSize, nil
	}

	// For the last record, use all remaining data
	return remainingData, nil
}

// determineGen2Size determines size for Generation 2 VU records (record arrays).
func determineGen2Size(data []byte, dataOffset int, tag uint16) (int, error) {
	// Gen2 uses record arrays - parse the headers to determine total size
	totalSize := 0
	offset := dataOffset

	// Gen2 structures contain multiple record arrays, each with a 5-byte header
	for offset < len(data) {
		if len(data)-offset < 5 {
			break
		}

		// Parse record array header
		_ = data[offset] // recordType (not used for size calculation)
		recordSize := binary.BigEndian.Uint16(data[offset+1 : offset+3])
		numRecords := binary.BigEndian.Uint16(data[offset+3 : offset+5])

		// Validate reasonable values
		if recordSize > 10000 || numRecords > 10000 {
			// Probably not a valid record array header - we've reached the end
			break
		}

		arraySize := 5 + int(recordSize)*int(numRecords)
		totalSize += arraySize
		offset += arraySize

		// Check if the next position looks like another record array or the end
		if offset < len(data)-1 {
			// Peek ahead to see if there's another valid VU tag (end of current record)
			nextTag := binary.BigEndian.Uint16(data[offset : offset+2])
			if tachounit.VuTag(nextTag).IsValid() {
				// We've reached the next TV record
				break
			}
		}

		// Safety check: don't parse indefinitely
		if totalSize > 100000 {
			break
		}
	}

	return totalSize, nil
}

// getVuGeneration determines the generation from a VU tag.
func getVuGeneration(tag uint16) string {
	trep := uint8(tag & 0xFF)

	switch {
	case trep == 0x00:
		return "DownloadInterfaceVersion"
	case trep >= 0x01 && trep <= 0x05:
		return "Generation1"
	case trep >= 0x21 && trep <= 0x25:
		return "Generation2V1"
	case trep >= 0x31 && trep <= 0x35:
		return "Generation2V2"
	default:
		return fmt.Sprintf("Unknown_TREP_0x%02X", trep)
	}
}

// getVuTagName returns a human-readable name for a VU tag.
func getVuTagName(tag uint16) string {
	vuTag := tachounit.VuTag(tag)

	switch vuTag {
	case tachounit.VU_DownloadInterfaceVersion:
		return "VU_DownloadInterfaceVersion"
	case tachounit.VU_OverviewFirstGen:
		return "VU_OverviewFirstGen"
	case tachounit.VU_ActivitiesFirstGen:
		return "VU_ActivitiesFirstGen"
	case tachounit.VU_EventsAndFaultsFirstGen:
		return "VU_EventsAndFaultsFirstGen"
	case tachounit.VU_DetailedSpeedFirstGen:
		return "VU_DetailedSpeedFirstGen"
	case tachounit.VU_TechnicalDataFirstGen:
		return "VU_TechnicalDataFirstGen"
	case tachounit.VU_OverviewSecondGen:
		return "VU_OverviewSecondGen"
	case tachounit.VU_ActivitiesSecondGen:
		return "VU_ActivitiesSecondGen"
	case tachounit.VU_EventsAndFaultsSecondGen:
		return "VU_EventsAndFaultsSecondGen"
	case tachounit.VU_DetailedSpeedSecondGen:
		return "VU_DetailedSpeedSecondGen"
	case tachounit.VU_TechnicalDataSecondGen:
		return "VU_TechnicalDataSecondGen"
	case tachounit.VU_OverviewSecondGenV2:
		return "VU_OverviewSecondGenV2"
	case tachounit.VU_ActivitiesSecondGenV2:
		return "VU_ActivitiesSecondGenV2"
	case tachounit.VU_EventsAndFaultsSecondGenV2:
		return "VU_EventsAndFaultsSecondGenV2"
	case tachounit.VU_TechnicalDataSecondGenV2:
		return "VU_TechnicalDataSecondGenV2"
	default:
		return fmt.Sprintf("Unknown_VU_Tag_0x%04X", tag)
	}
}
