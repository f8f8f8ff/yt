package dl

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
	Dry   bool
}

func (d *Downloader) download(url string, extraflags ...string) (string, error) {
	defaultflags := []string{
		"--quiet",
		"--print", "filename",
		"--no-mtime",
	}
	flags := append(d.Flags, extraflags...)
	if !d.Dry {
		flags = append(flags, "--no-simulate")
	}
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

func (d *Downloader) Video(url string) ([]string, error) {
	p, err := d.download(url)
	if err != nil {
		return nil, err
	}
	return splitMultipleFiles(p), nil
}

func (d *Downloader) Audio(url string) ([]string, error) {
	flags := []string{
		"-x",
		"--audio-format", "mp3",
		"--audio-quality", "0",
	}
	p, err := d.download(url, flags...)
	if err != nil {
		return nil, err
	}
	paths := splitMultipleFiles(p)
	for i, f := range paths {
		m := idPathRegexp.FindStringSubmatch(p)
		// trim original format suffix
		f = strings.TrimSuffix(f, m[2])
		paths[i] = f + "mp3"
	}
	return paths, nil
}

func splitMultipleFiles(files string) []string {
	return strings.Split(strings.TrimSuffix(files, "\n"), "\n")
}
