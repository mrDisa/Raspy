import { useEffect, useState } from "react";

import {
  getWeekSchedule,
  getNextWeekSchedule,
} from "../../api/schedule";

import type { Lesson, ScheduleDay } from "../../types/api";

interface WeekScheduleProps {
  onBack: () => void;
}

type Week = "current" | "next";

function WeekSchedule({ onBack }: WeekScheduleProps) {
  const [week, setWeek] = useState<Week>("current");
  const [schedule, setSchedule] = useState<ScheduleDay[]>([]);
  const [nextWeekAvailable, setNextWeekAvailable] = useState(false);

  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function loadWeek() {
      try {
        setLoading(true);
        setError(null);

        const data =
          week === "current"
            ? await getWeekSchedule()
            : await getNextWeekSchedule();

        setSchedule(data);
      } catch (err) {
        setError(
          err instanceof Error
            ? err.message
            : "Не удалось загрузить расписание",
        );
      } finally {
        setLoading(false);
      }
    }

    loadWeek();
  }, [week]);

  useEffect(() => {
    async function checkNextWeek() {
      try {
        const data = await getNextWeekSchedule();
        setNextWeekAvailable(data.length > 0);
      } catch {
        setNextWeekAvailable(false);
      }
    }

    checkNextWeek();
  }, []);

  return (
    <section className="schedule-page">
      <header className="week-header">
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
            Расписание на неделю
          </div>
        </div>
      </header>

      <div className="week-switcher">
        <button
          className={week === "current" ? "active" : ""}
          onClick={() => setWeek("current")}
        >
          Эта неделя
        </button>

        <button
          className={week === "next" ? "active" : ""}
          disabled={!nextWeekAvailable}
          onClick={() => setWeek("next")}
        >
          Следующая
        </button>
      </div>

      <div className="schedule-heading">
        <span className="eyebrow">
          {week === "current"
            ? "ЭТА НЕДЕЛЯ"
            : "СЛЕДУЮЩАЯ НЕДЕЛЯ"}
        </span>

        <h1>
          {schedule.length > 0
            ? formatWeekTitle(schedule)
            : "Расписание"}
        </h1>
      </div>

      {loading && (
        <div className="list-state">
          Загружаем расписание...
        </div>
      )}

      {error && (
        <div className="form-error">
          {error}
        </div>
      )}

      {!loading && !error && (
        <div className="week-list">
          {schedule.length === 0 ? (
            <div className="empty-state">
              <div className="empty-icon">☕</div>

              <h2>
                Расписание недоступно
              </h2>

              <p>
                На эту неделю пока нет занятий.
              </p>
            </div>
          ) : (
            schedule.map((day) => (
              <WeekDay
                key={day.date}
                day={day}
              />
            ))
          )}
        </div>
      )}

      <button
        className="week-back-button"
        onClick={onBack}
      >
        ← К расписанию
      </button>
    </section>
  );
}

function WeekDay({
  day,
}: {
  day: ScheduleDay;
}) {
  return (
    <section className="week-day">
      <div className="week-day-header">
        <span className="week-day-name">
          {formatDayName(day.date)}
        </span>

        <span className="week-day-date">
          {formatShortDate(day.date)}
        </span>
      </div>

      {day.list.length === 0 ? (
        <div className="week-day-empty">
          Пар нет
        </div>
      ) : (
        <div className="week-day-lessons">
          {day.list.map((lesson) => (
            <div
              className="week-lesson"
              key={`${day.date}-${lesson.number}-${lesson.timeStart}`}
            >
              <div className="week-lesson-time">
                {lesson.timeStart}
              </div>

              <div className="week-lesson-info">
                <strong>
                  {lesson.disciplines}
                </strong>

                <LessonMeta lesson={lesson} />
              </div>
            </div>
          ))}
        </div>
      )}
    </section>
  );
}

function LessonMeta({
  lesson,
}: {
  lesson: Lesson;
}) {
  const values = [
    lesson.types,
    lesson.auditorium,
    lesson.corpus,
  ].filter(Boolean);

  if (values.length === 0) {
    return null;
  }

  return (
    <div className="lesson-meta">
      {values.map((value) => (
        <span key={value}>
          {value}
        </span>
      ))}
    </div>
  );
}

function formatDayName(date: string) {
  const parsed = new Date(date);

  return parsed
    .toLocaleDateString("ru-RU", {
      weekday: "long",
    })
    .replace(/^./, (char) => char.toUpperCase());
}

function formatShortDate(date: string) {
  const parsed = new Date(date);

  return parsed.toLocaleDateString("ru-RU", {
    day: "numeric",
    month: "long",
  });
}

function formatWeekTitle(schedule: ScheduleDay[]) {
  const first = new Date(schedule[0].date);
  const last = new Date(
    schedule[schedule.length - 1].date,
  );

  const firstText = first.toLocaleDateString(
    "ru-RU",
    {
      day: "numeric",
      month: "long",
    },
  );

  const lastText = last.toLocaleDateString(
    "ru-RU",
    {
      day: "numeric",
      month: "long",
    },
  );

  return `${firstText} — ${lastText}`;
}

export default WeekSchedule;