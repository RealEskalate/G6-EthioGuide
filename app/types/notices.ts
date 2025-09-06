export interface Notice {
  id: string;
  organization_id: string;
  title: string;
  content: string;
  tags: string[];
  created_at: string;
  updated_at: string;
}

export interface NoticesResponse {
  data: Notice[];
  total: number;
  page: number;
  limit: number;
}
