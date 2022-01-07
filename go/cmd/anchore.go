package cmd

import "fmt"

func Anchore(enabler string) {

	fmt.Println(enabler)
	// #     fmt.Println(
	// #     fmt.Println( "Anchore Security Scan... "
	// #     enabler=$@
	// #
	// #     redirect_all echo $(pwd)
	// #     cd ./anchore
	// #
	// #     filename_anchore=""
	// #
	// #     cmd='.enablers[] | select(.name == "'${enabler}'") | .image'
	// #     images=$(jq -r "${cmd}" enablers.json)
	// #
	// #     # From the images, we need to iterate for the different values
	// #     for i in ${images//,/ }
	// #     do
	// #         # call your procedure/other scripts here below
	// #         redirect_all echo "$i"
	// #
	// #         extension="$(date +%Y%m%d_%H%M%S)-anchore.json"
	// #
	// #         # Extract the name of the docker image
	// #         short_name=$(echo $i | awk -F '/' '{print $2}' | awk -F ':' '{print $1}')
	// #         redirect_all echo "$short_name"
	// #
	// #         filename=$(echo "$enabler" | awk  -v a="$extension" -v b="$short_name" '{print $0"-"b"-"a}')
	// #
	// #
	// #         redirect_all docker pull $i
	// #
	// #         # Step 4: Add the FIWARE Enabler to be analysed
	// #         redirect_all docker-compose -f docker-compose-anchore.yaml exec api anchore-cli image add $i
	// #
	// #         # Step 5: Wait until the analysis is finished (it needs some time)
	// #         redirect_all docker-compose -f docker-compose-anchore.yaml exec api anchore-cli image wait $i
	// #
	// #         # Step 6: Get the list of vulnerabilities
	// #         redirect_stderr docker-compose -f docker-compose-anchore.yaml exec api anchore-cli --json image vuln $i all > ${filename}
	// #
	// #         # Just to finish, send the data to the nexus instance
	// #         redirect_all curl -fsSL -u ${user}':'${password} --upload-file ${filename}  https://nexus.lab.fiware.org/repository/security/check/${enabler}/sast_anchore/${filename}
	// #         # http -a ${user}:${password} https://nexus.lab.fiware.org/repository/security/check/${enabler}/sast_anchore/${filename} < ${filename}
	// #
	// #         # Keep the name of the file to send afterward the email to the owner
	// #         filename_anchore=$filename_anchore$(pwd)"/"$filename","
	// #     done
	// #
	// #     # We need to remove the last "," from the filename string(s)
	// #     filename_anchore=${filename_anchore::-1}
	// #
	// #     cd ..
}
