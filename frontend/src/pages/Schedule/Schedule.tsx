import WeekSchedule from "../WeekSchedule/WeekSchedule";

import { useEffect, useMemo, useState } from "react";

import { getTodaySchedule, getTomorrowSchedule } from "../../api/schedule";
import type { ScheduleDay, User, Lesson } from "../../types/api";

interface ScheduleProps {
  user: User;
}

type Day = "today" | "tomorrow";

function Schedule({ user }: ScheduleProps) {
  const [day, setDay] = useState<Day>("today");
  const [schedule, setSchedule] = useState<ScheduleDay | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showWeek, setShowWeek] = useState(false);

  useEffect(() => {
    async function loadSchedule() {
      try {
        setLoading(true);
        setError(null);

        const data =
          day === "today"
            ? await getTodaySchedule()
            : await getTomorrowSchedule();

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

    loadSchedule();
  }, [day]);

  const currentLesson = useMemo(() => {
    if (day !== "today" || !schedule) {
      return null;
    }

    return schedule.list.find((lesson) => isLessonCurrent(lesson)) ?? null;
  }, [day, schedule]);

  const nextLesson = useMemo(() => {
    if (day !== "today" || !schedule) {
      return null;
    }

    return (
      schedule.list.find((lesson) => isLessonUpcoming(lesson)) ?? null
    );
  }, [day, schedule]);

  const lessonsFinished = useMemo(() => {
    if (day !== "today" || !schedule || schedule.list.length === 0) {
      return false;
    }

    return schedule.list.every((lesson) => isLessonFinished(lesson));
  }, [day, schedule]);

  if (showWeek) {
    return (
        <WeekSchedule
        onBack={() => setShowWeek(false)}
        />
    );
  }

  return (
    <section className="schedule-page">
      <header className="schedule-top">
        <div>
          <div className="schedule-brand">Raspy</div>

          <div className="schedule-group">
            Группа #{user.GroupID} ·{" "}
            {user.Subgroup === 0
              ? "без подгруппы"
              : `${user.Subgroup}-я подгруппа`}
          </div>
        </div>

        <button className="settings-button">⚙</button>
      </header>

      <div className="day-switcher">
        <button
          className={day === "today" ? "active" : ""}
          onClick={() => setDay("today")}
        >
          Сегодня
        </button>

        <button
          className={day === "tomorrow" ? "active" : ""}
          onClick={() => setDay("tomorrow")}
        >
          Завтра
        </button>
      </div>

      <div className="schedule-heading">
        <span className="eyebrow">
          {day === "today" ? "СЕГОДНЯ" : "ЗАВТРА"}
        </span>

        <h1>
          {schedule
            ? formatDate(schedule.date)
            : day === "today"
              ? "Сегодня"
              : "Завтра"}
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

      {!loading && !error && schedule && (
        <>
          {day === "today" && schedule.list.length > 0 && (
            <TodayOverview
              currentLesson={currentLesson}
              nextLesson={nextLesson}
              lessonsFinished={lessonsFinished}
            />
          )}

          {schedule.list.length === 0 ? (
            <div className="empty-state">
              <div className="empty-icon">☕</div>

              <h2>Сегодня занятий нет</h2>

            </div>
          ) : (
            <div className="schedule-section">
              <div className="lesson-list">
                {schedule.list.map((lesson) => (
                  <LessonCard
                    key={`${lesson.number}-${lesson.timeStart}`}
                    lesson={lesson}
                    isCurrent={day === "today" && currentLesson === lesson}
                    isToday={day === "today"}
                  />
                ))}
              </div>
            </div>
          )}

          <button
            className="week-button"
            onClick={() => setShowWeek(true)}
            >
            <span>Расписание на неделю</span>
            <span>→</span>
          </button>
        </>
      )}
    </section>
  );
}

interface TodayOverviewProps {
  currentLesson: Lesson | null;
  nextLesson: Lesson | null;
  lessonsFinished: boolean;
}

function TodayOverview({
  currentLesson,
  nextLesson,
  lessonsFinished,
}: TodayOverviewProps) {
  if (lessonsFinished) {
    return (
      <section className="today-overview finished">
        <div className="overview-label">
          СЕГОДНЯ
        </div>

        <h2>Пары на сегодня закончились</h2>
        
      </section>
    );
  }

  return (
    <section className="today-overview">
      {currentLesson ? (
        <div className="overview-block current">
          <div className="overview-label">
            <span className="status-dot" />
            СЕЙЧАС ИДЁТ
          </div>

          <h2>{currentLesson.disciplines}</h2>

          <div className="overview-time">
            {currentLesson.timeStart} — {currentLesson.timeEnd}
          </div>

          <LessonMeta lesson={currentLesson} />
        </div>
      ) : (
        <div className="overview-block break">
          <div className="overview-label">
            СЕЙЧАС
          </div>

          <h2>Перерыв</h2>

          <p>Сейчас занятий нет.</p>
        </div>
      )}

      {nextLesson && (
        <div className="overview-block next">
          <div className="overview-label">
            СЛЕДУЮЩАЯ
          </div>

          <h2>{nextLesson.disciplines}</h2>

          <div className="overview-time">
            {nextLesson.timeStart} — {nextLesson.timeEnd}
          </div>

          <LessonMeta lesson={nextLesson} />
        </div>
      )}
    </section>
  );
}

interface LessonCardProps {
  lesson: Lesson;
  isCurrent: boolean;
  isToday: boolean;
}

function LessonCard({
  lesson,
  isCurrent,
  isToday,
}: LessonCardProps) {
  const finished =
    isToday && !isCurrent && isLessonFinished(lesson);

  return (
    <article
      className={[
        "lesson-card",
        isCurrent ? "current" : "",
        finished ? "finished" : "",
      ]
        .filter(Boolean)
        .join(" ")}
    >
      <div className="lesson-time">
        <strong>{lesson.timeStart}</strong>
        <span>{lesson.timeEnd}</span>
      </div>

      <div className="lesson-info">
        <div className="lesson-number">
          Пара {lesson.number}
        </div>

        <h2>{lesson.disciplines}</h2>

        <LessonMeta lesson={lesson} />
      </div>

      <div className="lesson-status">
        {isCurrent ? "Сейчас" : finished ? "✓" : ""}
      </div>
    </article>
  );
}

function LessonMeta({ lesson }: { lesson: Lesson }) {
  return (
    <div className="lesson-meta">
      {lesson.types && (
        <span>{lesson.types}</span>
      )}

      {lesson.auditorium && (
        <span>{lesson.auditorium}</span>
      )}

      {lesson.corpus && (
        <span>{lesson.corpus}</span>
      )}
    </div>
  );
}

function parseTime(time: string) {
  const [hours, minutes] = time
    .split(":")
    .map(Number);

  return hours * 60 + minutes;
}

function getCurrentMinutes() {
  const now = new Date();

  return now.getHours() * 60 + now.getMinutes();
}

function isLessonCurrent(lesson: Lesson) {
  const now = getCurrentMinutes();

  const start = parseTime(lesson.timeStart);
  const end = parseTime(lesson.timeEnd);

  return now >= start && now < end;
}

function isLessonUpcoming(lesson: Lesson) {
  const now = getCurrentMinutes();

  const start = parseTime(lesson.timeStart);

  return start > now;
}

function isLessonFinished(lesson: Lesson) {
  const now = getCurrentMinutes();

  const end = parseTime(lesson.timeEnd);

  return now >= end;
}

function formatDate(date: string) {
  const parsed = new Date(date);

  if (Number.isNaN(parsed.getTime())) {
    return date;
  }

  return parsed.toLocaleDateString("ru-RU", {
    weekday: "long",
    day: "numeric",
    month: "long",
  });
}

export default Schedule;