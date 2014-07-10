package libdevdevignore
import "testing"
import "os"

func TestLibdevdevIgnore(t *testing.T) {
    ignore_dir := os.Getenv("ignore");
    os.Chdir(ignore_dir);
    err := Init(t.Log, ".gitignore");
    if err != nil {
        t.Log(err);
        t.Fail();
        return;
    }
}
