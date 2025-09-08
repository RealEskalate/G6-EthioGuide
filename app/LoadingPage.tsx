import React from 'react'
import Image from 'next/image'

const LoadingPage = () => {
  return (
    <div className="flex items-center justify-center h-screen">
      <Image 
        src="/images/ethio-guide-boomerang-low-size.gif" 
        alt="Loading..." 
        width={100} 
        height={100} 
        priority 
      />
    </div>
  )
}

export default LoadingPage
