package ytdl

import "testing"

func TestDownloader_Video(t *testing.T) {
	d := &Downloader{
		cmd: "yt-dlp",
		flags: []string{
			"--simulate",
		},
	}
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		d       *Downloader
		args    args
		want    string
		wantErr bool
	}{
		{
			"should be ok",
			d,
			args{"https://www.youtube.com/watch?v=sCNj0WMBkrs"},
			"The Epic Battle： Jesus vs Cyborg Satan [sCNj0WMBkrs].webm",
			false,
		},
		{
			"should fail",
			d,
			args{"https://www.youtube.com/watch?v=sCNj0WMBkr"},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.d
			got, err := d.Video(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("Downloader.Video() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Downloader.Video() = %v, want %v", got, tt.want)
			}
		})
	}
}
