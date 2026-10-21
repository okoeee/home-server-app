
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
WHERE current_price.trade_date = ?;

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
