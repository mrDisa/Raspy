export interface User {
  ID: number;
  TelegramID: number;
  GroupID: number | null;
  Subgroup: number;
  NotificationsEnabled: boolean;
  CreatedAt: string;
  UpdatedAt: string;
}

export interface Group {
  ID: number;
  ExternalID: string;
  Name: string;
}

export interface Lesson {
  disciplines: string;
  types: string;
  timeStart: string;
  timeEnd: string;
  number: number;
  auditorium: string;
  corpus: string;
  teachers: string;
  subgroup: number;
}

export interface ScheduleDay {
  date: string;
  list: Lesson[];
}