export interface Procedure {
  id?: string;
  _id?: string;
  procedureId?: string;
  uuid?: string;

  title?: string;
  name?: string;
  procedureTitle?: string;

  description?: string;
  summary?: string;

  organization?: string;
  organization_id?: string;

  tags?: string[];
  category?: string;

  created_at?: string;
  updated_at?: string;
  createdAt?: string;
  updatedAt?: string;
}

export interface ProceduresResponse {
  data: Procedure[];
  total: number;
  page: number;
  limit: number;
}
