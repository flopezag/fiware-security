package cmd

import (
	//"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"

	"github.com/zricethezav/gitleaks/v8/config"
	"github.com/zricethezav/gitleaks/v8/detect"
	gl "github.com/zricethezav/gitleaks/v8/git"
	"github.com/zricethezav/gitleaks/v8/report"

	git "github.com/go-git/go-git/v5"
)

func Gitleaks(enabler_repository, filename string) {
	// Get gitloeaks config
	viper.SetConfigType("toml")
	viper.ReadConfig(strings.NewReader("./config/gitleaks.toml"))

	var (
		vc       config.ViperConfig
		findings []report.Finding
		err      error
	)

	filename = filename + "_gitleaks.json"

	fmt.Println("Audit git repository for secrets... ")
	fmt.Println("    Enabler repository: ", enabler_repository)
	fmt.Println("    Output file: ", filename)

	// Change to the Gitleaks folder to execute the analysis
	err = os.Chdir("./Gitleaks")
	CheckIfError(err)

	// Clone the given repository to the given directory
	fmt.Println("\n    git clone ", enabler_repository)

	_, err = git.PlainClone(".", false, &git.CloneOptions{URL: enabler_repository})

	CheckIfError(err)

	viper.Unmarshal(&vc)
	cfg, err := vc.Translate()
	if err != nil {
		fmt.Println(err)
		fmt.Println("Failed to load config")
		os.Exit(1)
	}

	if cfg.Path == "" {
		cfg.Path = filepath.Join("../config", "gitleaks.toml")
	}

	fmt.Println(cfg)

	start := time.Now()

	files, err := gl.GitLog(".", "")
	if err != nil {
		fmt.Println(err)
		fmt.Println("Failed to get git log")
		os.Exit(1)
	}

	fmt.Println(files)

	options := detect.Options{Verbose: false, Redact: false}

	findings = detect.FromGit(files, cfg, options)

	if len(findings) != 0 {
		fmt.Println("leaks found: ", len(findings))
	} else {
		fmt.Println("no leaks found")
	}

	source := "."
	findings, err = detect.FromFiles(source, cfg, options)
	if err != nil {
		fmt.Println("Failed to scan files")
	}

	fmt.Println("scan completed in ", time.Since(start), " seconds")

	writeJson(findings, filename)

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

	// Delete the cloned repository

	// Return to the original folder
	err = os.Chdir("..")
	CheckIfError(err)
}
