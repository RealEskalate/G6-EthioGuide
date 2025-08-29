import * as React from "react";

export interface TextareaProps extends React.TextareaHTMLAttributes<HTMLTextAreaElement> {}

export const Textarea = React.forwardRef<HTMLTextAreaElement, TextareaProps>(
  ({ className, ...props }, ref) => {
    return (
      <textarea
        ref={ref}
        className={`block w-full rounded-lg border border-gray-200 px-3 py-2 text-gray-900 shadow-sm focus:outline-none focus:ring-2 focus:ring-[#3A6A8D] focus:border-transparent transition-all resize-none ${className || ""}`}
        {...props}
      />
    );
  }
);

Textarea.displayName = "Textarea";
