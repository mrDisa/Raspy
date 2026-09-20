import asyncio
import os
from urllib.parse import parse_qsl, urlencode, urlsplit, urlunsplit

from aiogram import Bot, Dispatcher, Router
from aiogram.filters import Command, CommandStart
from aiogram.types import (
    BotCommand,
    InlineKeyboardButton,
    InlineKeyboardMarkup,
    Message,
    ReplyKeyboardRemove,
    WebAppInfo,
)
from dotenv import load_dotenv


load_dotenv()

BOT_TOKEN = os.getenv("BOT_TOKEN")
WEB_APP_URL = os.getenv("WEB_APP_URL")


if not BOT_TOKEN:
    raise RuntimeError("BOT_TOKEN is not set")

if not WEB_APP_URL:
    raise RuntimeError("WEB_APP_URL is not set")

if not WEB_APP_URL.startswith("https://"):
    raise RuntimeError("WEB_APP_URL must use HTTPS")


router = Router()


def web_app_url(screen: str | None = None) -> str:
    """Build a Mini App URL while preserving query parameters from .env."""
    if not screen:
        return WEB_APP_URL

    parts = urlsplit(WEB_APP_URL)
    query = dict(parse_qsl(parts.query, keep_blank_values=True))
    query["screen"] = screen
    return urlunsplit(
        (parts.scheme, parts.netloc, parts.path, urlencode(query), parts.fragment)
    )


def mini_app_button(text: str, screen: str | None = None) -> InlineKeyboardButton:
    return InlineKeyboardButton(text=text, web_app=WebAppInfo(url=web_app_url(screen)))


def launch_keyboard(screen: str | None = None) -> InlineKeyboardMarkup:
    labels = {
        None: "📅 Открыть расписание",
    }
    return InlineKeyboardMarkup(
        inline_keyboard=[[mini_app_button(labels[screen], screen)]],
    )


@router.message(CommandStart())
async def start(message: Message):
    await message.answer(
        "Меню обновлено.",
        reply_markup=ReplyKeyboardRemove(),
    )
    await message.answer(
        "Привет! 👋\n\n"
        "Raspy — расписание колледжа без рутины.\n\n"
        "Выбери группу один раз в приложении — мы запомним её для "
        "следующих входов.",
        reply_markup=launch_keyboard(),
    )


@router.message(Command("schedule"))
async def open_schedule(message: Message):
    await message.answer(
        "Открой расписание на сегодня, завтра или неделю.",
        reply_markup=launch_keyboard(),
    )


@router.message(Command("help"))
async def help_command(message: Message):
    await message.answer(
        "Команды Raspy:\n"
        "/schedule — расписание\n"
        "/start — показать меню",
        reply_markup=launch_keyboard(),
    )


async def main():
    bot = Bot(token=BOT_TOKEN)

    dp = Dispatcher()
    dp.include_router(router)

    await bot.delete_webhook(drop_pending_updates=True)
    await bot.set_my_commands(
        [
            BotCommand(command="start", description="Показать главное меню"),
            BotCommand(command="schedule", description="Открыть расписание"),
            BotCommand(command="help", description="Помощь"),
        ]
    )

    print("Bot started")

    await dp.start_polling(bot)


if __name__ == "__main__":
    asyncio.run(main())
