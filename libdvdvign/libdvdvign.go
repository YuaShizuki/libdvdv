package libdvdvign

/*
* This module builds a ignore file list during initialization. Ounce initialized, one
* can call check to check if a file is suppose to be ignored or not. This package
* reads the .libdvdvignore file present in the project's main directory. If it cant
* find a .libdvdvignore file, it creates a new one and adds shell globs patterns from 
* a .gitignore file if present.
* 

* Rules for ignoring.
*   -> lines starting with '#' are treated as comments
*   -> lines starting with '/' are globs for the main project directory.
*   -> lines ending with '/' will match only directories.
*   -> lines starting with '!' are negated from the ignore list.
*   -> all other lines are treated as shell glob patterns.
*   -> binary executable are ignored by default.
*
* please refer to http://en.wikipedia.org/wiki/Glob_(programming) for 
* understanding shell glob patterns.
*/

//import "os"
import "../libdvdvutil"
import "io/ioutil"

/* All modules have a log function*/
var libdvdvign_log func(a ...interface{});

/*Ignore file list*/
var ignore [][]string

/* Data structure to represent parsed ignore file, based on the rules mentioned */
/* above*/
var shell_globs []string;
var main_dir_shell_globs []string;
var dir_shell_globs []string
var not_shell_globs [3][]string

/*Ignore file message*/
var ignore_file_message string =
`# files specified hear are intentionally untracked files that libdevdev would ignore.
# Rules for ignoring.
#   -> lines starting with '#' are treated as comments
#   -> lines starting with '/' are globs for the main project directory.
#   -> lines ending with '/' will match only directories.
#   -> lines starting with '!' are negated from the ignore list.
#   -> all other lines are treated as shell glob patterns.
#   -> binary executables are ignored by default
# please refer to http://en.wikipedia.org/wiki/Glob_(programming) for 
# understanding shell glob patterns.
.*
*.o
`

func Init(log func(a ...interface{})) error {
    libdvdvign_log = log;
    libdvdvign_log("initializing libdvdv-ignore in current directory");
    /*Detect if there is a .libdvdvignore file in the current dir*/
    if libdvdvutil.PathExist(".libdvdvignore") {
        lines, err := ioutil.ReadFile(".libdvdvignore");
        if err != nil {
            libdvdvign_log("error-> unable to read .libdvdvignore file");
            return err;
        }
        parseIgnoreFile(string(lines));
    } else {
        if libdvdvuitl.PathExist(".gitignore") {
            lines, err := ioutil.ReadFile(".gitignore");
            if err != nil {
                libdvdvign_log("error-> unable to read .gitignore file");
                return err;
            }
            ioutil.WriteFile(".libdvdvignore", lines);
            parseIgnoreFile(string(lines));
        } else {
            lines := []byte("")
        }
    }
}


