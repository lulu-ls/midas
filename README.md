# wechat midas pay SDK for golang
### 微信虚拟支付 文档地址：https://docs.qq.com/doc/DVUN0QWJja0J5c2x4，大家注意看下文档版本是否和自己需要的一样 这个对接的是 v2
https://github.com/lulu-ls/midas

### 简介（目前只用到了两个接口，其他的接口如果有需要和我联系，我尽快加上）
| 接口  | 状态                    |
|---------------:|:---------|
| 获取余额   | 已完成          |
| 扣除代币   | 已完成          |
| 撤销扣除   | 未完成          |
| 赠送代币   | 未完成          |


## 安装
go get github.com/lulu-ls/midas

## 使用
```go
import (
	"errors"
  	"fmt"
	"github.com/lulu-ls/midas"
)	
```
```go
	vir := midas.NewMidas("")
	//vir.IsSandbox = true
	vir.OfferId = "" // 配置
	vir.ZoneId = "" // 配置
	vir.Env = 0 // 是否是正式环境

	vir.SessKey = "" // 用户的 session key
	vir.AccToken = "" // access token

	// 查询游戏币余额
	res, err := vir.GetBalance(&midas.BalanceRequest{
		CommonRequest: midas.CommonRequest{
		OpenId: "",
		//UserIp: f_ctx.GetIp(ctx), // 获取 ip
	},
	})
	if err != nil {
		if errors.Is(err, midas.ErrorIsSessionKeyInvalid) { // 判断是否 session key 失效，需要用户重新登录 401

		}
		fmt.Println(err)
	}

	fmt.Println(res)

	// 扣除代币
	res, err := vir.Pay(&PayRequest{
		CommonRequest: CommonRequest{
			OpenId: "",
		},
		Amount:  100,
		BillNo:  "",
		PayItem: "",
		Remark:  "测试支付",
	})
	if err != nil {
		if errors.Is(err, ErrorIsSessionKeyInvalid) { // 判断是否 session key 失效，需要用户重新登录 401

		}
		fmt.Println(err)
	}

	fmt.Println(res)

```

## 联系方式
QQ群: 312533472

## 文档
暂无，可以看参数注释很详细

## 授权(LICENSE)
[wechat is licensed under the Apache Licence, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0.html)
