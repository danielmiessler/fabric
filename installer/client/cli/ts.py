from dotenv import load_dotenv
from pydub import AudioSegment
from openai import OpenAI
from groq import Groq
import os
import argparse
import tempfile
import requests

class Whisper:
    def __init__(self, args, env_file="~/.config/fabric/.env"):
        """        
        Initialize the Whisper class with the provided arguments and environment file.

        Args:
            args: The arguments for initialization.
            env_file: The path to the environment file (default is "~/.config/fabric/.env").

        Returns:
            None
        """
        env_file = os.path.expanduser(env_file)
        load_dotenv(env_file)
        self.model = getattr(args, 'model', None)
        if not self.model:
            self.model = os.environ.get('DEFAULT_TS_MODEL', None)
            self.model = "whisper-1" if self.model is None else self.model

        self.openai_key = os.environ.get("OPENAI_API_KEY")
        self.groq_key = os.environ.get("GROQ_API_KEY")
        if not self.openai_key and not self.groq_key:
            print("Error: No API keys found in environment variables.")
            exit(1)

        try:
            self.openai_client = OpenAI(api_key=self.openai_key) if self.openai_key else None
        except Exception as e:
            print(f"Error instantiating OpenAI client: {e}")
            self.openai_client = None

        try:
            self.groq_client = Groq(api_key=self.groq_key) if self.groq_key else None
        except Exception as e:
            print(f"Error instantiating Groq client: {e}")
            self.groq_client = None

        self.openai_models, self.groq_models = self.fetch_available_models()
        if not self.openai_models and not self.groq_models:
            print("No models available. Check your API keys setup.")
            exit(1)

        use_openai = self.model in self.openai_models
        use_groq = self.model in self.groq_models

        if not use_openai and not use_groq:
            print(f"The selected model '{self.model}' is not available. Use `ts --listmodels` to check the available models")
            exit(1)
        self.client = self.openai_client if use_openai else self.groq_client

        self.whole_response = []

    def fetch_available_models(self):
        openai_list = []
        groq_list = []

        
        if self.openai_client:
            try:
                openai_list = [model.id.strip()
                          for model in self.openai_client.models.list().data 
                          if 'whisper' in model.id]
            except Exception as e:
                print(f"Error fetching OpenAI models: {e}")

        GROQ_API_URL = "https://api.groq.com/openai/v1/models"
        if self.groq_client:
            try:
                headers = {"Authorization": f"Bearer {self.groq_key}"}
                response = requests.get(GROQ_API_URL, headers=headers)
                if response.status_code == 200:
                    groq_list = [model['id'] for model in response.json()['data'] if 'whisper' in model['id']]
            except Exception as e:
                print(f"Error fetching Groq models: {e}")

        return openai_list, groq_list


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
            segment (AudioSegment): The segment audio to be transcribed.

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
            with tempfile.NamedTemporaryFile(delete=True, suffix=".mp3") as f:
                segment.export(f.name, format="mp3")
                with open(f.name, "rb") as audio_file:
                    response = self.client.audio.transcriptions.create(
                        model=self.model,
                        file=audio_file
                    )
                self.whole_response.append(response.text)

        except Exception as e:
            print(f"Error: {e}")

    def process_file(self, audio_file):
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

            segments = self.split_audio(audio_file)
            for i, segment in enumerate(segments):
                self.process_segment(segment)
            print(' '.join(self.whole_response))

        except Exception as e:
            print(f"Error: {e}")
    
    def list_models(self):
        print("OpenAI Models:")
        for model in self.openai_models:
            print(f" - {model}")

        print("\nGroq Models:")
        for model in self.groq_models:
            print(f" - {model}")


def main():
    parser = argparse.ArgumentParser(description="Transcribe an audio file.")
    parser.add_argument(
        "audio_file", nargs="?", help="The path to the audio file to be transcribed.")
    parser.add_argument(
        "--model", "-m", help="Select the model to use")
    parser.add_argument(
        "--listmodels", help="List all available models", action="store_true")
    args = parser.parse_args()

    if args.listmodels:
        whisper = Whisper(args)
        whisper.list_models()
    elif args.audio_file:
        whisper = Whisper(args)
        whisper.process_file(args.audio_file)
    else:
        parser.print_help()


if __name__ == "__main__":
    main()
