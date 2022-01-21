package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/evilmonkeyinc/jsonpath"
)

var (
	// Arch build identifier
	Arch string = ""
	// Command is the expected name of the build
	Command string = "jsonpath"
	// OS build identifier
	OS string = ""
	// Version build identifier
	Version string = "dev"

	errSelectorNotSpecified error = fmt.Errorf("selector not specified. expected -selector option or passed as argument")
	errJSONDataNotSpecified error = fmt.Errorf("json data not specified. expected -jsondata, or -input options or passed as argument")
)

const (
	cmdHelp    string = "help"
	cmdVersion string = "version"
)

func main() {
	flagset := flag.NewFlagSet("", flag.ContinueOnError)
	selectorPtr := flagset.String("selector", "", "a valid JSONPath selector")
	jsondataPtr := flagset.String("jsondata", "", "the json data to parse)")
	inputPtr := flagset.String("input", "", "optional path to a file that includes the json data to query")
	outputPtr := flagset.String("output", "", "optional path to a file that the result ")

	if err := flagset.Parse(os.Args[1:]); err != nil {
		if !errors.Is(err, flag.ErrHelp) {
			outputError(errSelectorNotSpecified)
			return
		}
		printHelp(flagset)
		return
	}

	args := flagset.Args()
	firstArg := ""
	if len(args) > 0 {
		firstArg = args[0]
	}

	switch firstArg {
	case cmdHelp:
		printHelp(flagset)
		break
	case cmdVersion:
		fmt.Printf("version %s %s/%s\n", Version, OS, Arch)
		break
	default:
		nextArg := 0

		selector := *selectorPtr
		if selector == "" && len(args) > nextArg {
			selector = args[nextArg]
			nextArg++
		}

		if selector == "" {
			outputError(errSelectorNotSpecified)
			return
		}

		compiled, err := jsonpath.Compile(selector)
		if err != nil {
			outputError(err)
			return
		}

		jsondata := *jsondataPtr
		if jsondata == "" && len(args) > nextArg {
			jsondata = args[nextArg]
			nextArg++
		}

		if jsondata == "" && *inputPtr != "" {
			bytes, err := loadFileContents(*inputPtr)
			if err != nil {
				outputError(err)
				return
			}
			jsondata = string(bytes)
		}
		if jsondata == "" {
			outputError(errJSONDataNotSpecified)
			return
		}

		result, err := compiled.QueryString(jsondata)
		if err != nil {
			outputError(err)
			return
		}

		jsonOutput, err := json.Marshal(result)
		if err != nil {
			outputError(err)
			return
		}

		if *outputPtr != "" {
			if err := saveFileContents(*outputPtr, jsonOutput); err != nil {
				outputError(err)
			}
			os.Exit(0)
			return
		}

		fmt.Printf("%s\n", string(jsonOutput))
		break
	}
	os.Exit(0)
}

func printHelp(flagset *flag.FlagSet) {
	fmt.Printf("%s is a tool for querying json data using a JSONPath selector\n\n", Command)
	fmt.Printf("Usage:\n\n")
	fmt.Printf("  %s '[selector]' '[jsondata]'\n", Command)
	fmt.Printf("\nExample:\n\n")
	fmt.Printf(`  %s '$[*].key' '[{"key":"show this"},{"key":"and this"},{"other":"but not this"}]'`+"\n", Command)
	fmt.Printf(`  > ["show this","and this"]` + "\n")
	fmt.Printf("\nOptions:\n\n")
	flagset.PrintDefaults()
	os.Exit(0)
}

func outputError(err error) {
	fmt.Printf("failed: %s\n", err.Error())
	os.Exit(1)
}

func loadFileContents(filename string) ([]byte, error) {
	jsonFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}
	return byteValue, nil
}

func saveFileContents(filename string, content []byte) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = file.Write(content); err != nil {
		return err
	}
	return nil
}
