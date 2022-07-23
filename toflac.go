package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sync/errgroup"
)

//convert dsf/ape/wav file to flac
func FileToFlac(shntool, ffmpeg, input string, verbose bool) {
	//prepare shell command
	var cmd string
	ext := filepath.Ext(input)
	if ext == ".dsf" {
		out := replaceExtToFlac(input)
		cmd = fmt.Sprintf("%s -i \"%s\" -af \"lowpass=24000, volume=6dB\" -sample_fmt s32 -ar 96000 \"%s\"", ffmpeg, input, out)
	} else {
		cmd = fmt.Sprintf("%s conv -o flac \"%s\" -d \"%s\"", shntool, input, filepath.Dir(input))
	}
	//convert file
	err, out, errout := Shellout(cmd)
	execVerbose(err, out, errout, verbose)
}

//convert dsf/ape/wav files to flac
func DirToFlac(shntool, ffmpeg, dir string, parallel uint, concat, remove, verbose bool) (dirs []string, err error) {
	//search files
	list, err := getFilesFromDir(dir, ".ape", ".wav", ".dsf")
	if len(list) == 0 {
		if !concat {
			return dirs, fmt.Errorf("audio files not found")
		}
	} else {

		pathes, dirs := prepareFiles(list, true)

		StartSpinner()
		g, _ := errgroup.WithContext(context.Background())
		g.SetLimit(int(parallel))
		for _, path := range dirs {
			path := path
			for _, file := range pathes[path] {
				input := fmt.Sprintf("%s/%s", path, file)
				ext := filepath.Ext(file)
				g.Go(func() error {
					//prepare shell command
					var cmd string
					if ext == ".dsf" {
						out := replaceExtToFlac(input)
						cmd = fmt.Sprintf("%s -i \"%s\" -af \"lowpass=24000, volume=6dB\" -sample_fmt s32 -ar 96000 \"%s\"", ffmpeg, input, out)
					} else {
						cmd = fmt.Sprintf("%s conv -o flac \"%s\" -d \"%s\"", shntool, input, path)
					}
					cmdVerbose(cmd, verbose)
					//convert file
					err, stdout, errout := Shellout(cmd)
					execVerbose(err, stdout, errout, verbose)
					if err != nil {
						return err
					}
					//remove sourse
					if remove {
						err = os.Remove(input)
					}
					return err
				})
			}
		}
		err = g.Wait()
		StopSpinner()
		return dirs, err
	}
	return dirs, nil
}
