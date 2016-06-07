package archiver

import (
	"archive/zip"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type (
	Archiver interface {
		Create(path string) error
		AddBytes(path string, b []byte) error
		Add(path string) error
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
func (z *Zip) AddBytes(path string, b []byte) error {
	f, err := z.writer.Create(path)
	if err != nil {
		return err
	}

	_, err = f.Write(b)

	return err
}

func (z *Zip) Add(pred string) error {
	return readFiles(path.Dir(pred), path.Base(pred), func(p string, i os.FileInfo, r io.Reader) error {
		zipfile, err := z.writer.Create(stripRootDir(p, path.Dir(pred)))
		if err != nil {
			return err
		}
		if _, err := io.Copy(zipfile, r); err != nil {
			return err
		}

		return nil
	})
	return nil
}

//Close closes the zipfile
func (z *Zip) Close() error {
	return z.writer.Close()
}

//strips the root directory of a path
func stripRootDir(p string, rootDir string) string {
	if len(p) > 0 && p[:1] != "." {
		return strings.Replace(p, rootDir, "", 1)
	}

	return p
}

//readFiles a helper function that opens and calls a function on every file
//correspondig to the path and predicate
func readFiles(dir string, predicate string, rfunc readFileFunc) error {
	return filepath.Walk(dir, func(p string, i os.FileInfo, err error) error {
		if i != nil && !i.IsDir() {
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
