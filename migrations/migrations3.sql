CREATE EXTENSION IF NOT EXISTS pg_cron;

SELECT cron.schedule(
    'cleanup-refresh-tokens',
    '17 3 * * *',
    $$DELETE FROM refresh_tokens WHERE expires_at < NOW()$$
);
