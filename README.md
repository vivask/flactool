# flactool

Multithreaded, batch tool for converting, concatenating and splitting audio files in flac, ape, wav formats


To use the program, you must first install the following packages:
- Monkey's Audio Codec; https://github.com/fernandotcl/monkeys-audio;
- shntool;
- cuetools;
- sox.

# Uasage:

    flactool [OPTION] 
-  -c concat all flac files in dir to one flac file
-  -s split flac or ape files in dir
-  -d "path"
-  -f "file"
-  -h help
-  -n output file name - number
-  -p Num core, default 4 (default 4)
-  -r rename ape file before convert
-  -R remove source after operation
-  -v verbose

# Examples:
1. Convert all ape files from ~/apedir (with subdirectories) to flac

    flactool -d ~/apedir 

2. All ape and wav files from the current directory (with subdirectories) are split with cue (if there is a cue file with a name similar to ape or wav) with conversion to flac. With the subsequent removal of the original ape and wav files

    flactool -s -R

# Build 
To build from source code you need make and Docker, then run the following commands:

    git clone https://github.com/vivask/flactool.git
    cd flactool
    install.sh