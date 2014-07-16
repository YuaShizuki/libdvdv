package libdvdvign

import "strings"
import "os"
import "testing"
import "io/ioutil"
import "../libdvdvutil"

func TestBuildIgnoreFile(t *testing.T) {
    //content of .gitignore
    t1 := `.swp
    .ram
    .cat`;

    d , err := ioutil.TempDir(".", "lddTest");
    if err != nil {
        t.Fatal("err-> ", err);
    }
    os.Chdir(d);
    LibdvdvLog = t.Log;
    //part-1
    if err := BuildIgnoreFile(); err != nil {
        t.Fatal("error-> ", err);
    }
    if exist,_ := libdvdvutil.PathExist(".libdvdvignore"); !exist {
        t.Fatal("error-> unable to build file");
    }
    os.Remove(".libdvdvignore");
    if exist,_ := libdvdvutil.PathExist(".libdvdvignore"); exist {
        t.Fatal("error-> unable to delete file");
    }
    //part-2
    err = ioutil.WriteFile(".gitignore",[]byte(t1), 0644)
    if err != nil {
        t.Fatal("error-> ", err);
    }
    if err := BuildIgnoreFile(); err != nil {
        t.Fatal("error-> ", err);
        return;
    }
    new_file , err := ioutil.ReadFile(".libdvdvignore");
    if err != nil {
        t.Fatal("error-> ", err);
    }
    if string(new_file) != t1 {
        t.Fatal("error-> .gitignore != .libdvdvignore");
    }
    os.Chdir("../");
}

func TestParseIgnoreFile(t *testing.T) {
    var globs *Ignore_shell_globs;
    d , err := ioutil.TempDir(".", "lddTest");
    if err := os.Chdir(d); err != nil {
        t.Fatal("error-> ", err);
    }
    t1 := []string {
        "*.txt",
        "/foo/bar/",
        "foos/dark/",
        "!*.txt",
        "!/foo/bar/",
        "!foos/dark/"};
    err = ioutil.WriteFile(".libdvdvignore", []byte(strings.Join(t1, "\n")), 0644);
    if err != nil {
        t.Fatal("error-> ", err);
    }

    if globs = ParseIgnoreFile(); globs == nil {
        t.Fatal("error-> unable to parese ignore file");
    }
    if (len(globs.Sg_simple) != 1) || (globs.Sg_simple[0] != t1[0]) {
        t.Fatal("error-> error parsing ", t1[0]);
    }
    if (len(globs.Sg_main) != 1) || (globs.Sg_main[0] != t1[1]) {
        t.Fatal("error-> error parsing ", t1[1]);
    }
    if (len(globs.Sg_dir) != 1) || (globs.Sg_dir[0] != t1[2]) {
        t.Fatal("error-> error parsing ", t1[2])
    }
    for _, v := range globs.Sg_not {
        if len(v) != 1 {
            t.Fatal("error-> error parsing not pattern");
        }
    }
    if globs.Sg_not[0][0] != t1[3][1:len(t1[3])] {
        t.Fatal("error-> error parsing ", t1[3]);
    }
    if globs.Sg_not[1][0] != t1[5][1:len(t1[5])] {
        t.Fatal("error-> error parsing ", t1[5]);
    }
    if globs.Sg_not[2][0] != t1[4][1:len(t1[4])] {
        t.Fatal("error-> error parsing ", t1[4]);
    }
    os.Chdir("../");
}

func TestIgnoreList(t *testing.T) {
    d, err := ioutil.TempDir(".", "lddTest");
    if err != nil {
        t.Fatal("error-> while creating temp directory", err);
    }
    if err = os.Chdir(d); err != nil {
        t.Fatal("error-> unable to change directory: ", err);
    }
    tmp := []string{
        "/foo/*",
        "!foo/bar/",
        "*.r",
        "car/",
        "/foo2/*",
        "!zar/"};

    fs := make(map[string][]byte);
    fs[".libdvdvignore"] = []byte(strings.Join(tmp, "\n"));
    fs["foo/f.t"] = []byte{1};
    fs["foo/f2.t"] = []byte{1}
    fs["foo/bar/fb.t"] = []byte{0, 1};
    fs["foo/bar/fb2.t"] = []byte{0, 1};
    fs["m.t"] = []byte{0, 1};
    fs["m.r"] = []byte{1};
    fs["foo2/f20.t"] = []byte{1};
    fs["foo2/car/ft.t"] = []byte{0,1};
    fs["foo2/zar/"] = []byte{0, 1};
    if err = libdvdvutil.CreateFiles(fs); err != nil {
        t.Fatal("error-> libdvdvutil.CreateFiles returning :", err);
    }
    if err = BuildIgnoreList(ParseIgnoreFile()); err != nil {
        t.Fatal("error-> unable to build ignore list: ", err);
    }
    if d , err = os.Getwd(); err != nil {
        t.Fatal("error-> unable to get pwd: ", err);
    }
    for k, v := range fs {
        if k == ".libdvdvignore" {
            continue;
        }
        ck := Check(d+"/"+k);
        if ((ck != nil) && (v[0] == 0)) || ((ck == nil) && (v[0] == 1)){
            t.Log("error-> failed for file=",k);
            t.Fail();
        }
    }
}

/*
* This is not a test, but only a clean up function
*/
func TestCleanup(t *testing.T){
    if err := libdvdvutil.ForTestCleanupTemp(t); err != nil {
        t.Fatal("error-> _________ UNABLE TO CLEAN _________");
    }
}




