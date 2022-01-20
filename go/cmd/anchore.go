package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func Anchore(enabler, filename string) {
	var out bytes.Buffer
	var stderr bytes.Buffer

	filename = filename + "_anchore.json"

	fmt.Println("Anchore Security Scan... ")
	fmt.Println("    Docker image: ", enabler)
	fmt.Println("    Output file: ", filename)

	// Change to the Anchore folder to execute the analysis
	err := os.Chdir("./Anchore")
	CheckIfError(err)

	// Step 1: Add the FIWARE Enabler to be analysed
	// docker-compose -f docker-compose-anchore.yaml exec api anchore-cli image add $enabler
	fmt.Print("    Adding the FIWARE Enabler to be analysed... ")
	cmd := exec.Command(absPathDockerCompose, "-f", "docker-compose-anchore.yaml", "exec", "-T", "api", "anchore-cli", "image", "add", enabler)
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		os.Exit(-1)
	} else {
		fmt.Println("Success")
		fmt.Printf("         Result:\n%8v\n", out.String())
	}

	// Step 2: Wait until the analysis is finished (it needs some time)
	// docker-compose -f docker-compose-anchore.yaml exec -T api anchore-cli image wait $enabler
	fmt.Print("    Waiting until the analysis is finished... ")
	out.Reset()
	stderr.Reset()
	cmd = exec.Command(absPathDockerCompose, "-f", "docker-compose-anchore.yaml", "exec", "-T", "api", "anchore-cli", "image", "wait", enabler)
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		os.Exit(-1)
	} else {
		fmt.Println("Success")
		fmt.Printf("         Result:\n%8v\n", out.String())
	}

	// Step 3: Get the list of vulnerabilities
	// docker-compose -f docker-compose-anchore.yaml exec -T api anchore-cli --json image vuln $enabler all

	fmt.Print("    Getting the list of vulnerabilities... ")
	out.Reset()
	stderr.Reset()
	cmd = exec.Command(absPathDockerCompose, "-f", "docker-compose-anchore.yaml", "exec", "-T", "api", "anchore-cli", "--json", "image", "vuln", enabler, "all")
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		os.Exit(-1)
	} else {
		fmt.Println("Success")
		fmt.Printf("         Result:\n%8v\n", out.String())

		// Step 4: Save the out into filename
		err = ioutil.WriteFile(filename, out.Bytes(), 0644)

		if err != nil {
			log.Fatal(err)
		}
	}

	// Just to finish, send the data to the nexus instance
	// redirect_all curl -fsSL -u ${user}':'${password} --upload-file ${filename}  https://nexus.lab.fiware.org/repository/security/check/${enabler}/sast_anchore/${filename}
	// # http -a ${user}:${password} https://nexus.lab.fiware.org/repository/security/check/${enabler}/sast_anchore/${filename} < ${filename}
	// #

	// Keep the name of the file to send afterward the email to the owner
	// filename_anchore=$filename_anchore$(pwd)"/"$filename","

	// #     # We need to remove the last "," from the filename string(s)
	// #     filename_anchore=${filename_anchore::-1}
	// #

	// Return to the original folder
	err = os.Chdir("..")
	CheckIfError(err)
}
