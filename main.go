// main.go

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/git719/utl"
	"gopkg.in/yaml.v3"
	"os"
	"runtime"
	"strings"
)

const (
	prgname = "jy"
	prgver  = "1.2.0"
)

func PrintUsage() {
	fmt.Printf(prgname + " JSON|YAML converter v" + prgver + "\n" +
		"    FILENAME  Convert given File from JSON to YAML or vice-versa\n" +
		"              You can also pipe the file into the program\n" +
		"    -v        Print this usage page\n")
	os.Exit(0)
}

func isGitBashOnWindows() bool {
	return runtime.GOOS == "windows" && strings.HasPrefix(os.Getenv("MSYSTEM"), "MINGW")
}

func isStdinEmpty() bool {
	// Git Bash on Windows handles input redirection differently than other shells. When a program
	// is run without any input or arguments, it still treats the input as if it were piped from an
	// empty stream, causing the program to consider it as piped input and hang. This works around that.
	if isGitBashOnWindows() {
		stat, err := os.Stdin.Stat()
		if err != nil || stat.Size() == 0 {
			return true
		}
	}
	return false
}

func main() {
	var buf bytes.Buffer

	stat, _ := os.Stdin.Stat() // Check if anything was piped in

	if (stat.Mode()&os.ModeCharDevice) == 0 && !isStdinEmpty() {
		// Processing piped input
		_, err := buf.ReadFrom(os.Stdin)
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
		byteString := []byte(buf.String())

		// If JSON then convert to YAML, or vice-versa
		var objRaw interface{}
		// Because JSON is essentially a subset of YAML, we have to check JSON first
		// As an interesting aside, see https://news.ycombinator.com/item?id=31406473
		_ = json.Unmarshal(byteString, &objRaw) // See if it's JSON
		if objRaw == nil {                      // Ok, it's NOT JSON
			_ = yaml.Unmarshal(byteString, &objRaw) // See if it's YAML
			if objRaw == nil {
				utl.Die("Piped input is neither JSON nor YAML\n")
			}
			utl.PrintJson(objRaw) // Print YAML as JSON
			os.Exit(0)
		}
		utl.PrintYaml(objRaw) // Print JSON as YAML
		os.Exit(0)
	} else if len(os.Args) == 2 {
		// Processing arguments
		switch os.Args[1] { // We only care/check for ONE argument
		case "-v":
			// To explicitly print the usage
			PrintUsage()
		default:
			// Or a potential JSON/YAML file to convert
			filePath := os.Args[1]
			if utl.FileUsable(filePath) {
				objRaw, _ := utl.LoadFileJson(filePath)
				if objRaw == nil { // If NOT JSON
					objRaw, _ = utl.LoadFileYaml(filePath) // See if it's YAML
					if objRaw == nil {
						utl.Die("File is neither JSON nor YAML\n")
					}
					utl.PrintJson(objRaw) // Print YAML as JSON
					os.Exit(0)
				}
				utl.PrintYaml(objRaw) // Print JSON as YAML
				os.Exit(0)
			} else {
				utl.Die("File is unusable\n")
			}
		}
	} else {
		PrintUsage()
	}
}
