TRUNCATE TABLE verification_codes CASCADE;
TRUNCATE TABLE users CASCADE;

INSERT INTO users (
    username,
    name,
    email,
    password_hash,
    email_verified,
    created_at,
    updated_at
)
SELECT 
    'user' || generate_series::TEXT,
    'Test User ' || generate_series,
    'user' || generate_series::TEXT || '@example.com',
    '$2a$10$ZKzGMGqU9Lz2.Nb0MYxgXu.KVJcZp.RXm5hPtzYE3ZBEdvnkM5NFi',
    CASE WHEN random() > 0.5 THEN true ELSE false END,
    NOW() - (random() * interval '90 days'),
    NOW() - (random() * interval '30 days')
FROM generate_series(1, 50);

INSERT INTO verification_codes (
    user_id,
    code,
    type,
    expires_at,
    created_at
)
SELECT 
    id,
    substring(md5(random()::text) from 1 for 6),
    'email_verification',
    NOW() + interval '24 hours',
    NOW()
FROM users 
WHERE NOT email_verified;
