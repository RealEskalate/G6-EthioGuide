export interface Procedure {
  id: string;
  name: string;
}

export default interface Notice {
  id: string;                    // maps to id
  organization_id: string;       // maps to organization_id
  title: string;                 // maps to title
  content: string;               // maps to content
  tags: string[];                // maps to tags array
  created_at: string;            // maps to created_at
  updated_at: string;            // maps to updated_at
}
