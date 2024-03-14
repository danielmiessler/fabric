import argparse
import os
import sys
from datetime import datetime

from dotenv import load_dotenv

DEFAULT_CONFIG = "~/.config/fabric/.env"
PATH_KEY = "FABRIC_OUTPUT_PATH"
FM_KEY = "FABRIC_FRONTMATTER_TAGS"
DATE_FORMAT = "%Y-%m-%d"
load_dotenv(os.path.expanduser(DEFAULT_CONFIG))


def main(tag, tags, silent, fabric):
    out = os.getenv(PATH_KEY)
    if out is None:
        print(f"'{PATH_KEY}' not set in {DEFAULT_CONFIG} or in your environment.")
        sys.exit(1)

    out = os.path.expanduser(out)

    if not os.path.isdir(out):
        print(f"'{out}' does not exist. Create it and try again.")
        sys.exit(1)

    if not out.endswith("/"):
        out += "/"

    if len(sys.argv) < 2:
        print(f"'{sys.argv[0]}' takes a single argument to tag your summary")
        sys.exit(1)

    yyyymmdd = datetime.now().strftime(DATE_FORMAT)
    target = f"{out}{yyyymmdd}-{tag}.md"

    # don't clobber existing files- add an incremented number to the end instead
    would_clobber = True
    inc = 0
    while would_clobber:
        if inc > 0:
            target = f"{out}{yyyymmdd}-{tag}-{inc}.md"
        if os.path.exists(target):
            inc += 1
        else:
            would_clobber = False

    # YAML frontmatter stubs for things like Obsidian
    # Prevent a NoneType ending up in the tags
    frontmatter_tags = ""
    if fabric:
        frontmatter_tags = os.getenv(FM_KEY)

    with open(target, "w") as fp:
        if frontmatter_tags or len(tags) != 0:
            fp.write("---\n")
            now = datetime.now().strftime(f"{DATE_FORMAT} %H:%M")
            fp.write(f"generation_date: {now}\n")
            fp.write(f"tags: {frontmatter_tags} {tag} {' '.join(tags)}\n")
            fp.write("---\n")

        # function like 'tee' and split the output to a file and STDOUT
        for line in sys.stdin:
            if not silent:
                print(line, end="")
            fp.write(line)


def cli():
    parser = argparse.ArgumentParser(
        description=(
            'save: a "tee-like" utility to pipeline saving of content, '
            "while keeping the output stream intact. Can optionally generate "
            '"frontmatter" for PKM utilities like Obsidian via the '
            '"FABRIC_FRONTMATTER" environment variable'
        )
    )
    parser.add_argument(
        "stub",
        nargs="?",
        help=(
            "stub to describe your content. Use quotes if you have spaces. "
            "Resulting format is YYYY-MM-DD-stub.md by default"
        ),
    )
    parser.add_argument(
        "-t,",
        "--tag",
        required=False,
        action="append",
        default=[],
        help=(
            "add an additional frontmatter tag. Use this argument multiple times"
            "for  multiple tags"
        ),
    )
    parser.add_argument(
        "-n",
        "--nofabric",
        required=False,
        action="store_false",
        help="don't use the fabric tags, only use tags from --tag",
    )
    parser.add_argument(
        "-s",
        "--silent",
        required=False,
        action="store_true",
        help="don't use STDOUT for output, only save to the file",
    )
    args = parser.parse_args()

    if args.stub:
        main(args.stub, args.tag, args.silent, args.nofabric)
    else:
        parser.print_help()


if __name__ == "__main__":
    cli()
