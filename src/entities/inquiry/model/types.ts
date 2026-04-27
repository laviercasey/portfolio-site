export type InquiryType = 'freelance' | 'fulltime' | 'collaboration' | 'other';
export type InquiryStatus = 'new' | 'read' | 'replied' | 'archived';

export interface Inquiry {
  id: string;
  name: string;
  email: string;
  company?: string;
  telegram?: string;
  type: InquiryType;
  budget?: string;
  message: string;
  status: InquiryStatus;
  adminNotes?: string;
  createdAt: string;
  updatedAt: string;
}
