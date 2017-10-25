package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/cytec/releaseparser"
	"github.com/ttacon/chalk"
)

var renameFiles = flag.Bool("files", false, "rename files as well")
var renameDirs = flag.Bool("rename", true, "rename directorys")
var help = flag.Bool("help", false, "echo usage command")
var testrun = flag.Bool("test", false, "only echo renmaed values but dont rename")
var mformat = flag.String("mformat", "{{.Title}}{{if .Year}} ({{.Year}}){{end}}", "format for movie direcotries")
var tvformat = flag.String("tvformat", "{{.Title}} S{{.Season}}E{{.Episode}}", "format for Series direcotries")
var formats = flag.Bool("formats", false, "print available formats for the renamer")

var videoextensions = []string{".mkv", ".mp4", ".avi", ".m4v", ".webm", ".flv", ".mov", ".wmv"}

func isVideoExtension(s string) bool {
	mimet := mime.TypeByExtension(s)

	if strings.Contains(mimet, "video") {
		return true
	}
	for _, a := range videoextensions {
		if a == s {
			return true
		}
	}
	return false
}

func main() {

	flag.Parse()

	if *help {
		fmt.Printf("usage: %s direcotry\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *formats {
		var examplename = "Release.Name.Uncut.2010.German.Dubbed.AC3.BluRay.1080p.x264-GroupName"
		example := releaseparser.Parse(examplename)
		fmt.Printf(chalk.Bold.TextStyle("available formats for example releasename: %s\n\n"), examplename)
		fmt.Printf("{{.Input}}\t => \t %s\n", example.Input)
		fmt.Printf("{{.Type}}\t => \t %s\n", example.Type)
		fmt.Printf("{{.Title}}\t => \t %s\n", example.Title)
		fmt.Printf("{{.Year}}\t => \t %d\n", example.Year)
		fmt.Printf("{{.Source}}\t => \t %s\n", example.Source)
		fmt.Printf("{{.Resolution}}\t => \t %s\n", example.Resolution)
		fmt.Printf("{{.Codec}}\t => \t %s\n", example.Codec)
		fmt.Printf("{{.Language}}\t => \t %s\n", example.Language)
		fmt.Printf("{{.Audio}}\t => \t %s\n", example.Audio)
		fmt.Printf("{{.Group}}\t => \t %s\n", example.Group)
		fmt.Printf("{{.Season}}\t => \t %d\n", example.Season)
		fmt.Printf("{{.Episode}}\t => \t %d\n", example.Episode)
		fmt.Printf("{{.Uncut}}\t => \t %t\n", example.Uncut)
		fmt.Printf("\nfor a full list of available formatters see: https://godoc.org/github.com/cytec/releaseparser#Release\n")
		os.Exit(1)
	}

	if len(flag.Args()) <= 0 {
		fmt.Printf("usage: %s direcotry\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *testrun {
		fmt.Printf(chalk.Red.Color("NOTICE: running in test mode, no actual renaming is done\n"))
	}

	mformatPointer := *mformat
	tvformatPointer := *tvformat
	mtemplate := template.Must(template.New("moviename").Parse(mformatPointer))
	tvtemplate := template.Must(template.New("moviename").Parse(tvformatPointer))

	var mydirs = flag.Args()
	// var currentDir = mydirs[0]
	//

	for _, currentDir := range mydirs {
		files, err := ioutil.ReadDir(currentDir)
		if err != nil {
			panic(err)
		}

		for _, f := range files {
			tpl := bytes.Buffer{}

			if f.IsDir() {
				r := releaseparser.Parse(f.Name())
				if r.Title != "" {

					err = nil
					if r.Type == "movie" {
						err = mtemplate.Execute(&tpl, r)
					} else if r.Type == "tvshow" {
						err = tvtemplate.Execute(&tpl, r)
						err = errors.New("tv releases are not supported")
						fmt.Printf(chalk.Red.Color("%s tv releases are not supported, skipping\n"), f.Name())
					}

					if err == nil {

						fmt.Printf("rename %s => %s\n", chalk.Yellow.Color(r.Input), chalk.Green.Color(tpl.String()))

						oldDirName := filepath.Join(currentDir, f.Name())
						newDirName := filepath.Join(currentDir, tpl.String())
						if *renameFiles {
							// fmt.Printf("checking dir '%s' for files to rename\n", f.Name())
							//check for movie file in current dir and rename it to dirname
							files, err := ioutil.ReadDir(oldDirName)
							if err != nil {
								panic(err)
							}
							for _, fn := range files {
								if !fn.IsDir() {
									fpath := filepath.Join(currentDir, f.Name(), fn.Name())
									ext := filepath.Ext(fn.Name())
									fnameNew := strings.Trim(f.Name(), " ") + ext
									fpathNew := filepath.Join(currentDir, f.Name(), fnameNew)
									if isVideoExtension(ext) && !strings.Contains(strings.ToLower(fpath), "sample") {
										// fmt.Printf("%s seems to be a non sample video file.", fn.Name())

										fmt.Printf("rename '%s' => '%s'\n", chalk.Yellow.Color(fn.Name()), chalk.Green.Color(fnameNew))
										if !*testrun {
											os.Rename(fpath, fpathNew)
										}
									}
								}
							}
						}

						if !*testrun {
							os.Rename(oldDirName, newDirName)
						}

					}
				}
			}
		}
	}
}
