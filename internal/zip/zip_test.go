package zip

import (
	"os"
	"testing"
)

func TestCreate(t *testing.T) {
	type args struct {
		files []*os.File
	}
	files, _ := getFiles("testdata/test.mp3", "testdata/test2.mp3")
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		{
			"testdata/1.zip",
			args{files},
			"",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, err := os.Create(tt.name)
			if err != nil {
				t.Error(err)
				return
			}
			defer w.Close()
			err = create(w, tt.args.files...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
