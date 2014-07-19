package libdvdvcodestate

import "../libdvdvutil"
import "testing"

func TestBuildStateDB(t *testing.T) {
    if err := BuildStateDB("../"); err != nil {
        t.Fatal(err);
    }
    if exist,_ := libdvdvutil.PathExist("../.libdvdv/state.db"); !exist {
        t.Fatal("error-> unable to build state database");
    }
}

