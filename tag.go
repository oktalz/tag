package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	input "github.com/tcnksm/go-input"
)

const programName = `  ____ _ _     _                  
 / ___(_) |_  | |_ __ _  __ _ ___ 
| |  _| | __| | __/ _` + "`" + ` |/ _` + "`" + ` / __|
| |_| | | |_  | || (_| | (_| \__ \
 \____|_|\__|  \__\__,_|\__, |___/
                        |___/     

`

func main() {
	fmt.Println(programName)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	dir, oserr := filepath.Abs(filepath.Dir(os.Args[0]))
	if oserr != nil {
		log.Fatal(oserr)
	}
	log.Println(dir)

	lastTagMaster := bash("git tag --merged | tail -1")
	//log.Println("last tag", lastTagMaster)
	lastTagOverall := bash("git tag | tail -1")
	//log.Println("last tag overall", lastTagOverall)
	lastTag := lastTagMaster

	var newTag string
	if strings.Contains(lastTagMaster, "-ee") {
		log.Println("Enterprise mode")
		var tag string
		if lastTagMaster != lastTagOverall {
			log.Println("EE: new upstream tag:", lastTagOverall)
			tag = lastTagOverall + "-ee0"
		} else {
			tag = lastTagMaster
		}
		dataL := strings.Split(tag, "e")
		lastIndex := dataL[len(dataL)-1]
		i2, err := strconv.ParseInt(lastIndex, 10, 64)
		if err == nil {
			lastIndex = strconv.Itoa(int(i2 + 1))
		}
		//log.Println(dataL)
		newTag = fmt.Sprintf("%see%s", dataL[0], lastIndex)
	} else {
		dataL := strings.Split(lastTag, ".")
		lastIndex := dataL[len(dataL)-1]
		i2, err := strconv.ParseInt(lastIndex, 10, 64)
		if err == nil {
			lastIndex = strconv.Itoa(int(i2 + 1))
		}
		newTag = fmt.Sprintf("%s.%s.%s", dataL[0], dataL[1], lastIndex)
	}
	fmt.Println(lastTagMaster, lastTagOverall, "=>", newTag)

	ui := &input.UI{
		Writer: os.Stdout,
		Reader: os.Stdin,
	}
	inputTag, err := ui.Ask("Define new tag:", &input.Options{
		Default:   newTag,
		Required:  true,
		Loop:      true,
		HideOrder: true,
	})
	if err != nil {
		return
	}

	logs := bash(fmt.Sprintf("git log %s..HEAD --pretty=format:\"%%h - %%s (%%aI) <%%an>\"", lastTag))
	fmt.Println("##########################################################")
	fmt.Println(logs)
	fmt.Println("##########################################################")
	shell("git", "tag", "-a", "-m", logs, inputTag)
}

func shell(name string, params ...string) string {
	out, err := exec.Command(name, params...).Output()
	if err != nil {
		log.Fatalf("exec.Command failed with %s, %s, %v\n", err, name, params)
	}
	return string(out)
}

func bash(cmd string) string {
	log.Println("EXECUTING", cmd)
	result := exec.Command("bash", "-c", cmd)
	out, err := result.CombinedOutput()
	if err != nil {
		log.Fatalf("exec.Command) failed with %s, %s\n", err, cmd)
	}
	res := string(out)
	if strings.HasSuffix(res, "\n") {
		res = res[:len(res)-1]
	}
	return res
}
