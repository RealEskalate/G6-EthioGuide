import { Button } from "@/components/ui/button"

interface PaginationProps {
  page: number
  totalPages: number
  onPageChange: (page: number) => void
}

export default function Pagination({ page, totalPages, onPageChange }: PaginationProps) {

  return (
    <div className="flex items-center justify-between p-4 text-sm text-neutral">
      <span>Page {page} of {totalPages}</span>
      <div className="flex gap-2">
        <Button 
          variant="outline" 
          size="sm" 
          disabled={page === 1} 
          onClick={() => onPageChange(page - 1)}
        >
          Previous
        </Button>

        <Button 
          className="text-white bg-primary" 
          size="sm"
        >
          {page}
        </Button>

        <Button 
          variant="outline" 
          size="sm" 
          disabled={page === totalPages} 
          onClick={() => onPageChange(page + 1)}
        >
          Next
        </Button>
      </div>
    </div>
  )
}
