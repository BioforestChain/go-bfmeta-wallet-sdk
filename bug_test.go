package sdk_test

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	sdk "github.com/BioforestChain/go-bfmeta-wallet-sdk"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/req/broadcastTra"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/req/createTransferAsset"
)

var bugSdkClient = sdk.NewLocalBCFWalletSDK(true)

var bugBCFSignUtil = bugSdkClient.NewBCFSignUtil("b")
var bugWallet = bugSdkClient.NewBCFWallet("35.213.66.234", 30003, "https://tracker.biw-meta.info/browser")

// var bugBCFSignUtil = bugSdkClient.NewBCFSignUtil("c")
// var bugWallet = bugSdkClient.NewBCFWallet("34.84.178.63", 19503, "https://qapmapi.pmchainbox.com/browser")

func sendTransactionBiw(t *testing.T, secret string, toAddr string, toAmount string) (success bool, err error) {
	bCFSignUtilCreateKeypair, _ := bugBCFSignUtil.CreateKeypair(secret)

	reqCreateTransferAsset := createTransferAsset.TransferAssetTransactionParams{
		TransactionCommonParamsWithRecipientId: createTransferAsset.TransactionCommonParamsWithRecipientId{
			TransactionCommonParams: createTransferAsset.TransactionCommonParams{
				PublicKey:        bCFSignUtilCreateKeypair.PublicKey,
				Fee:              "5000",
				ApplyBlockHeight: bugWallet.GetLastBlock().Result.Height,
			},
			RecipientId: toAddr, //钱包地址
		},
		Amount: toAmount,
	}
	reqCreateTransferAssetJson, _ := json.Marshal(reqCreateTransferAsset)
	log.Printf("reqCreateTransferAsset=%s", reqCreateTransferAssetJson)
	createTransferAssetResp, _ := bugWallet.CreateTransferAsset(reqCreateTransferAsset)
	if !createTransferAssetResp.Success {
		t.Errorf("createTransferAsset error=%v", createTransferAssetResp.Error)
		return false, nil
	}

	//// 3.3 生成签名
	detachedSign, _ := bugBCFSignUtil.DetachedSign(createTransferAssetResp.Result.Buffer.StringBuffer, bCFSignUtilCreateKeypair.SecretKey.StringBuffer)

	//// 3.4 bugWallet.BroadcastTransferAsset()
	req1 := broadcastTra.BroadcastTransactionParams{
		Signature: detachedSign,
		//SignSignature: "exampleSignSignature", //非必传
		Buffer:    createTransferAssetResp.Result.Buffer, //3.2 上面取得的buffer
		IsOnChain: true,
	}

	broadcastResult, err := bugWallet.BroadcastTransferAsset(req1)
	success = broadcastResult.Success

	return
}

func Test_bu(t *testing.T) {
	success, err := sendTransactionBiw(t, "ggggggggggggg", fmt.Sprintf("%sEAXDkaEJgWKMM61KYz2dYU1RfuxbB8Ma", bugBCFSignUtil.Prefix), "10000")
	log.Printf("success=%#v error=%#v", success, err)
}
