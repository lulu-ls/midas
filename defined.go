package midas

import "errors"

// 微信虚拟支付 文档地址：https://docs.qq.com/doc/DVUN0QWJja0J5c2x4

var (
	ErrorIsSessionKeyInvalid = errors.New("session key 过期")
)

const (
	wxApi      = "https://api.weixin.qq.com" // 微信接口地址
	signMethod = "hmac_sha256"               // 签名算法
)

type interType int

// 接口地址定义
const (
	getBalance      interType = iota // 获取余额
	cancelPay                        // 取消支付
	pay                              // 支付
	present                          // 赠送
	checkSessionKey                  // 检查session key 是否过期
)

var (
	apiPath = map[interType]string{
		getBalance:      "/wxa/game/getbalance",
		cancelPay:       "/wxa/game/cancelpay",
		pay:             "/wxa/game/pay",
		present:         "/wxa/game/present",
		checkSessionKey: "/wxa/checksession",
	}
)

// 常见错误码定义
const (
	CodeSuccess   = 0     // 请求成功
	CodeBusy      = -1    // 系统繁忙，此时请开发者稍候再试
	CodeSignature = 90010 // signature签名错误或用户登录态（session_key）已过期
	CodePaySign   = 90011 // pay_sig签名错误
	CodeParams    = 90018 // 参数错误，具体参数见errmsg描述
)

// CommonRequest 公共请求参数定义
type CommonRequest struct {
	OpenId  string `json:"openid"`   // 是 用户唯一标识符
	OfferId string `json:"offer_id"` // 是 支付应用ID（OfferId）
	Ts      int64  `json:"ts"`       // 是 当前UNIX时间戳（请尽可能确保时间准确），单位：秒 如：1668136271
	ZoneId  string `json:"zone_id"`  // 是 已发布的分区ID（MP-分区配置-分区ID）需要和env对应
	Env     int    `json:"env"`      // 是 环境配置 0：现网环境（也叫正式环境）1：沙箱环境
	UserIp  string `json:"user_ip"`  // 否 用户外网 ip
}

// CommonReply 公共返回参数定义
type CommonReply struct {
	ErrCode int    `json:"errcode"` // 错误码
	ErrMsg  string `json:"errmsg"`  // 错误信息
}

// BalanceRequest 获取余额请求参数
type BalanceRequest struct {
	CommonRequest
}

// BalanceReply 获取余额返回参数
type BalanceReply struct {
	CommonReply
	Balance        int  `json:"balance"`         // 游戏币总余额，包括现金充值和赠送部分。
	PresentBalance int  `json:"present_balance"` // 赠送账户的游戏币余额（原1.0的gen_balance）
	SumSave        int  `json:"sum_save"`        // 累计现金充值获得的游戏币数量（原1.0的save_amt）
	SumPresent     int  `json:"sum_present"`     // 累计赠送的游戏币数量（原1.0的present_sum）
	SumBalance     int  `json:"sum_balance"`     // 累计获得的游戏币数量，包括现金充值和赠送（原1.0的save_sum）
	SumCost        int  `json:"sum_cost"`        // 累计总消耗（即扣除）游戏币数量（原1.0的cost_sum）
	FirstSave      bool `json:"first_save"`      // 是否满足首充活动标记（原1.0的first_save）

}

// PayRequest 虚拟支付请求参数
type PayRequest struct {
	CommonRequest
	Amount  int    `json:"amount"`  // 是 扣除游戏币数量，需要大于0（原1.0的amt）
	BillNo  string `json:"bill_no"` // 是 扣除游戏币订单号，业务需要保证全局唯一，相同的订单号的多次请求不会重复扣除 长度不超过63，只能是数字、英文大小写字母及_-的组合 不能以下划线（_）开头（2.0新增约束）
	PayItem string `json:"payitem"` // 否 道具信息
	Remark  string `json:"remark"`  // 否 备注
}

// PayReply 虚拟支付返回参数
type PayReply struct {
	CommonReply
	BillNo            string `json:"bill_no"`             // 扣除游戏币订单号
	Balance           int    `json:"balance"`             // 扣款后的余额
	UsedPresentAmount int    `json:"used_present_amount"` // 本次扣的赠送币的数量（原1.0的used_gen_amt)
}
