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
*   -> lines starting with '!' are negated from the ignore list.
*   -> All other lines are treated as shell glob patterns.
*
* please refer to http://en.wikipedia.org/wiki/Glob_(programming) for 
* understanding shell glob patterns.
*/

import "os"

/* All modules have a log function*/
var libdvdvign_log func(a ...interface{});

func Init(log func(a ...interface{})) error {
    libdvdvign_log = log;
    libdvdvign("initializing libdvdv-ignore in current directory");

}
