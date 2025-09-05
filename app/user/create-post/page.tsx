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

  const handlePublish = async () => {
    if (!title.trim() || !content.trim()) return;
    try {
      await createDiscussion({
        title: title.trim(),
        content: content.trim(),
        tags: tag ? [tag] : undefined,
        procedures: [], // ensure empty array is sent
      }).unwrap();
      toast.success("Your discussion has been created successfully.");
      setTitle("");
      setContent("");
      setTag("");
    } catch (e: unknown) {
      const errorMessage =
        typeof e === "object" && e !== null && "data" in e && typeof (e as { data?: { message?: string } }).data?.message === "string"
          ? (e as { data?: { message?: string } }).data?.message
          : "Failed to publish. Please try again.";
      toast.error(errorMessage ?? "Failed to publish. Please try again.");
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <Toaster position="top-right" toastOptions={{ duration: 4000 }} />
      <div className="max-w-2xl mx-auto">
        {/* Header */}
        <div className="flex items-center gap-3 mb-6">
          <button
            type="button"
            className="p-2 rounded-full hover:bg-gray-200 transition-colors"
            onClick={() => router.push("/user/discussions")}
            aria-label="Back to discussions"
          >
            <ArrowLeft className="h-5 w-5 text-[#3A6A8D]" />
          </button>
          <h1 className="text-2xl font-bold text-gray-900">Create New Post</h1>
        </div>

        {/* Post Type Card */}
        <div className="bg-white rounded-xl p-6 mb-6 shadow-lg">
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
          </div>
        </div>

        {/* Title & Content Card */}
        <div className="bg-white rounded-xl p-6 mb-6 shadow-lg">
          <div className="mb-6">
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Title / Headline <span className="text-red-500">*</span>
            </label>
            <Input
              value={title}
              onChange={e => setTitle(e.target.value)}
              placeholder="Enter a clear, descriptive title..."
              maxLength={100}
              className="mb-1"
            />
            <div className="text-xs text-gray-400">{title.length}/100 characters</div>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Content <span className="text-red-500">*</span>
            </label>
            {/* Simple toolbar mockup */}
            <div className="flex items-center gap-2 mb-2 text-gray-400">
              {/* <span className="font-bold">B</span>
              <span className="italic">I</span>
              <span className="underline">U</span>
              <span className="">•</span>
              <span className="">@</span> */}
            </div>
            <Textarea
              value={content}
              onChange={e => setContent(e.target.value)}
              placeholder="Write your post content here... Use @ to mention users or organizations"
              maxLength={2000}
              className="mb-1"
              rows={6}
            />
            <div className="text-xs text-gray-400">{content.length}/2000 characters</div>
          </div>
        </div>

        {/* Actions */}
        <div className="flex gap-4">
          <Button
            className="flex items-center gap-2 px-6 py-2 bg-primary hover:bg-[#2d5470] text-white"
            disabled={isLoading || title.trim() === "" || content.trim() === ""}
            onClick={handlePublish}
          >
            <Send className="w-4 h-4" /> {isLoading ? "Publishing..." : "Publish Post"}
          </Button>
        </div>
      </div>
    </div>
  );
}

