
-- 账户表
CREATE TABLE IF NOT EXISTS accounts(
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  open_id VARCHAR(255) UNIQUE COMMENT '用户OpenID',
  create_time TIMESTAMP COMMENT '创建时间',
  amount NUMERIC(20,0) comment '账户金额',
  status int DEFAULT 0 COMMENT '账户状态 1.正常 0.异常'
);

-- 交易记录表
CREATE TABLE IF NOT EXISTS trades(
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  account_id BIGINT  COMMENT '账户ID'  ,
  in_out VARCHAR(255) COMMENT 'IN 收入 OUT 支出',
  create_time TIMESTAMP COMMENT '创建时间',
  changed_amount NUMERIC(20,0) COMMENT '变动金额(单位:分)',
  title VARCHAR(255) comment '交易标题',
  remark VARCHAR(255) COMMENT '交易备注'

)