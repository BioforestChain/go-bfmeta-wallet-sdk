package sdk_test

import (
	sdk "github.com/BioforestChain/go-bfmeta-wallet-sdk"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/req/account"
	accountAssetEntityReq "github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/req/accountAsset"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/req/address"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/req/assetDetails"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/req/assets"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/req/block"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/req/broadcast"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/req/broadcastTra"
	createAccountReq "github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/req/createAccount"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/req/createTransferAsset"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/req/generateSecretReq"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/req/pkgTranscaction"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/req/transactions"
	"log"
	"testing"
)

var sdkClient = sdk.NewLocalBCFWalletSDK(true)

// var sdkClient = sdk.NewBCFWalletSDK()
var wallet = sdkClient.NewBCFWallet("34.84.178.63", 19503, "https://qapmapi.pmchainbox.com/browser")

func TestSdk(t *testing.T) {
	//getAddressBalance
	p := address.Params{
		"cEAXDkaEJgWKMM61KYz2dYU1RfuxbB8Ma",
		"XXVXQ",
		"PMC",
	}
	balance := wallet.GetAddressBalance(p)
	log.Printf("balance= %#v\n", balance)

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

	//getAccountInfo
	accountInfoReq := account.GetAccountInfoParams{
		"cEAXDkaEJgWKMM61KYz2dYU1RfuxbB8Ma",
	}
	accountInfo := wallet.GetAccountInfo(accountInfoReq)
	//accountInfo= accountResp.GetAccountInfoRespResult{Success:true, Result:accountResp.GetAccountInfoResp{Address:"cEAXDkaEJgWKMM61KYz2dYU1RfuxbB8Ma", PublicKey:"4bda2c5366b10e709c560e846e4041d355446c910dd6238e418092af5736c227", SecondPublicK
	//ey:"", IsDelegate:false, IsAcceptVote:false, AccountStatus:0, EquityInfo:accountResp.EquityInfo{Round:0, Equity:"", FixedEquity:""}}}
	log.Printf("accountInfo= %#v\n", accountInfo)

	//getAccountAsset
	accountAssetReq := accountAssetEntityReq.GetAccountAssetParams{
		"cEAXDkaEJgWKMM61KYz2dYU1RfuxbB8Ma",
	}
	accountAsset := wallet.GetAccountAsset(accountAssetReq)
	//accountAssetResp.GetAccountAssetRespResult{Success:true, Result:accountAssetResp.GetAccountAssetResp{Address:"cEAXDkaEJgWKMM61KYz2dYU1RfuxbB8Ma", Assets:accountAssetResp.AssetsMap{"XXVXQ":map[string]accountAssetResp.AssetDetail{"PMC":accountAssetResp.AssetDetail{Sour
	//ceChainMagic:"XXVXQ", AssetType:"PMC", SourceChainName:"paymetachain", AssetNumber:"1789879447994549065"}, "USDM":accountAssetResp.AssetDetail{SourceChainMagic:"XXVXQ", AssetType:"USDM", SourceChainName:"paymetachain", AssetNumber:"4949328785323"}}}, ForgingRewards:"10000541419", VotingRewards:""}}
	log.Printf("accountAsset= %#v\n", accountAsset)

	assetsReq := assets.PaginationOptions{
		1,
		2,
		"USDM",
	}
	asset := wallet.GetAssets(assetsReq)
	//assetsResp.GetAssetsRespResult{Success:false, Result:assetsResp.GetAssetsResp{Page:0, PageSize:0, Total:0, HasMore:false, DataList:[]assetsResp.GetAssetsData(nil)}}
	log.Printf("asset= %#v\n", asset)

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

	//GetTransactionsParams
	//todo 是否需要测试 Signature
	reqTra := transactions.GetTransactionsParams{
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

	//defer sdkClient.Close()
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
}

func TestGetLastBlock(t *testing.T) {
	lastBlock := wallet.GetLastBlock()
	log.Printf("lastBlock= %#v\n", lastBlock)
}

func TestBroadcastCompleteTransaction(t *testing.T) {
	//broadcastCompleteTransaction
	reqBroadcastCompleteTransaction := broadcast.Params{
		"key": 123,
		//"key1": []string{"item1", "item2"},
	}
	bCTResp := wallet.BroadcastCompleteTransaction(reqBroadcastCompleteTransaction)
	log.Printf("bCTResp= %#v\n", bCTResp)
}

func TestCreateTransferAsset(t *testing.T) {
	reqCreateTransferAsset := createTransferAsset.TransferAssetTransactionParams{
		TransactionCommonParamsWithRecipientId: createTransferAsset.TransactionCommonParamsWithRecipientId{
			TransactionCommonParams: createTransferAsset.TransactionCommonParams{
				PublicKey:        "examplePublicKey",
				Fee:              "0.1",
				ApplyBlockHeight: 100,
				Remark: map[string]string{
					"note": "example transaction",
				},
				BinaryInfos: []createTransferAsset.KVStorageInfo{
					{
						Key: "exampleKey",
						FileInfo: createTransferAsset.FileInfo{
							Name: "exampleFile",
							Size: 1234,
						},
					},
				},
				Timestamp: 1622732931,
			},
			RecipientId: "exampleRecipientId",
		},
		SourceChainMagic: "exampleSourceChainMagic",
		SourceChainName:  "exampleSourceChainName",
		AssetType:        "exampleAssetType",
		Amount:           "10.0",
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

func TestBroadcastTransferAsset(t *testing.T) {
	req := broadcastTra.BroadcastTransactionParams{
		Signature:     "exampleSignature",
		SignSignature: "exampleSignSignature",
		Buffer:        "exampleBuffer",
		IsOnChain:     true,
	}
	broadcastTransferAsset, _ := wallet.BroadcastTransferAsset(req)
	log.Printf("broadcastTransferAsset= %#v\n", broadcastTransferAsset)
}
