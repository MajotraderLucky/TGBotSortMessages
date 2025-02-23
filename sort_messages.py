from telethon import TelegramClient, events
import datetime
import asyncio
import traceback
import os
from dotenv import load_dotenv

# Загружаем переменные окружения из .env файла
load_dotenv()

# Получаем API ID и API Hash из переменных окружения
api_id = int(os.getenv("API_ID"))
api_hash = os.getenv("API_HASH")

# Ключевые слова для важности
important_keywords = ["срочно", "важно", "error", "ошибка", "внимание", "Здравствуйте", "Сергей", "добрый"]

# Подключаемся к Telegram
client = TelegramClient("my_session", api_id, api_hash)

# Кешируем контакты и архивированные чаты
cached_contacts = set()
cached_archived_chats = set()

async def update_contacts():
    """Обновляем кеш контактов"""
    global cached_contacts
    try:
        dialogs = await client.get_dialogs(limit=200)
        cached_contacts = {dialog.entity.id for dialog in dialogs if dialog.is_user}
        print(f"🔄 Контакты обновлены: {len(cached_contacts)} пользователей")
    except Exception as e:
        await report_error("Ошибка при обновлении контактов", e)

async def update_archived_chats():
    """Обновляем список архивированных чатов (игнорируем их)"""
    global cached_archived_chats
    try:
        dialogs = await client.get_dialogs(archived=True, limit=200)
        cached_archived_chats = {dialog.entity.id for dialog in dialogs}
        print(f"🗂 Архивированные чаты обновлены: {len(cached_archived_chats)} чатов")
    except Exception as e:
        await report_error("Ошибка при обновлении архивированных чатов", e)

async def report_error(context, error):
    """Отправка ошибки в 'Избранное'"""
    error_msg = f"🚨 Ошибка в боте ({context}):\n{error}\n\n{traceback.format_exc()}"
    print(error_msg)
    with open("errors.log", "a", encoding="utf-8") as error_log:
        error_log.write(f"{datetime.datetime.now()} - {context}\n{error_msg}\n\n")
    try:
        await client.send_message("me", error_msg)
    except Exception as e:
        print(f"⚠ Ошибка при отправке уведомления об ошибке: {e}")

@client.on(events.NewMessage)
async def new_message_handler(event):
    """Обрабатываем только личные сообщения (PM), игнорируем каналы, группы и архив"""
    try:
        sender = await event.get_sender()

        if event.is_group or event.is_channel:
            print(f"🛑 Игнорируем сообщение из группы/канала: {event.chat_id}")
            return

        if sender.id in cached_archived_chats:
            print(f"📦 Игнорируем сообщение из архива: {sender.id}")
            return

        message_text = event.text.lower()
        timestamp = datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S")
        sender_name = sender.first_name if sender else "Неизвестный"

        is_known = sender.id in cached_contacts

        with open("messages.log", "a", encoding="utf-8") as log_file:
            log_file.write(f"[{timestamp}] {sender_name}: {event.text}\n")

        print(f"📩 Новое личное сообщение от {sender_name}: {event.text}")

        await event.mark_read()

        if any(keyword in message_text for keyword in important_keywords):
            print(f"⚠️ ВАЖНОЕ сообщение: {event.text}")
            await client.forward_messages("me", event.message)
            await client.send_message("me", f"⚠️ ВАЖНОЕ сообщение от {sender_name}: {event.text}")

        if not is_known:
            print(f"🚨 Неизвестный отправитель! {sender_name} ({sender.id})")
            with open("unknown_senders.log", "a", encoding="utf-8") as unknown_log:
                unknown_log.write(f"[{timestamp}] {sender_name} ({sender.id}): {event.text}\n")
            await client.send_message("me", f"🚨 Неизвестный отправитель: {sender_name} ({sender.id})\nСообщение: {event.text}")

    except Exception as e:
        await report_error("Ошибка в обработке личного сообщения", e)

async def main():
    """Основная функция, обновляет кеш контактов и архивированных чатов каждые 10 минут"""
    try:
        await update_contacts()
        await update_archived_chats()

        while True:
            await asyncio.sleep(600)
            await update_contacts()
            await update_archived_chats()

    except Exception as e:
        await report_error("Ошибка в основном цикле", e)

with client:
    client.loop.run_until_complete(main())
    client.run_until_disconnected()

