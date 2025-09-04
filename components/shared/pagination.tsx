interface PaginationProps {
  page: number
  limit: number
  total: number
  onPageChange: (page: number) => void
}

export function Pagination({ page, limit, total, onPageChange }: PaginationProps) {
  const totalPages = Math.ceil(total / limit)

  return (
    <div className="flex gap-2">
      <button 
        disabled={page === 1} 
        onClick={() => onPageChange(page - 1)}
      >
        Prev
      </button>

      <span>{page} / {totalPages}</span>

      <button 
        disabled={page === totalPages} 
        onClick={() => onPageChange(page + 1)}
      >
        Next
      </button>
    </div>
  )
}
