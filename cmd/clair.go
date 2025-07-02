package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

func Clair(enabler, filename string) string {
	var out bytes.Buffer
	var stderr bytes.Buffer

	start := time.Now()
	fmt.Println("\nStarting at: ", start)

	filename = filename + "_clair.json"

	fmt.Println("Clair CVE Security Scan... ")
	fmt.Println("    Docker image: ", enabler)
	fmt.Println("    Output file: ", filename)

	// Change to the Clair folder to execute the analysis
	err := os.Chdir("./Clair")
	CheckIfError(err)

	// Step 1:
	fmt.Print("    Security analysis of " + enabler + " image...")
	// TODO:     Security analysis of fiware/orion-ld:latest image...exit status 1: Creating network "clair_default" with the default driver
	cmd := exec.Command(absPathDockerCompose, "compose", "run", "--rm", "scanner", enabler)
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		os.Exit(-1)
	} else {
		fmt.Println("Success")
		fmt.Printf("         Result:\n%8v\n", out.String())

		// Step 2: Filtering the line from the result analysis
		result := out.Bytes()
		index := bytes.Index(result, []byte("latest: Pulling from arminc/clair-db"))
		if index != 0 {
			// We need to delete the 3 first lines of the result
			fmt.Println("        Deleting the 3 lines of the output")
		}

		// Step 3: Save the out into filename
		err = os.WriteFile(filename, result, 0644)

		if err != nil {
			log.Fatal(err)
		}
	}
	// #         # Filtering the line from the result analysis
	// #         line=$(grep 'latest: Pulling from arminc\/clair-db' ${filename})
	// #
	// #         # Just for the 1st time...
	// #         if [[ -n ${line} ]]; then
	// #             # Delete first 3 lines of the file due to the first time that it is executed
	// #             # it includes 3 extra no needed lines
	// #             sed -i '1,3 d' ${filename}
	// #         fi
	// #
	// #         # Just to finish, send the data to the nexus instance
	// #         redirect_all curl -fsSL -u ${user}':'${password} --upload-file ${filename}  https://nexus.lab.fiware.org/repository/security/check/${enabler}/sast_clair/${filename}
	// #
	// #         # Keep the name of the file to send afterward the email to the owner
	// #         filename_clair=$filename_clair$(pwd)"/"$filename","
	// #     done
	// #
	// #     # We need to remove the last "," from the filename string(s)
	// #     filename_clair=${filename_clair::-1}
	// # }
	// #

	fmt.Println("scan completed in ", time.Since(start), " seconds")

	// Return to the original folder
	err = os.Chdir("..")
	CheckIfError(err)

	return filename
}
