package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	workDir := kingpin.Flag("wd", "Working Directory").Short('w').Default("./").String()
	incPattern := kingpin.Flag("inc", "Included file(s)").Short('i').Default(".*").String()
	excPattern := kingpin.Flag("exc", "Excluded file(s)").Short('e').Default(`^\.{1,2}$`).String()
	outputFile := kingpin.Flag("out", "Output file name").Short('o').String()
	appendResult := kingpin.Flag("append", "Append output result").Short('a').Bool()

	var output string
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.Parse()

	dumOutput := *outputFile == ""

	if !dumOutput {
		fmt.Printf("Working Dir: %s\nInclude: %s\nExclude: %s\nOutput: %s\n", *workDir, *incPattern, *excPattern, *outputFile)
	}

	incRegex, err := regexp.Compile(*incPattern)
	checkError(err, "Include regex compile error")
	excRegex, err := regexp.Compile(*excPattern)
	checkError(err, "Exclude regex compile error")

	files, err := ioutil.ReadDir(*workDir)

	checkError(err, "Working dir read error")

	for _, f := range files {
		if !f.IsDir() && f.Name() != *outputFile && incRegex.MatchString(f.Name()) && !excRegex.MatchString(f.Name()) {
			data, err := readFile(f.Name(), *workDir)
			checkError(err, "Read file error")
			output = output + "\n" + string(data)
		}
	}

	if dumOutput {
		fmt.Println(output)
	} else {
		err = writeFile(*outputFile, *workDir, output, *appendResult)
		checkError(err, "Write file error")
	}

}

func readFile(filename string, dir string) (data []byte, err error) {
	path := filepath.Join(dir, filename)
	data, err = ioutil.ReadFile(path)
	return
}

func writeFile(filename string, dir string, data string, append bool) (err error) {
	path := filepath.Join(dir, filename)
	var f *os.File
	if append {
		f, err = os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	} else {
		f, err = os.OpenFile(path, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	}

	if err != nil {
		return
	}
	defer f.Close()

	_, err = f.WriteString(data)

	return
}

func checkError(err error, sufix string) {
	if err != nil {
		log.Fatalf("%s\nerror: %s", sufix, err.Error())
	}
}
