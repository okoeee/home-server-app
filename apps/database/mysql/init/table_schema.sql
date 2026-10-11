CREATE TABLE daily_prices (
    symbol VARCHAR(10) NOT NULL COMMENT '証券コード',
    trade_date DATE NOT NULL COMMENT '取引日',
    open_price DECIMAL(14, 4) NOT NULL COMMENT '調整済み始値',
    high_price DECIMAL(14, 4) NOT NULL COMMENT '調整済み高値',
    low_price DECIMAL(14, 4) NOT NULL COMMENT '調整済み安値',
    close_price DECIMAL(14, 4) NOT NULL COMMENT '調整済み終値',
    volume BIGINT UNSIGNED NOT NULL COMMENT '調整済み出来高',
    PRIMARY KEY (symbol, trade_date)
) COMMENT = '株価の日足';
