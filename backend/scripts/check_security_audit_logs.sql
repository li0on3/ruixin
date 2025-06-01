-- 查看最近的安全审计日志
SELECT 
    id,
    distributor_id,
    action,
    resource,
    JSON_EXTRACT(details, '$.card_code') as card_code,
    JSON_EXTRACT(details, '$.api') as api,
    JSON_EXTRACT(details, '$.authorized') as authorized,
    status,
    ip_address,
    created_at
FROM security_audit_logs
ORDER BY created_at DESC
LIMIT 20;

-- 统计未授权访问
SELECT 
    distributor_id,
    COUNT(*) as unauthorized_count,
    GROUP_CONCAT(DISTINCT resource) as accessed_cards
FROM security_audit_logs
WHERE action = 'UNAUTHORIZED_CARD'
GROUP BY distributor_id
ORDER BY unauthorized_count DESC;

-- 统计频率限制触发
SELECT 
    distributor_id,
    COUNT(*) as rate_limited_count,
    MAX(created_at) as last_limited_at
FROM security_audit_logs
WHERE action = 'RATE_LIMITED'
GROUP BY distributor_id
ORDER BY rate_limited_count DESC;

-- 查看特定分销商的审计日志
-- SELECT * FROM security_audit_logs WHERE distributor_id = 1 ORDER BY created_at DESC;

-- 查看今天的所有安全事件
SELECT 
    action,
    COUNT(*) as count
FROM security_audit_logs
WHERE DATE(created_at) = CURDATE()
GROUP BY action;