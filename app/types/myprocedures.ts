export interface UserProcedure {
  percent: number;
  procedureId: string;
  procedureTitle: string;
  status: "NOT_STARTED" | "IN_PROGRESS" | "COMPLETED";
  userProcedureId: string;
  // The following fields are optional for dummy data
  startDate?: string;
  estimatedCompletion?: string;
  documentsUploaded?: string;
  completedDate?: string;
  department?: string; // dummy only
  requirements?: string; // dummy only
  readyToStart?: string; // dummy only
  documentsRequired?: string; // dummy only
  statusColor?: string; // dummy only
  buttonText?: string; // dummy only
  buttonVariant?: "default" | "outline"; // dummy only
}

export interface ProceduresResponse {
  data: UserProcedure[];
  hasNext: boolean;
  limit: number;
  page: number;
  total: number;
}
