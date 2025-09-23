package main

import (
	"context"
	"fmt"
	"image/color"
	"os"
	"strings"

	"github.com/charmbracelet/fang"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/spf13/cobra"
	"github.com/way-platform/tachograph-go"
	"google.golang.org/protobuf/encoding/protojson"
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
		Use:   "tachograph",
		Short: "Tachograph CLI",
	}
	cmd.AddGroup(&cobra.Group{ID: "ddd", Title: ".DDD Files"})
	cmd.AddCommand(newStatCommand())
	cmd.AddCommand(newParseCommand())
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

			fileType := tachograph.InferFileType(data)
			fmt.Printf("Type: %s\n", fileType)

			fmt.Printf("Size: %d bytes\n", len(data))
			
			// TODO: Re-implement outline parsing using the new structure
			if fileType == tachograph.CardFileType {
				fmt.Println("\nCard file detected. Use 'parse' command for detailed parsing.")
			} else if fileType == tachograph.UnitFileType {
				fmt.Println("\nUnit file detected. Unit parsing not yet implemented.")
			}
		}

		return nil
	}
	return cmd
}

func newParseCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "parse <file1> [file2] [...]",
		Short:   "Parse .DDD files and output as protojson",
		GroupID: "ddd",
		Args:    cobra.MinimumNArgs(1),
	}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		for i, filename := range args {
			// Add separator between files if processing multiple files
			if i > 0 {
				fmt.Println()
			}

			data, err := os.ReadFile(filename)
			if err != nil {
				return fmt.Errorf("error reading file %s: %w", filename, err)
			}

			// Parse the file using our new unmarshal functionality
			file, err := tachograph.Unmarshal(data)
			if err != nil {
				return fmt.Errorf("error parsing file %s: %w", filename, err)
			}

			// Format as protojson
			marshaler := protojson.MarshalOptions{
				Indent:          "  ",
				UseProtoNames:   true,
				EmitUnpopulated: false,
			}
			jsonData, err := marshaler.Marshal(file)
			if err != nil {
				return fmt.Errorf("error marshaling to JSON for file %s: %w", filename, err)
			}

			// Print filename as comment if processing multiple files
			if len(args) > 1 {
				fmt.Printf("// File: %s\n", filename)
			}
			fmt.Println(string(jsonData))
		}

		return nil
	}
	return cmd
}

