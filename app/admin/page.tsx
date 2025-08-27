import React from 'react'
import AdminDashboard from '@/components/admin/AdminDashboard'
import LeftSidebar from '@/components/admin-org-components/LeftSidebar'

const page = () => {
  return (
    <div className='bg-gray-50'>
      {/* <LeftSidebar /> */}
      <AdminDashboard />
    </div>
  )
}

export default page