// main.go

package main

import (
	"encoding/json"
	"fmt"
	"github.com/git719/utl"
	goyaml "github.com/goccy/go-yaml"
	"github.com/gookit/color"
	"github.com/mattn/go-isatty"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"runtime"
	"strings"
)

const (
	prgname = "jy"
	prgver  = "1.2.9"
)

func printUsage() {
	fmt.Printf(prgname + " JSON|YAML converter v" + prgver + "\n" +
		"    FILENAME       Convert given file from JSON to YAML or vice-versa\n" +
		"                   You can also pipe the file into the program\n" +
		"    -c FILENAME    Prints given JSON or YAML file in color\n" +
		"    -v             Print this usage page\n")
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

func processPipedInput() {
	// Convert piped input from JSON to YAML or vice-versa, then print in color
	rawBytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading from stdin:", err)
	}

	// Remove color codes in piped input
	stringSansColor := color.ClearCode(string(rawBytes))
	rawBytes = []byte(stringSansColor)

	// JSON must be checked first because it is a subset of the YAML standard
	var rawObject interface{}
	_ = json.Unmarshal(rawBytes, &rawObject) // Is it JSON?
	if rawObject == nil {
		_ = yaml.Unmarshal(rawBytes, &rawObject) // Is it YAML?
		if rawObject == nil {
			utl.Die("Piped input is neither JSON nor YAML\n")
		}
		// It is YAML, print in colorized JSON
		jsonBytes, _ := goyaml.YAMLToJSON(rawBytes)
		jsonBytes2, _ := utl.JsonBytesReindent(jsonBytes, 2) // Two space indent
		utl.PrintJsonBytesColor(jsonBytes2)
	} else {
		// It is JSON, print in colorized YAML
		utl.PrintYamlColor(rawObject)
	}
}

func convertThenPrintInColor(filePath string) {
	// Convert given file from JSON to YAML or vice-versa, then print in color
	if !utl.FileUsable(filePath) {
		utl.Die("File is unusable\n")
	}
	// JSON must be checked first because it is a subset of the YAML standard
	rawObject, err := utl.LoadFileJson(filePath)
	if err == nil {
		// It's JSON, print in colorized YAML
		utl.PrintYamlColor(rawObject)
	} else {
		yamlBytes, err := utl.LoadFileYamlBytes(filePath)
		if err == nil {
			// It's YAML, print in colorized JSON
			jsonBytes, _ := goyaml.YAMLToJSON(yamlBytes)
			jsonBytes2, _ := utl.JsonBytesReindent(jsonBytes, 2) // Two space indent
			utl.PrintJsonBytesColor(jsonBytes2)
		} else {
			utl.Die("File is neither JSON nor YAML\n")
		}
	}
}

func printInColor(filePath string) {
	// Print given JSON or YAML file in color
	// JSON must be checked first because it is a subset of the YAML standard
	jsonBytes, err := utl.LoadFileYamlBytes(filePath)
	if err == nil {
		utl.PrintJsonBytesColor(jsonBytes) // Print colorized JSON
	} else {
		// Load as raw YAML byte slice that can include comments
		yamlBytes, err := utl.LoadFileYamlBytes(filePath)
		if err == nil {
			utl.PrintYamlBytesColor(yamlBytes) // Print colorized YAML
		} else {
			utl.Die("File is neither JSON nor YAML\n")
		}
	}
}

func processArgumentInput() {
	if len(os.Args) == 2 {
		switch os.Args[1] {
		case "-v":
			printUsage()
		default:
			convertThenPrintInColor(os.Args[1])
		}
	} else if len(os.Args) == 3 {
		switch os.Args[1] {
		case "-c":
			printInColor(os.Args[2])
		default:
			printUsage()
		}
	} else {
		printUsage()
	}
}

func main() {
	if hasPipedInput() {
		processPipedInput()
	} else {
		processArgumentInput()
	}
}
