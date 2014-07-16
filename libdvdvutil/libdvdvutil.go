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
import "bytes"

var libdvdvutil_log func(a ...interface{}) = func(a ...interface{}) {};

func Setup(log func(a ...interface{})) {
    libdvdvutil_log = log;
}

/*
* Checks if a file path exists
*/
func PathExist(path string) (bool, os.FileInfo) {
    info, err :=  os.Stat(path)
    if err != nil {
        if os.IsNotExist(err) {
            return false, info;
        }
        panic("panic-> cannot determine if path exists");
    }
    return true, info ;
}

/*
* Returns whether a path is a directory or a string.
* returns false if path does not exist
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
        if exist,_ := PathExist(path.Dir(k)); !exist {
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
        if exist,_ := PathExist(f[i]); exist {
            return errors.New("failed to remove file -> "+f[i]);
        }
    }
    return nil;
}

func IsBinaryExecutable(path string, finfo os.FileInfo) (bool, os.FileInfo) {
    var err error;
    var info os.FileInfo;
    if finfo != nil {
        info = finfo;
    } else {
        info, err = os.Stat(path);
        if err != nil {
            return false, nil;
        }
    }
    if info.IsDir() {
        return false, info;
    }
    perm := info.Mode().Perm();
    if ((perm % 2) != 0)  || (((perm >> 3)% 2) != 0) || (((perm >> 6)%2) != 0) {
        content, err := ioutil.ReadFile(path);
        if err != nil {
            return false, info;
        } else if bytes.Count(content, []byte{0}) > 1 {
            return true, info;
        }
    }
    return false, info;
}

/*
* This function is used mostely during testing, to delete temp files or folders
* created.
*/
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
