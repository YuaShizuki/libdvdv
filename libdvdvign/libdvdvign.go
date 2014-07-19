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
import "strings"
import "io/ioutil"
import "container/list"
import "path/filepath"
import "errors"

var LibdvdvLog func(a ...interface{}) = func(a ...interface{}){};

/*Ignore file list*/
var ignore *list.List;
var rm_frm_path int;

/* 
* Data structure to represent parsed ignore file, based on the rules mentioned 
* above 
*/
type Ignore_shell_globs struct {
    Sg_simple []string;
    Sg_dir []string;
    Sg_main []string;
    Sg_not [3][]string;
    ProjectDir string;
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

var standard_ignore_header string =
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
* Builds ".libdvdvignore" file in current directory.
*/
func BuildIgnoreFile(indir string) error {
    libdvdvignore_file_path := indir+"/.libdvdvignore";
    gitignore_file_path := indir+"/.gitignore";
    if exist,_ := libdvdvutil.PathExist(libdvdvignore_file_path); !exist {
        var lines []byte = nil;
        if exist2,_ := libdvdvutil.PathExist(gitignore_file_path); exist2 {
            var err error;
            lines, err = ioutil.ReadFile(gitignore_file_path);
            if err != nil {
                LibdvdvLog("error-> unable to read .gitignore file");
                return err;
            }
        } else {
            lines = []byte(standard_ignore_header);
        }
        if err := ioutil.WriteFile(libdvdvignore_file_path, lines, 0644); err != nil {
            LibdvdvLog(err);
            return err;
        }
    }
    return nil;
}

/*
* Determines shell glob pattern types. Mentioned in the rules above.
*/
func determine_pattern_type(s string) int {
    s_len := len(s);
    if (s_len == 0)  || (s[0] == '#') {
        return p_reject;
    } else if s[0] == '!' {
        return p_neg;
    } else if (s[0] != '/') && (s[s_len - 1] == '/') {
        return p_dir;
    } else if s[0] == '/' {
        return p_main;
    }
    return p_simple;
}

func ParseIgnoreFile(indir string) *Ignore_shell_globs {
    lines, err := ioutil.ReadFile(indir+"/.libdvdvignore");
    if err != nil {
        return nil;
    }
    line := strings.Split(string(lines), "\n");

    p := new(Ignore_shell_globs);
    p.Sg_simple = make([]string, 0, 25);
    p.Sg_dir = make([]string, 0, 25);
    p.Sg_main = make([]string, 0, 25);
    p.Sg_not = [3][]string {    make([]string,0, 15), make([]string,0,15),
                                make([]string,0,15)};
    p.ProjectDir = indir;
    for _,glob := range line {
        glob = strings.TrimSpace(glob);
        switch determine_pattern_type(glob) {
            case p_simple:
                p.Sg_simple = append(p.Sg_simple, glob);
            case p_dir:
                p.Sg_dir = append(p.Sg_dir, glob);
            case p_main:
                p.Sg_main = append(p.Sg_main, glob);
            case p_neg:
                switch determine_pattern_type(glob[1:len(glob)]) {
                    case p_simple:
                        p.Sg_not[0] = append(p.Sg_not[0],glob[1:len(glob)]);
                    case p_dir:
                        p.Sg_not[1] = append(p.Sg_not[1],glob[1:len(glob)]);
                    case p_main:
                        p.Sg_not[2] = append(p.Sg_not[2],glob[1:len(glob)]);
                    default:
                        continue;
                }
            default:
                continue;
        }
    }
    return p;
}

func BuildIgnoreList(globs *Ignore_shell_globs) error {
    if ignore != nil {
        ignore.Init();
    }
    if globs == nil {
        return errors.New("unknown error, shell globs empty");
    }
    rm_frm_path = len(globs.ProjectDir)+1;
    for _,s := range globs.Sg_main {
        match, err := filepath.Glob(globs.ProjectDir+s);
        if err != nil {
            return err;
        }
        for i, m := range match {
            if m[len(m)-1]  == '/' {
                match[i] = m[0:len(m)-1];
            }
        }
        appendToIgnore(match);
    }
    var err error;
    if err = buildIgnoreListDirWalk(globs.ProjectDir, globs); err == nil {
        err = negateFromIgnoreList(globs.ProjectDir, globs);
    } else {
        ignore.Init();
    }
    return err;
}

func buildIgnoreListDirWalk(path string, globs *Ignore_shell_globs) error {
    for _,s := range globs.Sg_simple {
        match, err := filepath.Glob(path+"/"+s);
        if err != nil {
            return err;
        }
        appendToIgnore(match);
    }
    for _,s := range globs.Sg_dir {
        match, err := filepath.Glob(path+"/"+s);
        if err != nil {
            return err;
        }
        for i, m := range match {
            match[i] = m[0:len(m)-1];
        }
        appendToIgnore(match);
    }
    finfo, err := ioutil.ReadDir(path);
    if err != nil  {
        return err;
    }
    for i := range finfo {
        path2 := path + "/" + finfo[i].Name();
        if finfo[i].IsDir() && (Check(path2) == nil) {
            if err := buildIgnoreListDirWalk(path2, globs); err != nil {
                return err;
            }
        }
    }
    return nil;
}

func negateFromIgnoreList(path string, globs *Ignore_shell_globs) error {
    for _, p := range globs.Sg_not[2] {
        match, err := filepath.Glob(path+p);
        if err != nil {
            return err;
        }
        for i, m := range match {
            if m[len(m)-1] == '/' {
                match[i] = m[0:len(m)-1];
            }
        }
        removeFromIgnore(match);
    }
    if globs.Sg_not[0] != nil || globs.Sg_not[1] != nil {
        err := negateFromIgnoreListDirWalk(path, globs);
        if err != nil {
            return err;
        }
    }
    return nil;
}

func negateFromIgnoreListDirWalk(path string, globs *Ignore_shell_globs) error {
    for _,s  := range globs.Sg_not[0] {
        match, err := filepath.Glob(path+"/"+s);
        if err != nil {
            return err;
        }
        removeFromIgnore(match);
    }
    for _,s := range globs.Sg_not[1] {
        match, err := filepath.Glob(path+"/"+s);
        if err != nil {
            return err;
        }
        for i, m := range match {
            match[i] = m[0:(len(m)-1)]
        }
        removeFromIgnore(match);
    }
    finfo, err := ioutil.ReadDir(path);
    if err != nil {
        return err;
    }
    for _,info := range finfo {
        path2 := path + "/" + info.Name();
        if info.IsDir() && (Check(path2) == nil) {
            if err := negateFromIgnoreListDirWalk(path2, globs); err != nil {
                return err;
            }
        }
    }
    return nil;
}

func appendToIgnore(str []string) {
    if ignore == nil {
        ignore = list.New();
    }
    for _,s := range str {
        ignore.PushBack(s[rm_frm_path:len(s)]);
    }
}

func removeFromIgnore(str []string) {
    for e := ignore.Front(); e != nil; {
        next := e.Next();
        for _,s := range str {
            if s[rm_frm_path:len(s)] == e.Value.(string) {
                ignore.Remove(e);
                break;
            }
        }
        e = next;
    }
}
/*
* Checks if a file should be ignored, if it returns a non nil pointer then YES
* else NO.
*/
func Check(s string) *list.Element {
    for e := ignore.Front(); e != nil ; e = e.Next() {
        if e.Value.(string) == s {
            return e;
        }
    }
    return nil;
}


