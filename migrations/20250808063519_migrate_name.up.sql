CREATE TABLE debt_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount BIGINT NOT NULL, -- + qo‘shildi, - kamaydi
    reason TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
