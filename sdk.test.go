package main

import (
	"bfmeta-wallet-bcf/entity/req/account"
	accountAssetEntityReq "bfmeta-wallet-bcf/entity/req/accountAsset"
	"bfmeta-wallet-bcf/entity/req/address"
	"bfmeta-wallet-bcf/entity/req/assetDetails"
	"bfmeta-wallet-bcf/entity/req/assets"
	"bfmeta-wallet-bcf/entity/req/broadcast"
	createAccountReq "bfmeta-wallet-bcf/entity/req/createAccount"
	"bfmeta-wallet-bcf/entity/req/generateSecretReq"
	"bfmeta-wallet-bcf/entity/req/transactions"
	"log"
)

func main() {
	sdk := newBCFWalletSDK()
	wallet := sdk.newBCFWallet("34.84.178.63", 19503, "https://qapmapi.pmchainbox.com/browser")
	//getAddressBalance
	p := address.Params{
		"cEAXDkaEJgWKMM61KYz2dYU1RfuxbB8Ma",
		"XXVXQ",
		"PMC",
	}
	balance := wallet.getAddressBalance(p)
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
	transactionsByBrowse, _ := wallet.getTransactionsByBrowser(req)
	log.Printf("transactionsByBrowse= %#v\n", transactionsByBrowse)

	//getAccountInfo
	accountInfoReq := account.GetAccountInfoParams{
		"cEAXDkaEJgWKMM61KYz2dYU1RfuxbB8Ma",
	}
	accountInfo := wallet.getAccountInfo(accountInfoReq)
	//accountInfo= accountResp.GetAccountInfoRespResult{Success:true, Result:accountResp.GetAccountInfoResp{Address:"cEAXDkaEJgWKMM61KYz2dYU1RfuxbB8Ma", PublicKey:"4bda2c5366b10e709c560e846e4041d355446c910dd6238e418092af5736c227", SecondPublicK
	//ey:"", IsDelegate:false, IsAcceptVote:false, AccountStatus:0, EquityInfo:accountResp.EquityInfo{Round:0, Equity:"", FixedEquity:""}}}
	log.Printf("accountInfo= %#v\n", accountInfo)

	//getAccountAsset
	accountAssetReq := accountAssetEntityReq.GetAccountAssetParams{
		"cEAXDkaEJgWKMM61KYz2dYU1RfuxbB8Ma",
	}
	accountAsset := wallet.getAccountAsset(accountAssetReq)
	//accountAssetResp.GetAccountAssetRespResult{Success:true, Result:accountAssetResp.GetAccountAssetResp{Address:"cEAXDkaEJgWKMM61KYz2dYU1RfuxbB8Ma", Assets:accountAssetResp.AssetsMap{"XXVXQ":map[string]accountAssetResp.AssetDetail{"PMC":accountAssetResp.AssetDetail{Sour
	//ceChainMagic:"XXVXQ", AssetType:"PMC", SourceChainName:"paymetachain", AssetNumber:"1789879447994549065"}, "USDM":accountAssetResp.AssetDetail{SourceChainMagic:"XXVXQ", AssetType:"USDM", SourceChainName:"paymetachain", AssetNumber:"4949328785323"}}}, ForgingRewards:"10000541419", VotingRewards:""}}
	log.Printf("accountAsset= %#v\n", accountAsset)

	assetsReq := assets.PaginationOptions{
		1,
		2,
		"USDM",
	}
	asset := wallet.getAssets(assetsReq)
	//assetsResp.GetAssetsRespResult{Success:false, Result:assetsResp.GetAssetsResp{Page:0, PageSize:0, Total:0, HasMore:false, DataList:[]assetsResp.GetAssetsData(nil)}}
	log.Printf("asset= %#v\n", asset)

	//getAllAccountAsset
	allAccountAssetReq := accountAssetEntityReq.GetAllAccountAssetReq{
		Filter: map[string]string{
			"assetType": "USDM",
		},
	}
	allAccountAsset := wallet.getAllAccountAsset(allAccountAssetReq)
	//accountAssetResp.GetAllAccountAssetRespResult{Success:true, Result:accountAssetResp.GetAllAccountAssetResp{map[string]map[string]string{"cCET2Sxt2LPDhx44wxJ9uhkpviKNrSacvE":map[string]string{}}...
	log.Printf("allAccountAsset= %#v\n", allAccountAsset)

	//getAssetDetails
	assetDetailsReq := assetDetails.Req{
		"USDM",
	}
	assetDetails := wallet.getAssetDetails(assetDetailsReq)
	//assetDetailsResp.GetAssetDetailsRespResult{Success:true, Result:assetDetailsResp.GetAssetDetailsResp{AssetInfo:assetDetailsResp.AssetInfo{Asset:assetDetailsResp.Asset{AssetType:"USDM", ApplyAddress:"cKFyTV2yNmCxdsnoLSbT25zKTYVa4kHv1e", GenesisAddress:"cEAXDkaEJgWKMM6
	//1KYz2dYU1RfuxbB8Ma", SourceChainName:"paymetachain", IssuedAssetPrealnum:"1000000000000000000", RemainAssetPrealnum:"1000000000000000000", FrozenMainAssetPrealnum:"100032000498000", PublishTime:1715325195000, SourceChainMagic:"XXVXQ"}, AddressQty:237}, IconURL:"https://bfm-fonts-cdn.oss-cn-hongkong.a
	//liyuncs.com/meta-icon/pmc/icon-USDM.png"}}
	log.Printf("assetDetails= %#v\n", assetDetails)

	lastBlock := wallet.getLastBlock()

	log.Printf("lastBlock= %#v\n", lastBlock)

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

	transactionsResp := wallet.getTransactions(reqTra)
	log.Printf("transactionsResp= %#v\n", transactionsResp)

	//generateSecret
	reqGenSecret := generateSecretReq.GenerateSecretParams{
		Lang: "en",
	}
	secretResp := wallet.generateSecret(reqGenSecret)
	log.Printf("secretResp= %#v\n", secretResp)

	//createAccount
	reqCreateAccount := createAccountReq.CreateAccountReq{
		Secret: "xxxxxxxxxxxxxxxxxxxxxxx",
	}
	createAccountResp := wallet.createAccount(reqCreateAccount)
	log.Printf("createAccountResp= %#v\n", createAccountResp)

	//broadcastCompleteTransaction
	reqBroadcastCompleteTransaction := broadcast.Params{
		"key":  123,
		"key1": []string{"item1", "item2"},
	}
	bCTResp := wallet.broadcastCompleteTransaction(reqBroadcastCompleteTransaction)
	log.Printf("bCTResp= %#v\n", bCTResp)

	defer sdk.Close()
}
