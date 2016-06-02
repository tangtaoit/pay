
-- 应用信息
CREATE TABLE IF NOT EXISTS app(

  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  open_id VARCHAR(255) COMMENT '用户的OPENID',
  app_id VARCHAR(255) UNIQUE COMMENT '应用ID',
  app_key VARCHAR(255) COMMENT '应用KEY',
  app_name VARCHAR(255) COMMENT '应用名称',
  app_desc VARCHAR(1000) COMMENT '应用描述',
  status int COMMENT '应用状态 0.待审核 1.已审核',
  create_time TIMESTAMP COMMENT '创建时间',
  update_time TIMESTAMP COMMENT '更新时间'

) CHARACTER SET utf8;

-- 账户表
CREATE TABLE IF NOT EXISTS accounts(
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  app_id VARCHAR(255) COMMENT '应用ID',
  open_id VARCHAR(255) UNIQUE COMMENT '用户OpenID',
  create_time TIMESTAMP COMMENT '创建时间',
  amount NUMERIC(20,0) comment '账户金额',
  update_time  TIMESTAMP COMMENT '更新时间',
  password VARCHAR(255) COMMENT '账户密码',
  status int DEFAULT 0 COMMENT '账户状态 1.正常 0.异常'
) CHARACTER SET utf8;

-- 账户记录
CREATE TABLE IF NOT EXISTS accounts_record(

  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  trade_no VARCHAR(255)   COMMENT '交易号',
  app_id VARCHAR(255) COMMENT '应用ID',
  open_id VARCHAR(255) COMMENT '用户openID',
  account_id BIGINT  COMMENT '账户ID',
  create_time TIMESTAMP COMMENT '创建时间',
  amount_before NUMERIC(20,0) COMMENT '金额变动前',
  amount_after NUMERIC(20,0) COMMENT '金额变动后',
  changed_amount NUMERIC(20,0) COMMENT '变动金额(单位:分)'


) CHARACTER SET utf8;

-- 预付款
CREATE TABLE IF NOT EXISTS imprest(

  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  code VARCHAR(255) UNIQUE COMMENT '付款代号',
  open_id VARCHAR(255) COMMENT '用户ID',
  create_time TIMESTAMP COMMENT '创建时间',
  update_time  TIMESTAMP COMMENT '更新时间',
  type int COMMENT '付款码类型 1.普通类型',
  amount NUMERIC(20,0) COMMENT '金额',
  status int DEFAULT 0 COMMENT '状态 1.进行中 0.异常 2.完成'

) CHARACTER SET utf8;

-- 预付款记录
CREATE TABLE IF NOT EXISTS imprest_record(

  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  code VARCHAR(255) UNIQUE COMMENT '付款码代号',
  trade_no VARCHAR(255) COMMENT '交易编号',
  open_id VARCHAR(255) COMMENT '用户ID',
  opr_open_id VARCHAR(255) COMMENT '操作者用户ID',
  create_time TIMESTAMP COMMENT '创建时间',
  update_time  TIMESTAMP COMMENT '更新时间',
  type int COMMENT '付款码类型 1.普通类型',
  amount NUMERIC(20,0) COMMENT '金额'

) CHARACTER SET utf8;




-- 交易表
CREATE TABLE IF NOT EXISTS trades(

  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  trade_no VARCHAR(255) UNIQUE  COMMENT '交易号' ,
  trade_type INT COMMENT '交易类型 1.充值 2.普通支出 3.预付款',
  out_trade_no VARCHAR(255) COMMENT '第三方系统中的交易号',
  out_trade_type INT COMMENT '第三方系统中的交易类型',
  notify_url VARCHAR(1000) COMMENT '交易通知地址',
  notify_status INT COMMENT '通知状态(1.通知成功 0.等待通知 2.通知失败)',
  app_id VARCHAR(255) COMMENT '应用ID',
  open_id VARCHAR(255) COMMENT '用户openID',
  create_time TIMESTAMP COMMENT '创建时间',
  update_time TIMESTAMP COMMENT '更新时间',
  changed_amount NUMERIC(20,0) COMMENT '变动金额(单位:分)',
  actual_amount  NUMERIC(20,0) COMMENT '实际支付金额',
  title VARCHAR(255) comment '交易标题',
  remark VARCHAR(1000) COMMENT '交易备注',
  no_once int COMMENT '是否需要一次性付清 0.不需要 1.需要',
  status INT COMMENT '状态 1.交易成功 0.待交易 2.交易错误 3.交易欠付'

) CHARACTER SET utf8;

-- 交易支付信息
CREATE TABLE IF NOT EXISTS trades_pay(

  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  open_id VARCHAR(255) COMMENT '付款人的ID',
  trade_no VARCHAR(255)  COMMENT '交易号',
  create_time TIMESTAMP COMMENT '创建时间',
  update_time TIMESTAMP COMMENT '更新时间',
  pay_type INT COMMENT '支付类型 0.支付宝 1.微信 2.账户余额',
  pay_amount NUMERIC(20,0) COMMENT '实际付款金额'

) CHARACTER SET utf8;

