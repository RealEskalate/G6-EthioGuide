export interface Category {
  id: string;
  organization_id?: string;
  parent_id?: string | null;
  title: string;
}

export interface CategoryListResponse {
  data: Category[];
  page: number;
  limit: number;
  total: number;
}