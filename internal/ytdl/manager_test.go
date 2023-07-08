package ytdl

import (
	"reflect"
	"testing"
)

func Test_getId(t *testing.T) {
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
	}
	for _, url := range urls {
		t.Run(url, func(t *testing.T) {
			got, err := getId(url)
			if err != nil {
				t.Errorf("getId() error = %v", err)
			}
			if len(got) != 11 {
				t.Errorf("getId() = %v", got)
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
		want    *file
		wantErr bool
	}{
		{
			".webm",
			args{"The Epic Battle： Jesus vs Cyborg Satan [sCNj0WMBkrs].webm"},
			&file{
				id:     "sCNj0WMBkrs",
				format: "webm",
				path:   "The Epic Battle： Jesus vs Cyborg Satan [sCNj0WMBkrs].webm",
			},
			false,
		},
		{
			".webm",
			args{"song [dQw4w9WgXcQ].mp3"},
			&file{
				id:     "dQw4w9WgXcQ",
				format: "mp3",
				path:   "song [dQw4w9WgXcQ].mp3",
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
