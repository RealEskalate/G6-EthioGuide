import Image from "next/image"

interface CustomIconProps {
  src: string
  alt: string
  className?: string
}

export function CustomIcon({ src, alt, className = "w-5 h-5" }: CustomIconProps) {
  return <Image src={src || "/placeholder.svg"} alt={alt} width={20} height={20} className={className} />
}
