package archiver

import (
	"fmt"
	"os"
	"testing"
	"time"
)

type ArchTest struct {
	AddBytesFileName string
	AddBytesContent  []byte
	Add              string
}

var ArchTypes = []Archiver{
	&Zip{},
}

var Tests = []ArchTest{
	{
		AddBytesFileName: "some-new-file.txt",
		AddBytesContent:  []byte("hi, i'm new here"),
		Add:              "./tests/loremipsum.txt",
	},
	{
		AddBytesFileName: "some-new-file2.txt",
		AddBytesContent:  []byte("hi, i'm new here too"),
		Add:              "./tests/*.go",
	},
}

func TestArchiver(t *testing.T) {
	if err := MakeTestOutputDir(); err != nil {
		t.Fatal(err)
	}

	for _, arch := range ArchTypes {
		arch.Create(fmt.Sprintf("_test/test-%d", time.Now().Unix()))
		for _, at := range Tests {
			//add bytes
			if err := arch.AddBytes(at.AddBytesFileName, at.AddBytesContent); err != nil {
				t.Fatalf("failed adding bytes to archive %+v: %s", arch, err)
			}

			//add
			if err := arch.Add(at.Add); err != nil {
				t.Fatalf("failed adding file to archive %+v: %s", arch, err)
			}
		}
		arch.Close()
	}
}

func MakeTestOutputDir() error {
	if err := os.RemoveAll("./_test"); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed cleaning test directory: %s", err)
	}

	return os.Mkdir("./_test", 0777)
}
