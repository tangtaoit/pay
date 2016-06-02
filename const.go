package main
const (
	ReturnCodeSuccess = "SUCCESS"
	ReturnCodeFail    = "FAIL"
)

const (
	ResultCodeSuccess = "SUCCESS"
	ResultCodeFail    = "FAIL"
)

const (
	//等待交易
	Trade_Status_Wait =0

	//交易成功
	Trade_Status_Success =1

	//交易失败
	Trade_Status_Fail =2

	//交易金额未满
	Trade_Status_NOFULL =3
)

//支付类型
const (
	//支付宝
	Pay_Type_AliPay = 1

	//微信支付
	Pay_Type_WXPAY =2

	//账户支付
	Pay_Type_Account =3
)


//订单类型
const  (

	//充值订单
	Trade_Type_Recharge = 1
	//购买
	Trade_Type_Buy =2
	//预付款
	Trade_Type_Imprest =3
)

//通知状态
const (
	//等待
	Notify_Status_Wait = 0

	//完成
	Notify_Status_Finish=1

	//错误
	Notify_Status_Error =2
)

const(
	//禁用
	Account_Status_Disable=0

	//启用
	Account_Status_Enable=1

)