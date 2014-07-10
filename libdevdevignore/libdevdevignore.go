package libdevdevignore
import "io/ioutil"
import "strings"
import "os"
import "path/filepath"
import "container/list"

/*
* this allows the flexibility to use fmt.Println in normal mode and t.Log in
* test enviornment.
*/
var consoleLog func (a ...interface{});

/*
* ign_glb_typ_sim (ignore glob type simple) is glob string that has no '/' its 
* suppose to executed by glob function in every directory.
* ign_glb_typ_dir (ignore glob type directory) is glob string that only matches
* with a directory, it also is suppose to be executed in every single directory
* ign_glb_typ_mn (ignore glob type main) is glob string that's suppose to be 
* executed form the main directory
* ign_glb_typ_nt (ignore glob type not) is globe string that start with '!'
* it reinclude the files if ignored
* ign_glb_typ_c (ignore globe type comment); these are suppose to be ignored;
* ign_glb_type_uk (ignore glob type unknown); unknown.
* http://git-scm.com/docs/gitignore has information on .gitignore files.
* .libdevdevignore follows the same rules.
*/
const (
    Ign_glb_typ_smpl = iota;
    Ign_glb_typ_dir = iota;
    Ign_glb_typ_mn = iota;
    Ign_glb_typ_nt = iota;
    Ign_glb_typ_c = iota;
    Ign_glb_typ_uk = iota;
)

/*
* a not glob pattern can be any of ign_glb_typ_smpl, ign_glb_typ_dir, 
* ign_glb_typ_mn. Hence this struct.
*/
type not_ignore struct {
    ign_glb_smpl []*string;
    ign_glb_dir []*string;
    ign_glb_mn []*string;
}

/*
* this structure is a direct map for the ignore file. This struct needs to be
* initialized if you plane to use BuildIgnoreFileList witout Init. Since Init
* loads a ignore file from the working directory and builds this struct 
* automaticaly.
*/
type Ignore struct {
    ign_glb_smpl *[]*string;
    ign_glb_dir *[]*string;
    ign_glb_mn *[]*string;
    ign_glb_nt *not_ignore;
}

/*
* initializes the module, with the specified ignore file. responsible for 
* building a array of all files to be excluded form the project. This module
* will load ignore_file from current working directory os.Getwd();
*/
func Init(log func(a ...interface{}), ignore_file string) error {
    consoleLog = log;
    consoleLog("initializing libdevdevignore module");
    ignore,err :=  ioutil.ReadFile(ignore_file);
    if err != nil {
        consoleLog("unable to open ",ignore_file," file");
        return err;
    }
    BuildIgnoreFileList(parseIgnoreFile(string(ignore)));
    return nil;
}

func ShouldIgnore(file_name, full_path string) bool {
    return false;
}

func BuildIgnoreFileList(ignore *Ignore) error {
    consoleLog("++++Ignore File List++++");
    pwd,_ := os.Getwd();
    ignore_list := list.New();
    //Handel ignore glob patterns for main directory.
    for _,ign := range *(ignore.ign_glb_mn) {
        l,_ := filepath.glob(pwd+ign);
       for i,_ := range l {
           if (len(l[i]) != 0) && (l[i][len(l[i]) - 1] == '/') {
               //ignore the patter ie just "/"
               if len(l[i]) == 1 {
                   continue;
               }
               ignore_list.PushBack(l[i][:len(l[i])-1]);
           }
           else {
               ignore_list.PushBack(l[i]);
           }
       }
    }
    //Now remove from ignore_list all '!' not ignore files.
    for e := ignore_list.Front(); e != nil ; e := e.Next() {
        file,_ := e.Value.(string);
        for _,ign := ignore.ign_glb_nt.ign_
    }
    //Now recursivly walk through the main directory and append
    return nil;
}

func parseIgnoreFile(content string) *Ignore {
    patterns := strings.Split(content, "\n");

    var ignore_globs Ignore;
    ign_glb_smpl := make([]*string, 0, len(patterns));
    ign_glb_dir := make([]*string, 0, len(patterns));
    ign_glb_mn := make([]*string, 0, len(patterns));
    var ign_glb_nt not_ignore;
    ign_glb_nt.ign_glb_smpl = make([]*string, 0, len(patterns));
    ign_glb_nt.ign_glb_dir = make([]*string, 0, len(patterns));
    ign_glb_nt.ign_glb_mn = make([]*string, 0, len(patterns));

    for _,pat := range patterns {
        pattern := strings.TrimSpace(pat);
        switch IgnoreGlobType(pattern) {
            case Ign_glb_typ_smpl:
                ign_glb_smpl = append(ign_glb_smpl, &pattern);
            case Ign_glb_typ_dir:
                ign_glb_dir = append(ign_glb_dir, &pattern);
            case Ign_glb_typ_mn:
                ign_glb_mn = append(ign_glb_mn, &pattern);
            case Ign_glb_typ_nt:
                pattern_nt := pattern[1:];
                switch IgnoreGlobType(pattern_nt) {
                    case Ign_glb_typ_smpl:
                        ign_glb_nt.ign_glb_smpl = append(ign_glb_nt.ign_glb_smpl, &pattern_nt);
                    case Ign_glb_typ_dir:
                        ign_glb_nt.ign_glb_dir = append(ign_glb_nt.ign_glb_dir, &pattern_nt);
                    case Ign_glb_typ_mn:
                        ign_glb_nt.ign_glb_mn = append(ign_glb_nt.ign_glb_mn, &pattern_nt);
                    default:
                        continue;
                }
            default:
                continue;
        }
    }
    //log all patterns to check every thing is working fine..
    log_str_ptr_array := func (s []*string) {
        if len(s) == 0 {
            consoleLog("[empty]");
            return;
        }
        for _,p := range s {
            consoleLog(*p);
        }
    }
    consoleLog("ign_glb_smpl =======>");
    log_str_ptr_array(ign_glb_smpl);
    consoleLog("ing_glb_dir =======>");
    log_str_ptr_array(ign_glb_dir);
    consoleLog("ign_glb_mn =======>");
    log_str_ptr_array(ign_glb_mn);
    consoleLog("ign_glb_nt.ign_glb_smpl =======>");
    log_str_ptr_array(ign_glb_nt.ign_glb_smpl);
    consoleLog("ign_glb_nt.ign_glb_dir =======>");
    log_str_ptr_array(ign_glb_nt.ign_glb_dir);
    consoleLog("ign_glb_nt.ign_glb_mn =======>");
    log_str_ptr_array(ign_glb_nt.ign_glb_mn);
    //Build the return Ignore structure.
    ignore_globs.ign_glb_smpl = &ign_glb_smpl;
    ignore_globs.ign_glb_dir = &ign_glb_dir;
    ignore_globs.ign_glb_mn = &ign_glb_mn;
    ignore_globs.ign_glb_nt = &ign_glb_nt;
    return &ignore_globs;
}

func IgnoreGlobType(str string) int {
    str_len := len(str);
    slash_count := strings.Count(str ,"/");
    if len(str) == 0 {
        return Ign_glb_typ_uk;
    } else if str[0] == '#' {
        return Ign_glb_typ_c;
    } else if str[0] == '!' {
        return Ign_glb_typ_nt;
    } else if str[0] == '\\' {
        return IgnoreGlobType(str[1:]);
    } else if slash_count  == 0 {
        return Ign_glb_typ_smpl;
    } else if (slash_count == 1) && (str[str_len - 1] == '/') {
        return Ign_glb_typ_dir;
    } else if slash_count > 0 {
        return Ign_glb_typ_mn;
    }
    return Ign_glb_typ_uk;
}


func logAllIgnoreGlobs() {
}
