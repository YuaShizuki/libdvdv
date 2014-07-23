package libdvdvupdate

/*
* This is the core update module, its responsible to build a tar.gz file 
* encapsulating all changes made on the project directory by master-dev.
* It can also function to untar the tar.gz file and update the project directory 
* of worker-dev to be in sync with master-dev. check the definition of 
* [master|worker]dev and project directory, on libdevdev.org
* this module would ignore all file pattern found in .libdvdvignore file.
*/
import "errors"

var LibdvdvLog func(a ...interface{}) = func(a ...interface{}){}

/*
* This the main function exposed by this module, it builds the tar.gz file. It 
* stores temporaray directory state in  sqlite3 database (state.db) located in
* .libdvdv directory.
*/
func BuildUpdate(dir string) error {
}





