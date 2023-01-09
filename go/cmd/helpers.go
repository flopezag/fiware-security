package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
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

func Error(format string, args ...interface{}) {
	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
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
	path, err := exec.LookPath("docker")

	if err != nil {
		fmt.Printf("didn't find 'docker' executable, need to execute docker compose\n")
		os.Exit(-1)
	} else {
		fmt.Printf("'docker' executable is in '%s'\n", path)
	}

	path = path + " compose"
	return path
}

func FindDockerCompose() {
	path, err := exec.LookPath("docker")
	if err != nil {
		fmt.Println(err)
		log.Fatal("We cannot find the 'docker' in the PATH variable.\n" +
			"It is needed to have installed and configured Docker in order to execute docker compose to run this security scan analysis")
	}

	path = path + " compose"
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

func FilenameFromUrl(component, uri string) string {
	var result string

	// Getting the current time
	t := time.Now()
	extension := fmt.Sprintf("%d%02d%02d_%02d%02d", t.Year(), int(t.Month()), t.Day(), t.Hour(), t.Minute())

	u, err := url.Parse(uri)
	if err != nil {
		panic(err)
	}

	fmt.Println(u.Path)
	res := strings.Split(u.Path, "/")

	result = component + "_" + res[len(res)-1] + "_" + extension

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

func deleteClonedFolder() {
	var files []string

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		res, _ := regexp.MatchString("_gitleaks.json", path)

		if path != "." && !res {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	// We need to delete the *_gitleaks.json and . files for the list of files to delete

	for _, file := range files {
		err = os.RemoveAll(file)
		if err != nil {
			return
		}
	}
}

func check_mandatory_commands() {
	// Firstly: The scanner need to be executed in linux system, basically due to the Docker-Bench-Security
	//        uses system commands calls
	if runtime.GOOS == "windows" {
		fmt.Println("Security Scan can't be executed on a windows machine")
		os.Exit(-1)
	}

	// Secondly: There are several commands tools needed to be installed to be used by Docker-Bench-Security
	fmt.Print("    Checking mandatory programs... ")

	programs := [10]string{"awk", "docker", "grep", "stat", "tee", "tail", "wc", "xargs", "truncate", "sed"}

	for i := 0; i < 10; i++ {
		cmd := exec.Command("bash", "-c", "command", "-v", programs[i])
		err := cmd.Run()
		if err != nil {
			fmt.Println("Required program not found: ", programs[i])
			os.Exit(-1)
		}
	}

	fmt.Println("Success")
}

func fileExists(fileName string) bool {
	// check if a file exist, it can return false if it is not exist or it is a directory
	info, err := os.Stat(fileName)
	if err != nil && !os.IsNotExist(err) {
		return false
	}

	if info != nil && err == nil {
		if !info.IsDir() {
			return true
		}
	}
	return false
}
