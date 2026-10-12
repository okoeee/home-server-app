
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
