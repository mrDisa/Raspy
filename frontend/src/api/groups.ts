import { apiRequest } from "./client";
import type { Group } from "../types/api";

export async function getGroups(query: string): Promise<Group[]> {
  const params = new URLSearchParams({
    q: query,
  });

  return apiRequest<Group[]>(
    `/api/v1/groups?${params.toString()}`,
  );
}