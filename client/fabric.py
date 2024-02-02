#!/usr/bin/env python3 

from utils import Remote, Standalone, Server
import argparse
import sys
import os
import shutil

script_directory = os.path.dirname(os.path.realpath(__file__))
config_file = os.path.join(script_directory, 'config.yaml')


if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description='An open source framework for augmenting humans using AI.')
    parser.add_argument('--text', '-t', help='Text to extract summary from')
    parser.add_argument(
        '--copy', '-c', help='Copy the response to the clipboard', action='store_true')
    parser.add_argument('--output', '-o', help='Save the response to a file',
                        nargs='?', const='analyzepaper.txt', default=None)
    parser.add_argument(
        '--stream', '-s', help='Use this option if you are piping output to another app. The output will not be streamed', action='store_true')
    parser.add_argument('--pattern', '-p', help='The pattern (prompt) to use')
    parser.add_argument(
        '--server', '-S', help='Server mode!!!', action='store_true')
    parser.add_argument('--domain', '-d', help='The domain to use for server mode')
    parser.add_argument('--port', '-P', help='The port to use for server mode')
    parser.add_argument('--apikey', '-a', help='Add an OpenAI key')

    args = parser.parse_args()
    home_holder = os.path.expanduser('~')
    config = os.path.join(home_holder, '.config', 'fabric')
    config_patterns_directory = os.path.join(config, 'patterns')
    env_file = os.path.join(config, '.env')
    if not os.path.exists(config):
        os.makedirs(config)
    if not os.path.exists(config_patterns_directory):
        source_patterns_directory = os.path.join(script_directory, 'server/app/chatgpt/patterns')
        shutil.copytree(source_patterns_directory, config_patterns_directory)
    if args.apikey:
        with open(env_file, 'w') as f:
            f.write(f'OPENAI_API_KEY={args.apikey}')
        # print the api key to the console
        print(f'OpenAI API key set to {args.apikey}')
        # quit
        sys.exit()
    if args.server:
        server = Server()
        if not args.domain and not args.port:
            server.run_server(domain='127.0.0.1', port='5000')
        elif args.domain and not args.port:
            server.run_server(domain=args.domain, port='5000')
        elif not args.domain and args.port:
            server.run_server(domain='127.0.0.1', port=args.port)
        else:
            server.run_server(domain=args.domain, port=args.port)
    else:
        if args.domain or args.port:
            parser.error('--domain and --port can only be used with --server')
        if args.text is not None:
            text = args.text
        else:
            text = sys.stdin.read()
        if os.path.exists(config_file):
            analyzer = Remote(args.pattern, args)
            analyzer.analyze(text, copy_to_clipboard=args.copy,
                             save_to_file=args.output)
            analyzer.disconnect_handler()
        else:
            standalone = Standalone(args, args.pattern)
            if args.stream:
                standalone.streamMessage(text)
            else:
                standalone.sendMessage(text)
