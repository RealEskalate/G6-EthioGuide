export interface DiscussionPost {
  ID: string;
  UserID: string;
  Title: string;
  Content: string;
  Procedures: string[] | null;
  Tags: string[];
  CreatedAt: string;
  UpdatedAt: string;
}

export interface DiscussionsList {
  posts: DiscussionPost[];
  total: number;
  page: number;
  limit: number;
}
