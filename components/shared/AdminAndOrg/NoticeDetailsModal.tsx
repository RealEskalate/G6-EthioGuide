"use client"

import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription } from "@/components/ui/dialog"
import { Badge } from "@/components/ui/badge"
import { format } from "date-fns"

interface Notice {
  id: string
  orgId: string
  title: string
  body: string
  pinned?: boolean
  effectiveFrom?: string
  effectiveTo?: string
  procedures: { id: string; name: string }[]
  createdAt: string
  updatedAt: string
}

interface NoticeDetailsModalProps {
  open: boolean
  onClose: () => void
  notice: Notice | null
}

export function NoticeDetailsModal({ open, onClose, notice }: NoticeDetailsModalProps) {
  if (!notice) return null

  return (
    <Dialog open={open} onOpenChange={onClose}>
      <DialogContent className="max-w-2xl">
        <DialogHeader>
          <DialogTitle>{notice.title}</DialogTitle>
          <DialogDescription>
            Notice details and related procedures
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4">
          <p className="text-sm text-muted-foreground">{notice.body}</p>

          {notice.pinned && <Badge variant="secondary">📌 Pinned</Badge>}

          <div className="grid grid-cols-2 gap-4 text-sm">
            {notice.effectiveFrom && (
              <p>
                <strong>Effective From:</strong>{" "}
                {format(new Date(notice.effectiveFrom), "PPP")}
              </p>
            )}
            {notice.effectiveTo && (
              <p>
                <strong>Effective To:</strong>{" "}
                {format(new Date(notice.effectiveTo), "PPP")}
              </p>
            )}
          </div>

          {notice.procedures?.length > 0 && (
            <div>
              <strong>Related Procedures:</strong>
              <ul className="list-disc list-inside text-sm mt-1">
                {notice.procedures.map(proc => (
                  <li key={proc.id}>{proc.name}</li>
                ))}
              </ul>
            </div>
          )}

          <div className="text-xs text-muted-foreground">
            <p>Created: {format(new Date(notice.createdAt), "PPP p")}</p>
            <p>Updated: {format(new Date(notice.updatedAt), "PPP p")}</p>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  )
}
