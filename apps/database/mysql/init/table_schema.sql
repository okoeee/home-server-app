CREATE TABLE stocks (
    symbol VARCHAR(5) NOT NULL COMMENT '銘柄を識別する5桁の証券コード',
    effective_date DATE NOT NULL COMMENT '銘柄情報の適用日',
    company_name VARCHAR(255) NOT NULL COMMENT '会社名',
    sector_33_code VARCHAR(4) NOT NULL COMMENT '33業種コード',
    sector_33_name VARCHAR(100) NOT NULL COMMENT '33業種名',
    market_code VARCHAR(4) NOT NULL COMMENT '市場区分コード',
    market_name VARCHAR(100) NOT NULL COMMENT '市場区分名',
    product_category_code VARCHAR(3) NOT NULL COMMENT '商品区分コード',
    is_listed BOOLEAN NOT NULL DEFAULT TRUE COMMENT '現在の上場状態',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '登録日時',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日時',
    PRIMARY KEY (symbol)
) COMMENT = '銘柄マスター';

CREATE TABLE daily_prices (
    symbol VARCHAR(5) NOT NULL COMMENT '銘柄を識別する5桁の証券コード',
    trade_date DATE NOT NULL COMMENT '取引日',
    open_price DECIMAL(14, 1) NOT NULL COMMENT '調整済み始値（小数第1位）',
    high_price DECIMAL(14, 1) NOT NULL COMMENT '調整済み高値（小数第1位）',
    low_price DECIMAL(14, 1) NOT NULL COMMENT '調整済み安値（小数第1位）',
    close_price DECIMAL(14, 1) NOT NULL COMMENT '調整済み終値（小数第1位）',
    volume DECIMAL(20, 1) UNSIGNED NOT NULL COMMENT '調整済み出来高（小数第1位）',
    PRIMARY KEY (symbol, trade_date)
) COMMENT = '株価の日足';
