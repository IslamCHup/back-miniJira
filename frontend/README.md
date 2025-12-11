# Frontend Mini-Jira

## Запуск Frontend

### Важно: Frontend должен запускаться через HTTP сервер, а не через file://

### Вариант 1: Python (рекомендуется)

```bash
cd frontend
python3 server.py
```

Или:

```bash
cd frontend
python server.py
```

Затем откройте в браузере: http://localhost:3000/index.html

### Вариант 2: Python HTTP Server (простой)

```bash
cd frontend
python3 -m http.server 3000
```

### Вариант 3: Bash скрипт

```bash
cd frontend
./server.sh
```

### Вариант 4: Node.js (если установлен)

```bash
cd frontend
npx http-server -p 3000
```

## Запуск Backend

**ВАЖНО:** Убедитесь, что backend сервер запущен и перезапущен после добавления CORS:

```bash
# В корне проекта
go run cmd/mini-jira/main.go
```

Сервер должен быть доступен на: http://localhost:8080

## Решение проблем

### "Ошибка соединения с сервером"

1. **Убедитесь, что backend запущен:**
   ```bash
   curl http://localhost:8080/auth/register
   ```
   Должен вернуть ошибку 400 (не 404 и не "connection refused")

2. **Убедитесь, что frontend запущен через HTTP:**
   - ❌ НЕ открывайте `file:///path/to/index.html`
   - ✅ Откройте `http://localhost:3000/index.html`

3. **Проверьте CORS:**
   ```bash
   curl -X OPTIONS http://localhost:8080/auth/register \
     -H "Origin: http://localhost:3000" \
     -H "Access-Control-Request-Method: POST" \
     -v
   ```
   Должны быть заголовки `Access-Control-Allow-Origin`

4. **Перезапустите backend** после изменений в коде

### SMTP ошибки при регистрации

Если видите ошибку "failed to send verification email", это нормально для разработки. Пользователь все равно создается в БД. Вы можете:
- Вручную подтвердить email через БД
- Или настроить SMTP правильно

