-- name: CreateOrUpdateCounterStatistics :batchexec
INSERT INTO counter_statistics (banner_id, timestamp_from, timestamp_to, COUNT)
VALUES ($1, $2, $3, $4)
ON CONFLICT (timestamp_from, timestamp_to) DO UPDATE
    SET COUNT = counter_statistics.count + EXCLUDED.count;

-- name: GetCounterStatistics :many
SELECT *
FROM counter_statistics
WHERE banner_id = $1
  AND timestamp_from >= $2
  AND timestamp_to <= $3;

-- name: CreateBanner :exec
INSERT INTO banners (name)
VALUES ($1);

-- name: GetBanner :one
SELECT id
FROM banners
WHERE id = $1;