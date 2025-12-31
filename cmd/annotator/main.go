package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "1.0.0"
var showVersion bool

var rootCmd = &cobra.Command{
	Use:   "annotator",
	Short: "Annotator scans OpenEdge 4GL class files for annotations",
	Long:  `Annotator recursively scans directories for .cls files and extracts annotations in JSON format.`,
	Run: func(cmd *cobra.Command, args []string) {
		if showVersion {
			fmt.Println(version)
			return
		}
		cmd.Help()
	},
}

func init() {
	rootCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "Print the version number")
	rootCmd.AddCommand(parseCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
