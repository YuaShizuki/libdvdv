package libdvdvutil
import "os"
import "testing"
import "io/ioutil"
import "path/filepath"

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

var test_is_dir_passed = false;
func TestIsDir(t *testing.T) {
    f, err := ioutil.TempFile(".", "lddTest");
    if err != nil {
        t.Fatal("error-> ", err);
    }
    d, err := ioutil.TempDir(".", "lddTest");
    if err != nil {
        t.Fatal("error-> ", err);
    }
    if IsDir(f.Name()) || !IsDir(d) {
        t.Fatal("error-> IsDir functioning error");
    }
    f.Close();
    test_is_dir_passed = true;
}

func TestCreateFiles(t *testing.T) {
    if !test_path_exist_passed {
        t.SkipNow();
    }
    wd, err := ioutil.TempDir(".", "lddTest");
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
    fs[".libdvdvignore"] = []byte{1};
    fs["foo/f.t"] = []byte{1};
    fs["foo/f2.t"] = []byte{1}
    fs["foo/bar/fb.t"] = []byte{0};
    fs["foo/bar/fb2.t"] = []byte{0};
    fs["m.t"] = []byte{0};
    fs["m.r"] = []byte{1};
    fs["foo2/f20.t"] = []byte{1};
    fs["foo2/tar/ft.t"] = []byte{1};
    fs["foo2/zar/"] = []byte{0};
    err = CreateFiles(fs);
    if err != nil {
        t.Fatal("error->", err);
    }
    for k,_ := range fs {
        if !PathExist(k) {
            t.Fatal("error-> path=",k," does not exist");
        }
    }
}

func TestRemoveFiles(t *testing.T) {
    if !test_is_dir_passed || !test_path_exist_passed {
        t.SkipNow();
    }
    wd, err := os.Getwd();
    if err != nil {
        t.Fatal("error-> directory structure might be unclean");
    }
    for i := 0; i < 10; i++ {
        if (i % 2) == 0 {
            ioutil.TempDir(".", "lddTest");
        } else {
            ioutil.TempFile(".", "lddTest");
        }
    }
    m, err := filepath.Glob(wd+"/lddTest*");
    if len(m) == 0 {
        t.Fatal("error-> unable to create directory structure, directory ",
                "structure corrupt");
    }
    if err := RemoveFiles("lddTest*", wd); err != nil {
        t.Fatal("error-> ", err);
    }
    for i := range m {
        if PathExist(m[i]) {
            t.Fatal("error-> RemoveFiles failed");
        }
    }
}

func TestIsBinaryExecutable(t *testing.T) {
    if !test_path_exist_passed {
        t.SkipNow();
    }
    if PathExist("/bin/ls") {
        if exe,_ := IsBinaryExecutable("/bin/ls"); !exe {
            t.Fatal("error-> test reporting /bin/ls is not binary executable");
        }
    }
    if PathExist("/etc/hosts") {
        if exe,_ := IsBinaryExecutable("/etc/hosts"); exe {
            t.Fatal("error-> test reporting /etc/hosts is binary executable");
        }
    }
}

func TestCleanup(t *testing.T) {
    if err := ForTestCleanupTemp(t); err != nil {
        t.Fatal("error-> _______ UNABLE TO CLEAN ________ ");
    }
}



