from dotenv import load_dotenv
from pydub import AudioSegment
from openai import OpenAI
import os
import argparse
import agentops

# Initialize AgentOps
AGENTOPS_API_KEY = os.getenv("AGENTOPS_API_KEY")
if not AGENTOPS_API_KEY:
    raise ValueError("AGENTOPS_API_KEY not found in environment variables")

agentops.init(AGENTOPS_API_KEY)

@agentops.track_agent(name='Whisper')
class Whisper:
    @agentops.record_function('__init__')
    def __init__(self):
        env_file = os.path.expanduser("~/.config/fabric/.env")
        load_dotenv(env_file)
        try:
            apikey = os.environ["OPENAI_API_KEY"]
            self.client = OpenAI()
            self.client.api_key = apikey
        except KeyError:
            print("OPENAI_API_KEY not found in environment variables.")
        except FileNotFoundError:
            print("No API key found. Use the --apikey option to set the key")
        self.whole_response = []

    @agentops.record_function('split_audio')
    def split_audio(self, file_path):
        audio = AudioSegment.from_file(file_path)
        segments = []
        segment_length_ms = 10 * 60 * 1000  # 10 minutes in milliseconds
        for start_ms in range(0, len(audio), segment_length_ms):
            end_ms = start_ms + segment_length_ms
            segment = audio[start_ms:end_ms]
            segments.append(segment)
        return segments

    @agentops.record_function('process_segment')
    def process_segment(self, segment):
        try:
            audio_file = open(segment, "rb")
            response = self.client.audio.transcriptions.create(
                model="whisper-1",
                file=audio_file
            )
            self.whole_response.append(response.text)
        except Exception as e:
            agentops.log_error(f"Error processing segment: {str(e)}")
            print(f"Error: {e}")

    @agentops.record_function('process_file')
    def process_file(self, audio_file):
        try:
            segments = self.split_audio(audio_file)
            for i, segment in enumerate(segments):
                segment_file_path = f"segment_{i}.mp3"
                segment.export(segment_file_path, format="mp3")
                self.process_segment(segment_file_path)
            print(' '.join(self.whole_response))
        except Exception as e:
            agentops.log_error(f"Error processing file: {str(e)}")
            print(f"Error: {e}")

@agentops.record_function('main')
def main():
    parser = argparse.ArgumentParser(description="Transcribe an audio file.")
    parser.add_argument(
        "audio_file", help="The path to the audio file to be transcribed.")
    args = parser.parse_args()
    whisper = Whisper()
    whisper.process_file(args.audio_file)

if __name__ == "__main__":
    try:
        main()
    except Exception as e:
        agentops.log_error(str(e))
    finally:
        agentops.end_session('Success')