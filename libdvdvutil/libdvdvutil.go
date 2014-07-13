package libdvdvutil
/*
* This package contains utility function used by all modules 
*/
import "os"
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

