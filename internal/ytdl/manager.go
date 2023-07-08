package ytdl

import (
	"errors"
	"regexp"
)

type file struct {
	id     string
	format string
	path   string
}

type Manager struct {
	files map[string]*file
}

// returns the path to a downloaded youtube video by its url
// downloads the video if not existent
func (m *Manager) GetVideo(url string) string {
	return url
}

var idUrlRegexp *regexp.Regexp
var idPathRegexp *regexp.Regexp

func init() {
	// https://stackoverflow.com/questions/3452546/how-do-i-get-the-youtube-video-id-from-a-url
	idUrlRegexp = regexp.MustCompile(`^.*(?:(?:youtu\.be\/|v\/|vi\/|u\/\w\/|embed\/|shorts\/)|(?:(?:watch)?\?v(?:i)?=|\&v(?:i)?=))([^#\&\?]*).*`)
	idPathRegexp = regexp.MustCompile(`\[([^#\&\?]*)\]\.(.*$)`)
}

func getId(url string) (string, error) {
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
