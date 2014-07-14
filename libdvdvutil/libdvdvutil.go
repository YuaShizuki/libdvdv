package libdvdvutil
/*
* This package contains utility function used by all modules 
*/
import "os"
import "io/ioutil"
import "path"
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
* Creates files and directory based on the map with path and and byte content. 
* directory should be writen as x/y/z/directory/ ie with a trailing '/'
*/
func CreateFiles(fs map[string][]byte) error {
    for k, v := range fs {
        if !PathExist(path.Dir(k)) {
            os.Mkdir(path.Dir(k), 0744);
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
