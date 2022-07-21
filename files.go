package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
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
func prepareFiles(files []string, one bool) (list FileList, keys []string) {
	var rm []string
	list = make(FileList)
	//copy from slice to map
	for _, path := range files {
		dir := filepath.Dir(path)
		name := filepath.Base(path)
		list[dir] = append(list[dir], name)
	}

	if !one {
		// mark dirs with one file
		for path, files := range list {
			if len(files) < 2 {
				rm = append(rm, path)
			}
		}
		// remove dirs with one file
		for _, r := range rm {
			delete(list, r)
		}
	}

	//create sort slice for map
	keys = make([]string, 0, len(list))
	for k := range list {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	rm = rm[:0]
	for _, path := range keys {
		fmt.Println(path)
		for _, file := range list[path] {
			fmt.Println(file)
		}
		fmt.Println()
		t := "y\n"
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Use this files [Yes/no]?")
		t, _ = reader.ReadString('\n')
		t = strings.ToLower(t)
		if !(t == "y\n" || t == "yes\n" || t == "\n") {
			rm = append(rm, path)
		}
	}

	//delete not used
	for _, r := range rm {
		delete(list, r)
	}

	//create sort slice for map
	keys = make([]string, 0, len(list))
	for k := range list {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return
}

// replace file extension to flac
func replaceExtToFlac(fName string) string {
	ext := filepath.Ext(fName)
	return fName[0:len(fName)-len(ext)] + ".flac"
}

// replace file extension to cue
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

// return name last dir
func getLastDir(path string) string {
	split := strings.Split(path, "/")
	cnt := len(split)
	if cnt > 0 {
		return split[cnt-1]
	}
	return path
}

// return previous dir file path
func getParentPath(path string) string {
	dir := filepath.Dir(path)
	parent := filepath.Dir(dir)
	if parent == "/" {
		return fmt.Sprintf("/%s", filepath.Base(path))
	}
	return fmt.Sprintf("%s/%s", parent, filepath.Base(path))
}
