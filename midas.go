package midas

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"time"
)

// 微信虚拟支付 文档地址：https://docs.qq.com/doc/DVUN0QWJja0J5c2x4

type Midas struct {
	AppKey   string // app_key为当前支付环境(env参数)对应的AppKey，从“mp-支付基础配置”处获取
	OfferId  string // 支付应用ID（OfferId）
	ZoneId   string // 已发布的分区ID（MP-分区配置-分区ID） 需要和env对应
	Env      int    // 是否为沙箱环境 0：现网环境（也叫正式环境）1：沙箱环境
	SessKey  string // 用户的 session key
	AccToken string // 该 app 的 access token
}

type Opt func(vo *Midas)

func NewMidas(appKey string, opts ...Opt) *Midas {
	vir := &Midas{
		AppKey: appKey,
	}

	for _, f := range opts {
		f(vir)
	}

	return vir
}

func WithOfferId(offerId string) Opt {
	return func(vo *Midas) {
		vo.OfferId = offerId
	}
}

func WithZoneId(zoneId string) Opt {
	return func(vo *Midas) {
		vo.ZoneId = zoneId
	}
}

func WithSessionKey(sessionKey string) Opt {
	return func(vo *Midas) {
		vo.SessKey = sessionKey
	}
}

func WithAccToken(accToken string) Opt {
	return func(vo *Midas) {
		vo.AccToken = accToken
	}
}

// GetBalance POST https://api.weixin.qq.com/wxa/game/getbalance?access_token=ACCESS_TOKEN&signature=SIGNATURE&sig_method=SIG_METHOD&pay_sig=PAY_SIGNATURE
// 获取余额
func (v *Midas) GetBalance(req *BalanceRequest) (*BalanceReply, error) {
	// 检查 session key
	if err := v.CheckSessionKey(req.CommonRequest.OpenId); err != nil {
		return nil, ErrorIsSessionKeyInvalid
	}

	var res BalanceReply

	req.ZoneId = v.ZoneId
	req.Env = v.Env
	req.OfferId = v.OfferId
	req.Ts = time.Now().Unix()

	if err := v.Request(getBalance, req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// Pay POST https://api.weixin.qq.com/wxa/game/pay?access_token=ACCESS_TOKEN&signature=SIGNATURE&sig_method=SIG_METHOD&pay_sig=PAY_SIGNATURE
// 支付
func (v *Midas) Pay(req *PayRequest) (*PayReply, error) {
	// 检查 session key
	if err := v.CheckSessionKey(req.CommonRequest.OpenId); err != nil {
		return nil, ErrorIsSessionKeyInvalid
	}

	var res PayReply

	req.ZoneId = v.ZoneId
	req.Env = v.Env
	req.OfferId = v.OfferId
	req.Ts = time.Now().Unix()

	if err := v.Request(getBalance, req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// CheckSessionKey 检查 session key 是否过期
func (v *Midas) CheckSessionKey(openId string) error {
	at := v.AccToken
	signature := v.Signature("")
	method := signMethod
	query := fmt.Sprintf("access_token=%v&signature=%v&openid=%v&sig_method=%v", at, signature, openId, method)

	client := resty.New()
	resp, err := client.R().
		EnableTrace().
		SetHeader("Accept", "application/json").
		SetQueryString(query).
		Get(wxApi + apiPath[checkSessionKey])

	if err != nil {
		return err
	}
	fmt.Println(string(resp.Body()))

	var code CommonReply
	if err := json.Unmarshal(resp.Body(), &code); err != nil {
		return err
	}

	if code.ErrCode != CodeSuccess {
		return errors.New(fmt.Sprintf("[%v] %v", code.ErrCode, code.ErrMsg))
	}

	return nil
}

// GetQuery 拼接 url 包含支付签名 用户登录状态签名
func (v *Midas) GetQuery(uri, data string) string {
	at := v.AccToken
	signature := v.Signature(data)
	paySign := v.PaySign(uri, data)
	method := signMethod

	return fmt.Sprintf("access_token=%v&signature=%v&sig_method=%v&pay_sig=%v", at, signature, method, paySign)
}

// PaySign 根据 url 和 data 获取支付签名
func (v *Midas) PaySign(url, data string) string {
	str := url + "&" + data
	return v.GenerateSha256(v.AppKey, str)
}

// Signature 根据 session key 获取登录状态签名
func (v *Midas) Signature(data string) string {
	return v.GenerateSha256(v.SessKey, data)
}

// GenerateSha256 sha256 签名
func (v *Midas) GenerateSha256(key, data string) string {
	// Create a new HMAC by defining the hash type and the key (as byte array)
	h := hmac.New(sha256.New, []byte(key))
	// Write Data to it
	h.Write([]byte(data))
	// Get result and encode as hexadecimal string
	return hex.EncodeToString(h.Sum(nil))
}

// Request 请求微信接口
func (v *Midas) Request(urlType interType, body, res interface{}) error {
	str, err := json.Marshal(body)
	if err != nil {
		return err
	}

	data := string(str)
	uri := apiPath[urlType]
	query := v.GetQuery(uri, data)

	client := resty.New()
	resp, err := client.R().
		EnableTrace().
		SetHeader("Accept", "application/json").
		SetQueryString(query).
		SetBody(str).
		SetResult(res).
		Post(wxApi + uri)

	if err != nil {
		return err
	}
	fmt.Println(string(resp.Body()))

	var code CommonReply
	if err := json.Unmarshal(resp.Body(), &code); err != nil {
		return err
	}

	if code.ErrCode != CodeSuccess {
		return errors.New(fmt.Sprintf("[%v] %v", code.ErrCode, code.ErrMsg))
	}

	return nil
}
