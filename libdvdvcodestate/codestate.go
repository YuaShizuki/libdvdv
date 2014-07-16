package libdvdvcodestate

/*
* This package is responsible to retrive the state of the file system.
* ie which new files are added, which files are removed and which files are 
* modified. Basic use case for this module is to initilize it in a directory
* then call GetState() to get new changes or to update new changes.
*/
var LibdvdvLog func(a ...interface{})(...interface{}) = func(a...interface{}){} 

func Build(path string) error {
}

func GetState() error {
}

func Update() error {

}

func Clean() error {
}
