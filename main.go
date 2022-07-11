package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/vivask/flactool/cmd"
)

const (
	SHNTOOL    = "shntool"
	SOX        = "sox"
	CUETAG     = "cuetag"
	SEARCHPATH = "/usr/bin /usr/sbin /usr/local/bin usr/local/sbin"
)

type Arguments map[string]interface{}

func parseArgs() Arguments {
	var fName, dirName string
	var numNameOut, concat, split, rename, remove, verbose, help bool
	var parallel uint
	flag.StringVar(&fName, "f", "", `-f "file"`)
	flag.StringVar(&dirName, "d", "", `-d "path"`)
	flag.UintVar(&parallel, "p", 4, "-p Num core, default 4")
	flag.BoolVar(&numNameOut, "n", false, "-n output file name - number")
	flag.BoolVar(&concat, "c", false, "-c concat all flac files in dir to one flac file")
	flag.BoolVar(&split, "s", false, "-s split flac or ape files in dir")
	flag.BoolVar(&rename, "r", false, "-r rename ape file before convert")
	flag.BoolVar(&remove, "R", false, "-R remove source after operation")
	flag.BoolVar(&verbose, "v", false, "-v verbose")
	flag.BoolVar(&help, "h", false, "-h help")
	flag.Parse()

	args := Arguments{}
	args["file"] = fName
	args["dir"] = dirName
	args["parallel"] = parallel
	args["outnum"] = numNameOut
	args["concat"] = concat
	args["split"] = split
	args["rename"] = rename
	args["remove"] = remove
	args["verbose"] = verbose
	args["help"] = help
	return args
}

func searchNeedPkg() (sox, shntool, cuetag string, err error) {
	pathes := strings.Split(SEARCHPATH, " ")
	for _, path := range pathes {
		fName := fmt.Sprintf("%s/%s", path, SOX)
		if _, exist := os.Stat(fName); exist == nil {
			sox = fName
		}
		fName = fmt.Sprintf("%s/%s", path, SHNTOOL)
		if _, exist := os.Stat(fName); exist == nil {
			shntool = fName
		}
		fName = fmt.Sprintf("%s/%s", path, CUETAG)
		if _, exist := os.Stat(fName); exist == nil {
			cuetag = fName
		}
	}
	if len(shntool) == 0 {
		err = fmt.Errorf("shntool not found. need install shntool")
	}
	if len(sox) == 0 {
		err = fmt.Errorf("sox not found. need install sox")
	}
	if len(cuetag) == 0 {
		err = fmt.Errorf("cuetag not found. need install cuetag")
	}

	return
}

func main() {
	args := parseArgs()
	help := args["help"].(bool)

	//usage
	if len(os.Args) > 1 && (os.Args[1] == "?" || help) {
		flag.PrintDefaults()
		return
	}

	//checking the availability of third-party packages
	sox, shntool, cuetag, err := searchNeedPkg()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//check for incompatible arguments
	verbose := args["verbose"].(bool)
	file := args["file"].(string)
	dir := args["dir"].(string)
	if len(file) != 0 && len(dir) != 0 {
		fmt.Println("cannot use flags -f with -d together")
		os.Exit(1)
	}

	//current directory
	if len(file) == 0 && len(dir) == 0 {
		dir, err = os.Getwd()
		if err != nil {
			fmt.Printf("error: %v\n", err)
			os.Exit(1)
		}
	}

	//split flac, ape, wav files according to cue by directories
	split := args["split"].(bool)
	parallel := args["parallel"].(uint)
	rename := args["rename"].(bool)
	remove := args["remove"].(bool)
	if split {
		err = cmd.SplitApeOrFlac(shntool, cuetag, dir, parallel, rename, remove, verbose)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return
	}

	//convert input file to flac
	if len(file) != 0 {
		cmd.FileToFlac(shntool, file, verbose)
		return
	}

	//converting input files by directories to flac
	worked := false
	concat := args["concat"].(bool)
	if len(dir) != 0 {
		outnum := args["outnum"].(bool)
		worked, err = cmd.DirToFlac(shntool, dir, parallel, outnum, concat, rename, remove, verbose)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	//concatenation of flac, ape, wav files by directories with conversion to flac
	if concat {
		err = cmd.ConcatFlacs(sox, dir, parallel, rename, remove, verbose, worked)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

}
