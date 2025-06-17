CREATE TABLE IF NOT EXISTS chats (
    id BIGSERIAL PRIMARY KEY,
    is_group BOOLEAN NOT NULL DEFAULT FALSE,
    chat_key TEXT UNIQUE
);

CREATE TABLE IF NOT EXISTS chat_members (
    chat_id BIGINT NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL,
    PRIMARY KEY (chat_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_chat_members_user_id ON chat_members(user_id);

CREATE INDEX IF NOT EXISTS idx_chat_members_chat_id ON chat_members(chat_id);

CREATE UNIQUE INDEX IF NOT EXISTS idx_chats_chat_key ON chats(chat_key);