package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"time"
)

func Docker_bench_security(compose, filename string) {
	var out bytes.Buffer
	var stderr bytes.Buffer

	start := time.Now()
	fmt.Println("Starting at: ", start)

	filename = filename + "_bench.json"

	fmt.Println("Docker Bench Security Scan... ")
	fmt.Println("    Docker compose: ", compose)
	fmt.Println("    Output file: ", filename)

	// Change to the Docker-Bench-Security folder to execute the analysis
	err := os.Chdir("./Docker-Bench-Security")
	CheckIfError(err)

	// Rename of the Docker-Bench-Security compose file
	err = os.Rename("docker-compose.yml", "docker-compose.old")
	if err != nil {
		fmt.Println(err)
	}

	// Download the docker compose file for the corresponding enabler. It could be docker-compose.yml or docker-compose-*.yml
	DownloadFile(compose, "docker-compose.yml")

	// Pull the docker-compose.yaml (enabler) content
	fmt.Print("Pulling Docker-Bench-Security content... ")
	cmd := exec.Command(absPathDockerCompose, "pull")
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

	// Launch the docker-compose.yml (enabler) content
	fmt.Print("Launching Docker-Bench-Security... ")
	cmd = exec.Command(absPathDockerCompose, "up", "&")
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

	// Wait 120 seconds to finish the previous command
	fmt.Printf("Sleeping 120 seconds... ")
	time.Sleep(time.Duration(120) * time.Second)
	fmt.Println("Done")

	// Execute the script
	fmt.Print("Executing the Docker-Bench-Security... ")
	cmd = exec.Command("./docker-bench-security.sh", "-c container_images,container_runtime,docker_security_operations")
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

	// Rename the result file
	err = os.Rename("docker-bench-security.sh.log.json", filename)
	if err != nil {
		fmt.Println(err)
	}

	// Step 1: Creating an instance of the GE ${enabler}
	// #     redirect_all echo "Creating an instance of the GE ${enabler}"
	// #     mv docker-compose.yml docker-compose.old
	// #     wget -q $compose
	// #     compose=$(echo $compose | awk -F '/' '{print $NF}')
	// #
	// #     # Could be possible that they are the same names (e.g. Orion)
	// #     mv $compose docker-compose.yml 2>/dev/null
	// #     redirect_all docker-compose pull
	// #     redirect_all docker-compose up &
	// #
	// #     # We need to wait until it is ready the docker-compose service
	// #     sleep 120
	// #
	// #     redirect_all ./docker-bench-security.sh  -c container_images,container_runtime,docker_security_operations
	// #
	// #     mv docker-bench-security.sh.log.json ${filename}
	// #

	// #     # Just to finish, send the data to the nexus instance
	// #     redirect_all curl -fsSL -u ${user}':'${password} --upload-file ${filename}  https://nexus.lab.fiware.org/repository/security/check/${enabler}/bench-security/${filename}
	// #
	// #     # Keep the name of the file to send afterward the email to the owner
	// #     filename_bench=$(pwd)"/"$filename
	// #
	// #     cd ..
	// #

	fmt.Println("scan completed in ", time.Since(start), " seconds")

	// Return to the original folder
	err = os.Chdir("..")
	CheckIfError(err)
}
