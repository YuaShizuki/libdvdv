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
    for i := range globs.Sg_not {
        if (len(globs.Sg_not[i]) != 1) || (globs.Sg_not[i][0] != t1[3+i][1:len(t1[3+i])]) {
            t.Fatal("error-> error parsing", globs.Sg_not[i][0]);
        }
    }
}

func TestCleanup(t *testing.T){
    t.Log("test-> final cleanup");
    wd := os.Getwd();
    for {
        if path.Match("lddTest*", path.Base(path.Dir(wd))) {
            wd = path.Dir(wd);
        }
    }

}




