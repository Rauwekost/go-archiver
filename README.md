Go-Archiver
-----------

A simple wrapper to create archives, for now only zip is avaiable will add tar and tar.gz later.

## Example

```go
package main

import "github.com/rauwekost/go-archiver"


func main() {
	//create a zip instance
	zip := new(archiver.Zip)
	
	//create the zip file on disk
	zip.Create("./testing.zip")
	
	//add a bytes to a file inside the zip
	zip.Add("somefile.txt", []byte("someting inside"))
	
	//recursive add dir
	zip.AddDir("./testdir")
	
	//add all files in dir matching predicate
	zip.AddAll(".", "*.go")
	
	//close the zip
	zip.Close()
}
```