package sdk_test

import (
	"log"
	"testing"
	"time"

	sdk "github.com/BioforestChain/go-bfmeta-wallet-sdk"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/jbase"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/req/account"
	accountAssetEntityReq "github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/req/accountAsset"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/req/address"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/req/assetDetails"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/req/assets"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/req/asymmetricDecrypt"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/req/block"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/req/broadcast"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/req/broadcastTra"
	createAccountReq "github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/req/createAccount"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/req/createTransferAsset"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/req/generateSecretReq"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/req/pkgTranscaction"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/req/transactions"
)

var sdkClient = sdk.NewLocalBCFWalletSDK(true)

// var sdkClient = sdk.NewBCFWalletSDK()
var wallet = sdkClient.NewBCFWallet("34.84.178.63", 19503, "https://qapmapi.pmchainbox.com/browser")

var bCFSignUtil = sdkClient.NewBCFSignUtil("c")

func TestSdk(t *testing.T) {
	//getAddressBalance
	p := address.Params{
		Address:   "cEAXDkaEJgWKMM61KYz2dYU1RfuxbB8Ma",
		Magic:     "XXVXQ",
		AssetType: "PMC",
	}
	balance := wallet.GetAddressBalance(p)
	log.Printf("balance= %#v\n", balance)

	//getAccountInfo
	accountInfoReq := account.GetAccountInfoParams{
		Address: "cEAXDkaEJgWKMM61KYz2dYU1RfuxbB8Ma",
	}
	accountInfo := wallet.GetAccountInfo(accountInfoReq)
	//accountInfo= accountResp.GetAccountInfoRespResult{Success:true, Result:accountResp.GetAccountInfoResp{Address:"cEAXDkaEJgWKMM61KYz2dYU1RfuxbB8Ma", PublicKey:"4bda2c5366b10e709c560e846e4041d355446c910dd6238e418092af5736c227", SecondPublicK
	//ey:"", IsDelegate:false, IsAcceptVote:false, AccountStatus:0, EquityInfo:accountResp.EquityInfo{Round:0, Equity:"", FixedEquity:""}}}
	log.Printf("accountInfo= %#v\n", accountInfo)
	//// 1. 获取区块高度
	//wallet.GetLastBlock()
	//// 2. 获取余额. 需要传地址 magic 币名
	//// 如果需要magic 就从 区块里面取magic
	////
	//wallet.GetAddressBalance()
	//// 这个函数只需要传地址,但是要自己解析里面的多个币名
	//wallet.GetAccountAsset()
	//// 3. 做转账
	//// 3.1 生成公私钥对
	//// 3.2
	//wallet.CreateTransferAsset()
	//// 3.3 签名 ts const signature = (await bfmetaSDK.bfchainSignUtil.detachedSign(bytes, keypair.secretKey)).toString("hex");
	//
	//// 3.4
	//wallet.BroadcastTransferAsset()

	defer sdkClient.Close()
}

func TestSdkAsset(t *testing.T) {
	//getAccountAsset
	accountAssetReq := accountAssetEntityReq.GetAccountAssetParams{
		Address: "cEAXDkaEJgWKMM61KYz2dYU1RfuxbB8Ma",
	}
	accountAsset := wallet.GetAccountAsset(accountAssetReq)
	//accountAssetResp.GetAccountAssetRespResult{Success:true, Result:accountAssetResp.GetAccountAssetResp{Address:"cEAXDkaEJgWKMM61KYz2dYU1RfuxbB8Ma", Assets:accountAssetResp.AssetsMap{"XXVXQ":map[string]accountAssetResp.AssetDetail{"PMC":accountAssetResp.AssetDetail{Sour
	//ceChainMagic:"XXVXQ", AssetType:"PMC", SourceChainName:"paymetachain", AssetNumber:"1789879447994549065"}, "USDM":accountAssetResp.AssetDetail{SourceChainMagic:"XXVXQ", AssetType:"USDM", SourceChainName:"paymetachain", AssetNumber:"4949328785323"}}}, ForgingRewards:"10000541419", VotingRewards:""}}
	log.Printf("accountAsset= %#v\n", accountAsset)

	//getAllAccountAsset
	allAccountAssetReq := accountAssetEntityReq.GetAllAccountAssetReq{
		Filter: map[string]string{
			"assetType": "USDM",
		},
	}
	allAccountAsset := wallet.GetAllAccountAsset(allAccountAssetReq)
	//accountAssetResp.GetAllAccountAssetRespResult{Success:true, Result:accountAssetResp.GetAllAccountAssetResp{map[string]map[string]string{"cCET2Sxt2LPDhx44wxJ9uhkpviKNrSacvE":map[string]string{}}...
	log.Printf("allAccountAsset= %#v\n", allAccountAsset)

	//getAssetDetails
	assetDetailsReq := assetDetails.Req{
		"USDM",
	}
	assetDetails := wallet.GetAssetDetails(assetDetailsReq)
	//assetDetailsResp.GetAssetDetailsRespResult{Success:true, Result:assetDetailsResp.GetAssetDetailsResp{AssetInfo:assetDetailsResp.AssetInfo{Asset:assetDetailsResp.Asset{AssetType:"USDM", ApplyAddress:"cKFyTV2yNmCxdsnoLSbT25zKTYVa4kHv1e", GenesisAddress:"cEAXDkaEJgWKMM6
	//1KYz2dYU1RfuxbB8Ma", SourceChainName:"paymetachain", IssuedAssetPrealnum:"1000000000000000000", RemainAssetPrealnum:"1000000000000000000", FrozenMainAssetPrealnum:"100032000498000", PublishTime:1715325195000, SourceChainMagic:"XXVXQ"}, AddressQty:237}, IconURL:"https://bfm-fonts-cdn.oss-cn-hongkong.a
	//liyuncs.com/meta-icon/pmc/icon-USDM.png"}}
	log.Printf("assetDetails= %#v\n", assetDetails)
}

func TestGetBlock(t *testing.T) {
	p := block.GetBlockParams{
		//test 不传
		//Signature: "abc123",
		Height:   10,
		Page:     1,
		PageSize: 20,
	}
	block := wallet.GetBlock(p)
	log.Printf("block= %#v\n", block)
	defer sdkClient.Close()
}

func TestGetLastBlock(t *testing.T) {
	lastBlock := wallet.GetLastBlock()
	//
	log.Printf("lastBlock= %#v\n", lastBlock)
}

// /
func TestSomeTransactionEvent(t *testing.T) {
	//GetTransactionsParams
	req := transactions.GetTransactionsParams{
		Signature:    "exampleSignature",
		Height:       12345,
		MinHeight:    10000,
		MaxHeight:    20000,
		SenderId:     "exampleSenderId",
		RecipientId:  "exampleRecipientId",
		Address:      "exampleAddress",
		Type:         []string{"transfer", "stake"},
		StorageValue: "exampleStorageValue",
		Page:         1,
		PageSize:     50,
		Sort:         1,
	}
	transactionsByBrowse, _ := wallet.GetTransactionsByBrowser(req)
	log.Printf("transactionsByBrowse= %#v\n", transactionsByBrowse)

	assetsReq := assets.PaginationOptions{
		1,
		2,
		"PMC",
	}
	asset := wallet.GetAssets(assetsReq)
	//assetsResp.GetAssetsRespResult{Success:false, Result:assetsResp.GetAssetsResp{Page:0, PageSize:0, Total:0, HasMore:false, DataList:[]assetsResp.GetAssetsData(nil)}}
	log.Printf("asset= %#v\n", asset)

	//GetTransactionsParams
	//目前测试 不传Signature
	reqTra := transactions.GetTransactionsParams{
		//Signature:    "exampleSignature",
		Height:       12345,
		MinHeight:    10000,
		MaxHeight:    20000,
		SenderId:     "exampleSenderId",
		RecipientId:  "exampleRecipientId",
		Address:      "exampleAddress",
		Type:         []string{"transfer", "stake"},
		StorageValue: "exampleStorageValue",
		Page:         1,
		PageSize:     50,
		Sort:         1,
	}

	transactionsResp := wallet.GetTransactions(reqTra)
	log.Printf("transactionsResp= %#v\n", transactionsResp)

	//generateSecret
	reqGenSecret := generateSecretReq.GenerateSecretParams{
		Lang: "en",
	}
	secretResp := wallet.GenerateSecret(reqGenSecret)
	log.Printf("secretResp= %#v\n", secretResp)

	//createAccount
	reqCreateAccount := createAccountReq.CreateAccountReq{
		Secret: "xxxxxxxxxxxxxxxxxxxxxxx",
	}
	createAccountResp := wallet.CreateAccount(reqCreateAccount)
	log.Printf("createAccountResp= %#v\n", createAccountResp)

}

// todo test params
// https://qatracker.pmchainbox.info/#/info/transaction-data/2d0cea07ab73be6bdab258f12e7e0aa22776a8b9dd7b130f33fdd8fce6534cb0e29bc8d4983d3564178ae4189eedba80a864bda1a4ceb8b197e530ef1774ea07
// params transaction 这个结构
// {"success":false,"error":{"message":"fromMagic in body is required","code":"001-00002"},"minFee":972}
func TestBroadcastCompleteTransaction(t *testing.T) {
	//broadcastCompleteTransaction
	reqBroadcastCompleteTransaction := broadcast.Params{
		"version":              1,
		"type":                 "PMC-PAYMETACHAIN-AST-02",
		"senderId":             "c6C9ycTXrPBu8wXAGhUJHau678YyQwB2Mn",
		"senderPublicKey":      "0d3c8003248cc4c71493dd67c0c433e75b7a191758df94fb0be5db2c6a94fecd",
		"fee":                  "100000",
		"timestamp":            31839601,
		"applyBlockHeight":     114208,
		"effectiveBlockHeight": 114258,
		"signature":            "2d0cea07ab73be6bdab258f12e7e0aa22776a8b9dd7b130f33fdd8fce6534cb0e29bc8d4983d3564178ae4189eedba80a864bda1a4ceb8b197e530ef1774ea07",
		"asset": map[string]interface{}{
			"transferAsset": map[string]interface{}{
				"sourceChainName":  "paymetachain",
				"sourceChainMagic": "XXVXQ",
				"assetType":        "PMC",
				"amount":           "185184",
			},
		},
		"rangeType": 0,
		"range":     []string{},
		"fromMagic": "xxx",
		"toMagic":   "zzz",
		"remark": map[string]string{
			"orderId": "110b45fafcb84cb7a1de7eef5a957855",
		},
		"recipientId":  "cFqv1tiifgYE6xbhZp43XxbZVJp363BWXt",
		"storageKey":   "assetType",
		"storageValue": "PMC",
		//"key1":       []string{"item1", "item2"},
	}
	bCTResp, _ := wallet.BroadcastCompleteTransaction(reqBroadcastCompleteTransaction)
	log.Printf("bCTResp= %#v\n", bCTResp)
}

// Se 03ac674216f3e15c761ee1a5e255f067953623c8b388b4459e13f978d7c846f4caf0f4c00cf9240771975e42b6672c88a832f98f01825dda6e001e2aab0bc0cc
// Pu caf0f4c00cf9240771975e42b6672c88a832f98f01825dda6e001e2aab0bc0cc

func TestBroadcastTransferAsset(t *testing.T) {
	//助记词
	bCFSignUtil_CreateKeypair, _ := bCFSignUtil.CreateKeypair(Secret)
	buffer := jbase.NewBase64StringBuffer("123456")
	sign, _ := bCFSignUtil.DetachedSign(buffer.StringBuffer, bCFSignUtil_CreateKeypair.SecretKey.StringBuffer)
	req := broadcastTra.BroadcastTransactionParams{
		Signature: sign,
		Buffer:    *buffer,
		IsOnChain: true,
	}
	broadcastTransferAsset, _ := wallet.BroadcastTransferAsset(req)
	log.Printf("broadcastTransferAsset= %#v\n", broadcastTransferAsset)
}

func TestCreateTransferAsset(t *testing.T) {
	reqCreateTransferAsset := createTransferAsset.TransferAssetTransactionParams{
		TransactionCommonParamsWithRecipientId: createTransferAsset.TransactionCommonParamsWithRecipientId{
			TransactionCommonParams: createTransferAsset.TransactionCommonParams{
				PublicKey:        *PubKey,
				Fee:              "0.1",
				ApplyBlockHeight: 100,
			},
			RecipientId: "exampleRecipientId",
		},
		SourceChainMagic: "exampleSourceChainMagic",
		SourceChainName:  "exampleSourceChainName",
		AssetType:        "exampleAssetType",
		Amount:           "100",
	}
	createTransferAssetResp, _ := wallet.CreateTransferAsset(reqCreateTransferAsset)
	log.Printf("createTransferAssetResp= %#v\n", createTransferAssetResp)
}

func TestPackageTransferAsset(t *testing.T) {
	req := pkgTranscaction.PackageTransacationParams{
		Signature: "exampleSignature",
		Buffer:    "exampleBuffer",
	}
	pkgTransferAssetResp, _ := wallet.PackageTransferAsset(req)
	log.Printf("pkgTransferAssetResp= %#v\n", pkgTransferAssetResp)
}

func TestBCFSignUtil_CreateKeypairBySecretKey(t *testing.T) {
	bCFSignUtil_CreateKeypairBySecretKey, _ := bCFSignUtil.CreateKeypairBySecretKey(jbase.NewHexStringBuffer("a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3a4465fd76c16fcc458448076372abf1912cc5b150663a64dffefe550f96feadd").StringBuffer)
	//bCFSignUtil_CreateKeypair= sdk.ResKeyPair{SecretKey:"a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3a4465fd76c16fcc458448076372abf1912cc5b150663a64dffefe550f96feadd", PublicKey:"a4465fd76c16fcc458448076372abf1912cc5b150663a64dffefe550f96feadd"}
	//--- PASS: TestBCFSignUtil_CreateKeypairBySecretKey (0.01s)
	log.Printf("bCFSignUtil_CreateKeypairBySecretKey= %#v\n", bCFSignUtil_CreateKeypairBySecretKey)
}
func TestBCFSignUtil_GetPublicKeyFromSecret(t *testing.T) {
	bCFSignUtil_GetPublicKeyFromSecret, _ := bCFSignUtil.GetPublicKeyFromSecret("123456")
	//bCFSignUtil_GetPublicKeyFromSecret= "0363649faf7a83d0bc0d9faa9c6a5efa8adc772190b8072210bc825895ca3570"
	log.Printf("bCFSignUtil_GetPublicKeyFromSecret= %#v\n", bCFSignUtil_GetPublicKeyFromSecret)
}

func TestBCFSignUtil_CreateSecondKeypair(t *testing.T) {
	var s = "a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3a4465fd76c16fcc458448076372abf1912cc5b150663a64dffefe550f96feadd"
	var ss = "12345678"
	got, _ := bCFSignUtil.CreateSecondKeypair(s, ss)
	//CreateSecondKeypair= sdk.ResKeyPair{SecretKey:"9d3292b245d0196e9e2ea7f636b25a84bf518c86ee2af87cb476f754dbf4407dbb3d939c1d91e95154c8ec5683e981865e0baa3cbaa25bd382f1bde5b693306d", PublicKey:"bb3d939c1d91e95154c8ec5683e981865e0baa3cbaa25bd382f1bde5b693306d"}
	log.Printf("CreateSecondKeypair= %#v\n", got)
}
func TestBCFSignUtil_GetSecondPublicKeyFromSecretAndSecondSecret(t *testing.T) {
	var s = "a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3a4465fd76c16fcc458448076372abf1912cc5b150663a64dffefe550f96feadd"
	var ss = "12345678"
	got, _ := bCFSignUtil.GetSecondPublicKeyFromSecretAndSecondSecret(s, ss)
	//GetSecondPublicKeyFromSecretAndSecondSecret= sdk.ResPubKeyPair{PublicKey:"bb3d939c1d91e95154c8ec5683e981865e0baa3cbaa25bd382f1bde5b693306d"}
	//--- PASS: TestBCFSignUtil_GetSecondPublicKeyFromSecretAndSecondSecret (0.02s)
	//PASS
	log.Printf("GetSecondPublicKeyFromSecretAndSecondSecret= %#v\n", got)
}
func TestBCFSignUtil_GetSecondPublicKeyStringFromSecretAndSecondSecret(t *testing.T) {
	var s = "a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3a4465fd76c16fcc458448076372abf1912cc5b150663a64dffefe550f96feadd"
	var ss = "12345678"
	got, _ := bCFSignUtil.GetSecondPublicKeyStringFromSecretAndSecondSecret(s, ss)
	//GetSecondPublicKeyFromSecretAndSecondSecret= "bb3d939c1d91e95154c8ec5683e981865e0baa3cbaa25bd382f1bde5b693306d"
	//--- PASS: TestBCFSignUtil_GetSecondPublicKeyStringFromSecretAndSecondSecret (0.02s)
	log.Printf("GetSecondPublicKeyFromSecretAndSecondSecret= %#v\n", got)
}

func TestBCFSignUtil_AsymmetricEncrypt(t *testing.T) {
	var msg = []byte("123")
	var decryptPK = []byte("801e19ac714803ca50d53ba802667adc99f82c21bf4b5dfbbfd0c4b766103af1cf6c56944124bd9f219b1910135469796b817fefe5abb01aabc8df8772495a02")
	var encryptSK = []byte("a4465fd76c16fcc458448076372abf1912cc5b150663a64dffefe550f96feadd")
	got, _ := bCFSignUtil.AsymmetricEncrypt(msg, decryptPK, encryptSK)
	//{"encryptedMessage":"117,56,228,3,87,171,27,39,24,162,27,204,28,18,218,165,44","nonce":"0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0"}
	log.Printf("AsymmetricEncrypt= %#v\n", got)
}
func TestBCFSignUtil_AsymmetricDecrypt(t *testing.T) {
	var decryptPK = []byte("801e19ac714803ca50d53ba802667adc99f82c21bf4b5dfbbfd0c4b766103af1cf6c56944124bd9f219b1910135469796b817fefe5abb01aabc8df8772495a02")
	var encryptSK = []byte("a4465fd76c16fcc458448076372abf1912cc5b150663a64dffefe550f96feadd")
	var req = asymmetricDecrypt.Req{
		DecryptSK:        decryptPK,
		Nonce:            []byte("1"),
		EncryptedMessage: []byte("117,56,228,3,87,171,27,39,24,162,27,204,28,18,218,165,44"),
		EncryptPK:        encryptSK,
	}
	got, _ := bCFSignUtil.AsymmetricDecrypt(req)
	log.Printf("AsymmetricDecrypt= %#v\n", got)
}

// checkSecondSecret
// 校验二次密码公钥是否正确
// Params:
// secret – 主密码
// secondSecret – 二次密码
// secondPublicKey – 二次密码公钥
func TestBCFSignUtil_CheckSecondSecret(t *testing.T) {
	var secret, secondSecret, secondPublicKey = "a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3a4465fd76c16fcc458448076372abf1912cc5b150663a64dffefe550f96feadd", "12345678", "bb3d939c1d91e95154c8ec5683e981865e0baa3cbaa25bd382f1bde5b693306d"
	got, _ := bCFSignUtil.CheckSecondSecret(secret, secondSecret, secondPublicKey)
	// CheckSecondSecret= true
	log.Printf("CheckSecondSecret= %#v\n", got)
}

// checkSecondSecret
// 校验二次密码公钥是否正确
// Params:
// secret – 主密码
// secondSecret – 二次密码
// secondPublicKey – 二次密码公钥
func TestBCFSignUtil_CheckSecondSecretV2(t *testing.T) {
	var secret, secondSecret, secondPublicKey = "a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3a4465fd76c16fcc458448076372abf1912cc5b150663a64dffefe550f96feadd", "12345678", "bb3d939c1d91e95154c8ec5683e981865e0baa3cbaa25bd382f1bde5b693306d"
	got, _ := bCFSignUtil.CheckSecondSecretV2(secret, secondSecret, secondPublicKey)
	// CheckSecondSecretV2= true
	log.Printf("CheckSecondSecretV2= %#v\n", got)

}

func TestBCFSignUtil_CreateSecondKeypairV2(t *testing.T) {
	var secret, secondSecret = "a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3a4465fd76c16fcc458448076372abf1912cc5b150663a64dffefe550f96feadd", "12345678"
	got, _ := bCFSignUtil.CreateSecondKeypairV2(secret, secondSecret)
	// createSecondKeypairV2= sdk.ResKeyPair{SecretKey:"645fc86050eaa146ee8c0117adfee3a7125580dd2978d1e6d4cbf35b8aa2b19e1bc79b077e3476354f845cf3879a1d9a6e3254f9866450ec5d6c00c83268319e", PublicKey:"1bc79b077e3476354f845cf3879a1d9a6e3254f9866450ec5d6c00c83268319e"}
	log.Printf("createSecondKeypairV2= %#v\n", got)
}

func TestBCFSignUtil_GetSecondPublicKeyFromSecretAndSecondSecretV2(t *testing.T) {
	var s = "a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3a4465fd76c16fcc458448076372abf1912cc5b150663a64dffefe550f96feadd"
	var ss = "12345678"
	got, _ := bCFSignUtil.GetSecondPublicKeyFromSecretAndSecondSecretV2(s, ss)
	//GetSecondPublicKeyFromSecretAndSecondSecretV2= sdk.ResPubKeyPair{PublicKey:"1bc79b077e3476354f845cf3879a1d9a6e3254f9866450ec5d6c00c83268319e"}
	log.Printf("GetSecondPublicKeyFromSecretAndSecondSecretV2= %#v\n", got)
}

func TestMultiSdk(t *testing.T) {
	for i := 0; i < 10; i++ {
		_singleSdk(t)
	}
	for i := 0; i < 20; i++ {
		time.Sleep(1 * time.Second)
		log.Printf("look~~ %d", i)
	}
}
func _singleSdk(t *testing.T) {
	var sdkClient = sdk.NewLocalBCFWalletSDK(true)
	var bCFSignUtil = sdkClient.NewBCFSignUtil("c")
	defer sdkClient.Close()

	var signature, _ = bCFSignUtil.DetachedSign(Msg.StringBuffer, SecretKey.StringBuffer)
	got, _ := bCFSignUtil.DetachedVerify(Msg.StringBuffer, signature.StringBuffer, PubKey.StringBuffer)
	log.Printf("DetachedVeriy= %#v\n", got)
}
