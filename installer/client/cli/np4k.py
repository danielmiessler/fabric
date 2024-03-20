import argparse
from newspaper import Article
import json
import time

class Np4k:
    def __init__(self, file_path=None, single_url=None, output_format='stdout'): 
        self.file_path = file_path
        self.single_url = single_url
        self.output_format = output_format.lower()
        self.articles_data = []
        self.urls = self.load_urls()

    def load_urls(self):
        '''Load URLs from a file or a single URL based on the input provided.'''
        urls = []
        if self.file_path:
            try:
                with open(self.file_path, 'r') as file:
                    urls = [url.strip() for url in file.readlines() if url.strip()]
            except FileNotFoundError:
                print(f'The file {self.file_path} was not found.')
            except Exception as e:
                print(f'Error reading from {self.file_path}: {e}')
        elif self.single_url:
            urls = [self.single_url]
        return urls

    def process_urls(self):
        '''Run newspaper4k against each URL and extract/produce metadata'''
        timestamp = int(time.time())
        output_filename = f'_output_{timestamp}.{"json" if self.output_format == "json" else "txt"}'

        for url in self.urls:
            if url:  # Check if URL is not empty
                try:
                    article_data = self.newspaper4k(url)
                    self.articles_data.append(article_data)
                    # Always print the article text to stdout.
                    print(article_data.get('text', 'No text extracted'))
                except Exception as e:
                    print(f'Error processing URL {url}: {e}')
                    continue

        # Write the extracted data to a file in the specified format if 'json' or 'kvp' is specified
        if self.output_format != 'stdout':
            if self.output_format == 'json':
                self.write_json(output_filename)
            else:  # 'kvp' format
                self.write_kvp(output_filename)


    def format_data(self, article_data, format_type):
        '''Formats the article data based on the specified format for terminal output'''
        if format_type == 'json':
            return json.dumps(article_data, ensure_ascii=False, indent=4)
        elif format_type == 'kvp':
            formatted_data = ""
            for key, value in article_data.items():
                if isinstance(value, list):
                    value = ', '.join(value)
                if isinstance(value, str):
                    value = value.replace('\n', '\\n')
                formatted_data += f"{key}: {value}\n"
            return formatted_data
        elif format_type == 'stdout':  # Only print the article text for stdout
            return article_data.get('text', 'No text extracted')

    def write_json(self, output_filename):
        try:
            with open(output_filename, 'w', encoding='utf-8') as f:
                json.dump(self.articles_data, f, ensure_ascii=False, indent=4)
            print(f'Successfully wrote extracted data to {output_filename}')
        except Exception as e:
            print(f'Error writing data to {output_filename}: {e}')

    def write_kvp(self, output_filename):
        try:
            with open(output_filename, 'w', encoding='utf-8') as f:
                for article in self.articles_data:
                    for key, value in article.items():
                        if isinstance(value, list):
                            value = ', '.join(value)
                        if isinstance(value, str):
                            value = value.replace('\n', '\\n')
                        f.write(f"{key}: {value}\n")
                    f.write("---\n")
                print(f'Successfully wrote extracted data to {output_filename}')
        except Exception as e:
            print(f'Error writing data to {output_filename}: {e}')

    def newspaper4k(self, url):
        article = Article(url, fetch_images=False)
        processed_article = {
            "title": "",
            "keywords": [],
            "tags": [],
            "authors": [],
            "summary": "",
            "text": "",
            "publish_date": "",
            "url": "", 
        }
        try:
            article.download()
            article.parse()
            article.nlp()
            
            processed_article["title"] = article.title or "Not Found"
            processed_article["keywords"] = article.keywords if article.keywords is not None else []
            processed_article["tags"] = list(article.tags) if article.tags is not None else []
            processed_article["authors"] = article.authors if article.authors is not None else ["Not Found"]
            processed_article["summary"] = article.summary or "Not Found"
            processed_article["text"] = article.text or "Not Found"
            processed_article["publish_date"] = article.publish_date.isoformat() if article.publish_date else "Not Found"
            processed_article["url"] = url

        except Exception as e:
            print(f'Failed to process article from {url}: {e}')
            raise e
        return processed_article

def parse_arguments():
    parser = argparse.ArgumentParser(description='Np4k is a helper to extract information from blogs or articles.')
    parser.add_argument('--url', type=str, help='A single URL to process.')
    parser.add_argument('--file', type=str, help='A file containing the list of URLs to process.')
    parser.add_argument('--output', type=str, choices=['stdout', 'kvp', 'json'], default='stdout', help='The file format to write the extracted data in. Default is stdout.')
    return parser.parse_args()

def main():
    args = parse_arguments()
    np4k = Np4k(file_path=args.file, single_url=args.url, output_format=args.output)
    np4k.process_urls()

if __name__ == "__main__":
    main()