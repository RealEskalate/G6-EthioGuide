export interface ChecklistItem {
  checklistID: string;
  done: boolean;
  text: string;
}

export type ChecklistStatus = "NOT_STARTED" | "IN_PROGRESS" | "COMPLETED";

export interface ChecklistResponse {
  checklists: ChecklistItem[];
  id: string; // userProcedureId
  percent: number;
  procedureId: string;
  status: ChecklistStatus;
  updatedAt: string; // ISO
  userId: string;
}
