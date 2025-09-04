interface DocumentRequired {
  name: string;
  templateUrl: string | null;
}

interface Fee {
  amount: number;
  currency: string;
  label: string;
}

interface Requirement {
  text: string;
}

interface Step {
  order: number;
  text: string;
}

interface Office {
  address: string;
  city: string;
  hours: string;
}

interface LanguageVersions {
  amId: string;
  enId: string;
}

interface ProcessingTime {
  minDays: number;
  maxDays: number;
}

export default interface ProcedureProp {
  id: string;
  orgId: string;
  title: string;
  summary: string;
  slug: string;
  tags: string[];
  verified: boolean;
  updatedAt: string;
  documentsRequired: DocumentRequired[];
  fees: Fee[];
  requirements: Requirement[];
  steps: Step[];
  offices: Office[];
  languageVersions: LanguageVersions;
  processingTime: ProcessingTime;
}
