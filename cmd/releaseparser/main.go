package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/cytec/releaseparser"
	"github.com/ttacon/chalk"
)

var help = flag.Bool("help", false, "echo usage command")
var jsonout = flag.Bool("json", false, "dont rename just parse the releases and output JSON")
var stdin = flag.Bool("stdin", false, "input from stdin most usefull with --json")

func printRelease(r *releaseparser.Release) {
	v := reflect.ValueOf(r).Elem()
	fmt.Printf("'%s' parsed to:\n", r.Input)
	for i := 0; i < v.NumField(); i++ {
		name := v.Type().Field(i).Name
		// fmt.Printf("%v", name)
		if name == "start" || name == "end" || name == "parts" {
			continue
		}
		value := v.Field(i).Interface()
		fname := v.Field(i).Type().Name()
		if fname == "bool" && value != false {
			fmt.Printf("\t%s:\t%t\n", name, value)
		} else if fname == "string" && value != "" {
			fmt.Printf("\t%s:\t%s\n", name, value)
		} else if fname == "int" && value != 0 {
			fmt.Printf("\t%s:\t%d\n", name, value)
		}
	}
}

func main() {

	flag.Parse()

	if *help {
		fmt.Printf("usage: %s direcotry\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	if len(flag.Args()) <= 0 && !*stdin {
		fmt.Printf("usage: %s direcotry\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	var output []*releaseparser.Release
	var mydirs = flag.Args()
	// var currentDir = mydirs[0]
	//

	if *stdin {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		r := releaseparser.Parse(text)

		if r.Title != "" {
			if !*jsonout {
				printRelease(r)
			} else {
				output = append(output, r)
			}
		}
	}

	for _, currentDir := range mydirs {

		if !*jsonout {
			fmt.Printf(chalk.Bold.TextStyle("scanning directory '%s' for releases\n"), currentDir)
		}
		files, err := ioutil.ReadDir(currentDir)
		if err != nil {
			panic(err)
		}

		for _, f := range files {
			name := f.Name()
			if !f.IsDir() {
				name = strings.TrimSuffix(name, filepath.Ext(f.Name()))
			}
			r := releaseparser.Parse(name)
			// append current parser result to output and go to next one
			if *jsonout {
				if r.Title != "" {
					output = append(output, r)
				}
				continue
			}
			if r.Title != "" {
				printRelease(r)
			}

		}
	}
	//return the json output
	if *jsonout {
		json, err := json.MarshalIndent(output, " ", " ")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(json))
	}
}
