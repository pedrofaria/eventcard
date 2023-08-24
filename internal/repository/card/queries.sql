-- name: GetCard :one

SELECT * FROM cards WHERE id = $1;

-- name: GetCardIdByExternalId :one

SELECT id FROM cards WHERE external_id = $1;

-- name: GetCardByExternalId :one

SELECT * FROM cards WHERE external_id = $1;

-- name: GetCardFull :one

SELECT
    c.*,
    b.amount as "balance"
FROM cards c
    INNER JOIN balances b ON b.card_id = c.id
WHERE c.id = $1;

-- name: GetCardFullByExternalId :one

SELECT
    c.*,
    b.amount as "balance"
FROM cards c
    INNER JOIN balances b ON b.card_id = c.id
WHERE c.external_id = $1;

-- name: GetCardBalance :one

SELECT amount FROM balances WHERE card_id = $1;

-- name: CreateCard :exec

INSERT INTO
    cards (
        id,
        external_id,
        name,
        enabled,
        created_at,
        updated_at
    )
VALUES ($1, $2, $3, $4, $5, $6);

-- name: UpdateEnabledCardByExternalId :exec

UPDATE cards SET enabled = $2 WHERE external_id = $1;

-- name: CreateBalance :exec

INSERT INTO
    balances (
        id,
        card_id,
        amount,
        created_at,
        updated_at
    )
VALUES ($1, $2, $3, $4, $5);