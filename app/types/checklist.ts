export interface ChecklistItem {
  id: string;
  title: string;
  description?: string;
  completed: boolean;
  order?: number;
  updatedAt?: string;
}

export type ChecklistStatus = "NOT_STARTED" | "IN_PROGRESS" | "COMPLETED";

// Client-side richer model
export interface UserProcedureChecklist {
  id: string; // checklist or userProcedureId
  procedureId: string;
  procedureTitle?: string;
  organizationName?: string;
  status?: ChecklistStatus;
  progress?: number; // derived if not provided
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
  items?: Array<Partial<ChecklistItem> & { id: string }>; // patch individual items
  status?: UserProcedureChecklist['status'];
}
