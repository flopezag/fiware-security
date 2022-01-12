package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
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
		log.Fatal("installing fortune is in your future")
	}
	fmt.Printf("fortune is available at %s\n", path)

	absPathDockerCompose = path
}

// Generate the filename corresponding to the image
func Filename() {
	// #         extension="$(date +%Y%m%d_%H%M%S)-anchore.json"
	// #
	// #         # Extract the name of the docker image
	// #         short_name=$(echo $i | awk -F '/' '{print $2}' | awk -F ':' '{print $1}')
	// #         redirect_all echo "$short_name"
	// #
	// #         filename=$(echo "$enabler" | awk  -v a="$extension" -v b="$short_name" '{print $0"-"b"-"a}')

}
