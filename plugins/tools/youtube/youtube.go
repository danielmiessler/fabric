package youtube

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/anaskhan96/soup"
	"github.com/danielmiessler/fabric/plugins"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func NewYouTube() (ret *YouTube) {

	label := "YouTube"
	ret = &YouTube{}

	ret.PluginBase = &plugins.PluginBase{
		Name:             label,
		SetupDescription: label + " - to grab video transcripts and comments",
		EnvNamePrefix:    plugins.BuildEnvVariablePrefix(label),
	}

	ret.ApiKey = ret.AddSetupQuestion("API key", true)

	return
}

type YouTube struct {
	*plugins.PluginBase
	ApiKey *plugins.SetupQuestion

	service *youtube.Service
}

func (o *YouTube) initService() (err error) {
	if o.service == nil {
		ctx := context.Background()
		o.service, err = youtube.NewService(ctx, option.WithAPIKey(o.ApiKey.Value))
	}
	return
}

func (o *YouTube) GetVideoId(url string) (ret string, err error) {
	if err = o.initService(); err != nil {
		return
	}

	pattern := `(?:https?:\/\/)?(?:www\.)?(?:youtube\.com\/(?:[^\/\n\s]+\/\S+\/|(?:v|e(?:mbed)?)\/|\S*?[?&]v=)|youtu\.be\/)([a-zA-Z0-9_-]{11})`
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(url)
	if len(match) > 1 {
		ret = match[1]
	} else {
		err = fmt.Errorf("invalid YouTube URL, can't get video ID")
	}
	return
}

func (o *YouTube) GrabTranscriptForUrl(url string, language string) (ret string, err error) {
	var videoId string
	if videoId, err = o.GetVideoId(url); err != nil {
		return
	}
	return o.GrabTranscript(videoId, language)
}

func (o *YouTube) GrabTranscript(videoId string, language string) (ret string, err error) {
	var transcript string
	if transcript, err = o.GrabTranscriptBase(videoId, language); err != nil {
		err = fmt.Errorf("transcript not available. (%v)", err)
		return
	}

	// Parse the XML transcript
	doc := soup.HTMLParse(transcript)
	// Extract the text content from the <text> tags
	textTags := doc.FindAll("text")
	var textBuilder strings.Builder
	for _, textTag := range textTags {
		textBuilder.WriteString(textTag.Text())
		textBuilder.WriteString(" ")
		ret = textBuilder.String()
	}
	return
}

func (o *YouTube) GrabTranscriptBase(videoId string, language string) (ret string, err error) {
	if err = o.initService(); err != nil {
		return
	}

	watchUrl := "https://www.youtube.com/watch?v=" + videoId
	var resp string
	if resp, err = soup.Get(watchUrl); err != nil {
		return
	}

	doc := soup.HTMLParse(resp)
	scriptTags := doc.FindAll("script")
	for _, scriptTag := range scriptTags {
		if strings.Contains(scriptTag.Text(), "captionTracks") {
			regex := regexp.MustCompile(`"captionTracks":(\[.*?\])`)
			match := regex.FindStringSubmatch(scriptTag.Text())
			if len(match) > 1 {
				var captionTracks []struct {
					BaseURL string `json:"baseUrl"`
				}

				if err = json.Unmarshal([]byte(match[1]), &captionTracks); err != nil {
					return
				}

				if len(captionTracks) > 0 {
					transcriptURL := captionTracks[0].BaseURL
					for _, captionTrack := range captionTracks {
						parsedUrl, error := url.Parse(captionTrack.BaseURL)
						if error != nil {
							err = fmt.Errorf("error parsing caption track")
						}
						parsedUrlParams, _ := url.ParseQuery(parsedUrl.RawQuery)
						if parsedUrlParams["lang"][0] == language {
							transcriptURL = captionTrack.BaseURL
						}
					}
					ret, err = soup.Get(transcriptURL)
					return
				}
			}
		}
	}
	err = fmt.Errorf("transcript not found")
	return
}

func (o *YouTube) GrabComments(videoId string) (ret []string, err error) {
	if err = o.initService(); err != nil {
		return
	}

	call := o.service.CommentThreads.List([]string{"snippet", "replies"}).VideoId(videoId).TextFormat("plainText").MaxResults(100)
	var response *youtube.CommentThreadListResponse
	if response, err = call.Do(); err != nil {
		log.Printf("Failed to fetch comments: %v", err)
		return
	}

	for _, item := range response.Items {
		topLevelComment := item.Snippet.TopLevelComment.Snippet.TextDisplay
		ret = append(ret, topLevelComment)

		if item.Replies != nil {
			for _, reply := range item.Replies.Comments {
				replyText := reply.Snippet.TextDisplay
				ret = append(ret, "    - "+replyText)
			}
		}
	}
	return
}

func (o *YouTube) GrabDurationForUrl(url string) (ret int, err error) {
	if err = o.initService(); err != nil {
		return
	}

	var videoId string
	if videoId, err = o.GetVideoId(url); err != nil {
		return
	}
	return o.GrabDuration(videoId)
}

func (o *YouTube) GrabDuration(videoId string) (ret int, err error) {
	var videoResponse *youtube.VideoListResponse
	if videoResponse, err = o.service.Videos.List([]string{"contentDetails"}).Id(videoId).Do(); err != nil {
		err = fmt.Errorf("error getting video details: %v", err)
		return
	}

	durationStr := videoResponse.Items[0].ContentDetails.Duration

	matches := regexp.MustCompile(`(?i)PT(?:(\d+)H)?(?:(\d+)M)?(?:(\d+)S)?`).FindStringSubmatch(durationStr)
	if len(matches) == 0 {
		return 0, fmt.Errorf("invalid duration string: %s", durationStr)
	}

	hours, _ := strconv.Atoi(matches[1])
	minutes, _ := strconv.Atoi(matches[2])
	seconds, _ := strconv.Atoi(matches[3])

	ret = hours*60 + minutes + seconds/60

	return
}

func (o *YouTube) Grab(url string, options *Options) (ret *VideoInfo, err error) {
	var videoId string
	if videoId, err = o.GetVideoId(url); err != nil {
		return
	}

	ret = &VideoInfo{}

	if options.Duration {
		if ret.Duration, err = o.GrabDuration(videoId); err != nil {
			err = fmt.Errorf("error parsing video duration: %v", err)
			return
		}

	}

	if options.Comments {
		if ret.Comments, err = o.GrabComments(videoId); err != nil {
			err = fmt.Errorf("error getting comments: %v", err)
			return
		}
	}

	if options.Transcript {
		if ret.Transcript, err = o.GrabTranscript(videoId, "en"); err != nil {
			return
		}
	}
	return
}

type Options struct {
	Duration   bool
	Transcript bool
	Comments   bool
	Lang       string
}

type VideoInfo struct {
	Transcript string   `json:"transcript"`
	Duration   int      `json:"duration"`
	Comments   []string `json:"comments"`
}

func (o *YouTube) GrabByFlags() (ret *VideoInfo, err error) {
	options := &Options{}
	flag.BoolVar(&options.Duration, "duration", false, "Output only the duration")
	flag.BoolVar(&options.Transcript, "transcript", false, "Output only the transcript")
	flag.BoolVar(&options.Comments, "comments", false, "Output the comments on the video")
	flag.StringVar(&options.Lang, "lang", "en", "Language for the transcript (default: English)")
	flag.Parse()

	if flag.NArg() == 0 {
		log.Fatal("Error: No URL provided.")
	}

	url := flag.Arg(0)
	ret, err = o.Grab(url, options)
	return
}
