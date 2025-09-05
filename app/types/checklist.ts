export interface ChecklistItem {
  id: string;
  title: string;
  description?: string;
  completed: boolean;
  order?: number;
  updatedAt?: string;
}

export interface UserProcedureChecklist {
  id: string; // checklist or userProcedureId
  procedureId: string;
  procedureTitle?: string;
  organizationName?: string;
  status?: 'NOT_STARTED' | 'IN_PROGRESS' | 'COMPLETED';
  progress?: number; // derived if not provided
  startedAt?: string;
  completedAt?: string;
  createdAt?: string;
  updatedAt?: string;
  items?: ChecklistItem[];
}

export interface PatchChecklistPayload {
  items?: Array<Partial<ChecklistItem> & { id: string }>; // patch individual items
  status?: UserProcedureChecklist['status'];
}
// Checklist related domain types (assumed shape; adjust when backend schema clarifies)
export interface ChecklistItem {
  id: string;
  title: string;
  description?: string;
  completed: boolean;
  // optional ordering
  order?: number;
  // backend may send timestamps
  updatedAt?: string;
}

export interface UserProcedureChecklist {
  id: string;            // userProcedure / checklist id
  procedureId: string;   // underlying procedure id
  procedureTitle?: string;
  organizationName?: string;
  status?: 'NOT_STARTED' | 'IN_PROGRESS' | 'COMPLETED';
  progress?: number;     // 0-100 (client derived if missing)
  startedAt?: string;
  completedAt?: string;
  createdAt?: string;
  updatedAt?: string;
  items?: ChecklistItem[];
}

export interface PatchChecklistPayload {
  items?: Array<Partial<ChecklistItem> & { id: string }>; // minimal patch model
  status?: UserProcedureChecklist['status'];
}
