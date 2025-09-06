// export interface DocumentRequired {
//   name: string;
//   templateUrl: string | null;
// }

// export interface Fee {
//   amount: number;
//   currency: string;
//   label: string;
// }

// export interface Requirement {
//   text: string;
// }

// export interface Step {
//   order: number;
//   text: string;
// }

// export interface Office {
//   address: string;
//   city: string;
//   hours: string;
// }

// export interface LanguageVersions {
//   amId: string;
//   enId: string;
// }

// export interface ProcessingTime {
//   minDays: number;
//   maxDays: number;
// }

// export default interface ProcedureProp {
//   // title: string;
//   // requirements: Requirement[];
//   // steps: Step[];
//   // fees: Fee[];
//   // processingTime: ProcessingTime;
//   id: string;
//   orgId: string;
//   title: string;
//   summary: string;
//   slug: string;
//   tags: string[];
//   verified: boolean;
//   updatedAt: string;
//   documentsRequired: DocumentRequired[];
//   fees: Fee[];
//   requirements: Requirement[];
//   steps: Step[];
//   offices: Office[];
//   languageVersions: LanguageVersions;
//   processingTime: ProcessingTime;
// }
export default interface ProcedureProp {
  id: string;
  name: string;
  organizationId: string;
  noticeIds: string[];
  content: {prerequisites: string[];
  steps: { [key: number]: string };
  result: string;}
  fees: {
    label: string;
    amount: number;
    currency: string;
  };
  processingTime: {
    minDays: number;
    maxDays: number;
  };
  createdAt: string;
}
