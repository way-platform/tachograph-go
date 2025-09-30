package card

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

func TestInferCardFileType(t *testing.T) {
	expectedCardTypes := map[string]cardv1.CardType{
		"driver":   cardv1.CardType_DRIVER_CARD,
		"workshop": cardv1.CardType_WORKSHOP_CARD,
		"control":  cardv1.CardType_CONTROL_CARD,
		"company":  cardv1.CardType_COMPANY_CARD,
	}
	if err := filepath.WalkDir("../../testdata/card", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if strings.ToLower(filepath.Ext(path)) != ".ddd" {
			return nil
		}
		// Extract card type from directory structure (testdata/card/driver/file.DDD)
		rel, _ := filepath.Rel("../../testdata/card", path)
		cardTypeName := strings.Split(rel, string(filepath.Separator))[0]
		expectedCardType, exists := expectedCardTypes[cardTypeName]
		if !exists {
			t.Fatalf("Unknown card type: %s", cardTypeName)
			return nil
		}
		t.Run(rel, func(t *testing.T) {
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("Failed to read file: %v", err)
			}
			rawCardFile, err := UnmarshalRawCardFile(data)
			if err != nil {
				t.Fatalf("Failed to unmarshal: %v", err)
			}
			inferredCardType := InferCardFileType(rawCardFile)
			if inferredCardType != expectedCardType {
				t.Errorf("Expected %s, got %s (%d records)",
					expectedCardType, inferredCardType, len(rawCardFile.GetRecords()))
			}
		})
		return nil
	}); err != nil {
		t.Fatalf("Failed to walk directory: %v", err)
	}
}
