package libdvdvutil
import "os"
import "testing"
import "io/ioutil"

var test_path_exist_passed = false;
func TestPathExist(t *testing.T) {
    if !PathExist("libdvdvutil_test.go") {
        t.FailNow();
    }
    if PathExist("dediejaijdjiejdaide.txt") {
        t.FailNow();
    }
    test_path_exist_passed = true;
}

func TestCreateFiles(t *testing.T) {
    if !test_path_exist_passed {
        t.SkipNow();
    }
    wd, err := ioutil.TempDir(".", "TestBuildFs");
    if err != nil {
        t.Log("test-> ",err);
        t.Fail();
        return;
    }
    if os.Chdir(wd) != nil {
        t.Log("test-> cannot change directory");
        t.Fail();
        return;
    }
    fs := make(map[string][]byte);
    fs["keshav.txt"] = []byte("working...");
    fs["darvin.txt"] = []byte("working...");
    fs["jack/issac.txt"] = []byte("free now ...");
    err = CreateFiles(fs);
    if err != nil {
        t.Fatal("error->", err);
    }
    for k,_ := range fs {
        if !PathExist(k) {
            t.Fatal("error-> path=",k," does not exist");
        }
    }
    os.Chdir("../")
    os.RemoveAll(wd);
}

