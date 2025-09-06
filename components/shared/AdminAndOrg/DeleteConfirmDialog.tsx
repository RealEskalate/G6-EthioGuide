"use client";
import { Trash2 } from "lucide-react";

import {
  AlertDialog,
  AlertDialogTrigger,
  AlertDialogContent,
  AlertDialogHeader,
  AlertDialogFooter,
  AlertDialogTitle,
  AlertDialogDescription,
  AlertDialogCancel,
  AlertDialogAction,
} from "@/components/ui/alert-dialog";
import { Button } from "@/components/ui/button";
import { ReactNode } from "react";

interface DeleteConfirmDialogProps {
  trigger?: ReactNode;
  title?: string;
  description?: string;
  confirmLabel?: string;
  onConfirm: () => void;
}

export default function DeleteConfirmDialog({
  title = "Are you sure?",
  description = "This action cannot be undone. This will permanently delete the item.",
  confirmLabel = "Delete",
  onConfirm,
}: DeleteConfirmDialogProps) {
  return (
    <AlertDialog>
      {/* Trigger button (or custom trigger) */}
      <AlertDialogTrigger>
        <Button className="bg-neutral-50 hover:bg-neutral">
          <Trash2 className="w-4 h-4 text-red-600 cursor-pointer" />
        </Button>
      </AlertDialogTrigger>

      <AlertDialogContent className="bg-white text-primary-dark">
        <AlertDialogHeader>
          <AlertDialogTitle>{title}</AlertDialogTitle>
          <AlertDialogDescription>{description}</AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel>Cancel</AlertDialogCancel>
          <AlertDialogAction
            onClick={onConfirm}
              className="bg-red-600 hover:bg-red-700 text-white font-bold"
          >
            {confirmLabel}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}
