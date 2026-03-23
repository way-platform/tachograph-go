package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/way-platform/tachograph-go"
	"google.golang.org/protobuf/encoding/protojson"
)

// NewCommand creates the root tachograph CLI command tree.
func NewCommand(opts ...Option) *cobra.Command {
	cfg := config{}
	for _, opt := range opts {
		opt(&cfg)
	}

	cmd := &cobra.Command{
		Use:   "tachograph",
		Short: "Tachograph CLI",
	}
	cmd.AddGroup(&cobra.Group{ID: "ddd", Title: ".DDD Files"})
	cmd.AddCommand(newParseCommand())
	cmd.AddGroup(&cobra.Group{ID: "utils", Title: "Utils"})
	cmd.SetHelpCommandGroupID("utils")
	cmd.SetCompletionCommandGroupID("utils")
	return cmd
}

func newParseCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "parse [file ...]",
		Short:   "Parse .DDD files",
		GroupID: "ddd",
		Args:    cobra.MinimumNArgs(1),
	}

	raw := cmd.Flags().Bool("raw", false, "Output raw intermediate format (skip semantic parsing)")
	authenticate := cmd.Flags().Bool("authenticate", false, "Authenticate signatures and certificates")
	strict := cmd.Flags().Bool("strict", true, "Error on unrecognized tags (default true)")
	preserveRawData := cmd.Flags().Bool("preserve-raw-data", true, "Store raw bytes for round-trip fidelity (default true)")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		for _, filename := range args {
			data, err := os.ReadFile(filename)
			if err != nil {
				return fmt.Errorf("error reading %s: %w", filename, err)
			}

			// Step 1: Unmarshal to raw format
			unmarshalOpts := tachograph.UnmarshalOptions{
				Strict: *strict,
			}
			rawFile, err := unmarshalOpts.Unmarshal(data)
			if err != nil {
				return fmt.Errorf("error parsing raw %s: %w", filename, err)
			}

			// Step 2: Optionally authenticate (works on raw files)
			if *authenticate {
				authOpts := tachograph.AuthenticateOptions{
					Mutate: true, // Mutate for CLI efficiency
				}
				rawFile, err = authOpts.Authenticate(ctx, rawFile)
				if err != nil {
					return fmt.Errorf("error authenticating %s: %w", filename, err)
				}
			}

			// Step 3: Output raw or parse to semantic format
			if *raw {
				// Output raw format (with or without authentication)
				fmt.Println(protojson.Format(rawFile))
			} else {
				// Parse to semantic format (authentication results are propagated)
				parseOpts := tachograph.ParseOptions{
					PreserveRawData: *preserveRawData,
				}
				file, err := parseOpts.Parse(rawFile)
				if err != nil {
					return fmt.Errorf("error parsing %s: %w", filename, err)
				}
				fmt.Println(protojson.Format(file))
			}
		}
		return nil
	}
	return cmd
}
