import re
from googleapiclient.discovery import build
from googleapiclient.errors import HttpError
from youtube_transcript_api import YouTubeTranscriptApi
from dotenv import load_dotenv
import os
import json
import isodate
import argparse


def get_video_id(url):
    # Extract video ID from URL
    pattern = r"(?:https?:\/\/)?(?:www\.)?(?:youtube\.com\/(?:[^\/\n\s]+\/\S+\/|(?:v|e(?:mbed)?)\/|\S*?[?&]v=)|youtu\.be\/)([a-zA-Z0-9_-]{11})"
    match = re.search(pattern, url)
    return match.group(1) if match else None


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
        video_response = (
            youtube.videos().list(id=video_id, part="contentDetails").execute()
        )

        # Extract video duration and convert to minutes
        duration_iso = video_response["items"][0]["contentDetails"]["duration"]
        duration_seconds = isodate.parse_duration(duration_iso).total_seconds()
        duration_minutes = round(duration_seconds / 60)

        # Get video transcript
        try:
            transcript_list = YouTubeTranscriptApi.get_transcript(video_id)
            transcript_text = " ".join([item["text"]
                                       for item in transcript_list])
            transcript_text = transcript_text.replace("\n", " ")
        except Exception as e:
            transcript_text = f"Transcript not available. ({e})"

        # Output based on options
        if options.duration:
            print(duration_minutes)
        elif options.transcript:
            print(transcript_text)
        else:
            # Create JSON object
            output = {"transcript": transcript_text,
                      "duration": duration_minutes}
            # Print JSON object
            print(json.dumps(output))
    except HttpError as e:

        print(
            f"Error: Failed to access YouTube API. Please check your YOUTUBE_API_KEY and ensure it is valid: {e}")


def main():
    parser = argparse.ArgumentParser(

        description='yt (video meta) extracts metadata about a video, such as the transcript and the video\'s duration. By Daniel Miessler.')
    parser.add_argument('url', nargs='?', help='YouTube video URL')
    parser.add_argument('--duration', action='store_true',
                        help='Output only the duration')
    parser.add_argument('--transcript', action='store_true',
                        help='Output only the transcript')
    parser.add_argument("url", nargs="?", help="YouTube video URL")
