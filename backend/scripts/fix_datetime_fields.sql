-- 修复卡片表中的时间字段问题
-- 将 '0000-00-00' 日期替换为 NULL

-- 修复 created_at 字段
UPDATE cards 
SET created_at = NOW() 
WHERE created_at = '0000-00-00 00:00:00' OR created_at = '0000-00-00' OR created_at IS NULL;

-- 修复 updated_at 字段  
UPDATE cards 
SET updated_at = NOW() 
WHERE updated_at = '0000-00-00 00:00:00' OR updated_at = '0000-00-00' OR updated_at IS NULL;

-- 修复 expired_at 字段
UPDATE cards 
SET expired_at = NULL 
WHERE expired_at = '0000-00-00 00:00:00' OR expired_at = '0000-00-00';

-- 修复 synced_at 字段
UPDATE cards 
SET synced_at = NULL 
WHERE synced_at = '0000-00-00 00:00:00' OR synced_at = '0000-00-00';

-- 修夏 used_at 字段
UPDATE cards 
SET used_at = NULL 
WHERE used_at = '0000-00-00 00:00:00' OR used_at = '0000-00-00';

-- 修复 reserved_at 字段
UPDATE cards 
SET reserved_at = NULL 
WHERE reserved_at = '0000-00-00 00:00:00' OR reserved_at = '0000-00-00';

-- 为没有过期时间的卡片设置默认过期时间（1年后）
UPDATE cards 
SET expired_at = DATE_ADD(NOW(), INTERVAL 1 YEAR)
WHERE expired_at IS NULL;

-- 检查修复结果
SELECT 'Cards with valid dates' as message, COUNT(*) as count FROM cards 
WHERE created_at > '1000-01-01' AND updated_at > '1000-01-01' AND expired_at IS NOT NULL;

-- 检查是否还有问题数据
SELECT 'Problem records' as message, COUNT(*) as count FROM cards 
WHERE created_at = '0000-00-00 00:00:00' OR updated_at = '0000-00-00 00:00:00' 
   OR expired_at = '0000-00-00 00:00:00' OR synced_at = '0000-00-00 00:00:00'
   OR used_at = '0000-00-00 00:00:00' OR reserved_at = '0000-00-00 00:00:00';