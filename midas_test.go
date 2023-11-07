package midas

import (
	"errors"
	"testing"
)

// 获取余额
func TestMidas_GetBalance(t *testing.T) {
	vir := NewMidas("")
	//vir.IsSandbox = true
	vir.OfferId = ""
	vir.ZoneId = ""
	vir.Env = 0

	vir.SessKey = ""
	vir.AccToken = ""

	common := CommonRequest{
		OpenId: "",
		//UserIp: f_ctx.GetIp(ctx), // 获取 ip
	}
	// 查询游戏币余额
	res, err := vir.GetBalance(&BalanceRequest{
		CommonRequest: common,
	})
	if err != nil {
		if errors.Is(err, ErrorIsSessionKeyInvalid) { // 判断是否 session key 失效，需要用户重新登录 401

		}
		t.Log(err)
	}

	t.Log(res)

}

func TestMidas_Pay(t *testing.T) {
	vir := NewMidas("")
	//vir.IsSandbox = true
	vir.OfferId = ""
	vir.ZoneId = ""
	vir.Env = 0

	vir.SessKey = ""
	vir.AccToken = ""

	// 查询游戏币余额
	res, err := vir.Pay(&PayRequest{
		CommonRequest: CommonRequest{
			OpenId: "",
		},
		Amount:  100,
		BillNo:  "123456798123456",
		PayItem: "2134141132",
		Remark:  "测试支付",
	})
	if err != nil {
		if errors.Is(err, ErrorIsSessionKeyInvalid) { // 判断是否 session key 失效，需要用户重新登录 401

		}
		t.Log(err)
	}

	t.Log(res)
}
