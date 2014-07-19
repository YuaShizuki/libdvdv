package libdvdvcodestate

/*
* Libdevdev Code State builds a sqlite database aka (state-DB) inside .libdvdv 
* directory. state-DB stores the most recent version of the directory state. 
* This module can be used to check the files that have been changed since last
* update. codestate rejects files that are mentioned in the .libdvdvignore file.
* hence libdvdvign is a dependancie.
*/

import "../libdvdvutil"
import "os"
import "database/sql"
import _ "github.com/mattn/go-sqlite3"
import "errors"
import "io/ioutil"

/* personal log function */
var LibdvdvLog func(a ...interface{}) = func(a...interface{}){};

/*
* Builds a database that represents overtime changes in files and directories of 
* "record". "record" is a file path, a directory to start recording changes in.
*/
func BuildStateDB(record string) error {
    libdvdv_dir := record+"/.libdvdv";
    state_db := libdvdv_dir+"/state.db";
    if exist,_ := libdvdvutil.PathExist(libdvdv_dir); !exist {
        if err := os.Mkdir(libdvdv_dir, 0744); err != nil {
            return err;
        }
    } else {
        if exist,_ := libdvdvutil.PathExist(state_db); exist {
            if err := os.Remove(state_db); err != nil {
                return err;
            }
        }
    }
    if err := UpdateStateDB(record); err != nil {
        return err;
    }
    return nil;
}


/*
* GetState() shows the changes since last update.
*/
func GetState() error {
    return nil;
}

/*
* Update() updates the state DB. 
*/
func UpdateStateDB(record string) error {
    conn, err := sql.Open("sqlite3", record+"/.libdvdv/state.db");
    if err != nil {
        return err;
    }
    defer conn.Close();
    row := conn.QueryRow("select name from sqlite_master where type=table "+
                            "and name=State;");
    if row != nil {
        var db_name string;
        err := row.Scan(&db_name);
        if (err != nil) || (len(db_name) == 0) {
            _, err = conn.Exec("create table State(file TEXT, modified INTEGER);");
            if err != nil {
                return err;
            }
        }
    }
    if err := enterDirState(record, len(record)+1, conn); err != nil {
        return err;
    }
    return nil;
}

func enterDirState(dir string, base int, conn *sql.DB) error {
    if conn == nil {
        return errors.New("error: unopend connection");
    }
    finfo, err := ioutil.ReadDir(dir);
    if err != nil {
        return err;
    }
    for i := range finfo {
        new_entry := dir + "/" + finfo[i].Name();
        if finfo[i].IsDir() {
            err := enterDirState(new_entry, base, conn);
            if err != nil {
                return err;
            }
        } else {
            _, err := conn.Exec("INSERT INTO State VALUES(?, ?)",
                                new_entry[base:len(new_entry)], finfo[i].ModTime().UnixNano());
            if err != nil {
                return err;
            }
        }
    }
    return nil;
}

/*
* ClearStateDB() clears the state-DB.
*/
func ClearStateDB() error {
    return nil;
}
