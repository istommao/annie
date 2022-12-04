package extractors

import (
	"fmt"

	"github.com/istommao/annie/downloader"
	"github.com/istommao/annie/request"
	"github.com/istommao/annie/utils"
)

// Universal download function
func Universal(url string) downloader.VideoData {
	fmt.Println()
	fmt.Println("annie doesn't support this URL right now, but it will try to download it directly")

	filename, ext := utils.GetNameAndExt(url)
	size := request.Size(url, url)
	urlData := downloader.URLData{
		URL:  url,
		Size: size,
		Ext:  ext,
	}
	format := map[string]downloader.FormatData{
		"default": {
			URLs: []downloader.URLData{urlData},
			Size: size,
		},
	}
	extractedData := downloader.VideoData{
		Site:    "Universal",
		Title:   utils.FileName(filename),
		Type:    request.ContentType(url, url),
		Formats: format,
	}
	extractedData.Download(url)
	return extractedData
}
