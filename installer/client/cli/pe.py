import pymupdf
import argparse

def main():
    parser = argparse.ArgumentParser()
    parser = argparse.ArgumentParser(
        description='pdfExtractor extracts all the text from a pdf. By Fabian Wieland.')
    parser.add_argument('filename')
    parser.add_argument("-p",'--page',  help='Select a specific page to extract')
    args = parser.parse_args()
    if args.filename is None:
        print("Error: No file provided.")
        return
    extract_content(args.filename, args.page)
   


def extract_content(filename,option):
    doc = pymupdf.open(filename)  
    text=""
    if option:
        try:
         print(doc.get_page_text(option))
        except:
         print("Error: {} is not a page in the provided document".format(option))
    else:
        for page in doc:
            text+=page.get_text()
        print(text)

if __name__ == "__main__":
    main()
