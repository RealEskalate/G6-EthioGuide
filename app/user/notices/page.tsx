"use client";

import { useState, useEffect, useMemo } from "react";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Calendar, FileText, Search } from "lucide-react";
import { useGetNoticesQuery } from "@/app/store/slices/noticesSlice";
import { useRouter } from "next/navigation";
import { motion, useReducedMotion } from "framer-motion";

type Notice = {
  id: string | number;
  title: string;
  content: string;
  tags?: string[];
  created_at: string;
  organization_id?: string;
};

export default function NoticesPage() {
  const [organization, setOrganization] = useState("all");
  const [expandedMap, setExpandedMap] = useState<Record<string | number, boolean>>({});
  const [searchInput, setSearchInput] = useState("");
  const [searchQuery, setSearchQuery] = useState("");

  // const router = useRouter();

  const { data: apiNotices, isLoading: noticesLoading, isError: noticesError } =
    useGetNoticesQuery({ page: 1, limit: 10 });

  useEffect(() => {
    if (apiNotices) console.log("Notices API response:", apiNotices);
  }, [apiNotices]);

  const notices = apiNotices?.data ?? [];

  const organizations = useMemo(() => {
    const set = new Set<string>();
    notices.forEach((n: Notice) => {
      if (n.organization_id) set.add(n.organization_id);
    });
    return Array.from(set);
  }, [notices]);

  const filteredNotices = useMemo(() => {
    const base =
      organization === "all"
        ? notices
        : notices.filter((n: Notice) => n.organization_id === organization);

    const q = searchQuery.trim().toLowerCase();
    if (!q) return base;

    const terms = q.split(/\s+/).filter(Boolean);
    if (!terms.length) return base;

    return base.filter((n: Notice) => {
      const haystack = [
        n.title,
        n.content,
        n.organization_id ?? "",
        ...(Array.isArray(n.tags) ? n.tags : []),
      ]
        .join(" ")
        .toLowerCase();
      return terms.every((t) => haystack.includes(t));
    });
  }, [organization, notices, searchQuery]);

  const prefersReducedMotion = useReducedMotion();

  const containerVariants = prefersReducedMotion
    ? { hidden: { opacity: 0 }, visible: { opacity: 1 } }
    : { hidden: { opacity: 0 }, visible: { opacity: 1, transition: { staggerChildren: 0.06, delayChildren: 0.08 } } };

  const itemVariants = prefersReducedMotion
    ? { hidden: { opacity: 0 }, visible: { opacity: 1 } }
    : {
        hidden: { opacity: 0, y: 10, scale: 0.985, filter: "blur(0.2px)" },
        visible: {
          opacity: 1,
          y: 0,
          scale: 1,
          filter: "none",
          transition: { type: "spring" as const, stiffness: 220, damping: 18, mass: 0.9 }
        },
      };

  const headerVariants = prefersReducedMotion
    ? { hidden: { opacity: 0 }, visible: { opacity: 1 } }
    : { hidden: { opacity: 0 }, visible: { opacity: 1, transition: { delay: 0.06 } } };

  const formatPublished = (iso?: string) => {
    if (!iso) return "Unknown";
    const d = new Date(iso);
    if (isNaN(d.getTime())) return "Unknown";
    return d.toLocaleString(undefined, {
      year: "numeric",
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  const tagPillClasses = (i: number) => {
    const styles = [
      "bg-green-50 text-green-700 border-green-200 hover:bg-green-100 hover:text-green-800",
      "bg-amber-50 text-amber-700 border-amber-200 hover:bg-amber-100 hover:text-amber-800",
      "bg-teal-50 text-teal-700 border-teal-200 hover:bg-teal-100 hover:text-teal-800",
      "bg-indigo-50 text-indigo-700 border-indigo-200 hover:bg-indigo-100 hover:text-indigo-800",
      "bg-emerald-50 text-emerald-700 border-emerald-200 hover:bg-emerald-100 hover:text-emerald-800",
      "bg-cyan-50 text-cyan-700 border-cyan-200 hover:bg-cyan-100 hover:text-cyan-800",
    ];
    return `cursor-default rounded-full ${styles[i % styles.length]}`;
  };

  return (
    <motion.div className="min-h-screen bg-gray-50 flex flex-col">
      {/* Header */}
      <motion.div
        className="flex items-center justify-between mb-8"
        variants={headerVariants}
        initial="hidden"
        animate="visible"
      >
        <div>
          <div className="flex items-center gap-3 mb-2">
            <FileText className="h-8 w-8 text-[#3A6A8D]" />
            <h1 className="text-3xl font-bold text-gray-900">Official Notices</h1>
          </div>
          <p className="text-gray-600">Get notices of different organizations</p>
        </div>
      </motion.div>

      {/* Search and Filters - match Community Discussions UI */}
      <Card className="p-4 mb-6">
        <div className="flex flex-col gap-4 w-full mb-2 sm:flex-row">
          <div className="relative flex-1 flex">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
            <input
              type="text"
              placeholder="Search notices..."
              value={searchInput}
              onChange={(e) => setSearchInput(e.target.value)}
              className="w-full pl-10 pr-4 py-2 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-[#3A6A8D] focus:border-transparent"
            />
            <Button
              type="button"
              className="ml-2 px-4 py-2 bg-[#3A6A8D] hover:bg-[#2d5470] text-white"
              onClick={() => setSearchQuery(searchInput)}
            >
              Search
            </Button>
          </div>
          <div className="flex gap-2 flex-none w-full sm:w-64">
            <Select value={organization} onValueChange={setOrganization}>
              <SelectTrigger className="w-full border-[#3A6A8D] text-[#3A6A8D] hover:bg-[#3A6A8D]/5 focus:ring-2 focus:ring-[#3A6A8D] focus:border-transparent">
                <SelectValue placeholder="All Organizations" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Organizations</SelectItem>
                {organizations.map((org) => (
                  <SelectItem key={org} value={org}>
                    {org}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        </div>
      </Card>

      {/* Notices List */}
      <div className="space-y-6">
        {noticesLoading && <div className="text-gray-600">Loading notices...</div>}
        {noticesError && <div className="text-red-600">Failed to load notices.</div>}
        {!noticesLoading &&
          !noticesError &&
          filteredNotices.map((notice: Notice, index: number) => {
            const isExpanded = !!expandedMap[notice.id];
            return (
              <motion.div
                key={notice.id}
                variants={itemVariants}
                initial="hidden"
                animate="visible"
                transition={{ delay: index * 0.05 + 0.06 }}
              >
                <Card className="bg-white p-6 hover:shadow-lg transition-all duration-300">
                  <CardContent className="p-0">
                    <div className="flex items-center gap-3 mb-2">
                      <h2 className="text-lg font-semibold text-gray-900 mr-2">{notice.title}</h2>
                      {!!notice.tags?.length && (
                        <Badge variant="outline" className={`text-xs ${tagPillClasses(index)}`}>
                          {notice.tags[0]}
                        </Badge>
                      )}
                    </div>

                    <p className={`text-gray-600 mb-2 ${isExpanded ? "" : "line-clamp-2"}`}>{notice.content}</p>

                    <div className="flex flex-wrap items-center gap-8 sm:gap-16 text-sm text-gray-500 mb-2">
                      <div className="flex items-center gap-1">
                        <Calendar className="w-4 h-4" />
                        <span>Published on {formatPublished(notice.created_at)}</span>
                      </div>
                      {notice.organization_id && (
                        <span>
                          Organization:{" "}
                          <span className="font-semibold text-gray-700">{notice.organization_id}</span>
                        </span>
                      )}
                    </div>

                    <div className="flex justify-end pt-2">
                      <Button
                        variant="ghost"
                        size="sm"
                        className="text-[#3A6A8D] hover:bg-[#3A6A8D]/10 hover:text-[#2d5470]"
                        onClick={() =>
                          setExpandedMap((prev) => ({ ...prev, [notice.id]: !prev[notice.id] }))
                        }
                      >
                        {isExpanded ? "View Less" : "View More"}
                      </Button>
                    </div>
                  </CardContent>
                </Card>
              </motion.div>
            );
          })}
        {!noticesLoading && !noticesError && filteredNotices.length === 0 && (
          <motion.div
            variants={itemVariants}
            initial="hidden"
            animate="visible"
            transition={{ delay: 0.3 }}
            className="text-gray-600"
          >
            No notices found.
          </motion.div>
        )}
      </div>
    </motion.div>
  );
}

