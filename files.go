package main

import (
	"os"
	"path/filepath"
)

type FileList map[string][]string

func getFilesFromDir(dir string, extSet ...string) (list []string, err error) {
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

func prepareFiles(files []string) (list FileList) {
	list = make(FileList)
	for _, path := range files {
		dir := filepath.Dir(path)
		name := filepath.Base(path)
		list[dir] = append(list[dir], name)
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

func isInclude(item string, set []string) bool {
	for _, i := range set {
		if item == i {
			return true
		}
	}
	return false
}

func getSplitFilesFromDir(dir string, extSet ...string) (list []string, err error) {
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
