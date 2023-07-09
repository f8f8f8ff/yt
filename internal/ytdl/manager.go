package ytdl

import (
	"errors"
	"regexp"
)

type File struct {
	id     string
	format string
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

// returns the path to a downloaded youtube video by its url
// downloads the video if not existent
func (m *Manager) GetVideo(url string) (string, error) {
	id, err := idFromUrl(url)
	if err != nil {
		return "", err
	}
	file := m.Files[id+"/v"]
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

func (m *Manager) DownloadVideo(url string) (*File, error) {
	path, err := m.Downloader.Video(url)
	if err != nil {
		return nil, err
	}
	f, err := fileFromPath(path)
	if err != nil {
		return nil, err
	}
	m.Files[f.id+"/v"] = f
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

func fileFromPath(path string) (*File, error) {
	match := idPathRegexp.FindStringSubmatch(path)
	if len(match) < 2 {
		return nil, errors.New("couldn't get id from path")
	}
	return &File{
		id:     match[1],
		path:   path,
		format: match[2],
	}, nil
}
