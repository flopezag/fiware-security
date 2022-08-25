/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"security/scan/gomail"

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
		var files []string

		initialize()

		if len(args) == 0 {
			// It means that the check operation has not specified FIWARE GE, therefore we scan all the
			// FIWARE GEs described in the configuration file (enablers.json)
			fmt.Println("Scanning all FIWARE Generic Enablers...")
		} else if len(args) == 1 {
			// We have received a specific FIWARE GE to scan
			ge := args[0]
			fmt.Println("FIWARE GE to scan: " + ge)
			ParseJSON()
			var images []string = Search(ge, "Image")

			for j := 0; j < len(images); j++ {
				out := Filename(ge, images[j])

				// Step 0: Pull the docker image
				fmt.Print("\nPulling image... ")
				err := exec.Command("docker", "pull", images[j]).Run()
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(-1)
				} else {
					fmt.Print("Success\n\n")
				}

				fmt.Println(out)

				// Step 1: Anchore and Clair scan image
				out = Anchore(images[j], out)
				files = append(files, out)

				out = Clair(images[j], out)
				files = append(files, out)
			}

			var repositories []string = Search(ge, "Repository")
			for j := 0; j < len(repositories); j++ {
				out := FilenameFromUrl(ge, repositories[j])
				fmt.Println(out)

				out = Gitleaks(repositories[j], out)
				files = append(files, out)
			}

			var compose []string = Search(ge, "Compose")
			out := FilenameFromUrl(ge, compose[0])
			fmt.Println(out)
			out = Docker_bench_security(compose[0], out)
			files = append(files, out)

			// Send the files by email
			SendMail(files)
		}

		clean()
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

func SendMail(files []string) {
	data := gomail.Data{
		ComponentName: "",
		Subject:       "[Security Analysis] Analysis of docker image: ",
		Body:          "",
		EmailTo:       "",
	}

	data.EmailTo = ""
	data.ComponentName = ""
	data.Subject = data.Subject + data.ComponentName

	body, err := gomail.GenerateBody(data)
	if err != nil {
		log.Fatalf("Error: %v", err)
	} else {
		data.Body = body
	}

	gomail.OAuthGmailService()

	status, err := gomail.SendEmailOAuth2(data)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	if status {
		log.Println("Email sent successfully using OAUTH")
	}

	status, err = gomail.SendEmailOAuth2WithAttachment(data, files)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	if status {
		log.Println("Email with attachments sent successfully using OAUTH")
	}
}
