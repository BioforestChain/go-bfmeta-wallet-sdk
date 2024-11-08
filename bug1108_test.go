package sdk_test

import (
	"fmt"
	"testing"

	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/jbase"
)

func TestBug1108(t *testing.T) {
	msg := "login"
	sigHex := "5954b23402e69ab7121dd6a9f70910620d54d581fd66e033405e02f9c259246e44fe7e804b365a08e87d4f73d04b20bfb01398588c8bedb5cc244478ae107800"
	publicKey := "68a75f6390a98a0bb23b6293f0b47b451ed5529b8a20839b3dfe75e80b739991"
	// 验证签名
	msgTmp := jbase.NewUtf8StringBuffer(msg)
	sigHexTmp := jbase.NewUtf8StringBuffer(sigHex)
	publicKeyTmp := jbase.NewUtf8StringBuffer(publicKey)
	fmt.Println(msg, sigHex, publicKey)
	res, err := bCFSignUtil.DetachedVerify(msgTmp.StringBuffer, sigHexTmp.StringBuffer, publicKeyTmp.StringBuffer)
	fmt.Println(res, err)
}
