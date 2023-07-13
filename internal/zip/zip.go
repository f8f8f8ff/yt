package zip

import (
	"archive/zip"
	"io"
	"os"
	"path"
)

func create(w io.Writer, files ...*os.File) error {
	zipWriter := zip.NewWriter(w)
	for _, f := range files {
		w, err := zipWriter.Create(path.Base(f.Name()))
		if err != nil {
			return err
		}
		_, err = io.Copy(w, f)
		if err != nil {
			return err
		}
		f.Close()
	}
	zipWriter.Close()
	return nil
}

func ZipFiles(w io.Writer, paths ...string) error {
	files, err := getFiles(paths...)
	if err != nil {
		return err
	}
	create(w, files...)
	return nil
}

func getFiles(paths ...string) ([]*os.File, error) {
	files := []*os.File{}
	for _, p := range paths {
		f, err := os.Open(p)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	return files, nil
}
