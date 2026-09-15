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

export async function updateNotifications(
  enabled: boolean,
): Promise<User> {
  return apiRequest<User>(
    "/api/v1/me/notifications",
    {
      method: "PUT",
      body: JSON.stringify({
        enabled,
      }),
    },
  );
}