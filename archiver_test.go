package archiver

import (
	"fmt"
	"os"
	"testing"
	"time"
)

type ArchTest struct {
	AddFileName     string
	AddContent      []byte
	AddFile         string
	AddDir          string
	AddAllDir       string
	AddAllPredicate string
}

var ArchTypes = []Archiver{
	&Zip{},
}

var Tests = []ArchTest{
	{
		AddFileName:     "some-new-file.txt",
		AddContent:      []byte("hi, i'm new here"),
		AddFile:         "./tests/loremipsum.txt",
		AddDir:          "./tests/dir",
		AddAllDir:       "./tests/dir2",
		AddAllPredicate: "*.go",
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
			if err := arch.Add(at.AddFileName, at.AddContent); err != nil {
				t.Fatalf("failed adding bytes to archive %+v: %s", arch, err)
			}

			//addfile
			if err := arch.AddFile(at.AddFile); err != nil {
				t.Fatalf("failed adding file to archive %+v: %s", arch, err)
			}

			//addDir
			if err := arch.AddDir(at.AddDir); err != nil {
				t.Fatalf("failed adding dir to archive %+v: %s", arch, err)
			}

			//addAll
			if err := arch.AddAll(at.AddAllDir, at.AddAllPredicate); err != nil {
				t.Fatalf("failed adding all with predicate %s, to archive %+v: %s",
					at.AddAllPredicate,
					arch,
					err,
				)
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
