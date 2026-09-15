import { useEffect, useMemo, useState } from "react";
import { getGroups } from "../../api/groups";
import type { Group } from "../../types/api";

interface SelectGroupProps {
  onSelect: (groupId: number, subgroup: number) => Promise<void>;
}

function SelectGroup({ onSelect }: SelectGroupProps) {
  const [groups, setGroups] = useState<Group[]>([]);
  const [search, setSearch] = useState("");
  const [selectedGroup, setSelectedGroup] = useState<Group | null>(null);
  const [subgroup, setSubgroup] = useState(0);

  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const query = search.trim();

    if (!query) {
        setGroups([]);
        return;
    }

    const timeout = setTimeout(async () => {
        try {
        setLoading(true);
        setError(null);

        const data = await getGroups(query);
        setGroups(data);
        } catch (err) {
        setError(
            err instanceof Error
            ? err.message
            : "Не удалось загрузить группы",
        );
        } finally {
        setLoading(false);
        }
    }, 300);

    return () => clearTimeout(timeout);
    }, [search]);

  const filteredGroups = useMemo(() => {
    const query = search.trim().toLowerCase();

    if (!query) {
      return [];
    }

    return groups
      .filter((group) =>
        group.Name.toLowerCase().includes(query),
      )
      .slice(0, 8);
  }, [groups, search]);

  function handleGroupClick(group: Group) {
    setSelectedGroup(group);
    setSearch("");
  }

  function handleChangeGroup() {
    setSelectedGroup(null);
    setSubgroup(0);
  }

  async function handleContinue() {
    if (!selectedGroup) {
      return;
    }

    try {
      setSaving(true);
      setError(null);

      await onSelect(selectedGroup.ID, subgroup);
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : "Не удалось сохранить группу",
      );
    } finally {
      setSaving(false);
    }
  }

  return (
    <section className="select-group">
      <div className="brand">
        <div className="brand-icon">R</div>

        <div>
          <h1>Raspy</h1>
          <p>Твоё расписание. Без рутины.</p>
        </div>
      </div>

      <div className="select-header">
        <span className="eyebrow">НАСТРОЙКА</span>

        <h2>Выбери свою группу</h2>

        <p>
          Мы запомним её, поэтому в следующий раз
          расписание откроется сразу.
        </p>
      </div>

      {!selectedGroup ? (
        <>
          <div className="search">
            <span>⌕</span>

            <input
              type="text"
              placeholder="Поиск группы"
              value={search}
              onChange={(event) => setSearch(event.target.value)}
              autoFocus
            />
          </div>

          {loading && (
            <div className="list-state">
              Загружаем группы...
            </div>
          )}

          {!loading && search.trim() && (
            <div className="group-list">
              {filteredGroups.length === 0 ? (
                <div className="list-state">
                  Ничего не найдено
                </div>
              ) : (
                filteredGroups.map((group) => (
                  <button
                    key={group.ID}
                    className="group-item"
                    onClick={() => handleGroupClick(group)}
                  >
                    <span>{group.Name}</span>
                    <span className="arrow">→</span>
                  </button>
                ))
              )}
            </div>
          )}
        </>
      ) : (
        <>
          <div className="selected-group">
            <div>
              <span className="selected-label">
                ГРУППА
              </span>

              <strong>{selectedGroup.Name}</strong>
            </div>

            <button
              className="change-button"
              onClick={handleChangeGroup}
            >
              Изменить
            </button>
          </div>

          <div className="subgroup">
            <div className="subgroup-header">
              <span>Подгруппа</span>

              <span className="subgroup-hint">
                Если у вас их несколько
              </span>
            </div>

            <div className="segmented">
              {[0, 1, 2].map((value) => (
                <button
                  key={value}
                  className={subgroup === value ? "active" : ""}
                  onClick={() => setSubgroup(value)}
                >
                  {value === 0 ? "Вся группа" : value}
                </button>
              ))}
            </div>
          </div>
        </>
      )}

      {error && (
        <div className="form-error">
          {error}
        </div>
      )}

      <button
        className="continue-button"
        disabled={!selectedGroup || saving}
        onClick={handleContinue}
      >
        {saving ? "Сохраняем..." : "Продолжить"}
      </button>
    </section>
  );
}

export default SelectGroup;