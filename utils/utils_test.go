package utils

import (
	"reflect"
	"testing"

	"github.com/istommao/annie/config"
)

func TestMatchOneOf(t *testing.T) {
	type args struct {
		patterns []string
		text     string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "normal test",
			args: args{
				patterns: []string{`aaa(\d+)`, `hello(\d+)`},
				text:     "hello12345",
			},
			want: []string{
				"hello12345", "12345",
			},
		},
		{
			name: "normal test",
			args: args{
				patterns: []string{`aaa(\d+)`, `bbb(\d+)`},
				text:     "hello12345",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MatchOneOf(tt.args.text, tt.args.patterns...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MatchOneOf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMatchAll(t *testing.T) {
	type args struct {
		pattern string
		text    string
	}
	tests := []struct {
		name string
		args args
		want [][]string
	}{
		{
			name: "normal test",
			args: args{
				pattern: `hello(\d+)`,
				text:    "hello12345hello123",
			},
			want: [][]string{
				{
					"hello12345", "12345",
				},
				{
					"hello123", "123",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MatchAll(tt.args.text, tt.args.pattern); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MatchAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileSize(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "normal test",
			args: args{
				filePath: "hello",
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := FileSize(tt.args.filePath); got != tt.want {
				t.Errorf("FileSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDomain(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "normal test",
			args: args{
				url: "http://www.aa.com",
			},
			want: "aa",
		},
		{
			name: "normal test",
			args: args{
				url: "https://aa.com",
			},
			want: "aa",
		},
		{
			name: "normal test",
			args: args{
				url: "aa.cn",
			},
			want: "aa",
		},
		{
			name: "normal test",
			args: args{
				url: "www.aa.cn",
			},
			want: "aa",
		},
		{
			name: ".com.cn test",
			args: args{
				url: "http://www.aa.com.cn",
			},
			want: "aa",
		},
		{
			name: "Universal test",
			args: args{
				url: "http://aa",
			},
			want: "Universal",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Domain(tt.args.url); got != tt.want {
				t.Errorf("Domain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLimitLength(t *testing.T) {
	type args struct {
		s      string
		length int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "normal test",
			args: args{
				s:      "你好 hello",
				length: 8,
			},
			want: "你好 hello",
		},
		{
			name: "normal test",
			args: args{
				s:      "你好 hello",
				length: 6,
			},
			want: "你好 ...",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LimitLength(tt.args.s, tt.args.length); got != tt.want {
				t.Errorf("LimitLength() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "normal test",
			args: args{
				name: "hello/world",
			},
			want: "hello world",
		},
		{
			name: "normal test",
			args: args{
				name: "hello:world",
			},
			want: "hello：world",
		},
		{
			name: "overly long strings test",
			args: args{
				name: "super 超级长超级长超级长超级长超级长超级长超级长超级长超级长超级长超级长超级长超级长超级长超级长超级长超级长超级长超级长超级长超级长超级长超级长超级长超级长", // length 81
			},
			want: "super 超级长超级长超级长超级长超级长超级长超级长超级长超级长超级长超级长超级长超级长超级长超级长超级长超级长超级长超级长超级长超级长超级长超级长超级...",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FileName(tt.args.name); got != tt.want {
				t.Errorf("FileName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilePath(t *testing.T) {
	type args struct {
		name   string
		ext    string
		escape bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "normal test",
			args: args{
				name:   "hello",
				ext:    "txt",
				escape: false,
			},
			want: "hello.txt",
		},
		{
			name: "normal test",
			args: args{
				name:   "hello:world",
				ext:    "txt",
				escape: true,
			},
			want: "hello：world.txt",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FilePath(tt.args.name, tt.args.ext, tt.args.escape); got != tt.want {
				t.Errorf("FilePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringInSlice(t *testing.T) {
	type args struct {
		str  string
		list []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "normal test",
			args: args{
				str: "hello",
				list: []string{
					"hello", "abc",
				},
			},
			want: true,
		},
		{
			name: "normal test",
			args: args{
				str: "123",
				list: []string{
					"hello", "abc",
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringInSlice(tt.args.str, tt.args.list); got != tt.want {
				t.Errorf("StringInSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetNameAndExt(t *testing.T) {
	type args struct {
		uri string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		{
			name: "normal test",
			args: args{
				uri: "https://img9.bcyimg.com/drawer/15294/post/1799t/1f5a87801a0711e898b12b640777720f.jpg/w650",
			},
			want:  "w650",
			want1: "jpeg",
		},
		{
			name: "normal test",
			args: args{
				uri: "https://img9.bcyimg.com/drawer/15294/post/1799t/1f5a87801a0711e898b12b640777720f.jpg",
			},
			want:  "1f5a87801a0711e898b12b640777720f",
			want1: "jpg",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetNameAndExt(tt.args.uri)
			if got != tt.want {
				t.Errorf("GetNameAndExt() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetNameAndExt() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestMd5(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "normal test",
			args: args{
				text: "123456",
			},
			want: "e10adc3949ba59abbe56e057f20f883e",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Md5(tt.args.text); got != tt.want {
				t.Errorf("Md5() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrintVersion(t *testing.T) {
	PrintVersion()
}

func TestShouldExtract(t *testing.T) {
	type args struct {
		format      string
		bestQuality string
	}
	tests := []struct {
		name          string
		args          args
		want          bool
		infoOnly      bool
		extractedData bool
		format        string
	}{
		{
			name: "InfoOnly test",
			args: args{
				format:      "1",
				bestQuality: "2",
			},
			want:     true,
			infoOnly: true,
		},
		{
			name: "ExtractedData test",
			args: args{
				format:      "1",
				bestQuality: "2",
			},
			want:          true,
			extractedData: true,
		},
		{
			name: "Format test",
			args: args{
				format:      "bd",
				bestQuality: "bd2",
			},
			want:   true,
			format: "bd",
		},
		{
			name: "Format test2",
			args: args{
				format:      "bd2",
				bestQuality: "bd2",
			},
			want:   false,
			format: "bd",
		},
		{
			name: "bestQuality test",
			args: args{
				format:      "bd2",
				bestQuality: "bd2",
			},
			want: true,
		},
		{
			name: "bestQuality test2",
			args: args{
				format:      "bd",
				bestQuality: "bd2",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.ExtractedData = tt.extractedData
			config.InfoOnly = tt.infoOnly
			config.Format = tt.format
			if got := ShouldExtract(tt.args.format, tt.args.bestQuality); got != tt.want {
				t.Errorf("ShouldExtract() = %v, want %v", got, tt.want)
			}
		})
	}
}
