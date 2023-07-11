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
// TODO only returns the path of the first video downloaded
func (m *Manager) Get(url, medium string) ([]string, error) {
	id, err := IdFromUrl(url)
	// if we cant get the id from the url, go ahead and continue
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
	log.Println("DOWNLOADING:", url, medium)
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
		m.Files[f.id] = f
		files = append(files, f)
	}
	return files, nil
}

var idUrlRegexp *regexp.Regexp
var idPathRegexp *regexp.Regexp

func init() {
	// https://stackoverflow.com/questions/3452546/how-do-i-get-the-youtube-video-id-from-a-url
	idUrlRegexp = regexp.MustCompile(`^.*(?:(?:youtu\.be\/|v\/|vi\/|u\/\w\/|embed\/|shorts\/)|(?:(?:watch)?\?v(?:i)?=|\&v(?:i)?=))([^#\&\?]*).*`)
	idPathRegexp = regexp.MustCompile(`\[([^#\&\?]*)\]\.(.*$)`)
}

var BadUrl = errors.New("couldn't get youtube id from url")

func IdFromUrl(url string) (string, error) {
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

func (m *Manager) SetDir(dir string) {
	m.Downloader.Dir = dir
}
