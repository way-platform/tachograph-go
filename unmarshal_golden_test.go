package tachograph

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/encoding/protojson"
)

var update = flag.Bool("update", false, "update golden files")

func TestUnmarshalFile_golden(t *testing.T) {
	if err := filepath.WalkDir("testdata", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".DDD") {
			return nil
		}
		goldenFile := strings.TrimSuffix(path, ".DDD") + ".json"
		t.Run(path, func(t *testing.T) {
			// Read and parse the DDD file
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("Failed to read DDD file %s: %v", path, err)
			}

			var actual string
			file, err := UnmarshalFile(data)
			if err != nil {
				// If parsing fails, use the error message as the golden content
				actual = "ERROR: " + err.Error()
			} else {
				// If parsing succeeds, use the JSON representation
				actual = protojson.Format(file)
			}

			if *update {
				if err := os.WriteFile(goldenFile, []byte(actual), 0o644); err != nil {
					t.Fatalf("Failed to write golden file %s: %v", goldenFile, err)
				}
				t.Logf("Updated golden file: %s", goldenFile)
				return
			}
			expected, err := os.ReadFile(goldenFile)
			if err != nil {
				if os.IsNotExist(err) {
					t.Fatalf("Golden file %s does not exist. Run with -update to create it.", goldenFile)
				}
				t.Fatalf("Failed to read golden file %s: %v", goldenFile, err)
			}
			if diff := cmp.Diff(string(expected), actual); diff != "" {
				t.Errorf("Golden file mismatch for %s (-expected +actual):\n%s", path, diff)
				t.Logf("To update the golden file, run: go test -update")
			}
		})
		return nil
	}); err != nil {
		t.Fatalf("Failed to walk testdata directory: %v", err)
	}
}
