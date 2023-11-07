package midas

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

// 测试签名是否和文档对应 文档地址：https://docs.qq.com/doc/DVUN0QWJja0J5c2x4
func TestSign(t *testing.T) {
	pb := `{"offer_id": "12345678", "openid": "oUrsfxxxxxxxxxx", "ts": 1668136271, "zone_id": "1", "env": 0}`
	appkey := "12345"
	uri := "/wxa/game/getbalance"

	data := uri + "&" + pb

	paySign := GenerateSha256(appkey, data)
	t.Log(paySign) // 11bac6388871d29c055c7d16fbe42e8d646855b666faf89b15c815218b1b23bd

	sessionKey := "9hAb/NEYUlkaMBEsmFgzig=="
	signature := GenerateSha256(sessionKey, pb)
	t.Log(signature) // 42fe1d3341fb1c8bd6f5014ba735ab04eacc80a2deb3ab4669eab4700b5b6729
}

func GenerateSha256(key, data string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
