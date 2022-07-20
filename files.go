package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type FileList map[string][]string

//Recursive search for files in a given directory by a list of extensions
func getFilesFromDir(dir string, extSet ...string) (list []string, err error) {
	StartSpinner()
	defer StopSpinner()
	err = filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			ext := filepath.Ext(path)
			if !info.IsDir() && isInclude(ext, extSet) {
				list = append(list, path)
			}
			return nil
		})
	return
}

//Generating a list of files in the FileList format
func prepareFiles(files []string) (list FileList) {
	list = make(FileList)
	for _, path := range files {
		dir := filepath.Dir(path)
		name := filepath.Base(path)
		list[dir] = append(list[dir], name)
	}

	var rm []string
	for path, files := range list {
		fmt.Println(path)
		for _, file := range files {
			fmt.Println(file)
		}
		fmt.Println()
		t := "y\n"
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Use this files?")
		fmt.Println("[Yes/no]?")
		t, _ = reader.ReadString('\n')
		t = strings.ToLower(t)
		if t == "n\n" || t == "no\n" {
			rm = append(rm, path)
		}
	}

	for _, r := range rm {
		delete(list, r)
	}

	return
}

func replaceExtToFlac(fName string) string {
	ext := filepath.Ext(fName)
	return fName[0:len(fName)-len(ext)] + ".flac"
}

func replaceExtToCue(fName string) string {
	ext := filepath.Ext(fName)
	return fName[0:len(fName)-len(ext)] + ".cue"
}

//Checking if an element is in a set
func isInclude(item string, set []string) bool {
	for _, i := range set {
		if item == i {
			return true
		}
	}
	return false
}

//Recursive search for files in a given directory according to the list of extensions if there is a cue file.
func getSplitFilesFromDir(dir string, extSet ...string) (list []string, err error) {
	StartSpinner()
	defer StopSpinner()
	err = filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			ext := filepath.Ext(path)
			if !info.IsDir() && isInclude(ext, extSet) {
				cue := replaceExtToCue(path)
				if _, exist := os.Stat(cue); exist == nil {
					list = append(list, path)
				}

			}
			return nil
		})
	return
}
