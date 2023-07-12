package dl

import (
	"reflect"
	"testing"
)

func TestManager_Get(t *testing.T) {
	m := Manager{
		Files: map[string]*File{},
		Downloader: Downloader{
			Cmd: "yt-dlp",
			Dir: "./tmp/ytdl",
			Dry: true,
		},
	}
	type args struct {
		url    string
		medium string
	}
	tests := []struct {
		name    string
		m       *Manager
		args    args
		want    []string
		wantErr bool
	}{
		{
			"single video",
			&m,
			args{"https://www.youtube.com/watch?v=sCNj0WMBkrs", "video"},
			[]string{"The Epic Battle： Jesus vs Cyborg Satan [sCNj0WMBkrs].webm"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.Get(tt.args.url, tt.args.medium)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_IdFromUrl(t *testing.T) {
	urls := []string{
		"https://youtube.com/shorts/dQw4w9WgXcQ?feature=share",
		"//www.youtube-nocookie.com/embed/up_lNV-yoK4?rel=0",
		"http://www.youtube.com/user/Scobleizer#p/u/1/1p3vcRhsYGo",
		"http://www.youtube.com/watch?v=cKZDdG9FTKY&feature=channel",
		"http://www.youtube.com/watch?v=yZ-K7nCVnBI&playnext_from=TL&videos=osPknwzXEas&feature=sub",
		"http://www.youtube.com/ytscreeningroom?v=NRHVzbJVx8I",
		"http://www.youtube.com/user/SilkRoadTheatre#p/a/u/2/6dwqZw0j_jY",
		"http://youtu.be/6dwqZw0j_jY",
		"http://www.youtube.com/watch?v=6dwqZw0j_jY&feature=youtu.be",
		"http://youtu.be/afa-5HQHiAs",
		"http://www.youtube.com/user/Scobleizer#p/u/1/1p3vcRhsYGo?rel=0",
		"http://www.youtube.com/watch?v=cKZDdG9FTKY&feature=channel",
		"http://www.youtube.com/watch?v=yZ-K7nCVnBI&playnext_from=TL&videos=osPknwzXEas&feature=sub",
		"http://www.youtube.com/ytscreeningroom?v=NRHVzbJVx8I",
		"http://www.youtube.com/embed/nas1rJpm7wY?rel=0",
		"http://www.youtube.com/watch?v=peFZbP64dsU",
		"http://youtube.com/v/dQw4w9WgXcQ?feature=youtube_gdata_player",
		"http://youtube.com/vi/dQw4w9WgXcQ?feature=youtube_gdata_player",
		"http://youtube.com/?v=dQw4w9WgXcQ&feature=youtube_gdata_player",
		"http://www.youtube.com/watch?v=dQw4w9WgXcQ&feature=youtube_gdata_player",
		"http://youtube.com/?vi=dQw4w9WgXcQ&feature=youtube_gdata_player",
		"http://youtube.com/watch?v=dQw4w9WgXcQ&feature=youtube_gdata_player",
		"http://youtube.com/watch?vi=dQw4w9WgXcQ&feature=youtube_gdata_player",
		"http://youtu.be/dQw4w9WgXcQ?feature=youtube_gdata_player",
		"/user/Scobleizer#p/u/1/1p3vcRhsYGo",
		"/watch?v=cKZDdG9FTKY&feature=channel",
		"/watch?v=yZ-K7nCVnBI&playnext_from=TL&videos=osPknwzXEas&feature=sub",
	}
	for _, url := range urls {
		t.Run(url, func(t *testing.T) {
			got, err := IdFromUrl(url)
			t.Logf("%v = %v", url, got)
			if err != nil {
				t.Errorf("IdFromUrl() error = %v", err)
			}
			if len(got) != 11 {
				t.Errorf("IdFromUrl() = %v", got)
			}
		})
	}
}

func Test_fileFromPath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    *File
		wantErr bool
	}{
		{
			".webm",
			args{"The Epic Battle： Jesus vs Cyborg Satan [sCNj0WMBkrs].webm"},
			&File{
				id:     "sCNj0WMBkrs",
				format: "webm",
				medium: MediumVideo,
				path:   "The Epic Battle： Jesus vs Cyborg Satan [sCNj0WMBkrs].webm",
				source: SourceYoutube,
			},
			false,
		},
		{
			".mp3",
			args{"song [dQw4w9WgXcQ].mp3"},
			&File{
				id:     "dQw4w9WgXcQ",
				format: "mp3",
				medium: MediumAudio,
				path:   "song [dQw4w9WgXcQ].mp3",
				source: SourceYoutube,
			},
			false,
		},
		{
			"mp3 soundcloud",
			args{"Theprivilegeoflovingyou [1444424830].mp3"},
			&File{
				id:     "1444424830",
				format: "mp3",
				medium: MediumAudio,
				path:   "Theprivilegeoflovingyou [1444424830].mp3",
				source: SourceSoundcloud,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fileFromPath(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("fileFromPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fileFromPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_LoadFiles(t *testing.T) {
	want := make(map[string]*File)
	want["dQw4w9WgXcQ/v"] = &File{
		id:     "dQw4w9WgXcQ",
		format: "webm",
		medium: MediumVideo,
		path:   "fake video [dQw4w9WgXcQ].webm",
		source: SourceYoutube,
	}
	want["NRHVzbJVx8I/a"] = &File{
		id:     "NRHVzbJVx8I",
		format: "mp3",
		medium: MediumAudio,
		path:   "unreal song [NRHVzbJVx8I].mp3",
		source: SourceYoutube,
	}
	want["1444424830/a"] = &File{
		id:     "1444424830",
		format: "mp3",
		medium: MediumAudio,
		path:   "Theprivilegeoflovingyou [1444424830].mp3",
		source: SourceSoundcloud,
	}
	got, err := LoadFiles("testdata")
	if err != nil {
		t.Errorf("LoadFiles() error = %v", err)
		return
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("LoadFiles() = %v, want %v", got, want)
	}
}

func Test_idFromName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name       string
		args       args
		wantId     string
		wantFormat string
		wantErr    bool
	}{
		{
			"normal",
			args{"The Epic Battle： Jesus vs Cyborg Satan [sCNj0WMBkrs].webm"},
			"sCNj0WMBkrs",
			"webm",
			false,
		},
		{
			"has [] in title",
			args{"Henry's Dress [07] Feathers [M3lOCJ--ikw].mp3"},
			"M3lOCJ--ikw",
			"mp3",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotId, gotFormat, err := IdFromName(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("idFromName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotId != tt.wantId {
				t.Errorf("idFromName() gotId = %v, want %v", gotId, tt.wantId)
			}
			if gotFormat != tt.wantFormat {
				t.Errorf("idFromName() gotFormat = %v, want %v", gotFormat, tt.wantFormat)
			}
		})
	}
}
