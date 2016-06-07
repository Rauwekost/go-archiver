package archiver

import (
	"archive/zip"
	"bufio"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type (
	Archiver interface {
		Create(path string) error
		Add(path string, b []byte) error
		AddFile(path string) error
		AddDir(path string) error
		AddAll(path, predicate string) error
		Close() error
	}

	Zip struct {
		writer *zip.Writer
	}

	readFileFunc func(p string, info os.FileInfo, file io.Reader) error
)

//Create creates a new zipfile
func (z *Zip) Create(path string) error {
	if strings.HasSuffix(path, ".zip") != true {
		path = path + ".zip"
	}
	file, err := os.Create(path)

	if err != nil {
		return err
	}
	z.writer = zip.NewWriter(file)

	return nil
}

//Add creates a file in the zip archive
func (z *Zip) Add(path string, b []byte) error {
	f, err := z.writer.Create(path)
	if err != nil {
		return err
	}

	_, err = f.Write(b)

	return err
}

//Add a file from path to the zip archive
func (z *Zip) AddFile(path string) error {
	file, err := os.Open(filepath.Join(path))
	if err != nil {
		return err
	}
	defer file.Close()
	fileReader := bufio.NewReader(file)

	zfile, err := z.writer.Create(path)
	if err != nil {
		return err
	}

	io.Copy(zfile, fileReader)
	return nil
}

//AddDir add a directory and it's content to the zip archive
func (z *Zip) AddDir(path string) error {
	return z.AddAll(path, "*")
}

//AddAll adds all files within a path and a search predicate
func (z *Zip) AddAll(path string, predicate string) error {
	return readFiles(path, predicate, func(p string, i os.FileInfo, r io.Reader) error {
		zfile, _ := z.writer.Create(p)
		if _, err := io.Copy(zfile, r); err != nil {
			return err
		}
		return nil
	})
}

//Close closes the zipfile
func (z *Zip) Close() error {
	return z.writer.Close()
}

//readFiles a helper function that opens and calls a function on every file
//correspondig to the path and predicate
func readFiles(p string, predicate string, rfunc readFileFunc) error {
	p = path.Clean(p)
	return filepath.Walk(p, func(p string, i os.FileInfo, err error) error {
		if !i.IsDir() {
			matched, err := filepath.Match(predicate, i.Name())
			if err != nil {
				return err
			}
			if !matched {
				return nil
			}
			r, err := os.Open(p)
			if err != nil {
				return err
			}
			defer r.Close()
			if err := rfunc(p, i, r); err != nil {
				return err
			}
		}
		return nil
	})
}
