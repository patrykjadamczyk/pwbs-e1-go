package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type PwbsInfo struct {
	version      string
	edition      string
	versionArray [4]int
}

var pwbs PwbsInfo = PwbsInfo{
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

type PWBSConfigFile struct {
	Commands map[string]string `json:"commands"`
}

func readJson(filename string) PWBSConfigFile {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error in reading File", err)
		return PWBSConfigFile{Commands: map[string]string{}}
	}
	jsonData := PWBSConfigFile{}
	jsonErr := json.Unmarshal(data, &jsonData)
	if jsonErr != nil {
		fmt.Println("Error in parsing JSON", jsonErr)
		return PWBSConfigFile{Commands: map[string]string{}}
	}
	return jsonData
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
	JsonData := readJson("pwbs.json")
	for _, arg := range args {
		baner := fmt.Sprintf(`Executing task "%v" ...`, arg)
		fmt.Println(baner)
		command := JsonData.Commands[arg]
		c := strings.SplitN(command, " ", 2)
		cmd, arguments := c[0], c[1]
		output := execute(cmd, arguments)
		fmt.Println(output)
		baner = fmt.Sprintf(`Finished task "%v" ...`, arg)
		fmt.Println(baner)
	}
}
