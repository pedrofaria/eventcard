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

-- name: GetDepositsByCardId :many

SELECT * FROM deposits WHERE card_id = $1 ORDER BY created_at DESC;

-- name: GetDepositByCardAndExternalId :one

SELECT * FROM deposits WHERE card_id = $1 AND external_id = $2;

-- name: UpdatePaidActiveDepositById :one

UPDATE deposits
SET paid = $2, updated_at = NOW()
WHERE
    id = $1
    AND cancelled = false RETURNING *;

-- name: CancelDepositById :one

UPDATE deposits
SET
    cancelled = true,
    updated_at = NOW()
WHERE
    id = $1
    AND cancelled = false RETURNING *;