import re
from googleapiclient.discovery import build
from googleapiclient.errors import HttpError
from youtube_transcript_api import YouTubeTranscriptApi
from dotenv import load_dotenv
import os
import json
import isodate
import argparse
import sys


def get_video_id(url):
    # Extract video ID from URL
    pattern = r"(?:https?:\/\/)?(?:www\.)?(?:youtube\.com\/(?:[^\/\n\s]+\/\S+\/|(?:v|e(?:mbed)?)\/|\S*?[?&]v=)|youtu\.be\/)([a-zA-Z0-9_-]{11})"
    match = re.search(pattern, url)
    return match.group(1) if match else None


def get_comments(youtube, video_id):
    comments = []

    try:
        # Fetch top-level comments
        request = youtube.commentThreads().list(
            part="snippet,replies",
            videoId=video_id,
            textFormat="plainText",
            maxResults=100  # Adjust based on needs
        )

        while request:
            response = request.execute()
            for item in response['items']:
                # Top-level comment
                topLevelComment = item['snippet']['topLevelComment']['snippet']['textDisplay']
                comments.append(topLevelComment)
                
                # Check if there are replies in the thread
                if 'replies' in item:
                    for reply in item['replies']['comments']:
                        replyText = reply['snippet']['textDisplay']
                        # Add incremental spacing and a dash for replies
                        comments.append("    - " + replyText)
            
            # Prepare the next page of comments, if available
            if 'nextPageToken' in response:
                request = youtube.commentThreads().list_next(
                    previous_request=request, previous_response=response)
            else:
                request = None

    except HttpError as e:
        print(f"Failed to fetch comments: {e}")

    return comments



def main_function(url, options):
    # Load environment variables from .env file
    load_dotenv(os.path.expanduser("~/.config/fabric/.env"))

    # Get YouTube API key from environment variable
    api_key = os.getenv("YOUTUBE_API_KEY")
    if not api_key:
        print("Error: YOUTUBE_API_KEY not found in ~/.config/fabric/.env")
        return

    # Extract video ID from URL
    video_id = get_video_id(url)
    if not video_id:
        print("Invalid YouTube URL")
        return

    try:
        # Initialize the YouTube API client
        youtube = build("youtube", "v3", developerKey=api_key)

        # Get video details
        video_response = youtube.videos().list(
            id=video_id, part="contentDetails").execute()

        # Extract video duration and convert to minutes
        duration_iso = video_response["items"][0]["contentDetails"]["duration"]
        duration_seconds = isodate.parse_duration(duration_iso).total_seconds()
        duration_minutes = round(duration_seconds / 60)

        # Get video transcript
        try:
            transcript_list = YouTubeTranscriptApi.get_transcript(video_id)
            transcript_text = " ".join([item["text"] for item in transcript_list])
            transcript_text = transcript_text.replace("\n", " ")
        except Exception as e:
            transcript_text = f"Transcript not available. ({e})"

        # Get comments if the flag is set
        comments = []
        if options.comments:
            comments = get_comments(youtube, video_id)

        # Output based on options
        if options.duration:
            print(duration_minutes)
        elif options.transcript:
            print(transcript_text)
        elif options.comments:
            print(json.dumps(comments, indent=2))
        else:
            # Create JSON object with all data
            output = {
                "transcript": transcript_text,
                "duration": duration_minutes,
                "comments": comments
            }
            # Print JSON object
            print(json.dumps(output, indent=2))
    except HttpError as e:
        print(f"Error: Failed to access YouTube API. Please check your YOUTUBE_API_KEY and ensure it is valid: {e}")


def main():
    parser = argparse.ArgumentParser(
        description='yt (video meta) extracts metadata about a video, such as the transcript, the video\'s duration, and now comments. By Daniel Miessler.')
    parser.add_argument('url', help='YouTube video URL')
    parser.add_argument('--duration', action='store_true', help='Output only the duration')
    parser.add_argument('--transcript', action='store_true', help='Output only the transcript')
    parser.add_argument('--comments', action='store_true', help='Output the comments on the video')

    args = parser.parse_args()

    if args.url is None:
        print("Error: No URL provided.")
        return

    main_function(args.url, args)

if __name__ == "__main__":
    main()
