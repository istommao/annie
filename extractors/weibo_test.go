package extractors

import (
	"testing"

	"github.com/istommao/annie/config"
	"github.com/istommao/annie/test"
)

func TestWeibo(t *testing.T) {
	config.InfoOnly = true
	tests := []struct {
		name string
		args test.Args
	}{
		{
			name: "normal test",
			args: test.Args{
				URL:   "https://m.weibo.cn/2815133121/G9VBqbsWM",
				Title: "当你超过25岁再去夜店……",
				Size:  3112080,
			},
		},
		{
			name: "fid url test",
			args: test.Args{
				URL:   "https://weibo.com/tv/v/Ga7XazXze?fid=1034:4a65c6e343dc672789d3ba49c2463c6a",
				Title: "看完更加睡不着了[二哈]",
				Size:  438757,
			},
		},
		{
			name: "title test",
			args: test.Args{
				URL:   "https://m.weibo.cn/status/4237529215145705",
				Title: `近日，日本视错觉大师、明治大学特任教授\"杉原厚吉的“错觉箭头“作品又引起世界人民的关注。反射，透视和视角的巧妙结合产生了这种惊人的幻觉：箭头向右？转过来还是向右？\n\n引用杉原教授的经典描述：“我们看外面的世界的方式——也就是我们的知觉——都是由大脑机制间接产生的，所以所有知觉在某`,
				Size:  1125984,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := Weibo(tt.args.URL)
			test.Check(t, tt.args, data)
		})
	}
}
