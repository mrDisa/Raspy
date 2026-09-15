import { useEffect, useState } from "react";

import { getMe } from "./api/user";
import type { User } from "./types/api";

import Schedule from "./pages/Schedule/Schedule";

function App() {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function loadUser() {
      try {
        setLoading(true);
        setError(null);

        const data = await getMe();

        setUser(data);
      } catch (err) {
        setError(
          err instanceof Error
            ? err.message
            : "Не удалось загрузить пользователя",
        );
      } finally {
        setLoading(false);
      }
    }

    loadUser();
  }, []);

  if (loading) {
    return (
      <div className="list-state">
        Загружаем...
      </div>
    );
  }

  if (error) {
    return (
      <div className="form-error">
        {error}
      </div>
    );
  }

  if (!user) {
    return (
      <div className="form-error">
        Пользователь не найден
      </div>
    );
  }

  return (
  <div className="app">
    <Schedule
      user={user}
      onUserUpdate={setUser}
    />
  </div>
);
}

export default App;