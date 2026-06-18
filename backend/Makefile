# คำสั่งลัดสำหรับงานที่ทำบ่อย (ใช้ได้ถ้ามี make ติดตั้ง)
# บน Windows ถ้าไม่มี make ใช้คำสั่งเต็มใน README ได้เลย

.PHONY: db-up db-down run build tidy

# เปิด postgres + adminer ด้วย docker
db-up:
	docker compose up -d

# ปิด database (ข้อมูลยังอยู่ใน volume)
db-down:
	docker compose down

# รัน backend
run:
	go run .

# build เป็น binary
build:
	go build -o bin/api .

tidy:
	go mod tidy
