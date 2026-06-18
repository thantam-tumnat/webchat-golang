# ---------- build stage ----------
# คอมไพล์ Go เป็น binary ตัวเดียว ในภาพที่มี toolchain ครบ
FROM golang:1.25-alpine AS build

WORKDIR /app

# คัดลอกเฉพาะ go.mod/go.sum ก่อน เพื่อ cache layer download dependency
# (ถ้า dependency ไม่เปลี่ยน Docker จะไม่ download ใหม่ทุกครั้งที่แก้โค้ด)
COPY go.mod go.sum ./
RUN go mod download

# คัดลอกโค้ดที่เหลือแล้ว build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/api

# ---------- run stage ----------
# ภาพเล็ก ๆ สำหรับรันจริง ใส่แค่ binary + ca-certificates (เผื่อต่อ DB ผ่าน SSL)
FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=build /server /server

EXPOSE 8080
CMD ["/server"]
