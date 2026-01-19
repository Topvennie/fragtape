-- name: HighlightSegmentCreate :one
INSERT INTO highlight_segments (highlight_id, start_tick, end_tick)
VALUES ($1, $2, $3)
RETURNING id;
