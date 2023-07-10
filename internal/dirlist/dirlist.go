package dirlist

import (
	"fmt"
	"os"
	"sort"
	"time"
)

type File struct {
	Name  string
	bytes int64
	time  time.Time
}

func (f *File) Size() string {
	const unit = 1000
	b := f.bytes
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp += 1
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}

func Dir(path string) ([]File, error) {
	osfiles, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	files := []File{}
	for _, file := range osfiles {
		if file.IsDir() {
			continue
		}
		info, err := file.Info()
		if err != nil {
			return nil, err
		}
		name := file.Name()
		size := info.Size()
		date := info.ModTime()
		files = append(files, File{
			Name:  name,
			bytes: size,
			time:  date,
		})
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].time.After(files[j].time)
	})
	return files, nil
}
