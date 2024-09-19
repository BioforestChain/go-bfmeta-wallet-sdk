package sdk_test

import (
	"log"
	"testing"

	sdk "github.com/BioforestChain/go-bfmeta-wallet-sdk"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/req/broadcastTra"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/req/createTransferAsset"
)

var bugSdkClient = sdk.NewLocalBCFWalletSDK(true)
var bugBCFSignUtil = bugSdkClient.NewBCFSignUtil("b")
var bugWallet = bugSdkClient.NewBCFWallet("35.213.66.234", 30003, "https://tracker.biw-meta.info/browser")

func sendTransactionBiw(secret string, toAddr string, toAmount string) (bool, error) {
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
	createTransferAssetResp, _ := bugWallet.CreateTransferAsset(reqCreateTransferAsset)

	//// 3.3 生成签名
	var s1 = []byte(createTransferAssetResp.Result.Buffer)
	var ss = []byte(bCFSignUtilCreateKeypair.SecretKey)
	detachedSign, _ := bugBCFSignUtil.DetachedSignToHex(s1, ss)

	//// 3.4 bugWallet.BroadcastTransferAsset()
	req1 := broadcastTra.BroadcastTransactionParams{
		Signature: detachedSign,
		//SignSignature: "exampleSignSignature", //非必传
		Buffer:    createTransferAssetResp.Result.Buffer, //3.2 上面取得的buffer
		IsOnChain: true,
	}

	var (
		err error
	)
	success, err := bugWallet.BroadcastTransferAsset(req1)

	return success.Success, err
}

func Test_bu(t *testing.T) {
	success, err := sendTransactionBiw("qaq", "bEAXDkaEJgWKMM61KYz2dYU1RfuxbB8Ma", "10000")
	log.Printf("success=%#v error=%#v", success, err)
}
