"use client";

import { useRouter } from "next/navigation";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Calendar, FileText, User } from "lucide-react";

import { notices } from "@/lib/noticesData";

export default function NoticeDetailPage({ params }: { params: { id: string } }) {
  const noticeId = Number(params.id);
  const notice = notices.find(n => n.id === noticeId);
  
  if (!notice) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <h2 className="text-2xl font-bold mb-4">Notice Not Found</h2>
          <button className="bg-[#3A6A8D] hover:bg-[#2d5470] text-white px-4 py-2 rounded" onClick={() => router.back()}>
            Back to Notices
          </button>
        </div>
      </div>
    );
  }

  const router = useRouter();
  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-3xl mx-auto space-y-6">
        {/* Back Arrow and Title */}
        <div className="flex items-center gap-2 mb-2">
          <button
            className="flex items-center gap-2 text-[#3A6A8D] hover:text-[#2d5470] font-medium"
            onClick={() => router.back()}
          >
            <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" /></svg>
            Back to Notices
          </button>
        </div>
        <h1 className="text-xl font-semibold text-gray-900 mb-2">{notice.title}</h1>
        {/* Card 1: Organization Info */}
        <Card className="bg-white p-4">
          <CardContent className="p-0">
            <div className="flex flex-wrap items-center gap-6">
              <span className="flex items-center gap-1 text-gray-500">
                <FileText className="w-4 h-4" />
                Organization: <span className="font-semibold text-gray-900">{notice.organization}</span>
              </span>
              <span className="flex items-center gap-1 text-gray-500">
                <Calendar className="w-4 h-4" /> Posted: <span className="font-semibold text-gray-900">{notice.published}</span>
              </span>
              <span className="flex items-center gap-1 text-gray-500">
                <User className="w-4 h-4" /> Author: <span className="font-semibold text-gray-900">{notice.author}</span>
              </span>
            </div>
          </CardContent>
        </Card>

        {/* Card 2: Main Notice Content */}
        <Card className="bg-white p-6">
          <CardContent className="p-0">
            <p className="text-gray-700 mb-4">Dear Team Members,</p>
            <p className="text-gray-700 mb-4">{notice.description}</p>
            {Array.isArray(notice.keyChanges) && notice.keyChanges.length > 0 && (
              <>
                <h2 className="font-semibold text-gray-900 mb-2">Key Policy Changes:</h2>
                <ul className="list-disc ml-6 mb-4 text-gray-700">
                  {notice.keyChanges.map((change, idx) => (
                    <li key={idx}>{change}</li>
                  ))}
                </ul>
              </>
            )}
            {Array.isArray(notice.securityUpdates) && notice.securityUpdates?.length > 0 && (
              <>
                <h2 className="font-semibold text-gray-900 mb-2">Security Protocol Updates:</h2>
                <ul className="list-disc ml-6 mb-4 text-gray-700">
                  {notice.securityUpdates?.map((update, idx) => (
                    <li key={idx}>{update}</li>
                  ))}
                </ul>
              </>
            )}
            <p className="text-gray-700 mt-4">For questions regarding these policy changes, please contact the HR department or attend one of our information sessions scheduled for next week. Your cooperation in implementing these changes is greatly appreciated.</p>
            <p className="text-gray-700 mt-4">Best regards, <span className="font-semibold">Sarah Johnson</span> HR Director,<br />TechCorp Solutions</p>
          </CardContent>
        </Card>

        {/* Card 3: Attachments */}
        {Array.isArray(notice.attachments) && notice.attachments?.length > 0 && (
          <Card className="bg-white p-6">
            <CardContent className="p-0">
              <h2 className="font-semibold text-gray-900 mb-4 flex items-center gap-2">
                <FileText className="w-5 h-5 text-[#3A6A8D]" /> Attachments
              </h2>
              <div className="space-y-4">
                {notice.attachments?.map((file, idx) => (
                  <div key={idx} className="flex items-center gap-4 p-4 bg-gray-50 rounded-lg">
                    <FileText className={`w-6 h-6 ${file.name.endsWith('.pdf') ? 'text-red-500' : 'text-blue-500'}`} />
                    <div className="flex-1">
                      <span className="font-medium text-gray-900">{file.name}</span>
                      <div className="text-xs text-gray-500">{file.size}</div>
                    </div>
                    <button className="text-[#3A6A8D] hover:underline">
                      <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16v2a2 2 0 002 2h12a2 2 0 002-2v-2M7 10l5 5m0 0l5-5m-5 5V4" /></svg>
                    </button>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        )}
      </div>
    </div>
  );
}
