package dl

import (
	"errors"
	"log"
	"os"
)

type File struct {
	id     string
	format string
	medium string
	path   string
}

type Manager struct {
	Downloader Downloader
	Files      map[string]*File
	Logger     *log.Logger
}

var DefaultManager Manager = Manager{
	Downloader: Downloader{
		Cmd:   "yt-dlp",
		Flags: []string{},
		Dir:   "./tmp/",
	},
	Files:  map[string]*File{},
	Logger: nil,
}

// returns the path to a downloaded youtube video by its url
// downloads the video if not existent
// TODO only returns the path of the first video downloaded
func (m *Manager) Get(url, medium string) ([]string, error) {
	m.log("requesting: %v as %v", url, medium)
	id, err := IdFromUrl(url)
	if errors.Is(err, BadUrl) {
		err = nil
	}
	if err != nil {
		return nil, err
	}
	switch medium {
	case "video":
		id += "/v"
	case "audio":
		id += "/a"
	default:
		return nil, BadMedium
	}
	file := m.Files[id]
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

func (m *Manager) DownloadVideo(url, medium string) ([]*File, error) {
	m.log("downloading %v as %v", url, medium)
	var (
		paths []string
		err   error
	)
	switch medium {
	case "video":
		paths, err = m.Downloader.Video(url)
	case "audio":
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
	m.Files[f.id] = f
}

// should return a file struct with "youtubeid/a" or /v for audio or video
func fileFromPath(path string) (*File, error) {
	id, format, err := IdFromName(path)
	if err != nil {
		return nil, err
	}
	var medium string
	switch format {
	case "webm", "mp4", "mov", "flv":
		id += "/v"
		medium = "video"
	case "mp3":
		id += "/a"
		medium = "audio"
	default:
		return nil, BadMedium
	}
	return &File{
		id:     id,
		path:   path,
		format: format,
		medium: medium,
	}, nil
}

func LoadFiles(dir string) (map[string]*File, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	m := make(map[string]*File)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		f, err := fileFromPath(file.Name())
		if err == nil {
			m[f.id] = f
		}
	}
	return m, nil
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
