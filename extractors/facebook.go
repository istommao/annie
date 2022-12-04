package extractors

import (
	"fmt"

	"github.com/istommao/annie/downloader"
	"github.com/istommao/annie/request"
	"github.com/istommao/annie/utils"
)

// Facebook download function
func Facebook(url string) downloader.VideoData {
	html := request.Get(url, url)
	title := utils.MatchOneOf(html, `<title id="pageTitle">(.+)</title>`)[1]

	format := map[string]downloader.FormatData{}
	var tag string
	for _, quality := range []string{"sd", "hd"} {
		tag = quality
		u := utils.MatchOneOf(
			html, fmt.Sprintf(`%s_src_no_ratelimit:"(.+?)"`, quality),
		)[1]
		size := request.Size(u, url)
		urlData := downloader.URLData{
			URL:  u,
			Size: size,
			Ext:  "mp4",
		}
		if quality == "hd" {
			tag = "default"
		}
		format[tag] = downloader.FormatData{
			URLs:    []downloader.URLData{urlData},
			Size:    size,
			Quality: quality,
		}
	}

	extractedData := downloader.VideoData{
		Site:    "Facebook facebook.com",
		Title:   utils.FileName(title),
		Type:    "video",
		Formats: format,
	}
	extractedData.Download(url)
	return extractedData
}
