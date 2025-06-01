-- 备份脚本：在重构前备份所有卡片相关数据
-- 执行时间：执行重构前
-- 使用方法：mysql -u root -p ruixin < backup_cards.sql

-- 1. 创建备份表
CREATE TABLE IF NOT EXISTS cards_backup_20250127 AS SELECT * FROM cards;
CREATE TABLE IF NOT EXISTS category_bindings_backup_20250127 AS SELECT * FROM category_bindings;
CREATE TABLE IF NOT EXISTS card_usage_logs_backup_20250127 AS SELECT * FROM card_usage_logs;

-- 2. 验证备份
SELECT 'cards_backup' as table_name, COUNT(*) as count FROM cards_backup_20250127
UNION ALL
SELECT 'category_bindings_backup', COUNT(*) FROM category_bindings_backup_20250127
UNION ALL
SELECT 'card_usage_logs_backup', COUNT(*) FROM card_usage_logs_backup_20250127;

-- 3. 导出当前卡片状态报告
SELECT 
    'Total Cards' as metric,
    COUNT(*) as value
FROM cards
UNION ALL
SELECT 
    'Active Cards',
    COUNT(*)
FROM cards 
WHERE status = 1
UNION ALL
SELECT 
    'Used Cards',
    SUM(used_count)
FROM cards
UNION ALL
SELECT 
    'Cards with Bindings',
    COUNT(DISTINCT category_id)
FROM category_bindings;