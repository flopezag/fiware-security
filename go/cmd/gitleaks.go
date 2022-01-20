package cmd

import (
	//"bytes"
	"fmt"
	"os"
)

func Gitleaks(enabler, filename string) {
	//var out bytes.Buffer
	//var stderr bytes.Buffer

	filename = filename + "_clair.json"

	fmt.Println("Audit git repository for secrets... ")
	fmt.Println("    Docker image: ", enabler)
	fmt.Println("    Output file: ", filename)

	// Change to the Clair folder to execute the analysis
	err := os.Chdir("./Gitleaks")
	CheckIfError(err)

	/*
			Get the last version of the software

		for linux:
			curl -s  https://api.github.com/repos/zricethezav/gitleaks/releases/latest |grep browser_download_url  |  cut -d '"' -f 4  | grep '\linux-amd64$'| wget -i -
		for Mac:
			curl -s  https://api.github.com/repos/zricethezav/gitleaks/releases/latest |grep browser_download_url  |  cut -d '"' -f 4  | grep '\darwin-amd64$'| wget -i -
		for Windows:
			If you’re a Windows user, download and install the gitleaks-windows-amd64.exe package.


		curl -s  https://api.github.com/repos/zricethezav/gitleaks/releases/latest | jq -r '.assets[] | select(.name == "gitleaks-darwin-amd64") | .browser_download_url' | wget -i -


		Move the software to the final name and give permissions

		mv gitleaks-darwin-amd64 gitleaks
		chmod 764 gitleaks


		Execute the analysis of data to a repo

		gitleaks --repo-url=https://github.com/gitleakstest/gronit -v


		Multiprocess
		CPU=$(cat /proc/cpuinfo | grep -ic ^processor)

		en mac: CPU=$(sysctl -a machdep.cpu.thread_count | awk '{print $2}')

		gitleaks --repo=https://github.com/jmutai/dotfiles --threads=$CPU



		Source: https://computingforgeeks.com/gitleaks-audit-git-repos-for-secrets/

	*/

	// Return to the original folder
	err = os.Chdir("..")
	CheckIfError(err)
}
