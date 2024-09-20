package sdk_test

import (
	"log"
	"testing"

	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/jbase"
	"github.com/stretchr/testify/assert"
)

var Msg = jbase.NewUtf8StringBuffer("utf8-1234")
var SecretKey = jbase.NewHexStringBuffer("03ac674216f3e15c761ee1a5e255f067953623c8b388b4459e13f978d7c846f4caf0f4c00cf9240771975e42b6672c88a832f98f01825dda6e001e2aab0bc0cc")
var PubKey = jbase.NewHexStringBuffer("caf0f4c00cf9240771975e42b6672c88a832f98f01825dda6e001e2aab0bc0cc")
var Secret = "1234"

func TestBCFSignUtil_CreateKeypair(t *testing.T) {
	bCFSignUtil_CreateKeypair, _ := bCFSignUtil.CreateKeypair(Secret)
	//{SecretKey:"03ac674216f3e15c761ee1a5e255f067953623c8b388b4459e13f978d7c846f4caf0f4c00cf9240771975e42b6672c88a832f98f01825dda6e001e2aab0bc0cc", PublicKey:"caf0f4c00cf9240771975e42b6672c88a832f98f01825dda6e001e2aab0bc0cc"}
	log.Printf("bCFSignUtil_CreateKeypair= %#v\n", bCFSignUtil_CreateKeypair)
	assert.Equal(t, SecretKey.Value, bCFSignUtil_CreateKeypair.SecretKey.Value)
	assert.Equal(t, PubKey.Value, bCFSignUtil_CreateKeypair.PublicKey.Value)
}

// 生成签名
func TestBCFSignUtil_DetachedSign(t *testing.T) {
	signature, err := bCFSignUtil.DetachedSign(Msg.StringBuffer, SecretKey.StringBuffer)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, "4a42fa5e984a54637d5e38dc0330551aa1a3c163c5e279fcab673ab058615225a8c71fd14f9340dcee25a1adfb15da72e959683a731d3cb8f6c5a1311350ee02", signature.Value)
}

func TestBCFSignUtil_DetachedVeriy(t *testing.T) {
	signature, err := bCFSignUtil.DetachedSign(Msg.StringBuffer, SecretKey.StringBuffer)
	if err != nil {
		panic(err)
	}
	verified, _ := bCFSignUtil.DetachedVerify(Msg.StringBuffer, signature.StringBuffer, PubKey.StringBuffer)

	assert.True(t, verified)
}

func TestBCFSignUtil_GetAddressFromPublicKey(t *testing.T) {
	var prefix = "c" //非必传
	address, _ := bCFSignUtil.GetAddressFromPublicKey(PubKey.StringBuffer, prefix)
	//cBUgBpP3mbJbVi7c9tYM8KJ7cd5Pgi5fmM
	log.Printf("AddressFromPublicKey= %#v\n", address)
	assert.Equal(t, "cBUgBpP3mbJbVi7c9tYM8KJ7cd5Pgi5fmM", address)
}

func TestBCFSignUtil_GetAddressFromSecret(t *testing.T) {
	address, _ := bCFSignUtil.GetAddressFromSecret(SecretKey.Value)
	//cZt3ajFJNZPC8zuQgAPvwcod5XPy8JS2S
	log.Printf("AddressFromSecret= %#v\n", address)
	assert.Equal(t, "cZt3ajFJNZPC8zuQgAPvwcod5XPy8JS2S", address)
}
