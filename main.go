// main.go

package main

import (
	"encoding/json"
	"fmt"
	"github.com/git719/utl"
	"github.com/mattn/go-isatty"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"runtime"
	"strings"
)

const (
	prgname = "jy"
	prgver  = "1.2.3"
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

func hasPipedInput() bool {
	stat, _ := os.Stdin.Stat() // Check if anything was piped in
	if isGitBashOnWindows() {
		// Git Bash on Windows handles input redirection differently than other shells. When a program
		// is run without any input or arguments, it still treats the input as if it were piped from an
		// empty stream, causing the program to consider it as piped input and hang. This works around that.
		if !isatty.IsCygwinTerminal(os.Stdin.Fd()) {
			return true
		}
	} else {
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			return true
		}
	}
	return false
}

func main() {
	if hasPipedInput() {
		//fmt.Fprintln(os.Stderr, "Processing piped input") // DEBUG
		buffer, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading from stdin:", err)
		}
		byteString := []byte(buffer)

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
		//fmt.Fprintln(os.Stderr, "Processing arguments") // DEBUG
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
