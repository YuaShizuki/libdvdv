package libdvdvign

import "strings"
import "os"
import "testing"
import "io/ioutil"
import "../libdvdvutil"
import "path"

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
    Setup(t.Log);
    //part-1
    if err := BuildIgnoreFile(); err != nil {
        t.Fatal("error-> ", err);
    }
    if !libdvdvutil.PathExist(".libdvdvignore") {
        t.Fatal("error-> unable to build file");
    }
    os.Remove(".libdvdvignore");
    if libdvdvutil.PathExist(".libdvdvignore") {
        t.Fatal("error-> unable to delete file");
    }
    //part-2
    err = ioutil.WriteFile(".gitignore",[]byte(t1), 0644)
    if err != nil {
        t.Fatal("error-> ", err);
    }
    Setup(t.Log);
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
    os.RemoveAll(d+"/");
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
    Setup(t.Log);
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
}

func TestCleanup(t *testing.T){
    libdvdvutil.Setup(t.Log);
    wd, err := os.Getwd();
    if err != nil {
        t.Fatal("error-> _________ UNABLE TO CLEAN _________", err);
    }
    for {
        var matched bool;
        if matched,_ = path.Match("lddTest*", path.Base(wd)); matched {
            wd = path.Dir(wd);
        } else {
            break;
        }
    }
    err = os.Chdir(wd);
    if err != nil {
        t.Fatal("error-> _________ UNABLE TO CLEAN _________");
    }
    err = libdvdvutil.RemoveFiles("lddTest*", wd);
    if err != nil {
        t.Fatal("error-> _________ UNABLE TO CLEAN _________")
    }
}




