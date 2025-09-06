export default interface NoticeUpdateProp {
  id: string;
  orgId: string;
  title: string;
  body: string;
  pinned: boolean;
  createdAt: string;      // ISO timestamp
  effectiveFrom: string;  // ISO timestamp
}
