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

import "../libdvdvutil"
import "io/ioutil"
import "errors"

/* All modules have a log function*/
var libdvdvign_log func(a ...interface{});

/*Ignore file list*/
var ignore [][]string

/* 
* Data structure to represent parsed ignore file, based on the rules mentioned 
* above 
*/
type ignore_shell_globs {
    sg_simple []*string;
    sg_dir []*string;
    sg_main []*string;
    sg_not [3][]*string;
};

/*Ignore file message*/
var ignore_file_message string =
`# files specified hear are intentionally untracked files that libdevdev would 
# ignore. 
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
        return buildIgnoreList(parseIgnoreLines(string(lines)));
    } else {
        var lines []byte = nil;
        if libdvdvutil.PathExist(".gitignore") {
            var err error;
            lines, err = ioutil.ReadFile(".gitignore");
            if err != nil {
                libdvdvign_log("error-> unable to read .gitignore file");
                return err;
            }
        } else {
            lines = []byte(ignore_file_message);
        }
        err := ioutil.WriteFile(".libdvdvignore", lines, 0644);
        if err != nil {
            libdvdvign_log(err);
            return err;
        }
        return buildIgnoreList(parseIgnoreLines(string(lines)));
    }
    return errors.New("error-> Unknown error");
}

/*
* The below constants help parses every single line according to the rules mentioned
* above.
*/
const (
    p_simple = iota;
    p_dir = iota;
    p_main = iota;
    p_neg = iota;
    p_reject = iota;
);

func determine_pattern_type (s *string) int {
    s_len := len(*s);
    if (*s)[0] == '#' {
        return p_reject;
    } else if (*s)[0] == '!' {
        return p_neg;
    } else if (*s)[0] != '/' && (*s)[s_len - 1] == '/' {
        return p_dir;
    } else if (*s)[0] == '/' {
        return p_main;
    }
    return p_simple;
}

func parseIgnoreLines(lines string) *ignore_shell_globs {
    line := strings.Split(lines, "\n");
    line_len := len(line);

    p := new(ignore_shell_globs);
    p.sg_simple = make([]*string, 0, 10);
    p.sg_dir = make([]*string, 0, 10);
    p.sg_main = make([]*string, 0, 10);
    p.sg_not = [3][]*string {   make([]*string,0, 5), make([]*string,0,5), 
                                make([]*string,0,5)};

    for i := 0; i < line_len; i++ {
        line[i] = strings.TrimSpace(line[i]);
        switch determine_pattern_type(&(line[i])) {
            case p_simple:
                p.pattern_simple = append(p.sg_simple,&(line[i]));
            case p_dir:
                p.pattern_dir = append(p.sg_dir, &(line[i]));
            case p_main:
                p.pattern_main = append(p.sg_main, &(line[i]));
            case p_neg:
                switch determine_pattern_type(&(line[1:len(line[i])])) {
                    case p_simple:
                        p.sg_not[0] = append(p.sg_not[0],&(line[1:len(line[i])]));
                    case p_dir:
                        p.sg_not[1] = append(p.sg_not[1],&(line[1:len(line[i])]));
                    case p_main:
                        p.sg_not[2] = append(p.sg_not[2], &(line[1:len(line[i])]));
                    case default:
                        continue;
                }
            case default:
                continue;
        }
    }
}

func buildIgnoreList() error {
    
}

func check(path string) error {
}
