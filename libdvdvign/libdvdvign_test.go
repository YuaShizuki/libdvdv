package libdvdvign
import "testing"
import "os"
import "../libdvdvutil"
import "bytes"
import "io/ioutil"

/*
* Pass the directory where libdvdvign module should be initialized as a 
* as a enviournment variable ign_dir
*/
func TestLibdvdvign(t *testing.T){
    dir := os.Getenv("ign_dir");
    if len(dir) == 0 {
        t.Log("error-> empty value for ign_dir enviornment variable");
        t.Fail();
        return;
    }
    if os.Chdir(dir) != nil {
        t.Log("error-> unable to change directory");
        t.Fail();
        return;
    }
    if Init(t.Log) != nil {
        t.Log("error-> init libdvdvign failed");
        t.Fail();
        return;
    }
    t.Log("test-> checking if .libdvdvignore file was built")
    if !libdvdvutil.PathExist(".libdvdvignore") {
        t.Fail();
        return;
    }
    if libdvdvutil.PathExist(".gitignore") {
        t.Log("test-> detected .gitignore file, comparing .libdvdvignore with",
                " .gitignore");
        gitignore, err1 := ioutil.ReadFile(".gitignore");
        libdvdvignore, err2 := ioutil.ReadFile(".libdvdvignore");
        if (err1 != nil)  || (err2 != nil) {
            t.Log("error-> ",err1,err2);
            t.Fail();
            return;
        }
        if !bytes.Equal(gitignore, libdvdvignore) {
            t.Log("test-> faild, .gitignore's got diffrent content than ",
                    ".libdvdvignore");
            t.Fail();
            return;
        }
        t.Log("test-> .gitignore content matches .libdvdvignore content");
    }
}
