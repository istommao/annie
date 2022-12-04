package extractors

import (
	"github.com/istommao/annie/downloader"
	"github.com/istommao/annie/parser"
	"github.com/istommao/annie/request"
)

// Pixivision download function
func Pixivision(url string) downloader.VideoData {
	html := request.Get(url, url)
	title, urls := parser.GetImages(url, html, "am__work__illust  ", nil)
	format := map[string]downloader.FormatData{
		"default": {
			URLs: urls,
			Size: 0,
		},
	}
	extractedData := downloader.VideoData{
		Site:    "pixivision pixivision.net",
		Title:   title,
		Type:    "image",
		Formats: format,
	}
	extractedData.Download(url)
	return extractedData
}
