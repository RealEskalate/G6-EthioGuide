export interface ChatCitation {
  id: string;
  type: "procedure" | string;
}

export interface ChatChecklistItem {
  checklistID: string;
  done: boolean;
  text: string;
}

export interface ChatGuideResponse {
  citations: ChatCitation[];
  response: string;
  source: "official" | string;
  verified: boolean;
  // added: optional checklists from backend
  checklists?: ChatChecklistItem[];
}

export interface PostChatRequest {
  message: string;
}

// added: history item type
export interface ChatHistoryItem {
  id: string;
  title: string;
  lastMessage: string;
  timestamp: string;
  messageCount: number;
}
