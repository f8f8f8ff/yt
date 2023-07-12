package dl

import (
	"os"
)

type Source int

const (
	SourceUnknown Source = iota
	SourceYoutube
	SourceSoundcloud
)

func (s Source) String() string {
	switch s {
	case SourceYoutube:
		return "youtube"
	case SourceSoundcloud:
		return "soundcloud"
	}
	return "unknown"
}

type Medium int

const (
	MediumUnknown Medium = iota
	MediumAudio
	MediumVideo
)

func (m Medium) String() string {
	switch m {
	case MediumAudio:
		return "audio"
	case MediumVideo:
		return "video"
	}
	return "unknown"
}

func (m Medium) Suffix() string {
	switch m {
	case MediumAudio:
		return "/a"
	case MediumVideo:
		return "/v"
	}
	return ""
}

type File struct {
	id     string
	format string
	medium Medium
	path   string
	source Source
}

func (f *File) Key() string {
	return f.id + f.medium.Suffix()
}

// should return a file struct with "youtubeid/a" or /v for audio or video
func fileFromPath(path string) (*File, error) {
	id, format, err := IdFromName(path)
	if err != nil {
		return nil, err
	}

	source := SourceFromId(id)

	// determine medium by file extension
	var medium Medium
	switch format {
	case "webm", "mp4", "mov", "flv":
		medium = MediumVideo
	case "mp3":
		medium = MediumAudio
	default:
		return nil, BadMedium
	}

	return &File{
		id:     id,
		format: format,
		medium: medium,
		path:   path,
		source: source,
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
			m[f.Key()] = f
		}
	}
	return m, nil
}
