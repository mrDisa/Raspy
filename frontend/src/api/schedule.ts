import { apiRequest } from "./client";
import type { ScheduleDay } from "../types/api";

export async function getTodaySchedule(): Promise<ScheduleDay> {
  return apiRequest<ScheduleDay>("/api/v1/schedule/today");
}

export async function getTomorrowSchedule(): Promise<ScheduleDay> {
  return apiRequest<ScheduleDay>("/api/v1/schedule/tomorrow");
}

export async function getWeekSchedule(): Promise<ScheduleDay[]> {
  return apiRequest<ScheduleDay[]>("/api/v1/schedule/week");
}

export async function getNextWeekSchedule(): Promise<ScheduleDay[]> {
  return apiRequest<ScheduleDay[]>("/api/v1/schedule/week/next");
}