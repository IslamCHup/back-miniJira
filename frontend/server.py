#!/usr/bin/env python3
"""
Простой HTTP сервер для запуска frontend
Использование: python3 server.py
Или: python server.py
"""

import http.server
import socketserver
import os
import sys

PORT = 3000

class MyHTTPRequestHandler(http.server.SimpleHTTPRequestHandler):
    def end_headers(self):
        # Добавляем CORS заголовки
        self.send_header('Access-Control-Allow-Origin', '*')
        self.send_header('Access-Control-Allow-Methods', 'GET, POST, OPTIONS')
        self.send_header('Access-Control-Allow-Headers', 'Content-Type')
        super().end_headers()

    def do_OPTIONS(self):
        self.send_response(200)
        self.end_headers()

if __name__ == "__main__":
    os.chdir(os.path.dirname(os.path.abspath(__file__)))
    
    with socketserver.TCPServer(("", PORT), MyHTTPRequestHandler) as httpd:
        print(f"Frontend сервер запущен на http://localhost:{PORT}")
        print(f"Откройте в браузере: http://localhost:{PORT}/index.html")
        print("Нажмите Ctrl+C для остановки")
        try:
            httpd.serve_forever()
        except KeyboardInterrupt:
            print("\nСервер остановлен")

