package ytdl

import (
	"errors"
	"net/http"
	"regexp"
)

type file struct {
	id     string
	format string
	path   string
}

type Manager struct {
	downloader Downloader
	files      map[string]*file
	Fs         managerFileSystem
}

// returns the path to a downloaded youtube video by its url
// downloads the video if not existent
func (m *Manager) GetVideo(url string) (string, error) {
	id, err := idFromUrl(url)
	if err != nil {
		return "", err
	}
	file := m.files[id+"/v"]
	if file == nil {
		// download video
		file, err = m.DownloadVideo(url)
		if err != nil {
			return "", err
		}
		return file.path, nil
	}
	return file.path, nil
}

func (m *Manager) DownloadVideo(url string) (*file, error) {
	path, err := m.downloader.Video(url)
	if err != nil {
		return nil, err
	}
	f, err := fileFromPath(path)
	if err != nil {
		return nil, err
	}
	m.files[f.id+"/v"] = f
	return f, nil
}

var idUrlRegexp *regexp.Regexp
var idPathRegexp *regexp.Regexp

func init() {
	// https://stackoverflow.com/questions/3452546/how-do-i-get-the-youtube-video-id-from-a-url
	idUrlRegexp = regexp.MustCompile(`^.*(?:(?:youtu\.be\/|v\/|vi\/|u\/\w\/|embed\/|shorts\/)|(?:(?:watch)?\?v(?:i)?=|\&v(?:i)?=))([^#\&\?]*).*`)
	idPathRegexp = regexp.MustCompile(`\[([^#\&\?]*)\]\.(.*$)`)
}

func idFromUrl(url string) (string, error) {
	match := idUrlRegexp.FindStringSubmatch(url)
	if len(match) < 2 {
		return "", errors.New("couldn't get id from url")
	}
	return match[1], nil
}

func fileFromPath(path string) (*file, error) {
	match := idPathRegexp.FindStringSubmatch(path)
	if len(match) < 2 {
		return nil, errors.New("couldn't get id from path")
	}
	return &file{
		id:     match[1],
		path:   path,
		format: match[2],
	}, nil
}

type managerFileSystem struct {
	fs http.FileSystem
}

func (mfs managerFileSystem) Open(name string) (http.File, error) {
	f, err := mfs.fs.Open(name)
	if err != nil {
		return nil, err
	}
	return f, nil
}
