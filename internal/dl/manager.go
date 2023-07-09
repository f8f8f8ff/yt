package dl

import (
	"errors"
	"log"
	"os"
	"regexp"
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
}

func NewManager() *Manager {
	return &Manager{
		Downloader: Downloader{
			Cmd:   "yt-dlp",
			Flags: []string{},
			Dir:   "./tmp/",
		},
		Files: map[string]*File{},
	}
}

var BadMedium = errors.New("unknown medium")

// returns the path to a downloaded youtube video by its url
// downloads the video if not existent
func (m *Manager) Get(url, medium string) (string, error) {
	id, err := idFromUrl(url)
	if err != nil {
		return "", err
	}
	switch medium {
	case "video":
		id += "/v"
	case "audio":
		id += "/a"
	default:
		return "", BadMedium
	}
	file := m.Files[id]
	if file == nil {
		file, err = m.DownloadVideo(url, medium)
		if err != nil {
			return "", err
		}
		return file.path, nil
	}
	return file.path, nil
}

func (m *Manager) DownloadVideo(url, medium string) (*File, error) {
	log.Println("DOWNLOADING:", url, medium)
	var (
		path string
		err  error
	)
	switch medium {
	case "video":
		path, err = m.Downloader.Video(url)
	case "audio":
		path, err = m.Downloader.Audio(url)
	default:
		return nil, BadMedium
	}
	if err != nil {
		return nil, err
	}
	f, err := fileFromPath(path)
	if err != nil {
		return nil, err
	}
	m.Files[f.id] = f
	return f, nil
}

var idUrlRegexp *regexp.Regexp
var idPathRegexp *regexp.Regexp

func init() {
	// https://stackoverflow.com/questions/3452546/how-do-i-get-the-youtube-video-id-from-a-url
	idUrlRegexp = regexp.MustCompile(`^.*(?:(?:youtu\.be\/|v\/|vi\/|u\/\w\/|embed\/|shorts\/)|(?:(?:watch)?\?v(?:i)?=|\&v(?:i)?=))([^#\&\?]*).*`)
	idPathRegexp = regexp.MustCompile(`\[([^#\&\?]*)\]\.(.*$)`)
}

var BadUrl = errors.New("couldn't get youtube id from url")

func idFromUrl(url string) (string, error) {
	match := idUrlRegexp.FindStringSubmatch(url)
	if len(match) < 2 {
		return "", BadUrl
	}
	return match[1], nil
}

var BadPath = errors.New("couldn't get youtube id from path")

// should return a file struct with "youtubeid/a" or /v for audio or video
func fileFromPath(path string) (*File, error) {
	match := idPathRegexp.FindStringSubmatch(path)
	if len(match) < 2 {
		return nil, BadPath
	}
	id := match[1]
	format := match[2]
	var medium string
	switch format {
	case "webm":
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

func (m *Manager) SetDir(dir string) {
	m.Downloader.Dir = dir
}
