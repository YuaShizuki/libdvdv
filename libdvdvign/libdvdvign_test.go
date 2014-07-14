package libdvdvign
import "os"
import "testing"
import "io/ioutil"
import "../libdvdvutil"

var t1 string = `.tmp
.temp
.swp
.exe`


func TestLibdvdvign_BuildIgnoreFile(t *testing.T) {
    d , err := ioutil.TempDir(".", "libdvdv-ignore-test");
    if err != nil {
        t.Log("test-> unable to create tmp directory");
        t.Log("err-> ", err);
        t.Fail();
        return;
    }
    os.Chdir(d);
    t.Log("test-> changed to directory ", d);
    Setup(t.Log);
    //part-1
    err = BuildIgnoreFile();
    if err != nil {
        t.Log("error-> ", err);
        t.Fail();
        return;
    }
    if !libdvdvutil.PathExist(".libdvdvignore") {
        t.Log("error-> unable to build file");
        t.Fail();
        return;
    }
    os.Remove(".libdvdvignore");
    if libdvdvutil.PathExist(".libdvdvignore") {
        t.Log("error-> unable to delete file");
        t.Fail();
        return;
    }
    //part-2
    err = ioutil.WriteFile(".gitignore",[]byte(t1), 0644)
    if err != nil {
        t.Log("error-> ", err);
        t.Fail();
    }
    Setup(t.Log);
    err = BuildIgnoreFile();
    if err != nil {
        t.Log("error-> ", err);
        t.Fail();
        return;
    }
    new_file , err := ioutil.ReadFile(".libdvdvignore");
    if err != nil {
        t.Log("error-> ", err);
        t.Fail();
        return;
    }
    if string(new_file) != t1 {
        t.Log("error-> new content is diffrent from .gitignore");
        t.Fail();
    }
    os.Chdir("../");
    os.RemoveAll(d+"/");
}
