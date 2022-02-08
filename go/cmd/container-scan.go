package cmd

import (
	"bytes"
	"context"
	// "encoding/json"
	"fmt"
	// "io/ioutil"
	"log"
	// "net/http"
	"os"
	"os/exec"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/go-git/go-git/v5"
	//"github.com/google/go-github/github"
)

var absPathDockerCompose string // global variable to store absolute path of filename.

/* # #!/bin/bash
# # set -xv
#
# user="{{ user }}"
# password="{{ password }}"
#
# export DOCKER_CLIENT_TIMEOUT=200
# export COMPOSE_HTTP_TIMEOUT=200
#
#
# function redirect_stderr() {
#     if [[ ${VERBOSE} -eq 1 ]]; then
#         "$@"
#     else
#         "$@" 2>/dev/null
#     fi
# }
#
# function redirect_all() {
#     if [[ ${VERBOSE} -eq 1 ]]; then
#         "$@"
#     else
#         "$@" 2>/dev/null >/dev/null
#     fi
# }
#
#
#
# function email() {
#   # This function receive 4 parameters
#   # FIWARE GE name
#   # Result file of the Docker Bench Security Scan
#   # Result files of the Clair scan
#   # Result files of the Anchore scan
#   # Return:
#   #             mail \
#   #                 -s <subject> \
#   #                 -t <destination od the email> \
#   #                 -b <content of the message> \
#   #                 -a <Clair test report> \
#   #                 -a <Anchore test report> \
#   #                 -a <Docker Bench Security report>
#
#
#   enabler_name=$1
#   email_owner=$(jq '.enablers[] | select(.name=="Orion") | .email' enablers.json)
#   bench_file=$2
#   clair_file=$3
#   anchore_file=$4
#
#
#
#   MESSAGE="Dear FIWARE GE owner,
#
#   As a result of the security analysis working group, we have identified
#   a set of possible security issues in your component, based on Clair,
#   Anchore and Docker Bench Security Analysis. We recommend that
#   you perform a thorough analysis of these results and take appropriate
#   actions to resolve the security issues found as soon as possible.
#
#   Thank you in advance for your cooperation...
#
#   Kind Regards,
#   Fernando"
#
#
#
#   SUBJECT="[Security Analysis] Analysis of docker image: "$enabler_name
#
#
#
#   echo "$MESSAGE" > /tmp/tmpfile.$$
#
#   email_command="/home/ubuntu/security-scan/common/mail
#                     -s \"$SUBJECT\"
#                     -t $email_owner
#                     -b /tmp/tmpfile.$$
#                     -a $bench_file"
#
#   # From the files associated to clair generate -a <file> for the different values
#   for i in ${clair_file//,/ }
#   do
#     email_command=$email_command"  -a $i"
#   done
#
#   # From the files associated to anchore generate -a <file> for the different values
#   for i in ${anchore_file//,/ }
#   do
#     email_command=$email_command"  -a $i"
#   done
#
#   echo
#   echo ${email_command}
# }
#
# function print_result {
#   file_clair=$1
#   file_bench=$2
#   file_anchore=$3
#
#   fmt.Println(
#   fmt.Println( "CVE Clair vulnerabilities"
#   fmt.Println(
#   for a in Unknown Negligible Low Medium High ;
#   do
#     data=$(more $file_clair | jq ".[].vulnerabilities[].severity | select (.==\"${a}\")" | wc -l)
#     fmt.Println( "    $a  $data"
#   done
#
#   TOTAL=$(more $file_clair  | jq '.[].vulnerabilities | length')
#   fmt.Println( "    TOTAL: ${TOTAL}"
#
#
#
#   fmt.Println(
#   fmt.Println(
#   fmt.Println( "CIS Docker Benchmark (security best practices)"
#   fmt.Println(
#   for a in "Container Images and Build File" "Container Runtime" "Docker Security Operations";
#   do
#     fmt.Println( "    $a"
#     for b in PASS INFO NOTE WARN;
#     do
#       data=$(more $file_bench | jq ".tests[] | select(.desc == \"${a}\") | .results[].result | select (.==\"${b}\")" | wc -l)
#       # data=$(more $file_bench | jq ".tests[].results[].result | select (.==\"${b}\")" | wc -l)
#       fmt.Println( "        $b  $data"
#     done
#   done
#   TOTAL=$(more $file_bench  | jq '.tests[].results | length' | awk '{sum+=$1} END{printf("%d\n",sum)}')
#
#   fmt.Println(
#   fmt.Println( "    TOTAL: ${TOTAL}"
#   fmt.Println(
#
#
#
#   fmt.Println(
#   fmt.Println(
#   fmt.Println( "Anchore Security Analysis"
#   fmt.Println(
#   for a in Unknown Negligible Low Medium High Critical;
#   do
#     data=$(more $file_anchore | jq ".vulnerabilities[].severity | select (.==\"${a}\")" | wc -l)
#     fmt.Println( "    $a  $data"
#   done
#
#   TOTAL=$(more $file_anchore  | jq '.vulnerabilities[].severity' | wc -l)
#   fmt.Println( "    TOTAL: ${TOTAL}"
#   fmt.Println(
# }
#
# PULL=0
# VERBOSE=0
#
# while getopts ":phv" opt; do
#     case ${opt} in
#         p)
#             PULL=1
#             ;;
#         v)
#             VERBOSE=1
#             ;;
#         \?)
#             echo "Invalid option: -$OPTARG" >&2
#             usage
#             ;;
#         h)
#             usage
#             ;;
#     esac
# done
# shift $(($OPTIND -1))
#
# init
#
# if [[ -n $1 ]]; then
#     security_analysis "$1"
#     docker_bench_security "$1"
#     anchore "$1"
#
#     email "$1" ${filename_bench} ${filename_clair} ${filename_anchore}
#     eval $email_command
#     rm /tmp/tmpfile.$$
#
#     print_result ${filename_clair} ${filename_bench} ${filename_anchore}
#
#     clean_docker
# else
#     for ge in `more enablers.json | jq .enablers[].name | sed 's/"//g'`
#     do
#       security_analysis ${ge}
#       docker_bench_security ${ge}
#       anchore ${ge}
#       fmt.Println(
#       fmt.Println(
#
#       email ${ge} ${filename_bench} ${filename_clair} ${filename_anchore}
#       eval $email_command
#       rm /tmp/tmpfile.$$
#
#       print_result ${filename_clair} ${filename_bench} ${filename_anchore}
#     done
#
#     clean_docker
# fi
#
# exit ${ret} */

func initialize() {
	var (
		out    bytes.Buffer
		stderr bytes.Buffer
		err    error
	)

	fmt.Println("\nInitialize the scan...")
	// Localize docker-compose program
	FindDockerCompose()

	// We want to create the corresponding folders for Clair, Ancoher, amd Docker Bench Security
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	Create_folder(path + "/Anchore")
	Create_folder(path + "/Docker-Bench-Security")
	Create_folder(path + "/Clair")
	Create_folder(path + "/Gitleaks")

	// Download the CVE_Severity_Scan docker compose yaml file in order to execute the Clair
	DownloadFile("https://raw.githubusercontent.com/flopezag/fiware-security/develop/Common/cve_severity_scan.yml", path+"/Clair/docker-compose.yml")

	// Change to the Clair directory
	err = os.Chdir("Clair")
	CheckIfError(err)

	// Pull the Clair content
	fmt.Print("Pulling Clair content... ")
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

	// Change to the initial directory
	err = os.Chdir("..")
	CheckIfError(err)

	// Download the enablers configuration file and copy the file to the three folders
	DownloadFile("https://raw.githubusercontent.com/flopezag/fiware-security/develop/Common/enablers.json", "enablers.json")
	Copy_file("enablers.json", "./Anchore/enablers.json")
	Copy_file("enablers.json", "./Clair/enablers.json")

	// Clone the given repository to the given directory
	Info("\ngit clone https://github.com/docker/docker-bench-security.git")

	_, err = git.PlainClone(path+"/Docker-Bench-Security", false, &git.CloneOptions{
		URL:      "https://github.com/docker/docker-bench-security.git",
		Progress: os.Stdout,
	})

	Info("%s\n", err)

	// Copy enablers.json to the Docker-Bench-Security folder and delete the file in the .. folder
	Copy_file("enablers.json", "./Docker-Bench-Security/enablers.json")
	err = os.Remove("enablers.json")
	CheckIfError(err)

	// Change to the Anchore directory
	err = os.Chdir("Anchore")
	CheckIfError(err)

	// Download the possible new version of the docker-compose.yaml file
	DownloadFile("https://engine.anchore.io/docs/quickstart/docker-compose.yaml", "docker-compose-anchore.yaml")

	// Start Anchore engine
	// #     redirect_all docker-compose -f docker-compose-anchore.yaml up -d
	fmt.Print("Starting Anchore engine... ")
	err = exec.Command(absPathDockerCompose, "-f", "docker-compose-anchore.yaml", "up", "-d").Run()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	} else {
		fmt.Println("Success")
	}

	// Verify service availability

	// Wait until the vulnarabilities dictionary is download
	fmt.Print("Waiting vulnerability dictionary downloads... ")
	err = exec.Command(absPathDockerCompose, "-f docker-compose-anchore.yaml exec api anchore-cli system wait").Run()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	} else {
		fmt.Print("Finished\n\n")
	}

	// Change to the Gitleaks directory
	/*
		err = os.Chdir("../Gitleaks")
		CheckIfError(err)

		*
				curl -s  https://api.github.com/repos/zricethezav/gitleaks/releases/latest | jq -r '.assets[] | select(.name == "gitleaks-darwin-amd64") | .browser_download_url' | wget -i -


			Move the software to the final name and give permissions

			mv gitleaks-darwin-amd64 gitleaks
			chmod 764 gitleaks
		*
		resp, err := http.Get("https://api.github.com/repos/zricethezav/gitleaks/releases/latest")

		if err != nil {
			log.Fatal(err)
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(body))
		fmt.Println()
		fmt.Println()

		// prepare result
		result := make(map[string]interface{})
		json.Unmarshal(body, &result)

		fmt.Println(result["assets"])
		*
			client := github.NewClient(nil)
			opt := &github.ListOptions{Page: 2, PerPage: 10}

			releases, rsp, err := client.Repositories.ListReleases(context.Background(), "zricethezav", "gitleaks", opt)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Printf("\n%+v\n", releases)
			fmt.Printf("\n%+v\n", rsp)

			client.Repositories
		*
	*/
	// Change to the original directory
	err = os.Chdir("..")
	CheckIfError(err)

	// os.Exit(0)
}

func clean() {
	// Stop/down the Anchore engine
	fmt.Println("\nClean up the Anchore docker-compose engine...")

	err := os.Chdir("Anchore")
	CheckIfError(err)

	exec.Command(absPathDockerCompose, "-f docker-compose-anchore.yaml down")

	fmt.Println()

	// Stop/down the Docker-Bench-Analysis engine
	// err = os.Chdir("../Docker-Bench-Security")
	// CheckIfError(err)

	fmt.Println("Clean up the Docker-Bench-Security docker-compose engine...")
	// exec.Command(absPathDockerCompose, "down")
	fmt.Println()

	// Stop/down the Clair engine
	// err = os.Chdir("../Clair")
	// CheckIfError(err)

	fmt.Println("Clean up the Clair docker-compose engine...")
	// exec.Command(absPathDockerCompose, "down")
	fmt.Println()

	// Going back to the original folder
	err = os.Chdir("..")
	CheckIfError(err)

	// Clean all docker images
	fmt.Println("Clean up all docker images...")
	// docker kill $(docker ps --all -q) 2>/dev/null
	// docker rm $(docker ps --all -q) 2>/dev/null
	// docker rmi $(docker images --all -q)
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Print("Stopping container ", container.ID[:10], "... ")
		if err := cli.ContainerStop(ctx, container.ID, nil); err != nil {
			panic(err)
		}
		fmt.Println("Success")
	}

	for _, container := range containers {
		fmt.Print("Removing container ", container.ID[:10], "... ")
		if err := cli.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{
			RemoveVolumes: true,
			RemoveLinks:   true,
			Force:         true,
		}); err != nil {
			panic(err)
		}
		fmt.Println("Success")
	}

	images, err := cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	for _, image := range images {
		fmt.Println(image.ID[7:])
	}

	// Removing images
	// for _, image := range images {
	// 	fmt.Print("Removing image ", image.ID[:10], "... ")
	//    if err := cli.ImageRemove(ctx, image.ID, nill); err !=nill {
	//        panic(err)
	//    }
	//  fmt.Println("Success")
	// }
}
