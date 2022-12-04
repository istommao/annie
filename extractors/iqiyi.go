package extractors

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/istommao/annie/downloader"
	"github.com/istommao/annie/parser"
	"github.com/istommao/annie/request"
	"github.com/istommao/annie/utils"
)

type vidl struct {
	M3utx      string `json:"m3utx"`
	Vd         int    `json:"vd"` // quality number
	ScreenSize string `json:"screenSize"`
}

type iqiyiData struct {
	Vidl []vidl `json:"vidl"`
}

type iqiyi struct {
	Code string    `json:"code"`
	Data iqiyiData `json:"data"`
}

var iqiyiFormats = []int{
	18, // 1080p
	5,  // 1072p, 1080p
	17, // 720p
	4,  // 720p
	21, // 504p
	2,  // 480p, 504p
	1,  // 336p, 360p
	96, // 216p, 240p
}

const iqiyiReferer = "https://www.iqiyi.com"

func getIqiyiData(tvid, vid string) iqiyi {
	t := time.Now().Unix() * 1000
	src := "76f90cbd92f94a2e925d83e8ccd22cb7"
	key := "d5fb4bd9d50c4be6948c97edd7254b0e"
	sc := utils.Md5(strconv.FormatInt(t, 10) + key + vid)
	info := request.Get(
		fmt.Sprintf(
			"http://cache.m.iqiyi.com/jp/tmts/%s/%s/?t=%d&sc=%s&src=%s",
			tvid, vid, t, sc, src,
		),
		iqiyiReferer,
	)
	var data iqiyi
	json.Unmarshal([]byte(info[len("var tvInfoJs="):]), &data)
	return data
}

// Iqiyi download function
func Iqiyi(url string) downloader.VideoData {
	html := request.Get(url, iqiyiReferer)
	tvid := utils.MatchOneOf(
		url,
		`#curid=(.+)_`,
		`tvid=([^&]+)`,
	)
	if tvid == nil {
		tvid = utils.MatchOneOf(
			html,
			`data-player-tvid="([^"]+)"`,
			`param\['tvid'\]\s*=\s*"(.+?)"`,
		)
	}
	vid := utils.MatchOneOf(
		url,
		`#curid=.+_(.*)$`,
		`vid=([^&]+)`,
	)
	if vid == nil {
		vid = utils.MatchOneOf(
			html,
			`data-player-videoid="([^"]+)"`,
			`param\['vid'\]\s*=\s*"(.+?)"`,
		)
	}
	doc := parser.GetDoc(html)
	title := strings.TrimSpace(doc.Find("h1 a").Text()) +
		strings.TrimSpace(doc.Find("h1 span").Text())
	if title == "" {
		title = doc.Find("title").Text()
	}
	videoDatas := getIqiyiData(tvid[1], vid[1])
	if videoDatas.Code != "A00000" {
		log.Fatal("Can't play this video")
	}
	format := map[string]downloader.FormatData{}
	var urlData downloader.URLData
	var size, totalSize int64
	for _, video := range videoDatas.Data.Vidl {
		urls := []downloader.URLData{}
		totalSize = 0
		for _, ts := range utils.M3u8URLs(video.M3utx) {
			size, _ = strconv.ParseInt(
				utils.MatchOneOf(ts, `contentlength=(\d+)`)[1], 10, 64,
			)
			// http://dx.data.video.qiyi.com -> http://data.video.qiyi.com
			urlData = downloader.URLData{
				URL:  strings.Replace(ts, "dx.data.video.qiyi.com", "data.video.qiyi.com", 1),
				Size: size,
				Ext:  "ts",
			}
			totalSize += size
			urls = append(urls, urlData)
		}
		format[strconv.Itoa(video.Vd)] = downloader.FormatData{
			URLs:    urls,
			Size:    totalSize,
			Quality: video.ScreenSize,
		}
	}
	// get best quality
	var videoData vidl
	for _, quality := range iqiyiFormats {
		for index, video := range videoDatas.Data.Vidl {
			if video.Vd == quality {
				videoData = videoDatas.Data.Vidl[index]
				break
			}
		}
		if videoData.M3utx != "" {
			break
		}
	}
	format["default"] = format[strconv.Itoa(videoData.Vd)]
	delete(format, strconv.Itoa(videoData.Vd))

	extractedData := downloader.VideoData{
		Site:    "爱奇艺 iqiyi.com",
		Title:   title,
		Type:    "video",
		Formats: format,
	}
	extractedData.Download(url)
	return extractedData
}
