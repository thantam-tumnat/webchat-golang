CREATE TABLE IF NOT EXISTS messages (
    id         BIGSERIAL PRIMARY KEY,
    room_id    BIGINT NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    user_id    BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content    TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- index ช่วยให้ query ข้อความในห้อง (เรียงตามเวลา) เร็วขึ้น
CREATE INDEX IF NOT EXISTS idx_messages_room_created ON messages (room_id, created_at DESC);
