package ytdl

import (
	"fmt"
	"os/exec"
	"strings"
)

// --no-simulate needs to be present in flags to enable downloading
type Downloader struct {
	cmd   string
	flags []string
}

func (d *Downloader) Video(url string) (string, error) {
	flags := append(d.flags, []string{
		"--quiet",
		"--print", "filename",
		url,
	}...)
	cmd := exec.Command(d.cmd, flags...)
	var (
		stdout strings.Builder
		stderr strings.Builder
	)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%s: stderr:\n%v", err, stderr.String())
	}
	path := strings.TrimSuffix(stdout.String(), "\n")
	return path, err
}
