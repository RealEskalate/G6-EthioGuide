"use client";
import { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Send, ArrowLeft } from "lucide-react";
import { useRouter } from "next/navigation";
import { useCreateDiscussionMutation } from "@/app/store/slices/discussionsSlice";
import { Toaster, toast } from "react-hot-toast";
import { useSession, SessionProvider } from "next-auth/react";

// Export a wrapper that provides SessionProvider above useSession
export default function CreatePostPage() {
  return (
    <SessionProvider>
      <CreatePostContent />
    </SessionProvider>
  );
}

function CreatePostContent() {
  const [tag, setTag] = useState("");
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  // added: optional procedure id
  const [procedureId, setProcedureId] = useState("");
  const router = useRouter();
  const [createDiscussion, { isLoading }] = useCreateDiscussionMutation();
  const { data: session } = useSession();
  interface SessionWithToken {
    accessToken?: string;
    [key: string]: unknown;
  }
  const accessToken = (session as SessionWithToken | null | undefined)?.accessToken;
  console.log("access token:", accessToken);

  // added: persist session token so RTK slices can read it
  useEffect(() => {
    if (accessToken) {
      try {
        localStorage.setItem("accessToken", accessToken);
      } catch {}
    }
  }, [accessToken]);

  // helper for simple progress bars
  const pct = (value: number, max: number) => Math.min(100, Math.round((value / max) * 100));

  const handlePublish = async () => {
    if (!title.trim() || !content.trim()) return;
    try {
      const safeProcedure = procedureId.trim();
      await createDiscussion({
        title: title.trim(),
        content: content.trim(),
        tags: tag ? [tag] : undefined,
        // include procedureId; send empty string if not provided (backend can ignore)
        procedureId: safeProcedure || "",
        procedures: [], // ensure empty array is sent
      }).unwrap();
      toast.success("Your discussion has been created successfully.");
      setTitle("");
      setContent("");
      setTag("");
      setProcedureId(""); // clear
    } catch (e: unknown) {
      const errorMessage =
        typeof e === "object" && e !== null && "data" in e && typeof (e as { data?: { message?: string } }).data?.message === "string"
          ? (e as { data?: { message?: string } }).data?.message
          : "Failed to publish. Please try again.";
      toast.error(errorMessage ?? "Failed to publish. Please try again.");
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 p-4 sm:p-6">
      {/* animations (scoped) */}
      <style jsx>{`
        .fade-in { animation: fadeIn 380ms ease-out both; }
        .fade-in-up { animation: fadeInUp 520ms ease-out both; }
        .pop-in { animation: popIn 360ms ease-out both; }
        .card-float { transition: transform .25s ease, box-shadow .25s ease; }
        .card-float:hover { transform: translateY(-3px); box-shadow: 0 12px 26px rgba(0,0,0,.06); }
        .btn-pop { transition: transform .2s ease; }
        .btn-pop:hover { transform: translateY(-1px) scale(1.01); }
        .btn-pop:active { transform: scale(0.98); }
        .focus-glow { transition: box-shadow .25s ease; }
        .focus-glow:focus { box-shadow: 0 0 0 4px rgba(58,106,141,.15); }
        @keyframes fadeIn { from { opacity: 0 } to { opacity: 1 } }
        @keyframes fadeInUp { from { opacity: 0; transform: translateY(12px) } to { opacity: 1; transform: translateY(0) } }
        @keyframes popIn { 0% { opacity: 0; transform: scale(.98) } 100% { opacity: 1; transform: scale(1) } }
      `}</style>

      <Toaster position="top-right" toastOptions={{ duration: 4000 }} />
      <div className="max-w-2xl mx-auto">
        {/* Header surface */}
        <div
          className="bg-white/90 border border-gray-100 rounded-xl p-4 sm:p-5 mb-6 shadow-sm fade-in-up"
          style={{ animationDelay: "40ms" }}
        >
          <div className="flex items-center gap-3">
            <button
              type="button"
              className="p-2 rounded-full hover:bg-gray-100 transition-colors"
              onClick={() => router.push("/user/discussions")}
              aria-label="Back to discussions"
            >
              <ArrowLeft className="h-5 w-5 text-[#3A6A8D]" />
            </button>
            <h1 className="text-2xl font-bold text-gray-900">Create New Post</h1>
          </div>
          <p className="text-sm text-gray-600 mt-1">Share helpful tips, ask questions, or start a discussion.</p>
        </div>

        {/* Post Type Card */}
        <div
          className="bg-white rounded-xl p-6 mb-6 shadow-sm border border-gray-100 fade-in-up card-float"
          style={{ animationDelay: "90ms" }}
        >
          <h2 className="text-lg font-semibold text-gray-800 mb-4">Post Type</h2>
          <div className="mb-2">
            <label className="block text-sm font-medium text-gray-700 mb-2">Tags</label>
            <Select value={tag} onValueChange={setTag}>
              <SelectTrigger className="w-full bg-white rounded-lg border border-gray-200 shadow-sm px-4 py-2 text-gray-900 focus:ring-2 focus:ring-[#3A6A8D] focus:border-transparent transition-all">
                <SelectValue placeholder="Select a tag" />
              </SelectTrigger>
              <SelectContent className="rounded-lg border border-gray-200 shadow-md bg-white">
                <SelectItem value="passport">passport</SelectItem>
                <SelectItem value="tax">tax</SelectItem>
                <SelectItem value="business">business</SelectItem>
                <SelectItem value="events">events</SelectItem>
              </SelectContent>
            </Select>
            <p className="text-xs text-gray-400 mt-1">Pick a tag that best fits your post.</p>
          </div>
          {/* added: optional procedure id input */}
          <div className="mt-4">
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Procedure (optional)
            </label>
              <Input
                value={procedureId}
                onChange={(e) => setProcedureId(e.target.value)}
                placeholder="Enter related procedure ID (if any)"
                className="border-gray-200 focus-glow"
                maxLength={80}
              />
            <p className="text-xs text-gray-400 mt-1">
              Link this discussion to a procedure. Leave blank if not applicable.
            </p>
          </div>
        </div>

        {/* Title & Content Card */}
        <div
          className="bg-white rounded-xl p-6 mb-6 shadow-sm border border-gray-100 fade-in-up card-float"
          style={{ animationDelay: "140ms" }}
        >
          <div className="mb-6">
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Title / Headline <span className="text-red-500">*</span>
            </label>
            <Input
              value={title}
              onChange={e => setTitle(e.target.value)}
              placeholder="Enter a clear, descriptive title..."
              maxLength={100}
              className="mb-2 focus-glow"
            />
            {/* visual progress for title */}
            <div className="h-1.5 w-full bg-gray-100 rounded-full overflow-hidden">
              <div
                className={`h-full rounded-full transition-all duration-300 ${pct(title.length, 100) > 85 ? "bg-red-300" : "bg-[#3A6A8D]"}`}
                style={{ width: `${pct(title.length, 100)}%` }}
              />
            </div>
            <div className="text-xs text-gray-400 mt-1">{title.length}/100 characters</div>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Content <span className="text-red-500">*</span>
            </label>
            {/* simple toolbar placeholder */}
            <div className="flex items-center gap-2 mb-2 text-gray-400">
              {/* ...optional toolbar icons (kept minimal to preserve color)... */}
              {/* <span className="font-bold">B</span><span className="italic">I</span><span className="underline">U</span> */}
            </div>
            <Textarea
              value={content}
              onChange={e => setContent(e.target.value)}
              placeholder="Write your post content here..."
              maxLength={2000}
              className="mb-2 focus-glow"
              rows={6}
            />
            {/* visual progress for content */}
            <div className="h-1.5 w-full bg-gray-100 rounded-full overflow-hidden">
              <div
                className={`h-full rounded-full transition-all duration-300 ${pct(content.length, 2000) > 90 ? "bg-red-300" : "bg-[#3A6A8D]"}`}
                style={{ width: `${pct(content.length, 2000)}%` }}
              />
            </div>
            <div className="text-xs text-gray-400 mt-1">{content.length}/2000 characters</div>
          </div>
        </div>

        {/* Actions */}
        <div className="flex flex-col sm:flex-row gap-3 fade-in" style={{ animationDelay: "180ms" }}>
          <Button
            className="flex items-center gap-2 px-6 py-2 bg-[#3A6A8D] hover:bg-[#2d5470] text-white shadow-sm w-full sm:w-auto btn-pop"
            disabled={isLoading || title.trim() === "" || content.trim() === ""}
            onClick={handlePublish}
          >
            <Send className="w-4 h-4" /> {isLoading ? "Publishing..." : "Publish Post"}
          </Button>
          <Button
            variant="outline"
            className="border-gray-300 hover:bg-gray-50 w-full sm:w-auto btn-pop"
            type="button"
            onClick={() => {
              // lightweight UX: allow user to clear quickly
              setTitle(""); setContent(""); setTag(""); setProcedureId("");
            }}
          >
            Clear
          </Button>
        </div>
      </div>
    </div>
  );
}


