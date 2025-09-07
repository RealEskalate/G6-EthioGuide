"use client";

import { useState, useEffect, useMemo } from "react";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Calendar, FileText, Search, Building2 } from "lucide-react";
import { useGetNoticesQuery } from "@/app/store/slices/noticesSlice";
import { motion, useReducedMotion } from "framer-motion";
import { useLazyGetOrgQuery } from "@/app/store/slices/orgsApi";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";

type Notice = {
  id: string | number;
  title: string;
  content: string;
  tags?: string[];
  created_at: string;
  organization_id?: string;
};

export default function NoticesPage() {
  const [expandedMap, setExpandedMap] = useState<Record<string | number, boolean>>({});
  const [searchInput, setSearchInput] = useState("");
  const [searchQuery, setSearchQuery] = useState("");
  const [orgFilter, setOrgFilter] = useState<string>("all");

  const { data: apiNotices, isLoading: noticesLoading, isError: noticesError } =
    useGetNoticesQuery({ page: 1, limit: 10 });

  useEffect(() => {
    if (apiNotices) console.log("Notices API response:", apiNotices);
  }, [apiNotices]);

  const notices = useMemo(() => apiNotices?.data ?? [], [apiNotices]);

  // Organization name fetching: cache orgId -> name locally in component state
  const [triggerGetOrg] = useLazyGetOrgQuery();
  const [orgNameMap, setOrgNameMap] = useState<Record<string, string>>({});

  useEffect(() => {
    const uniqueOrgIds = Array.from(
      new Set((notices as Notice[]).map((n) => n.organization_id).filter((v): v is string => typeof v === 'string' && !!v))
    );
    const missing = uniqueOrgIds.filter((id) => !(id in orgNameMap));
    if (!missing.length) return;
    let cancelled = false;
    (async () => {
      const entries: Array<[string, string]> = [];
      await Promise.all(
    missing.map(async (orgId) => {
          try {
      if (!orgId) return;
      const data = await triggerGetOrg(orgId, true).unwrap();
      entries.push([orgId, data.name]);
          } catch {
            // fallback: keep id as name on failure
      if (orgId) entries.push([orgId, orgId]);
          }
        })
      );
      if (!cancelled && entries.length) {
        setOrgNameMap((prev) => ({ ...prev, ...Object.fromEntries(entries) }));
      }
    })();
    return () => {
      cancelled = true;
    };
  }, [notices, orgNameMap, triggerGetOrg]);

  const filteredNotices = useMemo(() => {
    const base = notices;
    const q = searchQuery.trim().toLowerCase();
    if (!q) return base;

    const terms = q.split(/\s+/).filter(Boolean);
    if (!terms.length) return base;

  return base.filter((n: Notice) => {
      const haystack = [
        n.title,
        n.content,
        ...(Array.isArray(n.tags) ? n.tags : []),
    orgNameMap[n.organization_id || ""] || "",
      ]
        .join(" ")
        .toLowerCase();
      return terms.every((t) => haystack.includes(t));
    });
  }, [notices, searchQuery, orgNameMap]);

  const prefersReducedMotion = useReducedMotion();

  // containerVariants unused

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

  // Removed tag chips per request (no department/category text like "IT Operations")

  return (
    <motion.div className="min-h-screen w-full bg-gray-50 relative overflow-hidden flex flex-col">
      {/* subtle brand orbs */}
      <div className="absolute inset-0 overflow-hidden pointer-events-none">
        <div className="absolute -top-24 -right-24 w-56 h-56 rounded-full blur-3xl" style={{ background: 'radial-gradient(closest-side, rgba(167,179,185,0.10), rgba(167,179,185,0))' }} />
        <div className="absolute -bottom-28 -left-28 w-64 h-64 rounded-full blur-3xl" style={{ background: 'radial-gradient(closest-side, rgba(94,156,141,0.10), rgba(94,156,141,0))' }} />
      </div>
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

      {/* Search and Filters */}
      <Card className="bg-white/80 backdrop-blur-md rounded-2xl border border-[#e5e7eb] shadow-xl relative overflow-hidden p-4 mb-6">
        <div className="absolute inset-0 bg-gradient-to-r from-[#3a6a8d]/10 via-transparent to-[#5e9c8d]/10" />
        <div className="relative z-10 flex flex-col gap-4 w-full mb-2 sm:flex-row">
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
          {/* Org filter */}
      <div className="w-full sm:w-64">
            <Select value={orgFilter} onValueChange={setOrgFilter}>
        <SelectTrigger className="w-full">
                <SelectValue placeholder="All Organizations" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Organizations</SelectItem>
                {Array.from(
                  new Map(
                    (notices as Notice[])
                      .map((n) => n.organization_id)
                      .filter(Boolean)
                      .map((id) => [id!, orgNameMap[id!] || id!])
                  ).entries()
                ).map(([id, name]) => (
                  <SelectItem key={id} value={id}>
                    {name}
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
          filteredNotices
            .filter((n: Notice) => orgFilter === "all" || n.organization_id === orgFilter)
            .map((notice: Notice, index: number) => {
            const isExpanded = !!expandedMap[notice.id];
            return (
              <motion.div
                key={notice.id}
                variants={itemVariants}
                initial="hidden"
                animate="visible"
                transition={{ delay: index * 0.05 + 0.06 }}
              >
                <Card className="group bg-white/80 backdrop-blur-md rounded-2xl border border-[#e5e7eb] shadow-xl relative overflow-hidden ring-1 ring-transparent hover:ring-[#3a6a8d]/15 transition-all duration-300 transform-gpu hover:-translate-y-0.5 hover:shadow-2xl">
                  <div className="absolute inset-0 pointer-events-none">
                    <div className="absolute inset-0 bg-gradient-to-r from-[#3a6a8d]/5 via-transparent to-[#5e9c8d]/5" />
                  </div>
                  <CardContent className="relative z-10 p-6">
                    <div className="mb-3 flex items-center gap-3">
                      <div className="w-10 h-10 rounded-xl flex items-center justify-center" style={{ backgroundColor: '#e6f0f5' }}>
                        <FileText className="w-5 h-5" style={{ color: '#3a6a8d' }} />
                      </div>
                      <h2 className="text-lg font-semibold text-[#111827]">{notice.title}</h2>
                    </div>

                    <div className="flex flex-wrap items-center gap-x-6 gap-y-2 text-sm mb-3 text-[#4b5563]">
                      <span className="flex items-center gap-2">
                        <Calendar className="w-4 h-4" />
                        <span className="hidden sm:inline">Published:</span>
                        <span className="font-medium text-[#111827]">{formatPublished(notice.created_at)}</span>
                      </span>
                      <span className="flex items-center gap-2">
                        <Building2 className="w-4 h-4" />
                        <span className="hidden sm:inline">Organization:</span>
                        <span className="font-semibold text-[#111827]">{orgNameMap[notice.organization_id || ""] || "Organization"}</span>
                      </span>
                    </div>

                    <p className={`text-[#374151] mb-2 ${isExpanded ? "" : "line-clamp-2"}`}>{notice.content}</p>

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

