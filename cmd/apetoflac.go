package cmd

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/sync/errgroup"
)

func FileToFlac(shntool, input string, verbose bool) {
	task := fmt.Sprintf("%s conv -o flac \"%s\" -d \"%s\"", shntool, input, filepath.Dir(input))
	err, out, errout := Shellout(task)
	if verbose {
		if err != nil {
			log.Printf("error: %v\n", err)
		}
		fmt.Println("--- stdout ---")
		fmt.Println(out)
		fmt.Println("--- stderr ---")
		fmt.Println(errout)
	}
}

func DirToFlac(shntool, dir string, parallel uint, outnum, concat, rename, remove, verbose bool) (worked bool, err error) {
	worked = false
	list, err := getFilesFromDir(dir, ".ape", ".wav")
	if err != nil {
		return worked, fmt.Errorf("error: %w", err)
	}
	if len(list) == 0 {
		if !concat {
			return worked, fmt.Errorf("ape or wav files not found")
		}
	} else {

		pathes := prepareFiles(list)
		for path, files := range pathes {
			fmt.Println(path)
			for _, file := range files {
				fmt.Println(file)
			}
			fmt.Println()
		}

		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Convert this files to flac?")
		fmt.Println("[Yes/no]?")
		t, _ := reader.ReadString('\n')
		t = strings.ToLower(t)
		if t == "y\n" || t == "yes\n" || t == "\n" {
			StartSpinner()
			g, _ := errgroup.WithContext(context.Background())
			g.SetLimit(int(parallel))
			var mu sync.Mutex
			for path, files := range pathes {
				path := path
				for i, file := range files {
					i := i + 1
					input := fmt.Sprintf("%s/%s", path, file)
					newName := input
					if rename {
						newName = fmt.Sprintf("%s/%04d.ape", path, i+1)
						err := os.Rename(input, newName)
						if err != nil {
							return worked, fmt.Errorf("rename error: %w", err)
						}
					}
					g.Go(func() error {
						task := fmt.Sprintf("%s conv -o flac \"%s\" -d \"%s\"", shntool, newName, path)
						err, out, errout := Shellout(task)
						if verbose {
							if err != nil {
								log.Printf("error: %v\n", err)
							}
							fmt.Println("--- stdout ---")
							fmt.Println(out)
							fmt.Println("--- stderr ---")
							fmt.Println(errout)
						}

						if err == nil {
							mu.Lock()
							worked = true
							mu.Unlock()
							if remove {
								err = os.Remove(newName)
							}
						}
						return err
					})
				}
			}
			err = g.Wait()
			StopSpinner()
			return worked, err
		}
	}
	return worked, nil
}
