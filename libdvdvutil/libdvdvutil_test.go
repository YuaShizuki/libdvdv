package libdvdvutil
import "testing"

func TestLibdvdvutil(t *testing.T) {
    t.Log("testing-> libdvdvutil");
    if !PathExist("libdvdvutil_test.go") {
        t.Fail();
    }
    if PathExist("dediejaijdjiejdaide.txt") {
        t.Fail();
    }
}
