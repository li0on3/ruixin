-- 数据迁移脚本：将现有卡片数据迁移到新结构
-- 执行前请确保已备份数据
-- 使用方法：mysql -u root -p ruixin < migrate_cards.sql

-- 1. 检查是否已经执行过迁移
SELECT '开始检查数据库状态...' as message;

-- 2. 创建临时表存储迁移前的数据
CREATE TABLE IF NOT EXISTS cards_migration_temp AS 
SELECT * FROM cards WHERE 1=0;

-- 3. 备份现有数据到临时表
INSERT INTO cards_migration_temp SELECT * FROM cards;

-- 4. 添加新字段（如果不存在）
SET @exist := (SELECT COUNT(*) FROM information_schema.COLUMNS 
    WHERE TABLE_SCHEMA = DATABASE() 
    AND TABLE_NAME = 'cards' 
    AND COLUMN_NAME = 'batch_id');

SET @sql = IF(@exist = 0, 
    'ALTER TABLE cards ADD COLUMN batch_id INT UNSIGNED AFTER card_code',
    'SELECT "batch_id already exists"');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @exist := (SELECT COUNT(*) FROM information_schema.COLUMNS 
    WHERE TABLE_SCHEMA = DATABASE() 
    AND TABLE_NAME = 'cards' 
    AND COLUMN_NAME = 'price_id');

SET @sql = IF(@exist = 0, 
    'ALTER TABLE cards ADD COLUMN price_id INT UNSIGNED AFTER batch_id',
    'SELECT "price_id already exists"');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @exist := (SELECT COUNT(*) FROM information_schema.COLUMNS 
    WHERE TABLE_SCHEMA = DATABASE() 
    AND TABLE_NAME = 'cards' 
    AND COLUMN_NAME = 'cost_price');

SET @sql = IF(@exist = 0, 
    'ALTER TABLE cards ADD COLUMN cost_price DECIMAL(10,2) DEFAULT 0.00 AFTER price_id',
    'SELECT "cost_price already exists"');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @exist := (SELECT COUNT(*) FROM information_schema.COLUMNS 
    WHERE TABLE_SCHEMA = DATABASE() 
    AND TABLE_NAME = 'cards' 
    AND COLUMN_NAME = 'sell_price');

SET @sql = IF(@exist = 0, 
    'ALTER TABLE cards ADD COLUMN sell_price DECIMAL(10,2) DEFAULT 0.00 AFTER cost_price',
    'SELECT "sell_price already exists"');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @exist := (SELECT COUNT(*) FROM information_schema.COLUMNS 
    WHERE TABLE_SCHEMA = DATABASE() 
    AND TABLE_NAME = 'cards' 
    AND COLUMN_NAME = 'used_at');

SET @sql = IF(@exist = 0, 
    'ALTER TABLE cards ADD COLUMN used_at TIMESTAMP NULL AFTER status',
    'SELECT "used_at already exists"');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @exist := (SELECT COUNT(*) FROM information_schema.COLUMNS 
    WHERE TABLE_SCHEMA = DATABASE() 
    AND TABLE_NAME = 'cards' 
    AND COLUMN_NAME = 'order_id');

SET @sql = IF(@exist = 0, 
    'ALTER TABLE cards ADD COLUMN order_id INT UNSIGNED NULL AFTER used_at',
    'SELECT "order_id already exists"');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 5. 为现有数据创建默认批次
INSERT INTO card_batches (batch_no, price_id, cost_price, sell_price, total_count, used_count, imported_at, imported_by, description, created_at, updated_at)
SELECT 
    CONCAT('MIGRATED-', DATE_FORMAT(NOW(), '%Y%m%d-%H%i%s')) as batch_no,
    lp.id as price_id,
    7.50 as cost_price,  -- 默认成本价
    8.50 as sell_price,  -- 默认销售价
    COUNT(*) as total_count,
    SUM(CASE WHEN c.used_count > 0 THEN 1 ELSE 0 END) as used_count,
    NOW() as imported_at,
    1 as imported_by,  -- 默认管理员ID
    '数据迁移自动创建的批次' as description,
    NOW() as created_at,
    NOW() as updated_at
FROM cards c
LEFT JOIN luckin_prices lp ON lp.price_id = c.product_id
WHERE c.product_id IS NOT NULL
GROUP BY c.product_id, lp.id;

-- 6. 更新卡片数据
UPDATE cards c
JOIN luckin_prices lp ON lp.price_id = c.product_id
JOIN card_batches cb ON cb.price_id = lp.id AND cb.description = '数据迁移自动创建的批次'
SET 
    c.batch_id = cb.id,
    c.price_id = lp.id,
    c.cost_price = 7.50,
    c.sell_price = 8.50,
    c.status = CASE 
        WHEN c.used_count > 0 THEN 1  -- 已使用
        WHEN c.status = 0 THEN 0       -- 未使用
        ELSE 0                         -- 其他状态都改为未使用
    END,
    c.used_at = CASE 
        WHEN c.used_count > 0 THEN c.updated_at 
        ELSE NULL 
    END
WHERE c.product_id IS NOT NULL;

-- 7. 处理没有产品ID的卡片（创建一个特殊批次）
INSERT INTO card_batches (batch_no, price_id, cost_price, sell_price, total_count, used_count, imported_at, imported_by, description, created_at, updated_at)
SELECT 
    'MIGRATED-NO-PRODUCT' as batch_no,
    (SELECT id FROM luckin_prices LIMIT 1) as price_id,  -- 使用第一个价格配置
    0.00 as cost_price,
    0.00 as sell_price,
    COUNT(*) as total_count,
    SUM(CASE WHEN used_count > 0 THEN 1 ELSE 0 END) as used_count,
    NOW() as imported_at,
    1 as imported_by,
    '无产品ID的卡片迁移批次' as description,
    NOW() as created_at,
    NOW() as updated_at
FROM cards
WHERE product_id IS NULL OR product_id = 0
HAVING COUNT(*) > 0;

-- 更新无产品ID的卡片
UPDATE cards c
JOIN card_batches cb ON cb.batch_no = 'MIGRATED-NO-PRODUCT'
SET 
    c.batch_id = cb.id,
    c.price_id = cb.price_id,
    c.cost_price = 0.00,
    c.sell_price = 0.00,
    c.status = CASE 
        WHEN c.used_count > 0 THEN 1
        ELSE 0
    END,
    c.used_at = CASE 
        WHEN c.used_count > 0 THEN c.updated_at 
        ELSE NULL 
    END
WHERE c.product_id IS NULL OR c.product_id = 0;

-- 8. 删除旧字段（可选，建议先保留一段时间）
-- ALTER TABLE cards DROP COLUMN product_id;
-- ALTER TABLE cards DROP COLUMN daily_limit;
-- ALTER TABLE cards DROP COLUMN total_limit;
-- ALTER TABLE cards DROP COLUMN used_count;

-- 9. 添加索引优化查询性能
CREATE INDEX IF NOT EXISTS idx_cards_price_status ON cards(price_id, status);
CREATE INDEX IF NOT EXISTS idx_cards_batch_id ON cards(batch_id);
CREATE INDEX IF NOT EXISTS idx_cards_status ON cards(status);

-- 10. 输出迁移报告
SELECT '迁移完成报告' as message;

SELECT 
    'Total cards migrated' as metric,
    COUNT(*) as value
FROM cards
WHERE batch_id IS NOT NULL
UNION ALL
SELECT 
    'Batches created',
    COUNT(*)
FROM card_batches
WHERE batch_no LIKE 'MIGRATED-%'
UNION ALL
SELECT 
    'Cards marked as used',
    COUNT(*)
FROM cards
WHERE status = 1
UNION ALL
SELECT 
    'Cards marked as unused',
    COUNT(*)
FROM cards
WHERE status = 0;

-- 11. 验证数据完整性
SELECT 
    'Data integrity check' as check_item,
    CASE 
        WHEN COUNT(*) = 0 THEN 'PASS' 
        ELSE CONCAT('FAIL - ', COUNT(*), ' cards without batch_id')
    END as result
FROM cards
WHERE batch_id IS NULL
UNION ALL
SELECT 
    'Price ID assignment check',
    CASE 
        WHEN COUNT(*) = 0 THEN 'PASS' 
        ELSE CONCAT('FAIL - ', COUNT(*), ' cards without price_id')
    END
FROM cards
WHERE price_id IS NULL;

SELECT '迁移脚本执行完成！' as message;