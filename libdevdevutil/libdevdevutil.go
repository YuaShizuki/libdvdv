package libdevdevutil
import "os"


/*
* Checks if a file path exists
*/
func PathExist(path string) (bool, error) {
    _,err :=  os.Stat(path)
    if err != nil {
        if os.IsNotExist(err) {
            return false, nil;
        }
        return false, err;
    }
    return true, nil;
}


