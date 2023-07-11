package dl

import (
	"reflect"
	"testing"
)

func TestDownloader_Video(t *testing.T) {
	d := &Downloader{
		Cmd: "yt-dlp",
		Dir: "/tmp/ytdl",
		Dry: true,
	}
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		d       *Downloader
		args    args
		want    []string
		wantErr bool
	}{
		{
			"should be ok",
			d,
			args{"https://www.youtube.com/watch?v=sCNj0WMBkrs"},
			[]string{"The Epic Battle： Jesus vs Cyborg Satan [sCNj0WMBkrs].webm"},
			false,
		},
		{
			"should fail",
			d,
			args{"https://www.youtube.com/watch?v=sCNj0WMBkr"},
			nil,
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Downloader.Video() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDownloader_Audio(t *testing.T) {
	d := &Downloader{
		Cmd: "yt-dlp",
		Dir: "/tmp/ytdl",
		Dry: true,
	}
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		d       *Downloader
		args    args
		want    []string
		wantErr bool
	}{
		{
			"should be ok",
			d,
			args{"https://www.youtube.com/watch?v=sCNj0WMBkrs"},
			[]string{"The Epic Battle： Jesus vs Cyborg Satan [sCNj0WMBkrs].mp3"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.d
			got, err := d.Audio(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("Downloader.Audio() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Downloader.Audio() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDownloader_playlist(t *testing.T) {
	d := &Downloader{
		Cmd: "yt-dlp",
		Dir: "/tmp/ytdl",
		Dry: true,
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
			"video in playlist",
			d,
			args{"https://www.youtube.com/watch?v=kk-ksRk_a1U&list=PL53CCACFE7503DDF4"},
			`Henry's Dress [01] Definitely Nothing [kk-ksRk_a1U].webm
Henry's Dress [02] Title Fourthcoming [EnWAkV7ux9s].webm
Henry's Dress [03] Sally Wants [HixSmLYu_28].webm
Henry's Dress [04] (You're My) Radio One [kBJjPXMvcE4].webm
Henry's Dress [05] ＂A＂ is For Cribbage [EgCyRZr3pYk].webm
Henry's Dress [06] Three [IhQ43OVzOOw].webm
Henry's Dress [07] Feathers [M3lOCJ--ikw].webm
Henry's Dress [08] You Killed a Boy for Me [7n7gcOs8bTE].webm
Henry's Dress [06] Three [IhQ43OVzOOw].webm`,
			false,
		},
		{
			"https://www.youtube.com/playlist?list=PL53CCACFE7503DDF4",
			d,
			args{"https://www.youtube.com/watch?v=kk-ksRk_a1U&list=PL53CCACFE7503DDF4"},
			`Henry's Dress [01] Definitely Nothing [kk-ksRk_a1U].webm
Henry's Dress [02] Title Fourthcoming [EnWAkV7ux9s].webm
Henry's Dress [03] Sally Wants [HixSmLYu_28].webm
Henry's Dress [04] (You're My) Radio One [kBJjPXMvcE4].webm
Henry's Dress [05] ＂A＂ is For Cribbage [EgCyRZr3pYk].webm
Henry's Dress [06] Three [IhQ43OVzOOw].webm
Henry's Dress [07] Feathers [M3lOCJ--ikw].webm
Henry's Dress [08] You Killed a Boy for Me [7n7gcOs8bTE].webm
Henry's Dress [06] Three [IhQ43OVzOOw].webm`,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.d
			got, err := d.download(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("Downloader.download() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Downloader.download() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_splitMultipleFiles(t *testing.T) {
	files := `Henry's Dress [01] Definitely Nothing [kk-ksRk_a1U].webm
Henry's Dress [02] Title Fourthcoming [EnWAkV7ux9s].webm
Henry's Dress [03] Sally Wants [HixSmLYu_28].webm
Henry's Dress [04] (You're My) Radio One [kBJjPXMvcE4].webm
Henry's Dress [05] ＂A＂ is For Cribbage [EgCyRZr3pYk].webm
Henry's Dress [06] Three [IhQ43OVzOOw].webm
Henry's Dress [07] Feathers [M3lOCJ--ikw].webm
Henry's Dress [08] You Killed a Boy for Me [7n7gcOs8bTE].webm
Henry's Dress [06] Three [IhQ43OVzOOw].webm`
	want := []string{
		"Henry's Dress [01] Definitely Nothing [kk-ksRk_a1U].webm",
		"Henry's Dress [02] Title Fourthcoming [EnWAkV7ux9s].webm",
		"Henry's Dress [03] Sally Wants [HixSmLYu_28].webm",
		"Henry's Dress [04] (You're My) Radio One [kBJjPXMvcE4].webm",
		"Henry's Dress [05] ＂A＂ is For Cribbage [EgCyRZr3pYk].webm",
		"Henry's Dress [06] Three [IhQ43OVzOOw].webm",
		"Henry's Dress [07] Feathers [M3lOCJ--ikw].webm",
		"Henry's Dress [08] You Killed a Boy for Me [7n7gcOs8bTE].webm",
		"Henry's Dress [06] Three [IhQ43OVzOOw].webm",
	}
	got := splitMultipleFiles(files)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("splitMultipleFiles() = %v, want %v", got, want)
	}
}
