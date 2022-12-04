package extractors

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/istommao/annie/downloader"
	"github.com/istommao/annie/parser"
	"github.com/istommao/annie/request"
	"github.com/istommao/annie/utils"
)

type imageList struct {
	List []string `json:"@list"`
}

type tumblrImageList struct {
	Image imageList `json:"image"`
}

type tumblrImage struct {
	Image string `json:"image"`
}

func genURLData(url, referer string) (downloader.URLData, int64) {
	size := request.Size(url, referer)
	_, ext := utils.GetNameAndExt(url)
	data := downloader.URLData{
		URL:  url,
		Size: size,
		Ext:  ext,
	}
	return data, size
}

func tumblrImageDownload(url, html, title string) downloader.VideoData {
	jsonString := utils.MatchOneOf(
		html, `<script type="application/ld\+json">\s*(.+?)</script>`,
	)[1]
	var totalSize int64
	var urls []downloader.URLData
	if strings.Contains(jsonString, `"image":{"@list"`) {
		// there are two data structures in the same field(image)
		var imageList tumblrImageList
		json.Unmarshal([]byte(jsonString), &imageList)
		for _, u := range imageList.Image.List {
			urlData, size := genURLData(u, url)
			totalSize += size
			urls = append(urls, urlData)
		}
	} else {
		var image tumblrImage
		json.Unmarshal([]byte(jsonString), &image)
		urlData, size := genURLData(image.Image, url)
		totalSize = size
		urls = append(urls, urlData)
	}
	format := map[string]downloader.FormatData{
		"default": {
			URLs: urls,
			Size: totalSize,
		},
	}

	extractedData := downloader.VideoData{
		Site:    "Tumblr tumblr.com",
		Title:   title,
		Type:    "image",
		Formats: format,
	}
	extractedData.Download(url)
	return extractedData
}

func tumblrVideoDownload(url, html, title string) downloader.VideoData {
	videoURL := utils.MatchOneOf(html, `<iframe src='(.+?)'`)[1]
	if !strings.Contains(videoURL, "tumblr.com/video") {
		log.Fatal("annie doesn't support this URL right now")
	}
	videoHTML := request.Get(videoURL, url)
	realURL := utils.MatchOneOf(videoHTML, `source src="(.+?)"`)[1]
	urlData, size := genURLData(realURL, url)
	format := map[string]downloader.FormatData{
		"default": {
			URLs: []downloader.URLData{urlData},
			Size: size,
		},
	}
	extractedData := downloader.VideoData{
		Site:    "Tumblr tumblr.com",
		Title:   title,
		Type:    "video",
		Formats: format,
	}
	extractedData.Download(url)
	return extractedData
}

// Tumblr download function
func Tumblr(url string) downloader.VideoData {
	html := request.Get(url, url)
	// get the title
	doc := parser.GetDoc(html)
	title := parser.Title(doc)
	if strings.Contains(html, "<iframe src=") {
		// Video
		return tumblrVideoDownload(url, html, title)
	}
	// Image
	return tumblrImageDownload(url, html, title)
}
