package libdvdvcodestate

import "../libdvdvutil"
import "../libdvdvign"
import "testing"

func TestBuildStateDB(t *testing.T) {
    libdvdvign.LibdvdvLog = t.Log;
    libdvdvign.BuildIgnoreFile("../");
    if err := BuildStateDB("../"); err != nil {
        t.Fatal(err);
    }
    if exist,_ := libdvdvutil.PathExist("../.libdvdv/state.db"); !exist {
        t.Fatal("error-> unable to build state database");
    }
}

