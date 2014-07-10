package libdvdvign
import "os"
import "testing"
import "io/ioutil"
import "path/filepath"
import "strings"

var temp_log func(a ...interface{});
var libdvdvignore []string;
var libdvdvignore_result [][]string;

func check_should_ignore(path string) bool {
    for i,_ := range libdvdvignore_result {
        for _,str := range libdvdvignore_result[i] {
            if str == path {
                return true;
            }
        }
    }
    return false;
}


func recurse_glob(path string) error {
    /*get all globs from current directory*/
    for _,pattern := range libdvdvignore {
        p := strings.TrimSpace(pattern);
        if len(p) != 0 {
            var glob string;
            if p[0] == '/' {
                glob = path+p;
            } else {
                glob = path+"/"+p;
            }
            temp_log("tttt->",glob)
            temp, err := filepath.Glob(glob);
            if err != nil {
                return err;
            }
            libdvdvignore_result = append(libdvdvignore_result, temp);
        }
    }
    /* Recurese through all sub directories, but ignore the ones present */
    /* in libdvdvignore_result already */
    file_info, err := ioutil.ReadDir(path);
    if err != nil {
        return err;
    }
    for _,finfo := range file_info {
        if finfo.IsDir() && !check_should_ignore(path+"/"+finfo.Name()) {
            err := recurse_glob(path+"/"+finfo.Name());
            if err != nil {
                return err;
            }
        }
    }
    return nil;
}

func TestLibdvdvign(t *testing.T) {
    temp_log = t.Log;
    os.Chdir(os.Getenv("ign_path"));
    wd,_ := os.Getwd();
    f, err := ioutil.ReadFile(wd+"/.libdvdvignore");
    if err != nil {
        t.Log(err);
        t.Fail();
    }
    libdvdvignore_result = make([][]string, 20);
    libdvdvignore = strings.Split(string(f), "\n");
    if recurse_glob(wd) != nil {
        t.Log("error building glob pattern recursively");
        t.Fail();
    }
    for i,_ := range libdvdvignore_result {
        for _,str := range libdvdvignore_result[i] {
            t.Log("file-> ",str);
        }
    }
}



