package main

import "fmt"
import "os"
import "io/ioutil"
import "code.google.com/p/go-sqlite/go1/sqlite3"

var ignore []string;

func parseGitIgnore(str string) []string {
    return nil;
}

func pathExist(path string) bool {
    _, err := os.Stat(path);
    if err != nil {
        if os.IsNotExist(err) {
            return false;
        }
    }
    return true;
}

func isInIgnoreList(name string) bool {
    return false;
}

func isExecutable(mode os.FileMode) bool {
    if mode.IsDir() {
        return false;
    }
    perm := mode.Perm();
    if ((perm % 2) != 0)  || (((perm >> 3)% 2) != 0) || (((perm >> 6)%2) != 0) {
        return true;
    }
    return false;
}

func terminalLogStatus(c *sqlite3.Conn) {
    result, err := c.Query("SELECT * FROM status");
    if err != nil {
        fmt.Println("unable to make sqlite3 queries, recived error");
        fmt.Println(err);
        return;
    }
    for ; result.Next() == nil; {
        row_map := make(sqlite3.RowMap);
        result.Scan(row_map);
        fmt.Println(row_map);
    }
}

func enterDirInTable(dir string, c *sqlite3.Conn) error {
    file_info, err := ioutil.ReadDir("."+dir);
    if err != nil {
        fmt.Println("encounted unknown error");
        return err;
    }
    for i := 0; i < len(file_info); i++ {
        if isExecutable(file_info[i].Mode()) ||
            isInIgnoreList(file_info[i].Name()) {
            continue;
        }
        if file_info[i].IsDir() {
            err = enterDirInTable(dir+"/"+file_info[i].Name(), c);
            if err != nil {
                return err;
            }
        }
        sql_args := sqlite3.NamedArgs{"$time":file_info[i].ModTime().Unix(),
                                        "$file":dir+"/"+file_info[i].Name()};
        err = c.Exec("INSERT INTO status VALUES($time, $file)", sql_args);
        if err != nil {
            fmt.Println("encounterd unknown error");
            return err;
        }
    }
    return nil;
}

func initLibdevdev() {
    if !pathExist("./.libdevdev") {
        err :=  os.Mkdir(".libdevdev",0775);
        if err != nil {
            fmt.Println("unable to initialize, maybe try running sudo");
            fmt.Println(err);
            return;
        }
    } else {
        fmt.Println("libdevdev seems to be initialized in this directory, exec",
        "ute libdevdev update");
        return;
    }
    c,err := sqlite3.Open(".libdevdev/filestatus.db");
    if err != nil {
        fmt.Println("unable to build or open database file");
        fmt.Println(err);
        return;
    }
    err = c.Exec("CREATE TABLE status(time,file)");
    if err != nil {
        fmt.Println(err);
        c.Close();
        return;
    }
    err = enterDirInTable("", c);
    if err != nil {
        fmt.Println("unable to enter file modification timings");
        fmt.Println(err);
        c.Close();
        return;
    }
    fmt.Println("libdevdev perfectely initialized");
    c.Close();
}

func printStatus() {
    if !pathExist("./.libdevdev/filestatus.db") {
        fmt.Println("libdevdev is not initialized in this directory");
        return;
    }
    c,err := sqlite3.Open(".libdevdev/filestatus.db");
    if err != nil {
        fmt.Println("unable to connect to file status database, fatal error");
        fmt.Println(err);
        return;
    }
    terminalLogStatus(c);
    c.Close();
}

