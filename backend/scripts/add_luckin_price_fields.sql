-- 添加瑞幸原始价格字段到订单表
ALTER TABLE `orders` 
ADD COLUMN `luckin_price` DECIMAL(10,2) DEFAULT NULL COMMENT '瑞幸原始价格' AFTER `profit_amount`,
ADD COLUMN `luckin_cost_price` DECIMAL(10,2) DEFAULT NULL COMMENT '瑞幸原始成本价' AFTER `luckin_price`;

-- 添加索引方便查询价格差异
ALTER TABLE `orders` ADD INDEX `idx_price_diff` (`luckin_price`, `total_amount`);