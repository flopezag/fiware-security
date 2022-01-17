/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Scan analysis of a FIWARE Generic Enabler",
	Long: `Operation to scan the FIWARE Generic Enabler. If there is no arguments, the analysis will
	be over all the content defined in the enablers.json. In case that a specific FIWARE Generic 
	Enabler is specified, the analysis will be developed only on this component.`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			// It means that the check operation has not specified FIWARE GE, therefore we scan all the
			// FIWARE GEs described in the configuration file (enablers.json)
			fmt.Println("Scanning all FIWARE Generic Enablers...")
		} else if len(args) == 1 {
			// We have received a specific FIWARE GE to scan
			ge := args[0]
			fmt.Println("FIWARE GE to scan: " + ge)
			ParseJSON()
			var images []string

			images = Search(ge)

			for j := 0; j < len(images); j++ {
				out := Filename(ge, images[j])
				Anchore(images[j], out)
			}
		}

		// Check the arguments to see if we want to check all GEs (no data after check command or only
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// checkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
