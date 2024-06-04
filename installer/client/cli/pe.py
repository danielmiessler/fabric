import pymupdf
import argparse
import pymupdf4llm
import pathlib

def main():
    parser = argparse.ArgumentParser()
    parser = argparse.ArgumentParser(
        description='pdfExtractor extracts all the text from a pdf. By Fabian Wieland.')
    parser.add_argument('filename')
    parser.add_argument("-p", "--page", help='Select a specific page to extract')
    parser.add_argument("-m",'--markdown', action='store_true', help='Retrieve your document content in Markdown')
    parser.add_argument("-s",'--save', action='store_true', help='Saves markdownfile (created with -m) utf-8 encoded as output.md')
    args = parser.parse_args()
    if args.filename is None:
        print("Error: No file provided.")
        return
    extract_content(args.filename, args)
   


def extract_content(filename,options):
    doc = pymupdf.open(filename)  
    
    text=""
    if options.page:
        try:
         print(doc.get_page_text(options))
        except:
         print("Error: {} is not a page in the provided document".format(options)) 
    if (options.markdown):
        md_text = pymupdf4llm.to_markdown(filename)
        if(options.save):
            # Write the text to some file in UTF8-encoding
            pathlib.Path("output.md").write_bytes(md_text.encode())
        print(md_text)
    else:
        for page in doc:
            text+=page.get_text()
        print(md_text)
        

if __name__ == "__main__":
    main()
