import { useEffect, useState } from "react";

import { getMe } from "./api/user";
import type { User } from "./types/api";

function App() {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function loadUser() {
      try {
        const currentUser = await getMe();
        setUser(currentUser);
      } catch (err) {
        setError(
          err instanceof Error
            ? err.message
            : "Unknown error",
        );
      } finally {
        setLoading(false);
      }
    }

    loadUser();
  }, []);

  if (loading) {
    return <div>Загрузка...</div>;
  }

  if (error) {
    return <div>Ошибка: {error}</div>;
  }

  if (!user) {
    return <div>Пользователь не найден</div>;
  }

  return (
    <div>
      <h1>Raspy</h1>

      <p>Telegram ID: {user.TelegramID}</p>

      <p>
        Группа:{" "}
        {user.GroupID === null
          ? "Не выбрана"
          : user.GroupID}
      </p>

      <p>Подгруппа: {user.Subgroup}</p>
    </div>
  );
}

export default App;