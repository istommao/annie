package extractors

import (
	"testing"

	"github.com/istommao/annie/config"
	"github.com/istommao/annie/test"
)

func TestDouyu(t *testing.T) {
	config.InfoOnly = true
	tests := []struct {
		name string
		args test.Args
	}{
		{
			name: "normal test",
			args: test.Args{
				URL:   "https://v.douyu.com/show/l0Q8mMY3wZqv49Ad",
				Title: "每日撸报_每日撸报：有些人死了其实它还可以把你带走_斗鱼视频 - 最6的弹幕视频网站",
				Size:  10558080,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := Douyu(tt.args.URL)
			test.Check(t, tt.args, data)
		})
	}
}
