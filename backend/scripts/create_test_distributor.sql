-- 创建测试分销商数据
INSERT INTO distributors (
    name, 
    company_name, 
    contact_name, 
    phone, 
    email, 
    password, 
    api_key, 
    api_secret, 
    status, 
    balance, 
    frozen_amount, 
    credit_limit, 
    daily_order_limit, 
    monthly_order_limit,
    created_at,
    updated_at
) VALUES 
(
    '测试分销商A', 
    '测试科技有限公司A', 
    '张三', 
    '13800138001', 
    'testa@distributor.com', 
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- password: demo123
    'test-api-key-001', 
    'test-api-secret-001', 
    1, 
    1000.00, 
    0, 
    10000.00, 
    100, 
    3000,
    NOW(),
    NOW()
),
(
    '测试分销商B', 
    '测试科技有限公司B', 
    '李四', 
    '13800138002', 
    'testb@distributor.com', 
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- password: demo123
    'test-api-key-002', 
    'test-api-secret-002', 
    1, 
    2000.00, 
    0, 
    20000.00, 
    200, 
    5000,
    NOW(),
    NOW()
),
(
    '测试分销商C', 
    '测试科技有限公司C', 
    '王五', 
    '13800138003', 
    'testc@distributor.com', 
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- password: demo123
    'test-api-key-003', 
    'test-api-secret-003', 
    1, 
    500.00, 
    0, 
    5000.00, 
    50, 
    1500,
    NOW(),
    NOW()
);

-- 验证数据插入
SELECT id, name, email, balance, status FROM distributors WHERE email LIKE 'test%@distributor.com';