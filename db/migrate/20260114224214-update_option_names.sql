
-- +migrate Up
UPDATE options
SET
    name = CASE
        WHEN id = 1 THEN 'pull_request_review_wait_count'
        WHEN id = 2 THEN 'team_review_load'
        WHEN id = 5 THEN 'milestone_progress'
        ELSE name
    END
WHERE id IN (1, 2, 5);

UPDATE options 
SET deprecated = TRUE
WHERE id = 4;

-- +migrate Down
UPDATE options 
SET deprecated = FALSE
WHERE id = 4;

UPDATE options
SET
    name = CASE
        WHEN id = 1 THEN 'pr_review_pending_count'
        WHEN id = 2 THEN 'review_load_status'
        WHEN id = 5 THEN 'milestone_status'
        ELSE name
    END
WHERE id IN (1, 2, 5);
