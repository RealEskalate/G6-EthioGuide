export interface ProcedureProcessingTime { minDays?: number; maxDays?: number }
export interface ProcedureFee { amount: number; currency: string; label?: string }
export interface ProcedureStep {
  order: number
  text: string
  title?: string
  description?: string
  estimatedTime?: string
  time?: string
}
export interface ProcedureRequirement { text: string; optional?: boolean }
export interface ProcedureDocument { name: string; templateUrl?: string | null }
export interface Procedure {
  id: string
  // Some endpoints may return `name` instead of `title`; keep both for compatibility
  name?: string
  title: string
  // Some backend payloads may nest details inside a `content` object; keep loose typing
  content?: unknown
  slug?: string
  summary?: string
  tags?: string[]
  updatedAt?: string
  verified?: boolean
  views?: number
  likes?: number
  processingTime?: ProcedureProcessingTime
  fees?: ProcedureFee[]
  steps?: ProcedureStep[]
  requirements?: ProcedureRequirement[]
  documentsRequired?: ProcedureDocument[]
}
