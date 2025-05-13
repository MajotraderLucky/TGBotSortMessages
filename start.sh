#!/bin/bash
set -e

echo "🚀 Запускаем Telegram Message Sorter..."

# Пример запуска с параметрами
go run cmd/app/main.go \
  --messages=10 \
  --interval=120 \
  --formats=json,markdown,byuser \
  --debug=true

# Другие варианты запуска (закомментированы):
# Базовый запуск без дополнительных параметров:
# go run cmd/app/main.go

# Получение большего количества сообщений с меньшим интервалом:
# go run cmd/app/main.go --messages=20 --interval=30

# Сохранение только в определенных форматах:
# go run cmd/app/main.go --formats=json,structured

# Запуск в фоновом режиме (вывод в log файл):
# nohup go run cmd/app/main.go > app.log 2>&1 & 