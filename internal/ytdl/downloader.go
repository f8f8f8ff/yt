package ytdl

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// --no-simulate needs to be present in flags to enable downloading
type Downloader struct {
	Cmd   string
	Flags []string
	Dir   string
}

func (d *Downloader) download(url string, extraflags ...string) (string, error) {
	defaultflags := []string{
		"--quiet",
		"--print", "filename",
	}
	flags := append(d.Flags, extraflags...)
	flags = append(flags, defaultflags...)
	flags = append(flags, url)
	cmd := exec.Command(d.Cmd, flags...)
	var (
		stdout strings.Builder
		stderr strings.Builder
	)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if _, err := os.Stat(d.Dir); os.IsNotExist(err) {
		err = os.MkdirAll(d.Dir, os.ModePerm)
		if err != nil {
			return "", err
		}
	}
	cmd.Dir = d.Dir
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%s: stderr:\n%v", err, stderr.String())
	}
	path := strings.TrimSuffix(stdout.String(), "\n")
	return path, err
}

func (d *Downloader) Video(url string) (string, error) {
	return d.download(url)
}

func (d *Downloader) Audio(url string) (string, error) {
	flags := []string{
		"-x",
		"--audio-format",
		"mp3",
	}
	path, err := d.download(url, flags...)
	if err != nil {
		return path, err
	}
	m := idPathRegexp.FindStringSubmatch(path)
	// trim original format suffix
	path = strings.TrimSuffix(path, m[2])
	path = path + "mp3"
	return path, nil
}
