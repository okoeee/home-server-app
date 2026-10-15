-- name: UpsertStock :exec
INSERT INTO stocks (
    symbol,
    effective_date,
    company_name,
    sector_33_code,
    sector_33_name,
    market_code,
    market_name,
    product_category_code,
    is_listed
) VALUES (
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    TRUE
)
ON DUPLICATE KEY UPDATE
    effective_date = VALUES(effective_date),
    company_name = VALUES(company_name),
    sector_33_code = VALUES(sector_33_code),
    sector_33_name = VALUES(sector_33_name),
    market_code = VALUES(market_code),
    market_name = VALUES(market_name),
    product_category_code = VALUES(product_category_code),
    is_listed = TRUE;
