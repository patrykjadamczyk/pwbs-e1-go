package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// PwbsInfo : PWBS Info Structure
type PwbsInfo struct {
	version      string
	edition      string
	versionArray [4]int
}

var pwbs = PwbsInfo{
	version:      "0.9.1.1",
	edition:      "E1 GoLang",
	versionArray: [4]int{0, 9, 1, 0},
}

func main() {
	baner()
	programArguments := os.Args[1:]
	pwbsMain(programArguments)
}

func baner() {
	baner := fmt.Sprintf("PAiP Web Build System %v Edition %v", pwbs.version, pwbs.edition)
	fmt.Println(baner)
}

// PWBSConfigFile : PWBS Config File Structure
type PWBSConfigFile struct {
	Commands map[string]interface{} `json:"commands"`
}

func readJSON(filename string) PWBSConfigFile {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error in reading File", err)
		return PWBSConfigFile{Commands: map[string]interface{}{}}
	}
	JSONData := PWBSConfigFile{}
	jsonErr := json.Unmarshal(data, &JSONData)
	if jsonErr != nil {
		fmt.Println("Error in parsing JSON", jsonErr)
		return PWBSConfigFile{Commands: map[string]interface{}{}}
	}
	return JSONData
}

func execute(command string, args string) string {
	cmd := exec.Command(command, args)
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return string(out)
}

func pwbsMain(args []string) {
	JSONData := readJSON("pwbs.json")
	for _, arg := range args {
		baner := fmt.Sprintf(`Executing task "%v" ...`, arg)
		fmt.Println(baner)
		command := JSONData.Commands[arg]
		switch typedCommand := command.(type) {
		case string:
			c := strings.SplitN(typedCommand, " ", 2)
			cmd, arguments := c[0], c[1]
			output := execute(cmd, arguments)
			fmt.Println(output)
		case []interface{}:
			for _, tcwd := range typedCommand {
				switch typedCommandListItem := tcwd.(type) {
				case string:
					c := strings.SplitN(typedCommandListItem, " ", 2)
					cmd, arguments := c[0], c[1]
					output := execute(cmd, arguments)
					fmt.Println(output)
				default:
					fmt.Println("Unsupported typeof task")
					fmt.Printf("Task type: %T\n", typedCommandListItem)
					fmt.Printf("Task value: %v\n", typedCommandListItem)
				}
			}
		default:
			fmt.Println("Unsupported typeof task")
			fmt.Printf("Task type: %T\n", typedCommand)
			fmt.Printf("Task value: %v\n", typedCommand)
		}
		baner = fmt.Sprintf(`Finished task "%v" ...`, arg)
		fmt.Println(baner)
	}
}
