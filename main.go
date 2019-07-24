package main

import (
	// Package for encoding JSON from Config File
	"encoding/json"
	// Package for Stdout printing
	"fmt"
	// Package used for reading configuration file
	"io/ioutil"
	// Package used for reading cli arguments
	"os"
	// Package used for executing commands
	"os/exec"
	// Package with string utilities
	"strings"
)

// PwbsInfo : PWBS Info Structure
type PwbsInfo struct {
	// Version of PWBS
	version      string
	// Edition of PWBS
	edition      string
	// Version Array of PWBS
	versionArray [4]int
}

// PWBS Meta Data
var pwbs = PwbsInfo{
	version:      "0.9.1.2",
	edition:      "E1 GoLang",
	versionArray: [4]int{0, 9, 1, 2},
}

// Main Function of Program
func main() {
	// Show baner
	baner()
	// Get all cli arguments with task names to execute
	programArguments := os.Args[1:]
	// Start Main Logic Function
	pwbsMain(programArguments)
}

// Show baner function
func baner() {
	// Text Baner
	baner := fmt.Sprintf(
		"PAiP Web Build System %v Edition %v",
		pwbs.version,
		pwbs.edition,
	)
	// Print Text Baner
	fmt.Println(baner)
}

// PWBSConfigFile : PWBS Config File Structure
type PWBSConfigFile struct {
	// Commands map with keys which are tasknames and values that are strings or slices of strings
	Commands map[string]interface{} `json:"commands"`
}

// Read JSON File
func readJSON(filename string) PWBSConfigFile {
	// Read From Specified File
	data, err := ioutil.ReadFile(filename)
	// If error occured print it to stdout
	if err != nil {
		fmt.Println("Error in reading File", err)
		return PWBSConfigFile{Commands: map[string]interface{}{}}
	}
	// Make Variable for Unmarshaled JSON
	JSONData := PWBSConfigFile{}
	// Unmarshal json to variable
	jsonErr := json.Unmarshal(data, &JSONData)
	// If error occured print it to stdout
	if jsonErr != nil {
		fmt.Println("Error in parsing JSON", jsonErr)
		return PWBSConfigFile{Commands: map[string]interface{}{}}
	}
	// Return Unmarshalled JSON
	return JSONData
}

// Execute specified command and arguments for that command
func execute(command string, args string) string {
	// Execute command
	cmd := exec.Command(command, args)
	// Get output of command
	out, err := cmd.Output()
	// Return nothing on error
	if err != nil {
		return ""
	}
	// Convert output of command to string
	return string(out)
}

// PWBS Main Function
func pwbsMain(args []string) {
	// Read JSON from configuration file
	JSONData := readJSON("pwbs.json")
	// Loop over cli arguments that are processed as task to execute
	for _, arg := range args {
		// Start baner of task
		baner := fmt.Sprintf(`Executing task "%v" ...`, arg)
		fmt.Println(baner)
		// Get Command Value
		command := JSONData.Commands[arg]
		// Type Switch on Command Value
		switch typedCommand := command.(type) {
		case string:
			// Single Command Task
			c := strings.SplitN(typedCommand, " ", 2)
			cmd, arguments := c[0], c[1]
			output := execute(cmd, arguments)
			fmt.Println(output)
		case []interface{}:
			// Multi Command Task
			for _, tcwd := range typedCommand {
				switch typedCommandListItem := tcwd.(type) {
				case string:
					// Multi Command Task
					c := strings.SplitN(typedCommandListItem, " ", 2)
					cmd, arguments := c[0], c[1]
					output := execute(cmd, arguments)
					fmt.Println(output)
				default:
					// Make error
					fmt.Println("Unsupported typeof task")
					fmt.Printf("Task type: %T\n", typedCommandListItem)
					fmt.Printf("Task value: %v\n", typedCommandListItem)
				}
			}
		default:
			/// Make error
			fmt.Println("Unsupported typeof task")
			fmt.Printf("Task type: %T\n", typedCommand)
			fmt.Printf("Task value: %v\n", typedCommand)
		}
		// End of baner task
		baner = fmt.Sprintf(`Finished task "%v" ...`, arg)
		fmt.Println(baner)
	}
}
