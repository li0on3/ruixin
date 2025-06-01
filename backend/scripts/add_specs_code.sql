-- 添加specs_code字段到product_sku_mappings表
ALTER TABLE product_sku_mappings ADD COLUMN specs_code VARCHAR(100) COMMENT '完整的规格代码组合，如"0_0_0_1"';

-- 你可以手动执行这个SQL文件：
-- mysql -u your_username -p your_database < scripts/add_specs_code.sql