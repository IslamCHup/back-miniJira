#!/bin/bash
# Простой HTTP сервер для frontend (альтернатива Python)
# Использование: ./server.sh

PORT=3000

echo "Запуск frontend сервера на порту $PORT"
echo "Откройте в браузере: http://localhost:$PORT/index.html"
echo "Нажмите Ctrl+C для остановки"
echo ""

# Проверяем наличие Python
if command -v python3 &> /dev/null; then
    python3 -m http.server $PORT
elif command -v python &> /dev/null; then
    python -m http.server $PORT
else
    echo "Ошибка: Python не найден"
    echo "Установите Python или используйте другой HTTP сервер"
    exit 1
fi

