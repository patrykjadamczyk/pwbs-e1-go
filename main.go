package main

import (
	"encoding/json"
	"fmt"
	"io"
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
	version:      "0.9.1.0",
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

type JsonRecord struct {
	commands []string
}

func JsonDecode(r io.Reader) (x *JsonRecord, err error) {
	x = new(JsonRecord)
	if err = json.NewDecoder(r).Decode(x); err != nil {
		return
	}
	return
}

func readJson(filename string) *JsonRecord {
	file, err := os.Open(filename)
	if err != nil {
		return nil
	}
	data, jsonErr := JsonDecode(file)
	if jsonErr != nil {
		return nil
	}
	return data
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
	for i := 0; i < len(args); i++ {
		baner := fmt.Sprintf(`Executing task "%v" ...`, args[i])
		fmt.Println(baner)
		command := JsonData.commands[args[i]]
		c := strings.SplitN(command, " ", 2)
		cmd, args := c[0], c[1]
		output := execute(cmd, args)
		fmt.Println(output)
		baner = fmt.Sprintf(`Finished task "%v" ...`, args[i])
		fmt.Println(baner)
	}
}
