import { useEffect, useState } from "react";

import { getMe, updateGroup } from "./api/user";
import type { User } from "./types/api";

import SelectGroup from "./pages/SelectGroup/SelectGroup";
import Schedule from "./pages/Schedule/Schedule";

function App() {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function loadUser() {
      try {
        const data = await getMe();
        setUser(data);
      } catch (err) {
        setError(
          err instanceof Error ? err.message : "Не удалось загрузить пользователя",
        );
      } finally {
        setLoading(false);
      }
    }

    loadUser();
  }, []);

  async function handleGroupSelect(groupId: number, subgroup: number) {
    const updatedUser = await updateGroup(groupId, subgroup);
    setUser(updatedUser);
  }

  if (loading) {
    return <div className="app-state">Загрузка...</div>;
  }

  if (error) {
    return <div className="app-state error">{error}</div>;
  }

  if (!user) {
    return null;
  }

  return (
    <main className="app">
      {user.GroupID === null ? (
        <SelectGroup onSelect={handleGroupSelect} />
      ) : (
        <Schedule user={user} />
      )}
    </main>
  );
}

export default App;