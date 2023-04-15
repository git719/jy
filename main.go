// main.go

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/git719/utl"
	"gopkg.in/yaml.v3"
	"os"
)

const (
	prgname = "jy"
	prgver  = "1.1.4"
)

func PrintUsage() {
	fmt.Printf(prgname + " JSON & YAML converter v" + prgver + "\n" +
		"    FILENAME  Convert given File from JSON to YAML or vice-versa\n" +
		"              You can also pipe the file into the program\n" +
		"    -v        Print this usage page\n")
	os.Exit(0)
}

func main() {
	var buf bytes.Buffer

	_, isGitBash := os.LookupEnv("MSYSTEM")

	fileInfo, _ := os.Stdin.Stat() // Check if anything was piped in
	if (fileInfo.Mode()&os.ModeCharDevice) != 0 || (isGitBash && fileInfo.Size() <= 0) {
		// Nothing piped in, let's check for arguments
		switch len(os.Args[1:]) {
		case 1:
			switch os.Args[1] {
			case "-v":
				PrintUsage()
			default:
				filePath := os.Args[1]
				if utl.FileUsable(filePath) {
					objRaw, _ := utl.LoadFileJson(filePath)
					if objRaw == nil { // If NOT JSON
						objRaw, _ = utl.LoadFileYaml(filePath) // See if it's YAML
						if objRaw == nil {
							utl.Die("Neither JSON nor YAML.\n")
						}
						utl.PrintJson(objRaw) // Print YAML as JSON
						os.Exit(0)
					}
					utl.PrintYaml(objRaw) // Print JSON as YAML
					os.Exit(0)
				} else {
					utl.Die("File is unusable.\n")
				}
			}
		default:
			PrintUsage()
		}
	} else {
		// Read piped input into a buffer
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
				utl.Die("Neither JSON nor YAML.\n")
			}
			utl.PrintJson(objRaw) // Print YAML as JSON
			os.Exit(0)
		}
		utl.PrintYaml(objRaw) // Print JSON as YAML
		os.Exit(0)
	}
}
