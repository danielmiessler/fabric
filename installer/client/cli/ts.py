from dotenv import load_dotenv
from pydub import AudioSegment
from openai import OpenAI
from groq import Groq
import os
import argparse
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

        Raises:
            KeyError: If the "OPENAI_API_KEY" is not found in the environment variables.
            FileNotFoundError: If no API key is found in the environment variables.
        """
        env_file = os.path.expanduser(env_file)
        load_dotenv(env_file)
        self.model = getattr(args, 'model', None)
        if not self.model:
            self.model = os.environ.get('DEFAULT_TS_MODEL', None)
            self.model = "whisper-1" if self.model is None else self.model

        self.openai_key = os.environ.get("OPENAI_API_KEY")
        self.groq_key = os.environ.get("GROQ_API_KEY")

        self.openai_client = OpenAI(api_key=self.openai_key) if self.openai_key else None
        self.groq_client = Groq(api_key=self.groq_key) if self.groq_key else None

        openai_models, groq_models = self.fetch_available_models()
        print(f"{openai_models=}\n{groq_models=}")

        self.use_openai = self.model in openai_models
        self.use_groq = self.model in groq_models
        print(f"{ self.use_openai=}\n{self.use_groq=}")

        self.client = self.openai_client if self.use_openai else self.groq_client

        

        # try:
        #     apikey = os.environ["OPENAI_API_KEY"]
        #     self.client = OpenAI()
        #     self.client.api_key = apikey
        # except KeyError:
        #     print("OPENAI_API_KEY not found in environment variables.")

        # except FileNotFoundError:
        #     print("No API key found. Use the --apikey option to set the key")
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
    parser.add_argument(
        "--model", "-m", help="Select the model to use")
    args = parser.parse_args()
    whisper = Whisper(args)
    whisper.process_file(args.audio_file)


if __name__ == "__main__":
    main()
