package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

var (
	input      string
	inputFile  string
	outputFile string
	envPattern string
)

func init() {
	flag.StringVar(&input, "input", "", "input template")
	flag.StringVar(&inputFile, "file", "", "input template")
	flag.StringVar(&outputFile, "o", "", "output filename")
	flag.StringVar(&envPattern, "env-pattern", "${ENV}",
		"placeholder for env variable from input, formatted as <prefix>ENV<postfiv> where 'ENV' stands for env variable name")
}

func main() {
	flag.Parse()

	inputText, err := readInput()
	if err != nil {
		panic(err)
	}

	output := render(inputText, envPattern)

	if err = writeOutput(output); err != nil {
		panic(err)
	}
}

func render(text, envPattern string) string {
	// create regex for env placeholders
	holder := strings.Replace(envPattern, "ENV", ".*", -1)
	holderRegex := createRegex(holder)

	// create regex for placeholders' brackets
	brackets := strings.Replace(envPattern, "ENV", "|", -1)
	bracketsRegex := createRegex(brackets)

	// find all env vars placeholders in input
	foundEnvKeys := holderRegex.FindAllString(text, -1)

	subs := map[string]string{}
	for _, envKey := range foundEnvKeys {
		envName := bracketsRegex.ReplaceAllString(envKey, "")
		val := os.Getenv(envName)
		subs[envKey] = val
	}

	for k, v := range subs {
		text = strings.Replace(text, k, v, -1)
	}

	return text
}

func readInput() (string, error) {
	if inputFile != "" {
		bytes, err := ioutil.ReadFile(inputFile)
		if err != nil {
			return "", errors.Wrap(err, "Failed to read input file")
		}
		return string(bytes), nil
	}

	if input == "" {
		return "", errors.New("Please provide some input")
	}

	return input, nil
}

func writeOutput(text string) error {
	if outputFile != "" {
		err := ioutil.WriteFile(outputFile, []byte(text), 0644)
		if err != nil {
			return errors.Wrap(err, "Failed to save output file")
		}
		return nil
	}

	fmt.Println(text)
	return nil
}

func createRegex(pattern string) *regexp.Regexp {
	regexPattern := fmt.Sprintf("\\%s", pattern)
	return regexp.MustCompile(regexPattern)
}
