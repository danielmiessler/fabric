from .utils import Standalone, Update, Setup, Transcribe
import argparse
import sys
import os


script_directory = os.path.dirname(os.path.realpath(__file__))

def main():
    parser = argparse.ArgumentParser(
        description="An open source framework for augmenting humans using AI."
    )
    parser.add_argument("--text", "-t", help="Text to extract summary from")
    parser.add_argument(
        "--copy", "-c", help="Copy the response to the clipboard", action="store_true"
    )
    parser.add_argument(
        "--output",
        "-o",
        help="Save the response to a file",
        nargs="?",
        const="analyzepaper.txt",
        default=None,
    )
    parser.add_argument(
        "--stream",
        "-s",
        help="Use this option if you want to see the results in realtime. NOTE: You will not be able to pipe the output into another command.",
        action="store_true",
    )
    parser.add_argument(
        "--list", "-l", help="List available patterns", action="store_true"
    )
    parser.add_argument("--update", "-u", help="Update patterns", action="store_true")
    parser.add_argument("--pattern", "-p", help="The pattern (prompt) to use")
    parser.add_argument(
        "--setup", help="Set up your fabric instance", action="store_true"
    )
    parser.add_argument(
        "--model", "-m", help="Select the model to use (GPT-4 by default)", default="gpt-4-turbo-preview"
    )
    parser.add_argument(
        "--listmodels", help="List all available models", action="store_true"
    )
    parser.add_argument(
        "--youtube", "-y", help="video id for YouTube transcript"
    )

    

    args = parser.parse_args()
    home_holder = os.path.expanduser("~")
    config = os.path.join(home_holder, ".config", "fabric")
    config_patterns_directory = os.path.join(config, "patterns")
    env_file = os.path.join(config, ".env")
    if not os.path.exists(config):
        os.makedirs(config)
    if args.setup:
        Setup().run()
        sys.exit()
    if not os.path.exists(env_file) or not os.path.exists(config_patterns_directory):
        print("Please run --setup to set up your API key and download patterns.")
        sys.exit()
    if not os.path.exists(config_patterns_directory):
        Update()
        sys.exit()
    if args.update:
        Update()
        print("Your Patterns have been updated.")
        sys.exit()
    standalone = Standalone(args, args.pattern)
    if args.list:
        try:
            direct = os.listdir(config_patterns_directory)
            for d in direct:
                print(d)
            sys.exit()
        except FileNotFoundError:
            print("No patterns found")
            sys.exit()
    if args.listmodels:
        standalone.fetch_available_models()
        sys.exit()
    if args.text is not None:
        text = args.text
    else:
       if args.youtube is None:
           text = standalone.get_cli_input()
       else:
           text = Transcribe.youtube(args.youtube)    
    if args.stream:
        standalone.streamMessage(text)
    else:
        standalone.sendMessage(text)

if __name__ == "__main__":
    main()
