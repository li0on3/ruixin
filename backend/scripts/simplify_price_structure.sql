-- 简化价格结构迁移脚本
-- 1. 在卡片表添加瑞幸产品ID字段
ALTER TABLE `cards` 
ADD COLUMN `luckin_product_id` INT DEFAULT 6 COMMENT '瑞幸产品ID' AFTER `card_code`,
ADD INDEX `idx_luckin_product_id` (`luckin_product_id`);

-- 2. 从现有数据迁移产品ID
UPDATE cards c 
LEFT JOIN luckin_prices lp ON c.price_id = lp.id 
SET c.luckin_product_id = IFNULL(lp.price_id, 6);

-- 3. 更新字段注释，标记price_id为废弃
ALTER TABLE `cards` 
MODIFY COLUMN `price_id` BIGINT COMMENT '价格ID（已废弃，保留用于兼容）';

-- 4. 创建数据备份表（以防需要回滚）
CREATE TABLE IF NOT EXISTS `luckin_prices_backup_20250529` AS 
SELECT * FROM luckin_prices;

-- 5. 输出迁移统计
SELECT 
    COUNT(*) as total_cards,
    SUM(CASE WHEN luckin_product_id = 6 THEN 1 ELSE 0 END) as product_id_6,
    SUM(CASE WHEN luckin_product_id != 6 THEN 1 ELSE 0 END) as other_product_ids
FROM cards;

-- 注意：暂时不删除 luckin_prices 表和 price_id 字段，等系统稳定运行后再清理