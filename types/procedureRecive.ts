export default interface ProcedurePropCapital {
  ID: string;
  Name: string;
  Content: {
    Prerequisites: string[];
    Result: string;
    Steps: Record<string, string>;
  };
  Fees: {
    Label: string;
    Amount: number;
    Currency: string;
  };
  ProcessingTime: {
    MinDays: number;
    MaxDays: number;
  };
  CreatedAt?: string;
  GroupID?: string | null;
  OrganizationID?: string;
  NoticeIDs?: string[];
}
