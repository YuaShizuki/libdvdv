package libdvdvutil
/*
* This package contains utility function used by all modules 
*/

import "os"
import "io/ioutil"
import "path"
import "path/filepath"
import "errors"
import "testing"

var libdvdvutil_log func(a ...interface{}) = func(a ...interface{}) {};

func Setup(log func(a ...interface{})) {
    libdvdvutil_log = log;
}

/*
* Checks if a file path exists
*/
func PathExist(path string) bool {
    _,err :=  os.Stat(path)
    if err != nil {
        if os.IsNotExist(err) {
            return false;
        }
        panic("panic-> cannot determine if path exists");
    }
    return true;
}

/*
* Returns whether a path is a directory or a string.
*/
func IsDir(path string) bool {
    info, err := os.Stat(path);
    if err != nil {
        if os.IsNotExist(err) {
            return false;
        }
        panic("panic-> cannot use libraray function os.Stat ");
    }
    return info.IsDir();
}

/*
* Creates files and directory based on the map with path and and byte content. 
* directory should be writen as x/y/z/directory/ ie with a trailing '/'
*/
func CreateFiles(fs map[string][]byte) error {
    for k, v := range fs {
        if !PathExist(path.Dir(k)) {
            if err := os.MkdirAll(path.Dir(k), 0744); err != nil {
                return err;
            }
        }
        if k[len(k)-1] == '/' {
            continue;
        }
        err := ioutil.WriteFile(k, v, 0644);
        if err != nil {
            return err;
        }
    }
    return nil;
}

/*
* this func is used to remove files and directories.
* matching a glob pattern.
*/
func RemoveFiles(pattern string, dir string) error {
    f, err := filepath.Glob(dir+"/"+pattern);
    if err != nil {
      return err;
    }
    for i := range f {
        if IsDir(f[i]) {
            err = os.RemoveAll(f[i]);
        } else {
            err = os.Remove(f[i]);
        }
        if err != nil {
            return err;
        }
        if PathExist(f[i]) {
            return errors.New("failed to remove file -> "+f[i]);
        }
    }
    return nil;
}

func ForTestCleanupTemp(t *testing.T) error {
    wd, err := os.Getwd();
    if err != nil {
        return err;
    }
    for {
        if m,_ := path.Match("lddTest*", path.Base(wd)); m {
            wd = path.Dir(wd);
        } else {
            break;
        }
    }
    if err = os.Chdir(wd); err != nil {
        return err;
    }
    if err = RemoveFiles("lddTest*", wd); err != nil {
        return err;
    }
    return nil;
}
