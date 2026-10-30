
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

-- name: ListDailyPricesWithPreviousClose :many
SELECT
    current_price.symbol,
    stocks.company_name,
    current_price.trade_date,
    current_price.close_price,
    previous_price.trade_date AS previous_trade_date,
    previous_price.close_price AS previous_close_price
FROM daily_prices as current_price
         INNER JOIN daily_prices as previous_price
                    ON previous_price.symbol = current_price.symbol
                        AND previous_price.trade_date = (
                            SELECT MAX(candidate.trade_date)
                            FROM daily_prices AS candidate
                            WHERE candidate.symbol = current_price.symbol
                              AND candidate.trade_date < current_price.trade_date
                        )
         LEFT JOIN stocks ON stocks.symbol = current_price.symbol
WHERE current_price.trade_date = ?;

-- name: GetLatestDailyPriceTradeDate :one
SELECT trade_date
FROM daily_prices
ORDER BY trade_date DESC
LIMIT 1;

-- name: ListDailyVolumesByTradeDate :many
SELECT
    symbol,
    trade_date,
    volume
FROM daily_prices
WHERE trade_date = ?
ORDER BY symbol;

-- name: ListDailyVolumeHistory :many
WITH ranked_volumes AS (
    SELECT
        symbol,
        trade_date,
        volume,
        ROW_NUMBER() OVER (
            PARTITION BY symbol
            ORDER BY trade_date DESC
        ) AS history_number
    FROM daily_prices
    WHERE trade_date < ?
)
SELECT
    symbol,
    trade_date,
    volume
FROM ranked_volumes
WHERE history_number <= 20
ORDER BY symbol, trade_date DESC;

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
