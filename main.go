package main
import "bytes"
import "os"
import "fmt"
import "io/ioutil"
import "path/filepath"
import "strings"
import "./libdevdevutil"
import "sort"

/* 
*   master-dev command line.
*       libdevdev ignore 
*       libdevdev allocwrk --server=libdevdev.org
*       libdevdev pushupdate 
*       libdevdev status
*       libdevedv try {uuid} 
*       libdevdev accept
*       libdevedv untry
*       libdevdev wallet     
*   worker-dev command line.
*       libdevdev get {uuid} --server=libdevdev.org
*       libdevdev update 
*       libdevdev submit
*       libdevdev status
*   libdevdev help
*/

func main() {
   if len(os.Args) < 2 {
        fmt.Println("unknown command-linei, execute [$libdevdev help] for usage");
        return;
    }
    if os.Args[1] == "ignore" {
        if ans,_ := libdevdevutil.PathExist(".libdevdevignore"); ans == true {
            fmt.Println("libdevdev ignore file already built");
            return;
        }
        fmt.Println("building libdevdev ignore file (.libdevdevignore)");
        buildIgnore();
    } else if os.Args[1] == "allocwrk" {
        allocWork();
    }
    return;
}

/*
* build a ignore file (.libdevdevignore) and adds gitignore inputs.
*/
func buildIgnore() {
    var ignore_paterns bytes.Buffer;
    ignore_paterns.WriteString("#All binary executables are ignored by default\n");
    if ans,_ := libdevdevutil.PathExist(".gitignore"); ans != false {
        fmt.Println("appending git ignore inputs to .libdevdevignore");
        gitignore,err := ioutil.ReadFile(".gitignore");
        if err == nil {
            ignore_paterns.Write(gitignore);
        }
    }
    ignore_paterns.WriteString(".*\n.o\n.a\n");
    ioutil.WriteFile(".libdevdevignore",ignore_paterns.Bytes(), 0644);
    fmt.Println("ignoring files -\n");
    fmt.Println(string(ignore_paterns.Bytes()));
}

func allocWork() {
    if ans,_ := libdevdevutil.PathExist(".libdevdevignore"); ans == false {
        fmt.Println("requirs creation of a ignore file first, execute [$ libdevdev",
                    " ignore] first");
        return;
    }
    libdevdevignore, err := ioutil.ReadFile(".libdevdevignore");
    if err != nil {
        fmt.Println("unable to read .libdevdevignore program exiting");
        return;
    }
    ignore := buildIgnoreList(libdevdevignore);
    pwd_info,err := ioutil.ReadDir(".");
    for i := 0; i < len(pwd_info); i++ {
        for j := 0; j < len(ignore); j++ {
            matches := make([]string, 10);
            ans,err := glob(pwd_info[i].Name(), ignore[i], matches);
            if err != nil {
                fmt.Println("error compiling patern ", ignore[j]);
            }
            if len(ans) > 0 {
                fmt.Printf("%s ignore is -->",ignore[i]);
                fmt.Println(ans);
            }
        }
    }
}

func buildIgnoreList(content []byte) []string {

    str := strings.Split(string(content), "\n");
    ignore_list := make([]string, len(str));
    j := 0;
    for i := 0; i < len(str); i++ {
        if (len(str[i]) == 0) || (str[i][0] == '#') {
            continue;
        }
        ignore_list[j] = str[i];
        j++
    }
    return ignore_list;
}

func glob(dir, pattern string, matches []string) (m []string, e error) {
    m = matches;
    fi, err := os.Stat(dir);
    if err != nil {
        return
    }
    if !fi.IsDir() {
        return
    }
    d, err := os.Open(dir);
    if err != nil {
        return;
    }
    defer d.Close();
    names, err := d.Readdirnames(-1);
    if err != nil {
        return;
    }
    sort.Strings(names);
    for _, n := range names {
        matched, err := filepath.Match(pattern, n);
        if err != nil {
            return m,err;
        }
        if matched {
            m = append(m, filepath.Join(dir, n));
        }
    }
    return
}
