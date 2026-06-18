# Docker command 
#กันตัวเองลืม
รันแค่ Postgres (default)
```bash
docker compose up -d
```

รัน Postgres + backend (profile app)
```bash
docker compose --profile app up -d #เพราะตรงตั้ง profile ของ backend ใน compose ว่า app


