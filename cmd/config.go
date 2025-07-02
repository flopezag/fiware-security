package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"io"
)

/* Structure of the config data
{
    "enablers": [
        {
            "name": "...",
            "image": ["...", "...", ...],
            "compose": "...",
            "exclude": ["...", "...", ...],
            "email": "..."
        },
	...
	]
}
*/

// Enablers struct which contains an array of enablers
type Enablers struct {
	Enablers []Enabler `json:"enablers"`
}

// Enabler struct which contains a name, an image, a compose, an exclude, and an email
type Enabler struct {
	Name       string   `json:"name"`
	Image      []string `json:"image"`
	Repository []string `json:"repository"`
	Compose    string   `json:"compose"`
	Exclude    []string `json:"exclude"`
	Email      string   `json:"email"`
}

// we initialize our Enablers array
var enablers Enablers

func search(length int, f func(index int) bool) int {
	for index := 0; index < length; index++ {
		if f(index) {
			return index
		}
	}
	return -1
}

func Search(enabler, attribute string) []string {
	var result []string

	// Find the data of a specific FIWARE GE
	// Build a config map:
	idx := search(len(enablers.Enablers), func(index int) bool {
		return enablers.Enablers[index].Name == enabler
	})

	if idx < 0 {
		fmt.Println("not found")
		os.Exit(1)
	}

	if attribute == "Image" {
		result = enablers.Enablers[idx].Image
	} else if attribute == "Repository" {
		result = enablers.Enablers[idx].Repository
	} else if attribute == "Compose" {
		result = append(result, enablers.Enablers[idx].Compose)
	} else {
		fmt.Println("Attribute not found: " + attribute)
		os.Exit(1)
	}

	return result
}

func ParseJSON() {
	// Open our jsonFile
	jsonFile, err := os.Open("./config/enablers.json")

	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("    Successfully Opened enablers.json")

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened jsonFile as a byte array.
	byteValue, _ := io.ReadAll(jsonFile)

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &enablers)
}

func print_data() {
	// we iterate through every user within our users array and
	// print out the user Type, their name, and their facebook url
	// as just an example
	for i := 0; i < len(enablers.Enablers); i++ {
		fmt.Println("Enabler name: " + enablers.Enablers[i].Name)
		for j := 0; j < len(enablers.Enablers[i].Image); j++ {
			fmt.Println("Enabler image: " + enablers.Enablers[i].Image[j])
		}
		fmt.Println("Enabler compose: " + enablers.Enablers[i].Compose)
		for j := 0; j < len(enablers.Enablers[i].Exclude); j++ {
			fmt.Println("Enabler exclude: " + enablers.Enablers[i].Exclude[j])
		}
		fmt.Println("Enabler email: " + enablers.Enablers[i].Email + "\n")
	}
}
