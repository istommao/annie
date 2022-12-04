package extractors

import (
	"strings"

	"github.com/istommao/annie/downloader"
	"github.com/istommao/annie/request"
	"github.com/istommao/annie/utils"
)

// Weibo download function
func Weibo(url string) downloader.VideoData {
	if !strings.Contains(url, "m.weibo.cn") {
		statusID := utils.MatchOneOf(url, `weibo\.com/tv/v/([^\?/]+)`)[1]
		url = "https://m.weibo.cn/status/" + statusID
	}
	html := request.Get(url, url)
	title := utils.MatchOneOf(html, `"content2": "(.+?)",`)[1]
	realURL := utils.MatchOneOf(
		html, `"stream_url_hd": "(.+?)"`, `"stream_url": "(.+?)"`,
	)[1]
	size := request.Size(realURL, url)
	urlData := downloader.URLData{
		URL:  realURL,
		Size: size,
		Ext:  "mp4",
	}
	format := map[string]downloader.FormatData{
		"default": {
			URLs: []downloader.URLData{urlData},
			Size: size,
		},
	}
	extractedData := downloader.VideoData{
		Site:    "微博 weibo.com",
		Title:   title,
		Type:    "video",
		Formats: format,
	}
	extractedData.Download(url)
	return extractedData
}
