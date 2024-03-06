from .utils import Standalone, Update, Setup, Alias
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
        "--copy", "-C", help="Copy the response to the clipboard", action="store_true"
    )
    parser.add_argument(
        '--agents', '-a', choices=['trip_planner', 'ApiKeys'],
        help="Use an AI agent to help you with a task. Acceptable values are 'trip_planner' or 'ApiKeys'. This option cannot be used with any other flag."
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
    parser.add_argument(
        "--update", "-u", help="Update patterns", action="store_true")
    parser.add_argument("--pattern", "-p", help="The pattern (prompt) to use")
    parser.add_argument(
        "--setup", help="Set up your fabric instance", action="store_true"
    )
    parser.add_argument('--changeDefaultModel',
                        help="Change the default model. Your choice will be saved in ~/.config/fabric/.env). For a list of available models, use the --listmodels flag.")
    parser.add_argument(
        '--local', '-L', help="Use local LLM. Default is llama2", action="store_true")

    parser.add_argument(
        "--claude", help="Use Claude AI", action="store_true")

    parser.add_argument(
        "--model", "-m", help="Select the model to use (GPT-4 by default for chatGPT and llama2 for Ollama)", default="gpt-4-turbo-preview"
    )
    parser.add_argument(
        "--listmodels", help="List all available models", action="store_true"
    )
    parser.add_argument('--context', '-c',
                        help="Use Context file (context.md) to add context to your pattern", action="store_true")

    args = parser.parse_args()
    home_holder = os.path.expanduser("~")
    config = os.path.join(home_holder, ".config", "fabric")
    config_patterns_directory = os.path.join(config, "patterns")
    config_context = os.path.join(config, "context.md")
    env_file = os.path.join(config, ".env")
    if not os.path.exists(config):
        os.makedirs(config)
    if args.setup:
        Setup().run()
        Alias()
        sys.exit()
    if not os.path.exists(env_file) or not os.path.exists(config_patterns_directory):
        print("Please run --setup to set up your API key and download patterns.")
        sys.exit()
    if not os.path.exists(config_patterns_directory):
        Update()
        Alias()
        sys.exit()
    if args.changeDefaultModel:
        Setup().default_model(args.changeDefaultModel)
        print(f"Default model changed to {args.changeDefaultModel}")
        sys.exit()
    if args.agents:
        # Handle the agents logic
        if args.agents == 'trip_planner':
            from .agents.trip_planner.main import planner_cli
            tripcrew = planner_cli()
            tripcrew.ask()
            sys.exit()
        elif args.agents == 'ApiKeys':
            from .utils import AgentSetup
            AgentSetup().run()
            sys.exit()
    if args.update:
        Update()
        Alias()
        sys.exit()
    if args.context:
        if not os.path.exists(os.path.join(config, "context.md")):
            print("Please create a context.md file in ~/.config/fabric")
            sys.exit()
    standalone = None
    if args.local:
        standalone = Standalone(args, args.pattern, local=True)
    elif args.claude:
        standalone = Standalone(args, args.pattern, claude=True)
    else:
        standalone = Standalone(args, args.pattern)
    if args.list:
        try:
            direct = sorted(os.listdir(config_patterns_directory))
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
        text = standalone.get_cli_input()
    if args.stream and not args.context:
        standalone.streamMessage(text)
        sys.exit()
    if args.stream and args.context:
        with open(config_context, "r") as f:
            context = f.read()
            standalone.streamMessage(text, context=context)
        sys.exit()
    elif args.context:
        with open(config_context, "r") as f:
            context = f.read()
            standalone.sendMessage(text, context=context)
        sys.exit()
    else:
        standalone.sendMessage(text)
        sys.exit()


if __name__ == "__main__":
    main()
