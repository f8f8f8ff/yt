package dl

import (
	"fmt"
	"regexp"
)

func PartialUrlToYoutube(uri string) (string, error) {
	id, err := IdFromUrl(uri)
	if err != nil {
		return "", err
	}
	url := UrlFromId(id)
	return url, nil
}

func UrlFromId(id string) string {
	return fmt.Sprintf("https://youtube.com/watch?v=%v", id)
}

func IdFromUrl(url string) (string, error) {
	match := idUrlRegexp.FindStringSubmatch(url)
	if len(match) < 2 {
		return "", BadUrl
	}
	return match[1], nil
}

func IdFromName(name string) (id, format string, err error) {
	match := idPathRegexp.FindStringSubmatch(name)
	if len(match) < 2 {
		return "", "", BadPath
	}
	id = match[1]
	format = match[2]
	return id, format, nil
}

var idUrlRegexp *regexp.Regexp
var idPathRegexp *regexp.Regexp

func init() {
	// https://stackoverflow.com/questions/3452546/how-do-i-get-the-youtube-video-id-from-a-url
	idUrlRegexp = regexp.MustCompile(`^.*(?:(?:youtu\.be\/|v\/|vi\/|u\/\w\/|embed\/|shorts\/)|(?:(?:watch)?\?v(?:i)?=|\&v(?:i)?=))([^#\&\?]*).*`)
	idPathRegexp = regexp.MustCompile(`.+\[([^#\&\?]{11})\]\.(.*$)`)
}
