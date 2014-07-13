package libdvdvign

/*
* This module builds a ignore file list during initialization. Ounce initialized, one
* can call check to check if a file is suppose to be ignored or not. This package
* reads the .libdvdvignore file present in the project's main directory. If it cant
* find a .libdvdvignore file, it creates a new one and adds shell globs patterns 
* from a .gitignore file if present.
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
import "container/lists"
import "path/filesystem"
import "io/ioutil"

/* All modules have a log function*/
var libdvdvign_log func(a ...interface{}) = func(a ...interface{}) { };
/*Ignore file list*/
var ignore *container.List;

/* 
* Data structure to represent parsed ignore file, based on the rules mentioned 
* above 
*/
type Ignore_shell_globs {
    Sg_simple []*string;
    Sg_dir []*string;
    Sg_main []*string;
    Sg_not [3][]*string;
};

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

/* 
* All modules should have a setup function  
*/
func Setup(log func(a ...interface{})) error {
    libddvdvign_log = log;
    if ignore != nil {
        ignore.Init();
    } else {
        ignore := list.New();
    }
    return nil;
}

func BuildIgnoreFile() error {
    if !libdvdvutil.PathExist(".libdvdvignore") {
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
    }
    return nil;
}

func determine_pattern_type(s *string) int {
    s_len := len(*s);
    if (*s)[0] == '#' {
        return p_reject;
    } else if (*s)[0] == '!' {
        return p_neg;
    } else if ((*s)[0] != '/') && ((*s)[s_len - 1] == '/') {
        return p_dir;
    } else if (*s)[0] == '/' {
        return p_main;
    }
    return p_simple;
}

func ParseIgnoreFile() *Ignore_shell_globs {
    lines, err := ioutil.ReadFile(".libdvdvignore");
    if err != nil {
        return nil;
    }
    line := strings.Split(lines, "\n");
    line_len := len(line);

    p := new(ignore_shell_globs);
    p.Sg_simple = make([]*string, 0, 25);
    p.Sg_dir = make([]*string, 0, 25);
    p.Sg_main = make([]*string, 0, 25);
    p.Sg_not = [3][]*string {   make([]*string,0, 15), make([]*string,0,15),
                                make([]*string,0,15)};

    for i := 0; i < line_len; i++ {
        line[i] = strings.TrimSpace(line[i]);
        switch determine_pattern_type(&(line[i])) {
            case p_simple:
                p.Sg_simple = append(p.Sg_simple,&(line[i]));
            case p_dir:
                p.Sg_dir = append(p.Sg_dir, &(line[i]));
            case p_main:
                p.Sg_main = append(p.Sg_main, &(line[i]));
            case p_neg:
                switch determine_pattern_type(&(line[1:len(line[i])])) {
                    case p_simple:
                        p.Sg_not[0] = append(p.Sg_not[0],&(line[1:len(line[i])]));
                    case p_dir:
                        p.Sg_not[1] = append(p.Sg_not[1],&(line[1:len(line[i])]));
                    case p_main:
                        p.Sg_not[2] = append(p.Sg_not[2], &(line[1:len(line[i])]));
                    case default:
                        continue;
                }
            case default:
                continue;
        }
    }
}

func BuildIgnoreList(globs *Ignore_shell_globs) error {
    if globs == nil {
        return errors.New("unknown error, shell globs empty");
    }
    wd, err := os.Getwd();
    if err != nil {
        return err;
    }
    for _,p := range globs.Sg_main {
        match, err := filepath.Glob(wd+(*p));
        if err != nil {
            return err;
        }
        if len(match) != 0 {
            ignore.PushBack(match);
        }
    }
    err := buildIgnoreListDirWalk(wd, globs);
    if err == nil {
        negateFromIgnoreList(wd, globs);
    } else {
        ignore.Init();
    }
    return err;
}

func buildIgnoreListDirWalk(path string, globs *ignore_shell_globs) error {
    for _,p := range globs.Sg_simple {
        match, err := filepath.Glob(wd+"/"+(*p));
        if err != nil {
            return err;
        }
        if len(match) == 0 {
            continue;
        }
        ignore.PushBack(match);
    }
    for _,p := range golbs.Sg_dir {
        match, err := filepath.Glob(wd+"/"+(*p));
        if err != nil {
            return err;
        }
        if len(match) == 0 {
            continue;
        }
        for i := range match {
            match[i] = match[i][0:(len(match[i])-1)];
        }
        ignore.PushBack(match);
    }
    finfo, err := ioutil.ReadDir(path);
    if err != nil  {
        return err;
    }
    for i := range finfo {
        if finfo[i].IsDir() && (Check(wd+"/"+finfo[i].Name()) != nil) {
            err := buildIgnoreListDirWalk(wd+"/"+finfo[i].Name());
            if err != nil {
                return err;
            }
        }
    }
}

func Check(path string) *Element {
    for e := ignore.Front(); e != nil; e = e.Next() {
        ls,_ := e.Value.([]string);
        for i := range ls {
            if path == ls[i] {
                return e;
            }
        }
    }
    return nil;
}

func negateFromIgnoreList(path string, globs *Ignore_shell_globs) {
    for _, p := globs.Sg_not[2] {
         
    }
}




