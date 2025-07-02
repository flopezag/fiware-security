package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"time"
)

func Docker_bench_security(compose, filename string) string {
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
	cmd := exec.Command(absPathDockerCompose, "compose", "pull")
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
	cmd = exec.Command(absPathDockerCompose, "compose", "up", "-d")
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
	check := "container_images,container_runtime"
	exclude := "anchore_api_1,anchore_policy-engine_1,anchore_queue_1,anchore_analyzer_1,anchore_catalog_1,anchore_db_1,arminc/clair-db,arminc/clair-local-scan,quay.io/usr42/clair-container-scan"
	cmd = exec.Command("./docker-bench-security.sh", "-c", check, "-x", exclude, "-p")
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + out.String())
		os.Exit(-1)
	} else {
		fmt.Println("Success")
		// fmt.Printf("         Result:\n%8v\n", out.String())
	}

	// Change to the log folder where the result of the execution is produced
	err = os.Chdir("./log")
	CheckIfError(err)

	// Rename the result file
	err = os.Rename("docker-bench-security.log.json", filename)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("scan completed in ", time.Since(start), " seconds")

	// Return to the original folder
	err = os.Chdir("..")
	CheckIfError(err)

	err = os.Chdir("..")
	CheckIfError(err)

	return filename
}
