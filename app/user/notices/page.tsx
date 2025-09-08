// "use client";

// import { useState, useEffect, useMemo } from "react";
// import { Card, CardContent } from "@/components/ui/card";
// import { Button } from "@/components/ui/button";
// import { Calendar, FileText, Search, Building2 } from "lucide-react";
// import { useGetNoticesQuery } from "@/app/store/slices/noticesSlice";
// import { motion, useReducedMotion } from "framer-motion";
// import { useLazyGetOrgQuery } from "@/app/store/slices/orgsApi";
// import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";

// type Notice = {
//   id: string | number;
//   title: string;
//   content: string;
//   tags?: string[];
//   created_at: string;
//   organization_id?: string;
// };

// export default function NoticesPage() {
//   const [expandedMap, setExpandedMap] = useState<Record<string | number, boolean>>({});
//   const [searchInput, setSearchInput] = useState("");
//   const [searchQuery, setSearchQuery] = useState("");
//   const [orgFilter, setOrgFilter] = useState<string>("all");

//   const { data: apiNotices, isLoading: noticesLoading, isError: noticesError } =
//     useGetNoticesQuery({ page: 1, limit: 10 });

//   useEffect(() => {
//     if (apiNotices) console.log("Notices API response:", apiNotices);
//   }, [apiNotices]);

//   const notices = useMemo(() => apiNotices?.data ?? [], [apiNotices]);

//   // Organization name fetching: cache orgId -> name locally in component state
//   const [triggerGetOrg] = useLazyGetOrgQuery();
//   const [orgNameMap, setOrgNameMap] = useState<Record<string, string>>({});

//   useEffect(() => {
//     const uniqueOrgIds = Array.from(
//       new Set((notices as Notice[]).map((n) => n.organization_id).filter((v): v is string => typeof v === 'string' && !!v))
//     );
//     const missing = uniqueOrgIds.filter((id) => !(id in orgNameMap));
//     if (!missing.length) return;
//     let cancelled = false;
//     (async () => {
//       const entries: Array<[string, string]> = [];
//       await Promise.all(
//     missing.map(async (orgId) => {
//           try {
//       if (!orgId) return;
//       const data = await triggerGetOrg(orgId, true).unwrap();
//       entries.push([orgId, data.name]);
//           } catch {
//             // fallback: keep id as name on failure
//       if (orgId) entries.push([orgId, orgId]);
//           }
//         })
//       );
//       if (!cancelled && entries.length) {
//         setOrgNameMap((prev) => ({ ...prev, ...Object.fromEntries(entries) }));
//       }
//     })();
//     return () => {
//       cancelled = true;
//     };
//   }, [notices, orgNameMap, triggerGetOrg]);

//   const filteredNotices = useMemo(() => {
//     const base = notices;
//     const q = searchQuery.trim().toLowerCase();
//     if (!q) return base;

//     const terms = q.split(/\s+/).filter(Boolean);
//     if (!terms.length) return base;

//   return base.filter((n: Notice) => {
//       const haystack = [
//         n.title,
//         n.content,
//         ...(Array.isArray(n.tags) ? n.tags : []),
//     orgNameMap[n.organization_id || ""] || "",
//       ]
//         .join(" ")
//         .toLowerCase();
//       return terms.every((t) => haystack.includes(t));
//     });
//   }, [notices, searchQuery, orgNameMap]);

//   const prefersReducedMotion = useReducedMotion();

//   // containerVariants unused

//   const itemVariants = prefersReducedMotion
//     ? { hidden: { opacity: 0 }, visible: { opacity: 1 } }
//     : {
//         hidden: { opacity: 0, y: 10, scale: 0.985, filter: "blur(0.2px)" },
//         visible: {
//           opacity: 1,
//           y: 0,
//           scale: 1,
//           filter: "none",
//           transition: { type: "spring" as const, stiffness: 220, damping: 18, mass: 0.9 }
//         },
//       };

//   const headerVariants = prefersReducedMotion
//     ? { hidden: { opacity: 0 }, visible: { opacity: 1 } }
//     : { hidden: { opacity: 0 }, visible: { opacity: 1, transition: { delay: 0.06 } } };

//   const formatPublished = (iso?: string) => {
//     if (!iso) return "Unknown";
//     const d = new Date(iso);
//     if (isNaN(d.getTime())) return "Unknown";
//     return d.toLocaleString(undefined, {
//       year: "numeric",
//       month: "short",
//       day: "numeric",
//       hour: "2-digit",
//       minute: "2-digit",
//     });
//   };

//   // Removed tag chips per request (no department/category text like "IT Operations")

//   return (
//     <motion.div className="min-h-screen w-full bg-gray-50 relative overflow-hidden flex flex-col">
//       {/* subtle brand orbs */}
//       <div className="absolute inset-0 overflow-hidden pointer-events-none">
//         <div className="absolute -top-24 -right-24 w-56 h-56 rounded-full blur-3xl" style={{ background: 'radial-gradient(closest-side, rgba(167,179,185,0.10), rgba(167,179,185,0))' }} />
//         <div className="absolute -bottom-28 -left-28 w-64 h-64 rounded-full blur-3xl" style={{ background: 'radial-gradient(closest-side, rgba(94,156,141,0.10), rgba(94,156,141,0))' }} />
//       </div>
//       <div className="relative z-10 max-w-7xl w-full mx-auto px-4 sm:px-6 md:px-8">
//       {/* Header */}
//       <motion.div
//         className="bg-white border border-gray-100 rounded-xl p-4 sm:p-5 mb-4 sm:mb-6 shadow-sm"
//         variants={headerVariants}
//         initial="hidden"
//         animate="visible"
//       >
//         <div>
//           <div className="flex items-center gap-3 mb-1">
//             <FileText className="h-7 w-7" style={{ color: '#3a6a8d' }} />
//             <h1 className="text-xl leading-snug sm:text-2xl font-bold text-[#111827]">Official Notices</h1>
//           </div>
//           <p className="text-[#4b5563] text-sm">Get notices of different organizations</p>
//         </div>
//       </motion.div>

//       {/* Search and Filters */}
//   <Card className="bg-white border border-gray-100 rounded-xl p-4 mb-4 sm:mb-6 shadow-sm">
//         <div className="flex flex-col gap-4 w-full mb-2 sm:flex-row">
//           <div className="relative flex-1 flex">
//             <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
//             <input
//               type="text"
//               placeholder="Search notices..."
//               value={searchInput}
//               onChange={(e) => setSearchInput(e.target.value)}
//               className="w-full pl-10 pr-4 py-2 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-[#3A6A8D] focus:border-transparent"
//             />
//             <Button
//               type="button"
//               className="ml-2 px-4 py-2 bg-[#3A6A8D] hover:bg-[#2d5470] text-white"
//               onClick={() => setSearchQuery(searchInput)}
//             >
//               Search
//             </Button>
//           </div>
//           {/* Org filter */}
//           <div className="w-full sm:w-64">
//             <Select value={orgFilter} onValueChange={setOrgFilter}>
//         <SelectTrigger className="w-full border-[#3A6A8D] text-[#3A6A8D] focus:ring-2 focus:ring-[#3A6A8D] focus:border-transparent">
//                 <SelectValue placeholder="All Organizations" />
//               </SelectTrigger>
//               <SelectContent>
//                 <SelectItem value="all">All Organizations</SelectItem>
//                 {Array.from(
//                   new Map(
//                     (notices as Notice[])
//                       .map((n) => n.organization_id)
//                       .filter(Boolean)
//                       .map((id) => [id!, orgNameMap[id!] || id!])
//                   ).entries()
//                 ).map(([id, name]) => (
//                   <SelectItem key={id} value={id}>
//                     {name}
//                   </SelectItem>
//                 ))}
//               </SelectContent>
//             </Select>
//           </div>
//     </div>
//       </Card>

//       {/* Notices List */}
//       <div className="space-y-5 pb-8">
//         {noticesLoading && <div className="text-gray-600">Loading notices...</div>}
//         {noticesError && <div className="text-red-600">Failed to load notices.</div>}
//         {!noticesLoading &&
//           !noticesError &&
//           filteredNotices
//             .filter((n: Notice) => orgFilter === "all" || n.organization_id === orgFilter)
//             .map((notice: Notice, index: number) => {
//             const isExpanded = !!expandedMap[notice.id];
//             return (
//               <motion.div
//                 key={notice.id}
//                 variants={itemVariants}
//                 initial="hidden"
//                 animate="visible"
//                 transition={{ delay: index * 0.05 + 0.06 }}
//               >
//                 <Card className="group bg-white rounded-2xl border border-[#e5e7eb] shadow-xl relative overflow-hidden ring-1 ring-transparent hover:ring-[#3a6a8d]/20 transition-all duration-300 transform-gpu hover:-translate-y-0.5 hover:shadow-2xl p-3 sm:p-6">
//                   <div className="absolute inset-0 bg-gradient-to-r from-[#3a6a8d]/10 via-transparent to-[#5e9c8d]/10 opacity-0 group-hover:opacity-100 transition-opacity duration-300" />
//                   <CardContent className="relative z-10 p-0">
//                     <div className="mb-3 flex items-center gap-3">
//                       <div className="w-10 h-10 rounded-xl flex items-center justify-center" style={{ backgroundColor: '#e6f0f5' }}>
//                         <FileText className="w-5 h-5" style={{ color: '#3a6a8d' }} />
//                       </div>
//                       <h2 className="text-lg font-semibold text-[#111827]">{notice.title}</h2>
//                     </div>

//                     <div className="flex flex-wrap items-center gap-x-6 gap-y-2 text-sm mb-3 text-[#4b5563]">
//                       <span className="flex items-center gap-2">
//                         <Calendar className="w-4 h-4" />
//                         <span className="hidden sm:inline">Published:</span>
//                         <span className="font-medium text-[#111827]">{formatPublished(notice.created_at)}</span>
//                       </span>
//                       <span className="flex items-center gap-2">
//                         <Building2 className="w-4 h-4" />
//                         <span className="hidden sm:inline">Organization:</span>
//                         <span className="font-semibold text-[#111827]">{orgNameMap[notice.organization_id || ""] || "Organization"}</span>
//                       </span>
//                     </div>

//                     <p className={`text-[#374151] mb-2 ${isExpanded ? "" : "line-clamp-2"}`}>{notice.content}</p>

//                     <div className="flex justify-end pt-2">
//                       <Button
//                         variant="ghost"
//                         size="sm"
//                         className="text-[#3A6A8D] hover:bg-[#3A6A8D]/10 hover:text-[#2d5470]"
//                         onClick={() =>
//                           setExpandedMap((prev) => ({ ...prev, [notice.id]: !prev[notice.id] }))
//                         }
//                       >
//                         {isExpanded ? "View Less" : "View More"}
//                       </Button>
//                     </div>
//                   </CardContent>
//                 </Card>
//               </motion.div>
//             );
//           })}
//         {!noticesLoading && !noticesError && filteredNotices.length === 0 && (
//           <motion.div
//             variants={itemVariants}
//             initial="hidden"
//             animate="visible"
//             transition={{ delay: 0.3 }}
//             className="text-gray-600"
//           >
//             No notices found.
//           </motion.div>
//         )}
//       </div>
//   </div>
//     </motion.div>
//   );

// }


"use client";

import { useState, useEffect, useMemo } from "react";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Calendar, FileText, Search, Building2 } from "lucide-react";
import { useGetNoticesQuery } from "@/app/store/slices/noticesSlice";
// import { useRouter } from "next/navigation";
import { motion, useReducedMotion } from "framer-motion";
import { useLazyGetOrgQuery } from "@/app/store/slices/orgsApi";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { useTranslation } from "react-i18next";

type Notice = {
  id: string | number;
  title: string;
  content: string;
  tags?: string[];
  created_at: string;
  organization_id?: string;
};

export default function NoticesPage() {
  const { t } = useTranslation("user");
  const [expandedMap, setExpandedMap] = useState<Record<string | number, boolean>>({});
  const [searchInput, setSearchInput] = useState("");
  const [searchQuery, setSearchQuery] = useState("");
  const [orgFilter, setOrgFilter] = useState<string>("all");

  const { data: apiNotices, isLoading: noticesLoading, isError: noticesError } =
    useGetNoticesQuery({ page: 1, limit: 10 });

  useEffect(() => {
    if (apiNotices) console.log("Notices API response:", apiNotices);
  }, [apiNotices]);

  const notices = useMemo(() => {
    if (!Array.isArray(apiNotices?.data)) {
      console.error("Notices data is not an array:", apiNotices);
      return [];
    }
    return apiNotices.data.map((n, idx) => {
      if (!n) {
        console.error("Notice is undefined at index:", idx);
        return null;
      }
      return {
        ...n,
        title: n.title ?? t("notices.default_title"),
        content: n.content ?? "",
        tags: Array.isArray(n.tags) ? n.tags.map((t) => String(t)) : [],
        created_at: n.created_at ?? "",
        organization_id: n.organization_id ?? "",
      };
    }).filter((n): n is NonNullable<typeof n> => n !== null);
  }, [apiNotices, t]);

  const [triggerGetOrg] = useLazyGetOrgQuery();
  const [orgNameMap, setOrgNameMap] = useState<Record<string, string>>({});

  useEffect(() => {
    const uniqueOrgIds = Array.from(
      new Set(notices.map((n) => n.organization_id).filter((v): v is string => typeof v === 'string' && !!v))
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
            if (orgId) entries.push([orgId, t("notices.default_organization")]);
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
  }, [notices, orgNameMap, triggerGetOrg, t]);

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
        orgNameMap[n.organization_id || ""] || t("notices.default_organization"),
      ]
        .join(" ")
        .toLowerCase();
      return terms.every((t) => haystack.includes(t));
    });
  }, [notices, searchQuery, orgNameMap, t]);

  const prefersReducedMotion = useReducedMotion();

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
    if (!iso) return t("notices.unknown_date");
    const d = new Date(iso);
    if (isNaN(d.getTime())) return t("notices.unknown_date");
    return d.toLocaleString(undefined, {
      year: "numeric",
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  return (
    <motion.div className="min-h-screen w-full bg-gray-50 relative overflow-hidden flex flex-col">
      <div className="absolute inset-0 overflow-hidden pointer-events-none">
        <div className="absolute -top-24 -right-24 w-56 h-56 rounded-full blur-3xl" style={{ background: 'radial-gradient(closest-side, rgba(167,179,185,0.10), rgba(167,179,185,0))' }} />
        <div className="absolute -bottom-28 -left-28 w-64 h-64 rounded-full blur-3xl" style={{ background: 'radial-gradient(closest-side, rgba(94,156,141,0.10), rgba(94,156,141,0))' }} />
      </div>
      <div className="relative z-10 max-w-7xl w-full mx-auto px-4 sm:px-6 md:px-8">
        <motion.div
          className="bg-white border border-gray-100 rounded-xl p-4 sm:p-5 mb-4 sm:mb-6 shadow-sm"
          variants={headerVariants}
          initial="hidden"
          animate="visible"
        >
          <div>
            <div className="flex items-center gap-3 mb-1">
              <FileText className="h-7 w-7" style={{ color: '#3a6a8d' }} />
              <h1 className="text-xl leading-snug sm:text-2xl font-bold text-[#111827]">{t("notices.title")}</h1>
            </div>
            <p className="text-[#4b5563] text-sm">{t("notices.description")}</p>
          </div>
        </motion.div>

        <Card className="bg-white border border-gray-100 rounded-xl p-4 mb-4 sm:mb-6 shadow-sm">
          <div className="flex flex-col gap-4 w-full mb-2 sm:flex-row">
            <div className="relative flex-1 flex">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
              <input
                type="text"
                placeholder={t("notices.search.placeholder")}
                value={searchInput}
                onChange={(e) => setSearchInput(e.target.value)}
                className="w-full pl-10 pr-4 py-2 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-[#3A6A8D] focus:border-transparent"
              />
              <Button
                type="button"
                className="ml-2 px-4 py-2 bg-[#3A6A8D] hover:bg-[#2d5470] text-white"
                onClick={() => setSearchQuery(searchInput)}
              >
                {t("notices.search.button")}
              </Button>
            </div>
            <div className="w-full sm:w-64">
              <Select value={orgFilter} onValueChange={setOrgFilter}>
                <SelectTrigger className="w-full border-[#3A6A8D] text-[#3A6A8D] focus:ring-2 focus:ring-[#3A6A8D] focus:border-transparent">
                  <SelectValue placeholder={t("notices.filter.all")} />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">{t("notices.filter.all")}</SelectItem>
                  {Array.from(
                    new Map(
                      notices
                        .map((n) => n.organization_id)
                        .filter(Boolean)
                        .map((id) => [id!, orgNameMap[id!] || t("notices.default_organization")])
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

        <div className="space-y-5 pb-8">
          {noticesLoading && <div className="text-gray-600">{t("notices.loading")}</div>}
          {noticesError && <div className="text-red-600">{t("notices.errors.load_failed")}</div>}
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
                    <Card className="group bg-white rounded-2xl border border-[#e5e7eb] shadow-xl relative overflow-hidden ring-1 ring-transparent hover:ring-[#3a6a8d]/20 transition-all duration-300 transform-gpu hover:-translate-y-0.5 hover:shadow-2xl p-3 sm:p-6">
                      <div className="absolute inset-0 bg-gradient-to-r from-[#3a6a8d]/10 via-transparent to-[#5e9c8d]/10 opacity-0 group-hover:opacity-100 transition-opacity duration-300" />
                      <CardContent className="relative z-10 p-0">
                        <div className="mb-3 flex items-center gap-3">
                          <div className="w-10 h-10 rounded-xl flex items-center justify-center" style={{ backgroundColor: '#e6f0f5' }}>
                            <FileText className="w-5 h-5" style={{ color: '#3a6a8d' }} />
                          </div>
                          <h2 className="text-lg font-semibold text-[#111827]">{notice.title}</h2>
                        </div>

                        <div className="flex flex-wrap items-center gap-x-6 gap-y-2 text-sm mb-3 text-[#4b5563]">
                          <span className="flex items-center gap-2">
                            <Calendar className="w-4 h-4" />
                            <span className="hidden sm:inline">{t("notices.published")}</span>
                            <span className="font-medium text-[#111827]">{formatPublished(notice.created_at)}</span>
                          </span>
                          <span className="flex items-center gap-2">
                            <Building2 className="w-4 h-4" />
                            <span className="hidden sm:inline">{t("notices.organization")}</span>
                            <span className="font-semibold text-[#111827]">{orgNameMap[notice.organization_id || ""] || t("notices.default_organization")}</span>
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
                            {isExpanded ? t("notices.actions.view_less") : t("notices.actions.view_more")}
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
              {t("notices.empty")}
            </motion.div>
          )}
        </div>
      </div>
    </motion.div>
  );
}

