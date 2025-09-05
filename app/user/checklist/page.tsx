"use client"

import { useState, useEffect } from "react"
import {
  Download,
  Upload,
  FileText,
  AlertTriangle,
  Info,
  Calendar,
  CreditCard,
  CheckCircle, // added
  Circle, // added
} from "lucide-react"
import { Button } from "@/components/ui/button"
import { Checkbox } from "@/components/ui/checkbox"
import { Progress } from "@/components/ui/progress"
import { useSearchParams } from "next/navigation"
import { useGetChecklistQuery } from "@/app/store/slices/checklistSlice"
import { useGetMyProceduresQuery } from "@/app/store/slices/workspaceSlice"
import React from "react"
import type { ChecklistItem } from "@/app/types/checklist" // added

interface ApplicationState {
  step1: {
    completed: boolean
    formFilled: boolean
  }
  step2: {
    completed: boolean
    idUploaded: boolean
  }
  step3: {
    completed: boolean
    bankStatements: boolean
    proofOfFinancing: boolean
    businessPlan: boolean
  }
  step4: {
    completed: boolean
    inspectionScheduled: boolean
  }
  step5: {
    completed: boolean
    feesPaid: boolean
  }
}

export default function CityGovPortal() {
  const [applicationState, setApplicationState] = useState<ApplicationState>({
    step1: { completed: true, formFilled: true },
    step2: { completed: true, idUploaded: true },
    step3: { completed: false, bankStatements: false, proofOfFinancing: false, businessPlan: false },
    step4: { completed: false, inspectionScheduled: false },
    step5: { completed: false, feesPaid: false },
  })

  const [lastSaved, setLastSaved] = useState<Date | null>(null)

  useEffect(() => {
    const savedState = localStorage.getItem("business-license-progress")
    const savedDate = localStorage.getItem("business-license-last-saved")

    if (savedState) {
      try {
        setApplicationState(JSON.parse(savedState))
      } catch (error) {
        console.log("[v0] Error loading saved state:", error)
      }
    }

    if (savedDate) {
      setLastSaved(new Date(savedDate))
    }
  }, [])

  const calculateProgress = () => {
    const requirements = [
      applicationState.step1.formFilled,
      applicationState.step2.idUploaded, // Now dynamic based on state
      applicationState.step3.bankStatements,
      applicationState.step4.inspectionScheduled,
      applicationState.step5.feesPaid,
    ]

    const completedRequirements = requirements.filter(Boolean).length
    return {
      completed: completedRequirements,
      total: requirements.length,
      percentage: (completedRequirements / requirements.length) * 100,
    }
  }

  const progress = calculateProgress()

  const updateStepCompletion = (
    step: keyof ApplicationState,
    updates: Partial<ApplicationState[keyof ApplicationState]>,
  ) => {
    setApplicationState((prev) => {
      const newState = {
        ...prev,
        [step]: { ...prev[step], ...updates }, // Fixed to properly merge existing state with updates
      }

      if (step === "step1") {
        newState.step1.completed = newState.step1.formFilled
      } else if (step === "step2") {
        newState.step2.completed = newState.step2.idUploaded
      } else if (step === "step3") {
        newState.step3.completed = newState.step3.bankStatements
      } else if (step === "step4") {
        newState.step4.completed = newState.step4.inspectionScheduled
      } else if (step === "step5") {
        newState.step5.completed = newState.step5.feesPaid
      }

      return newState
    })
  }

  const saveProgress = () => {
    try {
      localStorage.setItem("business-license-progress", JSON.stringify(applicationState))
      localStorage.setItem("business-license-last-saved", new Date().toISOString())
      setLastSaved(new Date())
      alert("Progress saved successfully!")
    } catch (error) {
      console.log("[v0] Error saving progress:", error)
      alert("Error saving progress. Please try again.")
    }
  }

  const handleFileUpload = (stepKey: keyof ApplicationState, requirementKey: string) => {
    // Simulate file upload
    const input = document.createElement("input")
    input.type = "file"
    input.accept = ".pdf,.jpg,.jpeg,.png"
    input.onchange = (e) => {
      const file = (e.target as HTMLInputElement).files?.[0]
      if (file) {
        updateStepCompletion(stepKey, { [requirementKey]: true } as Record<string, boolean>)
        alert(`${file.name} uploaded successfully!`)
      }
    }
    input.click()
  }

  const scheduleInspection = () => {
    if (!applicationState.step3.completed) {
      alert("Please complete Step 3 before scheduling inspection.")
      return
    }

    // Simulate scheduling
    const confirmed = confirm("Schedule inspection for next available slot (January 15, 2024 at 10:00 AM)?")
    if (confirmed) {
      updateStepCompletion("step4", { inspectionScheduled: true })
      alert("Inspection scheduled successfully!")
    }
  }

  const processPayment = () => {
    if (!applicationState.step4.completed) {
      alert("Please complete Step 4 before making payment.")
      return
    }

    // Simulate payment
    const confirmed = confirm("Process payment of $150 for business license?")
    if (confirmed) {
      updateStepCompletion("step5", { feesPaid: true })
      alert("Payment processed successfully! Your license will be issued within 3-5 business days.")
    }
  }

  // Returns the next incomplete step number (1-based), or null if all complete
  const getNextStep = () => {
    if (!applicationState.step1.completed) return 1
    if (!applicationState.step2.completed) return 2
    if (!applicationState.step3.completed) return 3
    if (!applicationState.step4.completed) return 4
    if (!applicationState.step5.completed) return 5
    return null
  }

  const nextStep = getNextStep()

  // For scrolling to next step
  const handleContinue = () => {
    if (!nextStep) {
      alert("All steps are complete!")
      return
    }
    const el = document.getElementById(`step-${nextStep}`)
    if (el) {
      el.scrollIntoView({ behavior: "smooth", block: "center" })
    }
  }

  const searchParams = useSearchParams()
  const userProcedureId = searchParams.get("userProcedureId") || ""
  // read title passed from workspace (preferred)
  const paramProcedureTitle = searchParams.get("procedureTitle") || ""

  const {
    data: checklistData,
    error: checklistError,
    isLoading: checklistLoading,
  } = useGetChecklistQuery(userProcedureId, { skip: !userProcedureId })

  // fetch procedures to resolve the title as fallback
  // removed local ProceduresQueryParams interface; the hook already accepts { page?: number; limit?: number }
  const { data: proceduresResp } = useGetMyProceduresQuery({ page: 1, limit: 100 })

  // Define a type for procedure items
  interface ProcedureItem {
    userProcedureId?: string
    procedureId?: string
    procedureTitle?: string
    // add other fields as needed
  }

  // compute procedure title using userProcedureId or procedureId from checklist (fallback if param not present)
  const procedureTitle = React.useMemo(() => {
    const items: ProcedureItem[] = proceduresResp?.data ?? []
    const byUser = items.find((p) => p?.userProcedureId === userProcedureId)
    if (byUser?.procedureTitle) return byUser.procedureTitle
    const pid = checklistData?.procedureId
    if (pid) {
      const byProc = items.find((p) => p?.procedureId === pid)
      if (byProc?.procedureTitle) return byProc.procedureTitle
    }
    return null
  }, [proceduresResp, userProcedureId, checklistData])

  // final header title preference: param -> derived -> default
  const headerTitle = paramProcedureTitle || procedureTitle || "Business License Application"

  // add: AI checklist state populated from API
  const [aiChecklist, setAiChecklist] = useState<ChecklistItem[]>([])
  useEffect(() => {
    if (Array.isArray(checklistData?.checklists)) {
      setAiChecklist(checklistData.checklists)
    }
  }, [checklistData])

  return (
    <div className="min-h-screen bg-gray-50">
     
      <style jsx>{`
        .step-transition {
          transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
        }
        .checkbox-transition {
          transition: all 0.3s ease-in-out;
          transform-origin: center;
        }
        .checkbox-transition:hover {
          transform: scale(1.1);
        }
        .progress-transition {
          transition: all 0.6s cubic-bezier(0.4, 0, 0.2, 1);
        }
        .step-circle-transition {
          transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
          transform-origin: center;
        }
        .step-circle-transition:hover {
          transform: scale(1.1) rotate(5deg);
          box-shadow: 0 4px 20px rgba(94, 156, 141, 0.3);
        }
        .button-transition {
          transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
          transform-origin: center;
        }
        .button-transition:hover {
          transform: translateY(-2px) scale(1.02);
          box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
        }
        .button-transition:active {
          transform: translateY(0) scale(0.98);
        }
        .card-hover {
          transition: all 0.3s ease-in-out;
        }
        .card-hover:hover {
          transform: translateY(-2px);
          box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
        }
        .bounce-in {
          animation: bounceIn 0.6s cubic-bezier(0.68, -0.55, 0.265, 1.55);
        }
        @keyframes bounceIn {
          0% {
            transform: scale(0.3);
            opacity: 0;
          }
          50% {
            transform: scale(1.05);
          }
          70% {
            transform: scale(0.9);
          }
          100% {
            transform: scale(1);
            opacity: 1;
          }
        }
        .slide-in {
          animation: slideIn 0.5s ease-out;
        }
        @keyframes slideIn {
          from {
            transform: translateX(-20px);
            opacity: 0;
          }
          to {
            transform: translateX(0);
            opacity: 1;
          }
        }
        .pulse-success {
          animation: pulseSuccess 0.6s ease-in-out;
        }
        @keyframes pulseSuccess {
          0% {
            box-shadow: 0 0 0 0 rgba(34, 197, 94, 0.7);
          }
          70% {
            box-shadow: 0 0 0 10px rgba(34, 197, 94, 0);
          }
          100% {
            box-shadow: 0 0 0 0 rgba(34, 197, 94, 0);
          }
        }
      `}</style>
      <div className="max-w-7xl mx-auto px-4 sm:px-6 py-6 sm:py-8">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 lg:gap-8">
          <div className="lg:col-span-2 space-y-6">
            {/* Header Section */}
            <div className="bg-white rounded-lg p-6 mb-6 card-hover">
              <h1 className="text-2xl font-semibold text-[#2e4d57] mb-3">
                {headerTitle}
              </h1>
              <p className="text-[#a7b3b9] mb-6">
                Complete this step-by-step process to obtain your business license. This checklist will guide you
                through all required documentation and procedures.
              </p>

              <div className="flex items-center gap-6 text-sm">
                <div className="flex items-center gap-2">
                  <div className="w-2 h-2 bg-[#22c55e] rounded-full"></div>
                  <span className="text-[#a7b3b9]">Est. 2-3 weeks</span>
                </div>
                <div className="flex items-center gap-2">
                  <div className="w-2 h-2 bg-[#0075ff] rounded-full"></div>
                  <span className="text-[#a7b3b9]">Last saved: {lastSaved?.toLocaleDateString() || "Not saved yet"}</span>
                </div>
              </div>
            </div>

            {/* Progress Section */}
            <div className="bg-white rounded-lg p-6 mb-6 card-hover">
              <div className="flex items-center justify-between mb-4">
                <h2 className="text-lg font-medium text-[#2e4d57]">Your Progress</h2>
                {/* Prefer API percent, fallback to local progress */}
                {(() => {
                  const apiPercent = typeof checklistData?.percent === "number" ? checklistData.percent : null
                  const progressPercent = apiPercent ?? progress.percentage
                  return (
                    <span className="text-sm text-[#a7b3b9]">{`${Math.round(progressPercent)}% complete`}</span>
                  )
                })()}
              </div>

              {/* Progress bar uses API percent when available */}
              {(() => {
                const apiPercent = typeof checklistData?.percent === "number" ? checklistData.percent : null
                const progressPercent = apiPercent ?? progress.percentage
                return (
                  <div className="mb-4">
                    <Progress value={progressPercent} className="h-2 progress-transition" />
                  </div>
                )
              })()}

              <div className="flex justify-between text-xs text-[#a7b3b9]">
                <span>Started</span>
                <span>In Progress</span>
                <span>Complete</span>
              </div>
            </div>

            {/* Procedure Checklist */}
            <div className="bg-white rounded-lg p-6 card-hover">
              <h2 className="text-lg font-medium text-[#2e4d57] mb-6">Procedure Checklist</h2>

              {aiChecklist.length > 0 ? (
                <div className="space-y-4">
                  {aiChecklist.map((item) => (
                    <div key={item.checklistID} className="flex items-start gap-3">
                      <Checkbox
                        checked={item.done}
                        onCheckedChange={(checked) =>
                          setAiChecklist((prev) =>
                            prev.map((x) =>
                              x.checklistID === item.checklistID ? { ...x, done: Boolean(checked) } : x
                            )
                          )
                        }
                        className="h-5 w-5 data-[state=checked]:bg-[#3A6A8D] mt-0.5 checkbox-transition"
                      />
                      <span className={`text-base md:text-lg ${item.done ? "line-through text-[#a7b3b9]" : "text-[#2e4d57]"}`}>
                        {item.text}
                      </span>
                    </div>
                  ))}
                </div>
              ) : (
                <>
                  {/* Steps commented out as requested. The AI checklist will render here when available. */}
                  {/*
                  <div className="space-y-6">
                    // ...existing code...
                    // Step 1 - Now Interactive
                    // Step 2 - Now Fully Interactive
                    // Step 3 - In Progress
                    // Step 4 - Schedule Inspection
                    // Step 5 - Pay License Fees
                  </div>
                  */}
                </>
              )}
            </div>

            {/* Ready to Continue */}
            <div className="bg-white rounded-lg p-6 mt-6 card-hover">
              <h3 className="font-medium text-[#2e4d57] mb-2">
                {progress.completed === progress.total ? "Application Complete!" : "Ready to continue?"}
              </h3>
              <p className="text-sm text-[#a7b3b9] mb-4">
                {progress.completed === progress.total
                  ? "Your business license application is complete. You will receive your license within 3-5 business days."
                  : nextStep
                    ? `Complete Step ${nextStep} to move forward with your application.`
                    : "Complete the next step to move forward with your application."}
              </p>

              <div className="flex gap-3">
                <Button
                  variant="outline"
                  className="border-[#e5e7eb] text-[#2e4d57] hover:!bg-[#eff0f1] bg-transparent button-transition"
                  onClick={saveProgress}
                >
                  <FileText className="w-4 h-4 mr-2" />
                  Save Progress
                </Button>
                <Button
                  className={`bg-[#3A6A8D] hover:bg-[#2f5c81] text-white button-transition ${!nextStep ? "opacity-60 cursor-not-allowed" : ""}`}
                  onClick={handleContinue}
                  disabled={!nextStep}
                >
                  {nextStep ? `Continue Step ${nextStep}` : "All steps complete"}
                </Button>
              </div>
            </div>
          </div>

          {/* Sidebar */}
          <div className="space-y-6 order-first lg:order-none">
            {/* Required Documents */}
            <div className="bg-white rounded-lg p-5 sm:p-6 card-hover">
              <h3 className="font-medium text-[#2e4d57] mb-4">Required Documents</h3>

              {/* changed: do not render checklistData.checklists here anymore; keep static fallback */}
              {checklistLoading && (
                <div className="text-sm text-[#a7b3b9]">Loading documents...</div>
              )}

              {!checklistLoading && checklistError && (
                <div className="text-sm text-red-600">Failed to load documents.</div>
              )}

              {/* always show static list; AI items are shown in the Procedure Checklist card */}
              {!checklistLoading && !checklistError && (
                <div className="space-y-3">
                  {/* ...existing static Required Documents list... */}
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-3">
                      <div className="w-4 h-4 flex-shrink-0">
                        <svg width="16" height="16" viewBox="0 0 16 16" fill="none" xmlns="http://www.w3.org/2000/svg">
                          <g clipPath="url(#clip0_1_2794)">
                            <path
                              d="M0 2C0 0.896875 0.896875 0 2 0H7V4C7 4.55312 7.44688 5 8 5H12V9.5H5.5C4.39687 9.5 3.5 10.3969 3.5 11.5V16H2C0.896875 16 0 15.1031 0 14V2ZM12 4H8V0L12 4ZM5.5 11H6.5C7.46562 11 8.25 11.7844 8.25 12.75C8.25 13.7156 7.46562 14.5 6.5 14.5H6V15.5C6 15.775 5.775 16 5.5 16C5.225 16 5 15.775 5 15.5V14V11.5C5 11.225 5.225 11 5.5 11ZM6.5 13.5C6.91563 13.5 7.25 13.1656 7.25 12.75C7.25 12.3344 6.91563 12 6.5 12H6V13.5H6.5ZM9.5 11H10.5C11.3281 11 12 11.6719 12 12.5V14.5C12 15.3281 11.3281 16 10.5 16H9.5C9.225 16 9 15.775 9 15.5V11.5C9 11.225 9.225 11 9.5 11ZM10.5 15C10.775 15 11 14.775 11 14.5V12.5C11 12.225 10.775 12 10.5 12H10V15H10.5ZM13 11.5C13 11.225 13.225 11 13.5 11H15C15.275 11 15.5 11.225 15.5 11.5C15.5 11.775 15.275 12 15 12H14V13H15C15.275 13 15.5 13.225 15.5 13.5C15.5 13.775 15.275 14 15 14H14V15.5C14 15.775 13.775 16 13.5 16C13.225 16 13 15.775 13 15.5V13.5V11.5Z"
                              fill="#EF4444"
                            />
                          </g>
                          <defs>
                            <clipPath id="clip0_1_2794">
                              <path d="M0 0H16V16H0V0Z" fill="white" />
                            </clipPath>
                          </defs>
                        </svg>
                      </div>
                      <div>
                        <div className="text-sm font-medium text-[#2e4d57]">Business Registration Form</div>
                        <div className="text-xs text-[#a7b3b9]">PDF • 245 KB</div>
                      </div>
                    </div>
                    <div className="w-5 h-5 text-[#22c55e]">✓</div>
                  </div>

                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-3">
                      <div className="w-4 h-4 flex-shrink-0">
                        <svg width="12" height="16" viewBox="0 0 12 16" fill="none" xmlns="http://www.w3.org/2000/svg"/>
                          <g clipPath="url(#clip0_1_2808)">
                            <path
                              d="M2 0C0.896875 0 0 0.896875 0 2V14C0 15.1031 0.896875 16 2 16H10C11.1031 16 12 15.1031 12 14V5H8C7.44688 5 7 4.55312 7 4V0H2ZM8 0V4H12L8 0ZM2 8C2 7.73478 2.10536 7.48043 2.29289 7.29289C2.48043 7.10536 2.73478 7 3 7C3.26522 7 3.51957 7.10536 3.70711 7.29289C3.89464 7.48043 4 7.73478 4 8C4 8.26522 3.89464 8.51957 3.70711 8.70711C3.51957 8.89464 3.26522 9 3 9C2.73478 9 2.48043 8.89464 2.29289 8.70711C2.10536 8.51957 2 8.26522 2 8ZM6.75 9C6.91563 9 7.06875 9.08125 7.1625 9.21562L9.9125 13.2156C10.0188 13.3687 10.0281 13.5687 9.94375 13.7312C9.85938 13.8938 9.6875 14 9.5 14H6.75H5.5H4H2.5C2.31875 14 2.15313 13.9031 2.06562 13.7469C1.97812 13.5906 1.97813 13.3969 2.07188 13.2437L3.57188 10.7437C3.6625 10.5938 3.825 10.5 4 10.5C4.175 10.5 4.3375 10.5906 4.42812 10.7437L4.82812 11.4125L6.3375 9.21875C6.43125 9.08438 6.58437 9.00313 6.75 9.00313V9Z"
                            fill="#3B82F6"
                          />
                        </g>
                        <defs>
                          <clipPath id="clip0_1_2808">
                            <path d="M0 0H12V16H0V0Z" fill="white" />
                          </clipPath>
                        </defs>
                      </div>
                      <div>
                        <div className="text-sm font-medium text-[#2e4d57]">Photo ID Copy</div>
                        <div className="text-xs text-[#a7b3b9]">JPG/PDF • Max 5 MB</div>
                      </div>
                    </div>
                    <div className="w-5 h-5 text-[#22c55e]">✓</div>
                  </div>

                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-3">
                      <div className="w-4 h-4 flex-shrink-0">
                        <svg width="12" height="16" viewBox="0 0 12 16" fill="none" xmlns="http://www.w3.org/2000/svg"/>
                          <g clipPath="url(#clip0_1_2822)">
                            <path
                              d="M2 0C0.896875 0 0 0.896875 0 2V14C0 15.1031 0.896875 16 2 16H10C11.1031 16 12 15.1031 12 14V5H8C7.44688 5 7 4.55312 7 4V0H2ZM8 0V4H12L8 0ZM4.86562 7.81875L6 9.44063L7.13438 7.81875C7.37188 7.47813 7.84063 7.39687 8.17813 7.63438C8.51562 7.87188 8.6 8.34063 8.3625 8.67813L6.91563 10.75L8.36563 12.8188C8.60313 13.1594 8.52187 13.625 8.18125 13.8625C7.84062 14.1 7.375 14.0188 7.1375 13.6781L6 12.0562L4.86562 13.6781C4.62812 14.0188 4.15938 14.1 3.82188 13.8625C3.48438 13.625 3.4 13.1562 3.6375 12.8188L5.08437 10.75L3.63438 8.68125C3.39688 8.34062 3.47812 7.875 3.81875 7.6375C4.15937 7.4 4.625 7.48125 4.8625 7.82188L4.86562 7.81875Z"
                            fill="#22C55E"
                          />
                        </g>
                        <defs>
                          <clipPath id="clip0_1_2822">
                            <path d="M0 0H12V16H0V0Z" fill="white" />
                          </clipPath>
                        </defs>
                      </div>
                      <div>
                        <div className="text-sm font-medium text-[#2e4d57]">Financial Statements</div>
                        <div className="text-xs text-[#a7b3b9]">PDF/Excel • Max 10 MB</div>
                      </div>
                    </div>
                    <div className={`w-5 h-5 ${applicationState.step3.completed ? "text-[#22c55e]" : "text-[#ffb703]"}`}>
                      {applicationState.step3.completed ? "✓" : "●"}
                    </div>
                  </div>

                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-3">
                      <div className="w-4 h-4 flex-shrink-0">
                        <svg width="12" height="16" viewBox="0 0 12 16" fill="none" xmlns="http://www.w3.org/2000/svg">
                          <g clipPath="url(#clip0_1_2836)">
                            <path
                              d="M2 0C0.896875 0 0 0.896875 0 2V14C0 15.1031 0.896875 16 2 16H10C11.1031 16 12 15.1031 12 14V5H8C7.44688 5 7 4.55312 7 4V0H2ZM8 0V4H12L8 0ZM2.5 2H4.5C4.775 2 5 2.225 5 2.5C5 2.775 4.775 3 4.5 3H2.5C2.225 3 2 2.775 2 2.5C2 2.225 2.225 2 2.5 2ZM2.5 4H4.5C4.775 4 5 4.225 5 4.5C5 4.775 4.775 5 4.5 5H2.5C2.225 5 2 4.775 2 4.5C2 4.225 2.225 4 2.5 4ZM4.19375 11.9312C4.00312 12.5656 3.41875 13 2.75625 13H2.5C2.225 13 2 12.775 2 12.5C2 12.225 2.225 12 2.5 12H2.75625C2.97812 12 3.17188 11.8562 3.23438 11.6438L3.7 10.0969C3.80625 9.74375 4.13125 9.5 4.5 9.5C4.86875 9.5 5.19375 9.74063 5.3 10.0969L5.6625 11.3031C5.89375 11.1094 6.1875 11 6.5 11C6.99687 11 7.45 11.2812 7.67188 11.725L7.80937 12H9.5C9.775 12 10 12.225 10 12.5C10 12.775 9.775 13 9.5 13H7.5C7.30937 13 7.1375 12.8938 7.05312 12.725L6.77812 12.1719C6.725 12.0656 6.61875 12 6.50313 12C6.3875 12 6.27813 12.0656 6.22813 12.1719L5.95312 12.725C5.8625 12.9094 5.66563 13.0188 5.4625 13C5.25938 12.9812 5.08437 12.8406 5.02812 12.6469L4.5 10.9062L4.19375 11.9312Z"
                            fill="#A855F7"
                          />
                        </g>
                        <defs>
                          <clipPath id="clip0_1_2836">
                            <path d="M0 0H12V16H0V0Z" fill="white" />
                          </clipPath>
                        </defs>
                        </svg>
  
                      </div>
                      <div>
                        <div className="text-sm font-medium text-[#2e4d57]">Lease Agreement</div>
                        <div className="text-xs text-[#a7b3b9]">PDF • Max 5 MB</div>
                      </div>
                    </div>
                    <div className="w-5 h-5 text-[#a7b3b9]">○</div>
                  </div>
                </div>
              )}

              <Button variant="outline" className="w-full mt-4 border-[#e5e7eb] text-[#0075ff] bg-transparent">
                <Download className="w-4 h-4 mr-2" />
                Download All Templates
              </Button>
            </div>

            {/* Related Notices */}
            {/* <div className="bg-white rounded-lg p-6 card-hover">
              <h3 className="font-medium text-[#2e4d57] mb-4">Related Notices</h3>

              <div className="space-y-4">
                <div className="bg-[#fefce8] border border-[#ffb703] rounded-lg p-3">
                  <div className="flex items-start gap-2">
                    <AlertTriangle className="w-4 h-4 text-[#ffb703] mt-0.5" />
                    <div>
                      <div className="text-sm font-medium text-[#2e4d57]">New Fee Structure</div>
                      <div className="text-xs text-[#a7b3b9] mt-1">
                        Updated licensing fees effective January 2024. Review the new fee schedule before payment.
                      </div>
                      <div className="text-xs text-[#a7b3b9] mt-2">Dec 15, 2023</div>
                    </div>
                  </div>
                </div>

                <div className="bg-[#eff6ff] border border-[#0075ff] rounded-lg p-3">
                  <div className="flex items-start gap-2">
                    <Info className="w-4 h-4 text-[#0075ff] mt-0.5" />
                    <div>
                      <div className="text-sm font-medium text-[#2e4d57]">Holiday Processing Delays</div>
                      <div className="text-xs text-[#a7b3b9] mt-1">
                        Applications submitted during holiday season may take additional 2-7 business days.
                      </div>
                      <div className="text-xs text-[#a7b3b9] mt-2">Dec 20, 2023</div>
                    </div>
                  </div>
                </div>
              </div>
            </div> */}
          </div>
        </div>
      </div>

    </div>
  )
}

