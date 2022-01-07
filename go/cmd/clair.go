package cmd

import "fmt"

func Security_analysis(enabler string) {

	fmt.Println(enabler)
	// #     fmt.Println(
	// #     fmt.Println( "Clair CVE Security Scan... "
	// #
	// #     enabler=$@
	// #
	// #     redirect_all echo "Pulling from "$enabler"..."
	// #
	// #     # Get the list of docker images from the enabler.json taking into account the FIWARE GE name
	// #     cmd='.enablers[] | select(.name == "'${enabler}'") | .image'
	// #     images=$(jq -r "${cmd}" enablers.json)
	// #
	// #     filename_clair=""
	// #
	// #     # From the images, we need to iterate for the different values
	// #     for i in ${images//,/ }
	// #     do
	// #         # call your procedure/other scripts here below
	// #         redirect_all echo "$i"
	// #
	// #         redirect_all docker pull "$i"
	// #         redirect_all echo
	// #
	// #         labels=$(docker inspect --type=image "$i" 2>/dev/null | jq .[].Config.Labels)
	// #
	// #         if [[ ${PULL} -eq 1 ]];
	// #         then
	// #           redirect_all echo "Pulling Clair content ..."
	// #           redirect_all docker-compose pull
	// #           redirect_all echo
	// #         fi
	// #
	// #         redirect_all echo "Security analysis of "$i" image..."
	// #         extension="$(date +%Y%m%d_%H%M%S)-cve.json"
	// #         # filename=$(echo "$i" | awk -F '/' -v a="$extension" '{print $2 a}')
	// #         # enabler=$(echo "$i" | awk -F '/' '{print $2}')
	// #
	// #         # Extract the name of the docker image
	// #         short_name=$(echo $i | awk -F '/' '{print $2}' | awk -F ':' '{print $1}')
	// #         redirect_all echo "$short_name"
	// #
	// #         filename=$(echo "$enabler" | awk  -v a="$extension" -v b="$short_name" '{print $0"-"b"-"a}')
	// #
	// #         redirect_stderr docker-compose run --rm scanner "$i" > ${filename}
	// #         ret=$?
	// #         redirect_all echo
	// #
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
}
