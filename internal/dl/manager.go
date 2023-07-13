package dl

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"yt/internal/zip"
)

type Manager struct {
	Downloader Downloader
	Files      map[string]*File
	Logger     *log.Logger
	ZipDir     string
}

var DefaultManager Manager = Manager{
	Downloader: Downloader{
		Cmd:   "yt-dlp",
		Flags: []string{},
		Dir:   "./tmp/",
	},
	Files:  map[string]*File{},
	Logger: nil,
	ZipDir: "./tmp/zip/",
}

// returns the path to a downloaded youtube video by its url
// downloads the video if not existent
// TODO only returns the path of the first video downloaded
func (m *Manager) Get(url string, medium Medium) ([]string, error) {
	m.log("requesting: %v as %v", url, medium)
	id, err := IdFromUrl(url)
	if errors.Is(err, BadUrl) {
		err = nil
	}
	if err != nil {
		return nil, err
	}

	if medium == MediumUnknown {
		return nil, BadMedium
	}
	key := id + medium.Suffix()

	file := m.Files[key]
	paths := []string{}
	if file != nil {
		paths = append(paths, file.path)
		return paths, nil
	}
	// don't have video in the archive, download it
	files, err := m.DownloadVideo(url, medium)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		paths = append(paths, f.path)
	}
	return paths, nil
}

func (m *Manager) DownloadVideo(url string, medium Medium) ([]*File, error) {
	m.log("downloading %v as %v", url, medium)
	var (
		paths []string
		err   error
	)
	switch medium {
	case MediumVideo:
		paths, err = m.Downloader.Video(url)
	case MediumAudio:
		paths, err = m.Downloader.Audio(url)
	default:
		return nil, BadMedium
	}
	if err != nil {
		return nil, err
	}
	files := []*File{}
	for _, p := range paths {
		f, err := fileFromPath(p)
		if err != nil {
			return nil, err
		}
		m.addFile(f)
		files = append(files, f)
	}
	return files, nil
}

func (m *Manager) addFile(f *File) {
	m.log("adding file %v", f)
	m.Files[f.Key()] = f
}

func (m *Manager) Dir() string {
	return m.Downloader.Dir
}

func (m *Manager) log(format string, a ...any) {
	if m.Logger == nil {
		return
	}
	format = "dl.Manager: " + format
	m.Logger.Printf(format, a...)
}

func (m *Manager) Zip(files ...string) (string, error) {
	now := time.Now()
	outPath := fmt.Sprintf("%v_ytdl_%v_files.zip", now.Format("060102150405"), len(files))
	outPath = path.Join(m.ZipDir, outPath)
	m.log("zipping to %v: %v", outPath, files)
	w, err := os.Create(outPath)
	if err != nil {
		return "", err
	}
	defer w.Close()

	for i, f := range files {
		files[i] = path.Join(m.Dir(), f)
	}

	err = zip.ZipFiles(w, files...)
	if err != nil {
		return "", err
	}
	return outPath, nil
}
