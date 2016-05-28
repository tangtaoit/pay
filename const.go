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
	//普通订单
	Trade_Type_CommOrder =2
)