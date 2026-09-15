import { apiRequest } from "./client";
import type { ScheduleDay } from "../types/api";

export async function getTodaySchedule(): Promise<ScheduleDay> {
  return apiRequest<ScheduleDay>("/api/v1/schedule/today");
}

export async function getTomorrowSchedule(): Promise<ScheduleDay> {
  return apiRequest<ScheduleDay>("/api/v1/schedule/tomorrow");
}