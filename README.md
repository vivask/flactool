# flactool

To use the program, you must first install the following packages:
- Monkey's Audio Codec; https://github.com/fernandotcl/monkeys-audio;
- shntool;
- cuetools;

Uasage:

    flactool [OPTION] 
    -  -R remove source after operation
    -  -c concat all flac files in dir to one flac file
    -  -d "path"
    -  -f "file"
    -  -h help
    -  -n output file name - number
    -  -p Num core, default 4 (default 4)
    -   -r rename ape file before convert
    -   -s split flac or ape files in dir
    -   -v verbose

Examples:
1. Convert all ape files from ~/apefiles directory (with subdirectories) to flac

    flactool -d ~/apedir 

2. All ape and wav files from the current directory (with subdirectories) are split with cue (if there is a cue file with a name similar to ape or wav) with conversion to flac. With the subsequent removal of the original ape and wav files

    flactool -s -R