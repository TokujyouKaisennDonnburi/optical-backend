
-- +migrate Up

-- TODO: 期限切れレコードの自動削除を検討
-- 方法: Goアプリ内スケジューラ or 外部スケジューラ (GitHub Actions等)
-- クエリ: DELETE FROM calendar_invitations WHERE expires_at < NOW()

CREATE TABLE calendar_invitations (
    id UUID PRIMARY KEY,
    calendar_id UUID NOT NULL REFERENCES calendars(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    joined_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    token UUID NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_calendar_invitations_calendar_id ON calendar_invitations(calendar_id);

COMMENT ON COLUMN calendar_invitations.calendar_id IS '対象カレンダー無しで自動削除';
COMMENT ON COLUMN calendar_invitations.email IS '招待先メアド(記録用) ';
COMMENT ON COLUMN calendar_invitations.joined_user_id IS 'メアド違いで登録しても追跡するため';
COMMENT ON COLUMN calendar_invitations.token IS '個別トークンで招待者を識別';
COMMENT ON COLUMN calendar_invitations.expires_at IS 'トークンの有効期限';
COMMENT ON COLUMN calendar_invitations.used_at IS '参加日時';

-- +migrate Down
DROP TABLE IF EXISTS calendar_invitations;
