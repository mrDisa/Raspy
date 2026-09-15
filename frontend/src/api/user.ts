import { apiRequest } from "./client";
import type { User } from "../types/api";

export async function getMe(): Promise<User> {
  return apiRequest<User>("/api/v1/me");
}

export async function updateGroup(
  groupId: number,
  subgroup: number,
): Promise<User> {
  return apiRequest<User>("/api/v1/me/group", {
    method: "PUT",
    body: JSON.stringify({
      group_id: groupId,
      subgroup,
    }),
  });
}