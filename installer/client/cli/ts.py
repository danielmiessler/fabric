from dotenv import load_dotenv
from pydub import AudioSegment
from openai import OpenAI
import os
import argparse
try:
    import whisper as local_whisper
    import torch
    use_local_whisper = True
except ImportError:
    use_local_whisper = False


class Whisper:
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

    def split_audio(self, file_path):
        """
        Splits the audio file into segments of the given length.

        Args:
        - file_path: The path to the audio file.
        - segment_length_ms: Length of each segment in milliseconds.

        Returns:
        - A list of audio segments.
        """
        audio = AudioSegment.from_file(file_path)
        segments = []
        segment_length_ms = 10 * 60 * 1000  # 10 minutes in milliseconds
        for start_ms in range(0, len(audio), segment_length_ms):
            end_ms = start_ms + segment_length_ms
            segment = audio[start_ms:end_ms]
            segments.append(segment)

        return segments

    def process_segment(self, segment):
        """        Transcribe an audio file and print the transcript.

        Args:
            audio_file (str): The path to the audio file to be transcribed.

        Returns:
            None
        """

        try:
            # if audio_file.startswith("http"):
            #     response = requests.get(audio_file)
            #     response.raise_for_status()
            #     with tempfile.NamedTemporaryFile(delete=False) as f:
            #         f.write(response.content)
            #         audio_file = f.name
            audio_file = open(segment, "rb")
            response = self.client.audio.transcriptions.create(
                model="whisper-1",
                file=audio_file
            )
            self.whole_response.append(response.text)

        except Exception as e:
            print(f"Error: {e}")

    def process_file(self, audio_file, language="en", localOpenaiWhisper=False):
        """        Transcribe an audio file and print the transcript.

        Args:
            audio_file (str): The path to the audio file to be transcribed.

        Returns:
            None
        """

        if localOpenaiWhisper:
            try:
                if not use_local_whisper:
                    raise ImportError("Local `openai-whisper` model not found. Please install it using `pipx inject fabric openai-whisper`.")

                if audio_file.endswith(".mp3"):
                    # switch to wav format
                    audio = AudioSegment.from_mp3(audio_file)
                    audio_file = audio_file.replace(".mp3", ".wav")
                    audio.export(audio_file, format="wav")
                model = local_whisper.load_model(
                    "base",
                    device="cuda" if torch.cuda.is_available() else "cpu"
                )
                result = model.transcribe(
                    audio_file,
                    language=language,
                    # set False to display progress bar for debugging
                    # verbose=False
                )
                print(result["text"])

            except Exception as e:
                print(f"Error: {e}")

        else:
            try:
                # if audio_file.startswith("http"):
                #     response = requests.get(audio_file)
                #     response.raise_for_status()
                #     with tempfile.NamedTemporaryFile(delete=False) as f:
                #         f.write(response.content)
                #         audio_file = f.name

                segments = self.split_audio(audio_file)
                for i, segment in enumerate(segments):
                    segment_file_path = f"segment_{i}.mp3"
                    segment.export(segment_file_path, format="mp3")
                    self.process_segment(segment_file_path)
                print(' '.join(self.whole_response))

            except Exception as e:
                print(f"Error: {e}")


def main():
    parser = argparse.ArgumentParser(description="Transcribe an audio file.")
    parser.add_argument(
        "audio_file", help="The path to the audio file to be transcribed.")
    parser.add_argument('--lang', default='en', help='Language for the transcript (default: English). This is only used for the local model. (default: en)')
    parser.add_argument('--localOpenaiWhisper', action='store_true', help='Use Local `openai-whisper` model for transcription. (default: False)')
    args = parser.parse_args()
    whisper = Whisper()
    whisper.process_file(args.audio_file, args.lang, args.localOpenaiWhisper)


if __name__ == "__main__":
    main()
