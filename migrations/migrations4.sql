SELECT cron.schedule(
    'prune-job-run-details',
    '0 12 * * *',
    $$DELETE FROM cron.job_run_details WHERE end_time < now() - interval '7 days'$$
);
