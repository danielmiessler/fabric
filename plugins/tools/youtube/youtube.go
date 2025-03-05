package youtube

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

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

	normalizeRegex *regexp.Regexp
	service        *youtube.Service
}

func (o *YouTube) initService() (err error) {
	if o.service == nil {
		o.normalizeRegex = regexp.MustCompile(`[^a-zA-Z0-9]+`)
		ctx := context.Background()
		o.service, err = youtube.NewService(ctx, option.WithAPIKey(o.ApiKey.Value))
	}
	return
}

func (o *YouTube) GetVideoOrPlaylistId(url string) (videoId string, playlistId string, err error) {
	if err = o.initService(); err != nil {
		return
	}

	// Video ID pattern
	videoPattern := `(?:https?:\/\/)?(?:www\.)?(?:youtube\.com\/(?:live\/|[^\/\n\s]+\/\S+\/|(?:v|e(?:mbed)?)\/|(?:s(?:horts)\/)|\S*?[?&]v=)|youtu\.be\/)([a-zA-Z0-9_-]*)`
	videoRe := regexp.MustCompile(videoPattern)
	videoMatch := videoRe.FindStringSubmatch(url)
	if len(videoMatch) > 1 {
		videoId = videoMatch[1]
	}

	// Playlist ID pattern
	playlistPattern := `[?&]list=([a-zA-Z0-9_-]+)`
	playlistRe := regexp.MustCompile(playlistPattern)
	playlistMatch := playlistRe.FindStringSubmatch(url)
	if len(playlistMatch) > 1 {
		playlistId = playlistMatch[1]
	}

	if videoId == "" && playlistId == "" {
		err = fmt.Errorf("invalid YouTube URL, can't get video or playlist ID: '%s'", url)
	}
	return
}

func (o *YouTube) GrabTranscriptForUrl(url string, language string) (ret string, err error) {
	var videoId string
	var playlistId string
	if videoId, playlistId, err = o.GetVideoOrPlaylistId(url); err != nil {
		return
	} else if videoId == "" && playlistId != "" {
		err = fmt.Errorf("URL is a playlist, not a video")
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
		textBuilder.WriteString(strings.ReplaceAll(textTag.Text(), "&#39;", "'"))
		textBuilder.WriteString(" ")
		ret = textBuilder.String()
	}
	return
}

func (o *YouTube) GrabTranscriptWithTimestamps(videoId string, language string) (ret string, err error) {
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
		// Extract the start and duration attributes
		start := textTag.Attrs()["start"]
		dur := textTag.Attrs()["dur"]
		end := fmt.Sprintf("%f", parseFloat(start)+parseFloat(dur))
		// Format the timestamps
		startFormatted := formatTimestamp(parseFloat(start))
		endFormatted := formatTimestamp(parseFloat(end))
		text := strings.ReplaceAll(textTag.Text(), "&#39;", "'")
		textBuilder.WriteString(fmt.Sprintf("[%s - %s] %s\n", startFormatted, endFormatted, text))
	}
	ret = textBuilder.String()
	return
}

func parseFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func formatTimestamp(seconds float64) string {
	hours := int(seconds) / 3600
	minutes := (int(seconds) % 3600) / 60
	secs := int(seconds) % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, secs)
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
	var playlistId string
	if videoId, playlistId, err = o.GetVideoOrPlaylistId(url); err != nil {
		return
	} else if videoId == "" && playlistId != "" {
		err = fmt.Errorf("URL is a playlist, not a video")
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
	var playlistId string
	if videoId, playlistId, err = o.GetVideoOrPlaylistId(url); err != nil {
		return
	} else if videoId == "" && playlistId != "" {
		err = fmt.Errorf("URL is a playlist, not a video")
		return
	}

	ret = &VideoInfo{}

	if options.Metadata {
		if ret.Metadata, err = o.GrabMetadata(videoId); err != nil {
			err = fmt.Errorf("error getting video metadata: %v", err)
			return
		}
	}

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

	if options.TranscriptWithTimestamps {
		if ret.Transcript, err = o.GrabTranscriptWithTimestamps(videoId, "en"); err != nil {
			return
		}
	}

	return
}

// FetchPlaylistVideos fetches all videos from a YouTube playlist.
func (o *YouTube) FetchPlaylistVideos(playlistID string) (ret []*VideoMeta, err error) {
	if err = o.initService(); err != nil {
		return
	}

	nextPageToken := ""
	for {
		call := o.service.PlaylistItems.List([]string{"snippet"}).PlaylistId(playlistID).MaxResults(50)
		if nextPageToken != "" {
			call = call.PageToken(nextPageToken)
		}

		var response *youtube.PlaylistItemListResponse
		if response, err = call.Do(); err != nil {
			return
		}

		for _, item := range response.Items {
			videoID := item.Snippet.ResourceId.VideoId
			title := item.Snippet.Title
			ret = append(ret, &VideoMeta{videoID, title, o.normalizeFileName(title)})
		}

		nextPageToken = response.NextPageToken
		if nextPageToken == "" {
			break
		}

		time.Sleep(1 * time.Second) // Pause to respect API rate limit
	}
	return
}

// SaveVideosToCSV saves the list of videos to a CSV file.
func (o *YouTube) SaveVideosToCSV(filename string, videos []*VideoMeta) (err error) {
	var file *os.File
	if file, err = os.Create(filename); err != nil {
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write headers
	if err = writer.Write([]string{"VideoID", "Title"}); err != nil {
		return
	}

	// Write video data
	for _, record := range videos {
		if err = writer.Write([]string{record.Id, record.Title}); err != nil {
			return
		}
	}

	return
}

// FetchAndSavePlaylist fetches all videos in a playlist and saves them to a CSV file.
func (o *YouTube) FetchAndSavePlaylist(playlistID, filename string) (err error) {
	var videos []*VideoMeta
	if videos, err = o.FetchPlaylistVideos(playlistID); err != nil {
		err = fmt.Errorf("error fetching playlist videos: %v", err)
		return
	}

	if err = o.SaveVideosToCSV(filename, videos); err != nil {
		err = fmt.Errorf("error saving videos to CSV: %v", err)
		return
	}

	fmt.Println("Playlist saved to", filename)
	return
}

func (o *YouTube) FetchAndPrintPlaylist(playlistID string) (err error) {
	var videos []*VideoMeta
	if videos, err = o.FetchPlaylistVideos(playlistID); err != nil {
		err = fmt.Errorf("error fetching playlist videos: %v", err)
		return
	}

	fmt.Printf("Playlist: %s\n", playlistID)
	fmt.Printf("VideoId: Title\n")
	for _, video := range videos {
		fmt.Printf("%s: %s\n", video.Id, video.Title)
	}
	return
}

func (o *YouTube) normalizeFileName(name string) string {
	return o.normalizeRegex.ReplaceAllString(name, "_")

}

type VideoMeta struct {
	Id              string
	Title           string
	TitleNormalized string
}

type Options struct {
	Duration                 bool
	Transcript               bool
	TranscriptWithTimestamps bool
	Comments                 bool
	Lang                     string
	Metadata                 bool
}

type VideoInfo struct {
	Transcript string         `json:"transcript"`
	Duration   int            `json:"duration"`
	Comments   []string       `json:"comments"`
	Metadata   *VideoMetadata `json:"metadata,omitempty"`
}

type VideoMetadata struct {
	Id           string   `json:"id"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	PublishedAt  string   `json:"publishedAt"`
	ChannelId    string   `json:"channelId"`
	ChannelTitle string   `json:"channelTitle"`
	CategoryId   string   `json:"categoryId"`
	Tags         []string `json:"tags"`
	ViewCount    uint64   `json:"viewCount"`
	LikeCount    uint64   `json:"likeCount"`
}

func (o *YouTube) GrabMetadata(videoId string) (metadata *VideoMetadata, err error) {
	if err = o.initService(); err != nil {
		return
	}

	call := o.service.Videos.List([]string{"snippet", "statistics"}).Id(videoId)
	var response *youtube.VideoListResponse
	if response, err = call.Do(); err != nil {
		return nil, fmt.Errorf("error getting video metadata: %v", err)
	}

	if len(response.Items) == 0 {
		return nil, fmt.Errorf("no video found with ID: %s", videoId)
	}

	video := response.Items[0]
	viewCount := video.Statistics.ViewCount
	likeCount := video.Statistics.LikeCount

	metadata = &VideoMetadata{
		Id:           video.Id,
		Title:        video.Snippet.Title,
		Description:  video.Snippet.Description,
		PublishedAt:  video.Snippet.PublishedAt,
		ChannelId:    video.Snippet.ChannelId,
		ChannelTitle: video.Snippet.ChannelTitle,
		CategoryId:   video.Snippet.CategoryId,
		Tags:         video.Snippet.Tags,
		ViewCount:    viewCount,
		LikeCount:    likeCount,
	}
	return
}

func (o *YouTube) GrabByFlags() (ret *VideoInfo, err error) {
	options := &Options{}
	flag.BoolVar(&options.Duration, "duration", false, "Output only the duration")
	flag.BoolVar(&options.Transcript, "transcript", false, "Output only the transcript")
	flag.BoolVar(&options.TranscriptWithTimestamps, "transcriptWithTimestamps", false, "Output only the transcript with timestamps")
	flag.BoolVar(&options.Comments, "comments", false, "Output the comments on the video")
	flag.StringVar(&options.Lang, "lang", "en", "Language for the transcript (default: English)")
	flag.BoolVar(&options.Metadata, "metadata", false, "Output video metadata")
	flag.Parse()

	if flag.NArg() == 0 {
		log.Fatal("Error: No URL provided.")
	}

	url := flag.Arg(0)
	ret, err = o.Grab(url, options)
	return
}
