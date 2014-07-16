package libdvdvcodestate

/*
* Libdevdev Code State builds a sqlite database aka (state-DB) inside .libdvdv 
* directory. state-DB stores the most recent version of the directory state. 
* This module can be used to check the files that have been changed since last
* update. 
* Libdevdev Ignore is a dependancie.
*/

/* personal log function */
var LibdvdvLog func(a ...interface{})(...interface{}) = func(a...interface{}){};

/*
* Builds a database that represents overtime changes in files and directories of 
* "record". "record" is a file path, a directory to start recording changes in.
*/
func BuildStateDB(record string) error {
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
func UpdateStateDB() error {
    return nil;
}

/*
* ClearStateDB() clears the state database.
*/
func ClearStateDB() error {
    return nil;
}
