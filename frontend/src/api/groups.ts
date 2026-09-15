import { apiRequest } from "./client";
import type { Group } from "../types/api";

export async function getGroups(): Promise<Group[]> {
  return apiRequest<Group[]>("/api/v1/groups");
}