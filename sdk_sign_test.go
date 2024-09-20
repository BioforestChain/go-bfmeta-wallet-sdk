package sdk_test

import (
	"log"
	"testing"

	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/jbase"
)

var Msg = jbase.NewUtf8StringBuffer("utf8-1234")
var SecretKey = jbase.NewHexStringBuffer("03ac674216f3e15c761ee1a5e255f067953623c8b388b4459e13f978d7c846f4caf0f4c00cf9240771975e42b6672c88a832f98f01825dda6e001e2aab0bc0cc")
var PubKey = jbase.NewHexStringBuffer("caf0f4c00cf9240771975e42b6672c88a832f98f01825dda6e001e2aab0bc0cc")
var Secret = "1234"

func TestBCFSignUtil_CreateKeypair(t *testing.T) {
	bCFSignUtil_CreateKeypair, _ := bCFSignUtil.CreateKeypair(Secret)
	//{SecretKey:"03ac674216f3e15c761ee1a5e255f067953623c8b388b4459e13f978d7c846f4caf0f4c00cf9240771975e42b6672c88a832f98f01825dda6e001e2aab0bc0cc", PublicKey:"caf0f4c00cf9240771975e42b6672c88a832f98f01825dda6e001e2aab0bc0cc"}
	log.Printf("bCFSignUtil_CreateKeypair= %#v\n", bCFSignUtil_CreateKeypair)
}

/**
 * 签名并且转成 hex 字符串
 *
 * @param message
 * @param secretKey
 * @returns
 */
//signToString(message: Uint8Array, secretKey: Uint8Array): Promise<string>;
func TestBCFSignUtil_SignToString(t *testing.T) {
	got, _ := bCFSignUtil.SignToString(Msg.StringBuffer, SecretKey.StringBuffer)
	log.Printf("SignToString= %#v\n", got)
}

// 生成签名
func TestBCFSignUtil_DetachedSign(t *testing.T) {
	signature, err := bCFSignUtil.DetachedSign(Msg.StringBuffer, SecretKey.StringBuffer)
	if err != nil {
		panic(err)
	}
	log.Printf("DetachedSign= %#v\n", signature)
}

func TestBCFSignUtil_DetachedVeriy(t *testing.T) {
	signature, err := bCFSignUtil.DetachedSign(Msg.StringBuffer, SecretKey.StringBuffer)
	if err != nil {
		panic(err)
	}
	got, _ := bCFSignUtil.DetachedVeriy(Msg.StringBuffer, signature.StringBuffer, PubKey.StringBuffer)
	// true
	log.Printf("DetachedVeriy= %#v\n", got)
}

func TestBCFSignUtil_GetAddressFromPublicKey(t *testing.T) {
	var prefix = "c" //非必传
	got, _ := bCFSignUtil.GetAddressFromPublicKey(PubKey.StringBuffer, prefix)
	//cBUgBpP3mbJbVi7c9tYM8KJ7cd5Pgi5fmM
	log.Printf("AddressFromPublicKey= %#v\n", got)
}

func TestBCFSignUtil_GetAddressFromPublicKeyString(t *testing.T) {
	var prefix = "c" //非必传
	got, _ := bCFSignUtil.GetAddressFromPublicKeyString(PubKey.StringBuffer, prefix)
	//cBUgBpP3mbJbVi7c9tYM8KJ7cd5Pgi5fmM
	log.Printf("AddressFromPublicKeyString= %#v\n", got)
}

func TestBCFSignUtil_GetAddressFromSecret(t *testing.T) {
	got, _ := bCFSignUtil.GetAddressFromSecret(SecretKey.Value)
	//cZt3ajFJNZPC8zuQgAPvwcod5XPy8JS2S
	log.Printf("AddressFromSecret= %#v\n", got)
}
