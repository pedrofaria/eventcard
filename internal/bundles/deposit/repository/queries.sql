-- name: CreateDeposit :exec

INSERT INTO
    deposits (
        id,
        external_id,
        card_id,
        amount,
        paid,
        created_at,
        updated_at
    )
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: GetDepositsByExternalCardId :many

SELECT d.*
FROM cards c
    INNER JOIN deposits d ON d.card_id = c.id
WHERE c.external_id = $1;

-- name: GetDepositByExternalId :one

SELECT * FROM deposits WHERE external_id = $1;

-- name: UpdatePaidActiveDepositById :exec

UPDATE deposits SET paid = $2 WHERE id = $1 AND cancelled = false;

-- name: CancelDepositById :one

UPDATE deposits
SET cancelled = true
WHERE
    id = $1
    AND cancelled = false RETURNING *;