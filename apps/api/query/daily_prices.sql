
-- name: ListDailyPricesBySymbol :many
SELECT
    symbol,
    trade_date,
    open_price,
    high_price,
    low_price,
    close_price,
    volume
FROM daily_prices
WHERE symbol = ?
ORDER BY trade_date;

-- name: UpsertDailyPrice :exec
INSERT INTO daily_prices (
    symbol,
    trade_date,
    open_price,
    high_price,
    low_price,
    close_price,
    volume
) VALUES (
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?
)
ON DUPLICATE KEY UPDATE
    open_price = VALUES(open_price),
    high_price = VALUES(high_price),
    low_price = VALUES(low_price),
    close_price = VALUES(close_price),
    volume = VALUES(volume);
