package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var inputFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-scheme",
	Short: "Go-Scheme generates a json schema based on a json body",
	Long:  ``,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&inputFile, "input", "", "input file containing json reponse payload")
}
