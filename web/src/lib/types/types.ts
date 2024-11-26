export interface Post {
  title: string;
  date: string;
  description: string;
  tags: string[];
}

export interface PostMetadata extends Post {
  slug: string;
}