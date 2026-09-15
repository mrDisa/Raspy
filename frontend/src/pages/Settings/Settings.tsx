import { useEffect, useState } from "react";

import {
  updateGroup,
  updateNotifications,
} from "../../api/user";

import { getGroups } from "../../api/groups";

import type {
  Group,
  User,
} from "../../types/api";

interface SettingsProps {
  user: User;
  onBack: () => void;
  onUserUpdate: (user: User) => void;
}

function Settings({
  user,
  onBack,
  onUserUpdate,
}: SettingsProps) {
  const [loading, setLoading] =
    useState(false);

  const [error, setError] =
    useState<string | null>(null);

  const [showGroupPicker, setShowGroupPicker] =
    useState(false);

  const [groupQuery, setGroupQuery] =
    useState("");

  const [groups, setGroups] =
    useState<Group[]>([]);

  const [groupsLoading, setGroupsLoading] =
    useState(false);

  useEffect(() => {
    if (!showGroupPicker) {
      return;
    }

    if (groupQuery.trim().length < 2) {
      setGroups([]);
      return;
    }

    const timeout = setTimeout(
      async () => {
        try {
          setGroupsLoading(true);

          const data = await getGroups(
            groupQuery.trim(),
          );

          setGroups(data);
        } catch {
          setGroups([]);
        } finally {
          setGroupsLoading(false);
        }
      },
      300,
    );

    return () => {
      clearTimeout(timeout);
    };
  }, [
    groupQuery,
    showGroupPicker,
  ]);

  async function handleGroupChange(
    group: Group,
  ) {
    if (loading) {
      return;
    }

    try {
      setLoading(true);
      setError(null);

      const updatedUser =
        await updateGroup(
          group.ID,
          user.Subgroup,
        );

      onUserUpdate({
        ...updatedUser,
        Group: group,
      });

      setShowGroupPicker(false);
      setGroupQuery("");
      setGroups([]);
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : "Не удалось изменить группу",
      );
    } finally {
      setLoading(false);
    }
  }

  async function handleSubgroupChange(
    subgroup: number,
  ) {
    if (
      loading ||
      subgroup === user.Subgroup ||
      user.GroupID === null
    ) {
      return;
    }

    try {
      setLoading(true);
      setError(null);

      const updatedUser =
        await updateGroup(
          user.GroupID,
          subgroup,
        );

      onUserUpdate({
        ...updatedUser,
        Group: user.Group,
      });
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : "Не удалось изменить подгруппу",
      );
    } finally {
      setLoading(false);
    }
  }

  async function handleNotificationsToggle() {
    if (loading) {
      return;
    }

    try {
      setLoading(true);
      setError(null);

      const updatedUser =
        await updateNotifications(
          !user.NotificationsEnabled,
        );

      onUserUpdate({
        ...updatedUser,
        Group: user.Group,
      });
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : "Не удалось изменить настройки уведомлений",
      );
    } finally {
      setLoading(false);
    }
  }

  if (showGroupPicker) {
    return (
      <section className="settings-page">
        <header className="settings-header">
          <button
            className="back-button"
            onClick={() =>
              setShowGroupPicker(false)
            }
          >
            ←
          </button>

          <div>
            <div className="schedule-brand">
              Raspy
            </div>

            <div className="schedule-group">
              Выбор группы
            </div>
          </div>
        </header>

        <div className="group-picker">
          <input
            className="group-search"
            value={groupQuery}
            onChange={(event) =>
              setGroupQuery(
                event.target.value,
              )
            }
            placeholder="Введите группу"
            autoFocus
          />

          {groupsLoading && (
            <div className="list-state">
              Ищем группы...
            </div>
          )}

          {!groupsLoading &&
            groupQuery.trim().length >= 2 &&
            groups.length === 0 && (
              <div className="list-state">
                Группы не найдены
              </div>
            )}

          <div className="group-results">
            {groups.map((group) => (
              <button
                key={group.ID}
                className="group-result"
                disabled={loading}
                onClick={() =>
                  handleGroupChange(
                    group,
                  )
                }
              >
                {group.Name}
              </button>
            ))}
          </div>
        </div>
      </section>
    );
  }

  return (
    <section className="settings-page">
      <header className="settings-header">
        <button
          className="back-button"
          onClick={onBack}
        >
          ←
        </button>

        <div>
          <div className="schedule-brand">
            Raspy
          </div>

          <div className="schedule-group">
            Настройки
          </div>
        </div>
      </header>

      {error && (
        <div className="form-error">
          {error}
        </div>
      )}

      <div className="settings-list">
        <div className="settings-section">
          <div className="settings-section-title">
            Группа
          </div>

          <button
            className="settings-row"
            disabled={loading}
            onClick={() =>
              setShowGroupPicker(true)
            }
          >
            <div>
              <span className="settings-row-title">
                Группа
              </span>

              <span className="settings-row-value">
                {user.Group?.Name ??
                  "Не выбрана"}
              </span>
            </div>

            <span className="settings-row-arrow">
              →
            </span>
          </button>

          <div className="settings-section-title">
            Подгруппа
          </div>

          <div className="settings-subgroups">
            <button
              className={
                user.Subgroup === 0
                  ? "settings-subgroup active"
                  : "settings-subgroup"
              }
              disabled={loading}
              onClick={() =>
                handleSubgroupChange(0)
              }
            >
              Без подгруппы
            </button>

            <button
              className={
                user.Subgroup === 1
                  ? "settings-subgroup active"
                  : "settings-subgroup"
              }
              disabled={loading}
              onClick={() =>
                handleSubgroupChange(1)
              }
            >
              1-я подгруппа
            </button>

            <button
              className={
                user.Subgroup === 2
                  ? "settings-subgroup active"
                  : "settings-subgroup"
              }
              disabled={loading}
              onClick={() =>
                handleSubgroupChange(2)
              }
            >
              2-я подгруппа
            </button>
          </div>
        </div>

        <div className="settings-section">
          <div className="settings-section-title">
            Уведомления
          </div>

          <div className="settings-row">
            <div>
              <span className="settings-row-title">
                Изменения расписания
              </span>

              <span className="settings-row-value">
                Получать уведомления в Telegram
              </span>
            </div>

            <button
              className={[
                "settings-toggle",
                user.NotificationsEnabled
                  ? "active"
                  : "",
              ]
                .filter(Boolean)
                .join(" ")}
              disabled={loading}
              onClick={
                handleNotificationsToggle
              }
            >
              <span />
            </button>
          </div>
        </div>
      </div>
    </section>
  );
}

export default Settings;