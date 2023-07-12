package dl

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// --no-simulate needs to be present in flags to enable downloading
type Downloader struct {
	Cmd    string
	Flags  []string
	Dir    string
	Dry    bool
	Logger *log.Logger
}

func (d *Downloader) download(url string, extraflags ...string) (string, error) {
	d.log("downloading %v with flags %v", url, extraflags)
	defaultflags := []string{
		"--quiet",
		"--print", "filename",
		"--no-mtime",
		"-o", `"%(playlist_index&[{}] |)s%(uploader)s - %(title)s [%(id)s].%(ext)s"`,
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
		_, format, err := IdFromName(p)
		if err != nil {
			return nil, err
		}
		// trim original format suffix
		f = strings.TrimSuffix(f, format)
		paths[i] = f + "mp3"
	}
	return paths, nil
}

func (d *Downloader) log(format string, a ...any) {
	if d.Logger == nil {
		return
	}
	format = "dl.Downloader: " + format
	d.Logger.Printf(format, a...)
}

func splitMultipleFiles(files string) []string {
	return strings.Split(strings.TrimSuffix(files, "\n"), "\n")
}
