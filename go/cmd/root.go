/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	//	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "scan",
	Short: "Security scan analysis of a FIWARE Generic Enablers",
	Long: `This program searchs Docker Images vulnerabilities in the FIWARE Generic Enablers based on Anchore 
	and Clair tools and provide a set of best practices of a running instance of them based on a docker 
	compose file.

	If there is no arguments, the analysis will be over all the content defined in the enablers.json.
	In case that a specific FIWARE Generic Enabler is specified, the analysis will be developed only 
	on this component.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	/*
		Version:   "v0.1.0",
		Compiled:  time.Now(),
		Copyright: "(c) 2022 FIWARE Foundation, e.V.",
	*/
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.scan.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
