package youtube

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/istommao/annie/config"
	"github.com/istommao/annie/downloader"
	"github.com/istommao/annie/request"
	"github.com/istommao/annie/utils"
)

type args struct {
	Title  string `json:"title"`
	Stream string `json:"adaptive_fmts"`
	// not every page has `adaptive_fmts` field https://youtu.be/DNaOZovrSVo
	Stream2 string `json:"url_encoded_fmt_stream_map"`
}

type assets struct {
	JS string `json:"js"`
}

type youtubeData struct {
	Args   args   `json:"args"`
	Assets assets `json:"assets"`
}

const referer = "https://www.youtube.com"

var tokensCache = make(map[string][]string)

func getSig(sig, js string) string {
	url := fmt.Sprintf("https://www.youtube.com%s", js)
	tokens, ok := tokensCache[url]
	if !ok {
		tokens = getSigTokens(request.Get(url, referer))
		tokensCache[url] = tokens
	}
	return decipherTokens(tokens, sig)
}

func genSignedURL(streamURL string, stream url.Values, js string) string {
	var realURL, sig string
	if strings.Contains(streamURL, "signature=") {
		// URL itself already has a signature parameter
		realURL = streamURL
	} else {
		// URL has no signature parameter
		sig = stream.Get("sig")
		if sig == "" {
			// Signature need decrypt
			sig = getSig(stream.Get("s"), js)
		}
		realURL = fmt.Sprintf("%s&signature=%s", streamURL, sig)
	}
	return realURL
}

// Download YouTube main download function
func Download(uri string) {
	if !config.Playlist {
		youtubeDownload(uri)
		return
	}
	listID := utils.MatchOneOf(uri, `(list|p)=([^/&]+)`)[2]
	if listID == "" {
		log.Fatal("Can't get list ID from URL")
	}
	html := request.Get("https://www.youtube.com/playlist?list="+listID, referer)
	// "videoId":"OQxX8zgyzuM","thumbnail"
	videoIDs := utils.MatchAll(html, `"videoId":"([^,]+?)","thumbnail"`)
	for _, videoID := range videoIDs {
		u := fmt.Sprintf(
			"https://www.youtube.com/watch?v=%s&list=%s", videoID[1], listID,
		)
		youtubeDownload(u)
	}
}

func youtubeDownload(uri string) downloader.VideoData {
	vid := utils.MatchOneOf(
		uri,
		`watch\?v=([^/&]+)`,
		`youtu\.be/([^?/]+)`,
		`embed/([^/?]+)`,
		`v/([^/?]+)`,
	)
	if vid == nil {
		log.Fatal("Can't find vid")
	}
	videoURL := fmt.Sprintf(
		"https://www.youtube.com/watch?v=%s&gl=US&hl=en&has_verified=1&bpctr=9999999999",
		vid[1],
	)
	html := request.Get(videoURL, referer)
	ytplayer := utils.MatchOneOf(html, `;ytplayer\.config\s*=\s*({.+?});`)[1]
	var youtube youtubeData
	json.Unmarshal([]byte(ytplayer), &youtube)
	title := youtube.Args.Title

	format := extractVideoURLS(youtube, uri)

	extractedData := downloader.VideoData{
		Site:    "YouTube youtube.com",
		Title:   utils.FileName(title),
		Type:    "video",
		Formats: format,
	}
	extractedData.Download(uri)
	return extractedData
}

func extractVideoURLS(data youtubeData, referer string) map[string]downloader.FormatData {
	streams := strings.Split(data.Args.Stream, ",")
	if data.Args.Stream == "" {
		streams = strings.Split(data.Args.Stream2, ",")
	}
	var ext string
	var audio downloader.URLData
	format := map[string]downloader.FormatData{}

	bestQualityURL, _ := url.ParseQuery(streams[0])
	bestQualityItag := bestQualityURL.Get("itag")

	for _, s := range streams {
		stream, _ := url.ParseQuery(s)
		itag := stream.Get("itag")
		streamType := stream.Get("type")
		isAudio := strings.HasPrefix(streamType, "audio/mp4")

		if !isAudio && !utils.ShouldExtract(itag, bestQualityItag) {
			continue
		}

		quality := stream.Get("quality_label")
		if quality == "" {
			quality = stream.Get("quality") // for url_encoded_fmt_stream_map
		}
		if quality != "" {
			quality = fmt.Sprintf("%s %s", quality, streamType)
		} else {
			quality = streamType
		}
		if isAudio {
			// audio file use m4a extension
			ext = "m4a"
		} else {
			ext = utils.MatchOneOf(streamType, `(\w+)/(\w+);`)[2]
		}
		realURL := genSignedURL(stream.Get("url"), stream, data.Assets.JS)
		if !strings.Contains(realURL, "ratebypass=yes") {
			realURL += "&ratebypass=yes"
		}
		size := request.Size(realURL, referer)
		urlData := downloader.URLData{
			URL:  realURL,
			Size: size,
			Ext:  ext,
		}
		if ext == "m4a" {
			// Audio data for merging with video
			audio = urlData
		}
		format[itag] = downloader.FormatData{
			URLs:    []downloader.URLData{urlData},
			Size:    size,
			Quality: quality,
		}
	}

	format["default"] = format[bestQualityItag]
	delete(format, bestQualityItag)

	// `url_encoded_fmt_stream_map`
	if data.Args.Stream == "" {
		return format
	}

	// Unlike `url_encoded_fmt_stream_map`, all videos in `adaptive_fmts` have no sound,
	// we need download video and audio both and then merge them.
	// Another problem is that even if we add `ratebypass=yes`, the download speed still slow sometimes.

	// All videos here have no sound and need to be added separately
	for itag, f := range format {
		if strings.Contains(f.Quality, "video/") {
			f.Size += audio.Size
			f.URLs = append(f.URLs, audio)
			format[itag] = f
		}
	}
	return format
}
