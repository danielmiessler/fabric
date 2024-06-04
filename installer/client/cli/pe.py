import pymupdf
import argparse

def main():
    parser = argparse.ArgumentParser()
    parser = argparse.ArgumentParser(
        description='pdfExtractor extracts all the text from a pdf. By Fabian Wieland.')
    parser.add_argument('filename')
    args = parser.parse_args()
    doc = pymupdf.open(args.filename)  
    text=""
    for page in doc:
        text+=page.get_text()
    print(text)

if __name__ == "__main__":
    main()
