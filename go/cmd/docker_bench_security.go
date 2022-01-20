package cmd

import (
	//"bytes"
	"fmt"
	"os"
)

func Docker_bench_security(enabler, filename string) {
	//var out bytes.Buffer
	//var stderr bytes.Buffer

	filename = filename + "_clair.json"

	fmt.Println("Docker Bench Security Scan... ")
	fmt.Println("    Docker image: ", enabler)
	fmt.Println("    Output file: ", filename)

	// Change to the Clair folder to execute the analysis
	err := os.Chdir("./Docker-Bench-Security")
	CheckIfError(err)

	// Step 1: Creating an instance of the GE ${enabler}
	// #     redirect_all echo "Creating an instance of the GE ${enabler}"
	// #     cmd='.enablers[] | select(.name == "'${enabler}'") | .compose'
	// #     compose=$(jq -r "${cmd}" enablers.json)
	// #     redirect_all echo "Compose: $compose"
	// #     redirect_all echo
	// #
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
	// #     extension="$(date +%Y%m%d_%H%M%S)-bench.json"
	// #     filename="$@-$extension"
	// #     # enabler=$(echo "$@" | awk -F '/' '{print $2}')
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

	// Return to the original folder
	err = os.Chdir("..")
	CheckIfError(err)
}
