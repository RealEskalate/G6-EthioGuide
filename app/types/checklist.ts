export interface ChecklistItem {
  id: string;
  content: string;
  is_checked: boolean;
  type: string;
  user_procedure_id: string;
}

export type ChecklistStatus = "NOT_STARTED" | "IN_PROGRESS" | "COMPLETED";

// API Response models based on the provided specification
export interface CreateChecklistResponse {
  id: string;
  percent: number;
  procedure_id: string;
  status: string;
  updated_at: string;
  user_id: string;
}

export interface GetMyProceduresResponse {
  id: string;
  percent: number;
  procedure_id: string;
  status: string;
  updated_at: string;
  user_id: string;
}

export interface GetChecklistItemsResponse {
  content: string;
  id: string;
  is_checked: boolean;
  type: string;
  user_procedure_id: string;
}

// Client-side richer model for UI
export interface UserProcedureChecklist {
  id: string; // checklist or userProcedureId
  procedureId: string;
  procedureTitle?: string;
  organizationName?: string;
  status?: ChecklistStatus;
  progress?: number; // derived from percent
  startedAt?: string;
  completedAt?: string;
  createdAt?: string;
  updatedAt?: string;
  items?: ChecklistItem[];
}

// Backend response summary model
export interface ChecklistResponse {
  checklists: ChecklistItem[];
  id: string; // userProcedureId
  percent: number;
  procedureId: string;
  status: ChecklistStatus;
  updatedAt: string; // ISO
  userId: string;
}

export interface PatchChecklistPayload {
  is_checked?: boolean;
}
