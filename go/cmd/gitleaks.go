package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"

	"github.com/zricethezav/gitleaks/v8/config"
	"github.com/zricethezav/gitleaks/v8/detect"
	gl "github.com/zricethezav/gitleaks/v8/detect/git"
	"github.com/zricethezav/gitleaks/v8/report"

	git "github.com/go-git/go-git/v5"
)

func Gitleaks(enabler_repository, filename string) string {
	var (
		vc  config.ViperConfig
		err error

		findingsFromGit  []report.Finding
		findingsFromFile []report.Finding
		findings         []report.Finding
	)

	start := time.Now()
	fmt.Println("Starting at: ", start)

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
		return ""
	}

	if cfg.Path == "" {
		cfg.Path = filepath.Join("../config", "gitleaks.toml")
	}
	// TODO: Git log is executed twice, once here to see if it is valid and second in the scan analysis
	files, err := gl.NewGitLogCmd(".", "")
	if err != nil {
		fmt.Println(err)
		fmt.Println("Failed to get git log")
		return ""
	}

	fmt.Println(files)

	// Setup detector
	detector := detect.NewDetector(cfg)

	if detector.Config.Path == "" {
		detector.Config.Path = filepath.Join("../config", ".gitleaks.toml")
	}

	detector.Verbose = false
	detector.Redact = false
	var source string = "."
	var logOpts string = ""

	if fileExists(filepath.Join(source, ".gitleaksignore")) {
		detector.AddGitleaksIgnore(filepath.Join(source, ".gitleaksignore"))
	}

	// 1st: detect leaks from git logs
	findingsFromGit, err = detector.DetectGit(source, logOpts, detect.DetectType)
	if err != nil {
		// don't exit on error, just log it
		// log.Error().Msg(err.Error())
		Error(err.Error())
	}

	// 2nd: detect leaks from files
	findingsFromFile, err = detector.DetectFiles(source)
	if err != nil {
		// don't exit on error, just log it
		Error(err.Error())
	}

	findings = findingsFromGit
	for j := 0; j < len(findingsFromFile); j++ {
		findings = append(findings, findingsFromFile[j])
	}

	// log info about the scan
	if err == nil {
		Info("scan completed in %s", time.Since(start))
		if len(findings) != 0 {
			Warning("leaks found: %d", len(findings))
		} else {
			Info("no leaks found")
		}
	} else {
		Warning("partial scan completed in %s", time.Since(start))
		if len(findings) != 0 {
			Warning("%d leaks found in partial scan", len(findings))
		} else {
			Warning("no leaks found in partial scan")
		}
	}

	writeJson(findings, filename)

	// Delete the cloned repository
	path, _ := os.Getwd()
	fmt.Println(path)
	if filepath.Base(path) == "Gitleaks" {
		// We want to delete all the files and folders except the generated file with extension "_gitleaks.json"
		fmt.Print("    Removing cloning repository... ")
		deleteClonedFolder()
		fmt.Println("Success")
	}

	fmt.Println("scan completed in ", time.Since(start), " seconds")

	// Return to the original folder
	err = os.Chdir("..")
	CheckIfError(err)

	return filename
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
