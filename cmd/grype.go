package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

func Grype(enabler, filename string) string {
	var out bytes.Buffer
	var stderr bytes.Buffer

	start := time.Now()
	fmt.Print("\nStarting at: " + start.String())
	fmt.Println()

	filename = filename + "_grype.json"

	fmt.Println("  Grype Security Scan... ")
	fmt.Println("    Docker image: ", enabler)
	fmt.Println("    Output file: ", filename)

	// Change to the Grype folder to execute the analysis
	err := os.Chdir("./Grype")
	CheckIfError(err)

	// Step 1: Executing grype command
	// grype -o json -q --fail-on high --output json $enabler > $filename
	// grype fiware/idm:latest --scope all-layers -o json > ./results/keyrock.json

	fmt.Print("    Analysing the FIWARE Enabler... ")
	cmd := exec.Command("grype", enabler, "--scope", "all-layers", "-o", "json")
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		os.Exit(-1)
	} else {
		fmt.Println("Success")
		// fmt.Printf("         Result:\n%8v\n", out.String())
	}

	// Step 4: Save the out into filename
	filename = "../results/" + filename
	fmt.Print("    Saving the file: ")
	err = os.WriteFile(filename, out.Bytes(), 0644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Success")


	// Just to finish, send the data to the nexus instance
	// redirect_all curl -fsSL -u ${user}':'${password} --upload-file ${filename}  https://nexus.lab.fiware.org/repository/security/check/${enabler}/sast_anchore/${filename}
	// # http -a ${user}:${password} https://nexus.lab.fiware.org/repository/security/check/${enabler}/sast_anchore/${filename} < ${filename}
	// #

	// Keep the name of the file to send afterward the email to the owner
	// filename_anchore=$filename_anchore$(pwd)"/"$filename","

	// #     # We need to remove the last "," from the filename string(s)
	// #     filename_anchore=${filename_anchore::-1}
	// #

	fmt.Println("\nScan completed in ", time.Since(start), " seconds")

	// Return to the original folder
	err = os.Chdir("..")
	CheckIfError(err)

	return filename
}
