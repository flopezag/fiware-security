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
	var (
		vc               config.ViperConfig
		findingsFromGit  []report.Finding
		findingsFromFile []report.Finding
		findings         []report.Finding
		err              error
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
		return
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
		return
	}

	fmt.Println(files)

	options := detect.Options{Verbose: false, Redact: false}

	findingsFromGit = detect.FromGit(files, cfg, options)

	findingsFromFile, err = detect.FromFiles(".", cfg, options)
	if err != nil {
		fmt.Println("Failed to scan files")
	}

	findings = findingsFromGit
	for j := 0; j < len(findingsFromFile); j++ {
		findings = append(findings, findingsFromFile[j])
	}

	if len(findings) != 0 {
		fmt.Println("leaks found in Git: ", len(findings))
	} else {
		fmt.Println("no leaks found")
	}

	fmt.Println("scan completed in ", time.Since(start), " seconds")

	writeJson(findings, filename)

	// Delete the cloned repository
	path, _ := os.Getwd()
	fmt.Println(path)
	if filepath.Base(path) == "Gitleaks" {
		err = os.RemoveAll(".")
		if err != nil {
			fmt.Println("Error fatal: ", err)
		}
	}

	// Return to the original folder
	err = os.Chdir("..")
	CheckIfError(err)
}

func InitRules() {
	fmt.Println(os.Getwd())

	if _, err := os.Stat(filepath.Join("./config", "gitleaks.toml")); os.IsNotExist(err) {
		fmt.Println("No gitleaks config found in path ", filepath.Join(".", ".gitleaks.toml"), "using default gitleaks config")
		viper.SetConfigType("toml")
		viper.ReadConfig(strings.NewReader(config.DefaultConfig))
		return
	} else {
		fmt.Println("Using existing gitleaks config ", filepath.Join("./config", "gitleaks.toml"), "from `(--source)/.gitleaks.toml`")
	}

	viper.AddConfigPath(".")
	viper.SetConfigName(".gitleaks")
	viper.SetConfigType("toml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Unable to load gitleaks config, err: ", err)
	}

}
