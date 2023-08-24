-- name: CreateLedger :exec

INSERT INTO
    ledgers (
        id,
        card_id,
        reference,
        reference_id,
        amount,
        created_at
    )
VALUES ($1, $2, $3, $4, $5, $6);

-- name: IncreaseCardBalance :exec

UPDATE balances SET amount = amount + $1 WHERE card_id = $2;