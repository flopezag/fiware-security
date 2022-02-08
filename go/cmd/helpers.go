package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"time"

	"github.com/zricethezav/gitleaks/v8/report"
)

var Docker_Compose string // global variable to store absolute path of filename.

// Create a folder with 0755 permissions
func Create_folder(folder string) {
	_, err := os.Stat(folder)
	if os.IsNotExist(err) {
		err := os.Mkdir(folder, 0755)
		if err != nil {
			log.Fatal(err)
		}

	}

}

// CheckIfError should be used to naively panics if an error is not nil.
func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}

// Info should be used to describe the example commands that are about to run.
func Info(format string, args ...interface{}) {
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

// Warning should be used to display a warning
func Warning(format string, args ...interface{}) {
	fmt.Printf("\x1b[36;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

// DownloadFile will download a url and store it in local filepath.
// It writes to the destination file as it downloads it, without
// loading the entire file into memory.
func DownloadFile(url string, filepath string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func Copy_file(src string, dest string) {
	bytesRead, err := ioutil.ReadFile(src)

	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(dest, bytesRead, 0644)

	if err != nil {
		log.Fatal(err)
	}

}

func CheckDockerCompose() string {
	path, err := exec.LookPath("docker-compose")

	if err != nil {
		fmt.Printf("didn't find 'docker-compose' executable\n")
		os.Exit(-1)
	} else {
		fmt.Printf("'docker-compose' executable is in '%s'\n", path)
	}

	return path
}

func FindDockerCompose() {
	path, err := exec.LookPath("docker-compose")
	if err != nil {
		fmt.Println(err)
		log.Fatal("We cannot find the 'docker-compose' in the PATH variable.\n" +
			"It is needed to have installed and configured Docker Compose to run this security scan analysis")
	}

	absPathDockerCompose = path
}

// Generate the filename corresponding to the image
func Filename(component, filename string) string {
	var result string

	// Getting the current time
	t := time.Now()
	extension := fmt.Sprintf("%d%02d%02d_%02d%02d", t.Year(), int(t.Month()), t.Day(), t.Hour(), t.Minute())

	// Parse the docker image name to extract the name of the image without tag and organization
	var re = regexp.MustCompile(`([^:\/]*)\/([^:\/]*):?(.*)?`)

	match := re.FindAllStringSubmatch(filename, -1)

	result = component + "_" + match[0][2] + "_" + extension

	return result
}

func writeJson(findings []report.Finding, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	if len(findings) == 0 {
		findings = []report.Finding{}
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")
	return encoder.Encode(findings)
}
