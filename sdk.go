package sdk

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/req/block"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/req/broadcastTra"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/req/pkgTranscaction"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/resp/blockResp"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/resp/broadcastResultResp"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/resp/createTransferAssetResp"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/resp/pkgTranscactionResp"
	"io"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/req/account"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/req/accountAsset"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/req/address"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/req/assetDetails"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/req/assets"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/req/broadcast"
	createAccountReq "github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/req/createAccount"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/req/createTransferAsset"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/req/generateSecretReq"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/req/transactions"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/resp/accountAssetResp"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/resp/accountResp"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/resp/assetDetailsResp"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/resp/assetsResp"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/resp/broadcastResp"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/resp/createAccountResp"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/resp/generateSecretResp"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/resp/lastBlockResp"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/resp/transactionsResp"
)

type Result struct {
	Code    int    // 0 fail 1 success
	Message string //
}

type NodeProcess struct {
	ChannelMap sync.Map
	Cmd        *exec.Cmd
	Stdin      io.WriteCloser
	Stdout     io.ReadCloser
	Stderr     io.ReadCloser
	//NodeExec    func()
}

func newNodeProcess(cmd string, args []string, debug bool) *NodeProcess {
	command := exec.Command(cmd, args...)

	stdin, err := command.StdinPipe()
	if err != nil {
		//panic(err)
	}

	stdout, err := command.StdoutPipe()
	if err != nil {
		//panic(err)
	}

	stderr, err := command.StderrPipe()
	if err != nil {
		//panic(err)
	}
	nodeProcess := &NodeProcess{
		ChannelMap: sync.Map{},
		Cmd:        command,
		Stdin:      stdin,
		Stdout:     stdout,
		Stderr:     stderr,
	}

	err = nodeProcess.Cmd.Start()
	if err != nil {
		fmt.Println("Failed to start process:", err)
		panic(err)
	}

	//if err := cmd.Start(); err != nil {
	//	fmt.Println("Error starting command:", err)
	//	return
	//}
	// 在一个goroutine中读取stdout
	var readLines = func(name string, reader io.Reader, resultCode int) {
		scanner := bufio.NewReader(reader)
		for {
			line, err := scanner.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				fmt.Println(name+" read line error:", err)
				continue
			}
			if debug {
				fmt.Printf("line: %v\n", line)
			}
			parts := strings.SplitN(line, "Result ", 2)
			if len(parts) == 2 {
				//fmt.Printf("node: %v\n", parts[1])
				// 解析行并提取req_id和json_result
				result := parts[1]
				idIndex := strings.Index(result, " ")
				if idIndex == -1 {
					fmt.Println(name+" Invalid line format:", line)
					continue
				}
				reqIDStr := result[0:idIndex]
				resultData := result[idIndex+1:]

				// 将req_id转换为整数
				reqID, err := strconv.Atoi(reqIDStr)
				if err != nil {
					fmt.Println(name+" Invalid req_id:", reqIDStr)
					continue
				}

				// 从map中获取并移除通道
				if ch, ok := nodeProcess.ChannelMap.LoadAndDelete(reqID); ok {
					channel := ch.(chan Result)
					// 向通道发送结果
					channel <- Result{Code: resultCode, Message: resultData}
				} else {
					fmt.Println(name+" Channel not found for req_id:", reqID)
				}
			}
		}
		fmt.Println("node-process end")
	}
	go readLines("stdout", nodeProcess.Stdout, 1)
	go readLines("stderr", nodeProcess.Stderr, 0)

	return nodeProcess
}
func (p *NodeProcess) CloseProcess() error {
	//p.Stdin.Close()
	//p.Cmd.Wait()
	//return p.Cmd.Cancel()
	return nil
}

type BCFWalletSDK struct {
	debug       bool
	nodeProcess *NodeProcess
}

func (sdk *BCFWalletSDK) Close() error {
	return sdk.nodeProcess.Cmd.Process.Kill()
}

func NewLocalBCFWalletSDK(debug bool) BCFWalletSDK {
	// 启动Node.js进程
	nodeProcess := newNodeProcess("node", []string{"--no-warnings", "./sdk.js"}, debug)
	return BCFWalletSDK{nodeProcess: nodeProcess}
}
func NewBCFWalletSDK() BCFWalletSDK {
	var try = 0
	var repl string
	for try < 3 {
		var err error
		repl, err = exec.LookPath("bfcwallet-node-go-repl")
		if errors.Is(err, exec.ErrDot) {
			err = nil
		}
		if err != nil {
			fmt.Println("installing bfcwallet-node-go-repl...")
			install := exec.Command("npm", "i", "-g", "bfcwallet-node-go-repl")
			err = install.Run()
			if err != nil {
				log.Fatal(err)
			}
			continue
		}
		break
	}

	// 启动Node.js进程
	nodeProcess := newNodeProcess(repl, []string{}, false)

	return BCFWalletSDK{nodeProcess: nodeProcess}
}

var reqIdAcc = 0

func nodeExec[T any](nodeProcess *NodeProcess, jsCode string) (T, error) {
	var res T
	reqIdAcc += 1
	req_id := reqIdAcc
	channel := make(chan Result)
	nodeProcess.ChannelMap.Store(req_id, channel)

	var evalCode = fmt.Sprintf("await returnToGo(%d, async()=>%v)\r\n\n", req_id, jsCode)
	//fmt.Println("evalCode", evalCode)
	_, err := nodeProcess.Stdin.Write([]byte(evalCode))
	if err != nil {
		return res, err
	}

	result := <-channel
	if result.Code == 1 {
		err := json.Unmarshal([]byte(result.Message), &res)
		return res, err
	} else {
		return res, errors.New(result.Message)
	}
}

type BCFWallet struct {
	nodeProcess *NodeProcess
	walletId    string
}

func (sdk *BCFWalletSDK) NewBCFWallet(ip string, port int, browserPath string) *BCFWallet {
	bfcWalletId, _ := nodeExec[int](sdk.nodeProcess, `{
		const bfcwalletMap = (globalThis.bfcwalletMap??=new Map());
		globalThis.bfcwalletIdAcc ??= 0
		const id = globalThis.bfcwalletIdAcc++
		const bfcwallet = walletBcf.BCFWalletFactory({
			enable: true,
            host: [{ ip: "`+ip+`", port: `+strconv.Itoa(port)+` }],
            browserPath: "`+browserPath+`",
        });
		bfcwalletMap.set(id, bfcwallet)
		return id;
	}`)
	return &BCFWallet{nodeProcess: sdk.nodeProcess, walletId: strconv.Itoa(bfcWalletId)}
}

func (sdk *BCFWalletSDK) NewBCFWalletSignUtil(ip string, port int, browserPath string) *BCFWallet {
	//TODO
	bfcWalletId, _ := nodeExec[int](sdk.nodeProcess, `{
		const bfcwalletMap = (globalThis.bfcwalletMap??=new Map());
		globalThis.bfcwalletIdAcc ??= 0
		const id = globalThis.bfcwalletIdAcc++
		const bfcwallet = bfmetaSignUtil.BCFWalletFactory({
			enable: true,
            host: [{ ip: "`+ip+`", port: `+strconv.Itoa(port)+` }],
            browserPath: "`+browserPath+`",
        });
		bfcwalletMap.set(id, bfcwallet)
		return id;
	}`)
	return &BCFWallet{nodeProcess: sdk.nodeProcess, walletId: strconv.Itoa(bfcWalletId)}
}

type BalanceResult struct {
	Success bool `json:"success"`
	Result  struct {
		Amount string `json:"amount"`
	} `json:"result"`
	//Result interface{} `json:"result"`
}

// func (wallet *BCFWallet) GetAddressBalance(address, magic, assetType string) (res BalanceResult) {
func (wallet *BCFWallet) GetAddressBalance(req address.Params) (res BalanceResult) {
	script := fmt.Sprintf(`globalThis.bfcwalletMap.get(%s).getAddressBalance(%q, %q, %q)`, wallet.walletId, req.Address, req.Magic, req.AssetType)
	res, _ = nodeExec[BalanceResult](wallet.nodeProcess, script)
	return res
}

/// baseApi

func (wallet *BCFWallet) GetTransactionsByBrowser(req transactions.GetTransactionsParams) (resp transactionsResp.GetTransactionsByBrowserResult, err error) {
	reqData, err := json.Marshal(req)
	if err != nil {
		return resp, fmt.Errorf("failed to marshal request: %v", err)
	}
	script := fmt.Sprintf(`globalThis.bfcwalletMap.get(%s).getTransactionsByBrowser(%v)`, wallet.walletId, string(reqData))
	resp, _ = nodeExec[transactionsResp.GetTransactionsByBrowserResult](wallet.nodeProcess, script)
	return resp, nil
}

func (wallet *BCFWallet) GetAccountInfo(req account.GetAccountInfoParams) (resp accountResp.GetAccountInfoRespResult) {
	resp, _ = nodeExec[accountResp.GetAccountInfoRespResult](wallet.nodeProcess, `
		globalThis.bfcwalletMap.get(`+wallet.walletId+`).getAccountInfo("`+req.Address+`")
	`)
	return resp
}

func (wallet *BCFWallet) GetAccountAsset(req accountAsset.GetAccountAssetParams) (resp accountAssetResp.GetAccountAssetRespResult) {
	resp, _ = nodeExec[accountAssetResp.GetAccountAssetRespResult](wallet.nodeProcess, `
		globalThis.bfcwalletMap.get(`+wallet.walletId+`).getAccountAsset("`+req.Address+`")
	`)
	return resp
}

// script
func (wallet *BCFWallet) GetAssets(req assets.PaginationOptions) (resp assetsResp.GetAssetsRespResult) {
	script := fmt.Sprintf(`
        globalThis.bfcwalletMap.get(%q).getAssets(%d, %d, %q)
    `, wallet.walletId, req.Page, req.PageSize, req.AssetType)
	resp, _ = nodeExec[assetsResp.GetAssetsRespResult](wallet.nodeProcess, script)
	return resp
}

func (wallet *BCFWallet) GetAssetDetails(req assetDetails.Req) (resp assetDetailsResp.GetAssetDetailsRespResult) {
	script := fmt.Sprintf(`globalThis.bfcwalletMap.get(%s).getAssetDetails(%q)`, wallet.walletId, req.AssetType)
	resp, _ = nodeExec[assetDetailsResp.GetAssetDetailsRespResult](wallet.nodeProcess, script)
	return
}
func (wallet *BCFWallet) GetAllAccountAsset(req accountAsset.GetAllAccountAssetReq) (resp accountAssetResp.GetAllAccountAssetRespResult) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return
	}
	script := fmt.Sprintf(`globalThis.bfcwalletMap.get(%s).getAllAccountAsset(%v)`, wallet.walletId, string(jsonData))
	resp, _ = nodeExec[accountAssetResp.GetAllAccountAssetRespResult](wallet.nodeProcess, script)
	return
}

// / baseApis2
func (wallet *BCFWallet) GetBlock(req block.GetBlockParams) (resp blockResp.GetBlockResultResp) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return
	}
	script := fmt.Sprintf(`globalThis.bfcwalletMap.get(%s).sdk.api.basic.getBlock(%v)`, wallet.walletId, string(jsonData))
	resp, _ = nodeExec[blockResp.GetBlockResultResp](wallet.nodeProcess, script)
	return
}

func (wallet *BCFWallet) GetLastBlock() (resp lastBlockResp.GetLastBlockResultResp) {
	script := fmt.Sprintf(`globalThis.bfcwalletMap.get(%s).sdk.api.basic.getLastBlock()`, wallet.walletId)
	resp, _ = nodeExec[lastBlockResp.GetLastBlockResultResp](wallet.nodeProcess, script)
	return
}

func (wallet *BCFWallet) GetTransactions(req transactions.GetTransactionsParams) (resp transactionsResp.GetTransactionsResult) {
	reqData, err := json.Marshal(req)
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return
	}
	script := fmt.Sprintf(`globalThis.bfcwalletMap.get(%s).sdk.api.basic.getTransactions(%v)`, wallet.walletId, string(reqData))
	resp, _ = nodeExec[transactionsResp.GetTransactionsResult](wallet.nodeProcess, script)
	return resp
}

func (wallet *BCFWallet) GenerateSecret(req generateSecretReq.GenerateSecretParams) (resp generateSecretResp.GenerateSecretRespResult) {
	//reqData, err := json.Marshal(req)
	//if err != nil {
	//	fmt.Println("Error marshalling to JSON:", err)
	//	return
	//}
	script := fmt.Sprintf(`globalThis.bfcwalletMap.get(%s).sdk.api.basic.generateSecret(%q)`, wallet.walletId, req.Lang)
	resp, _ = nodeExec[generateSecretResp.GenerateSecretRespResult](wallet.nodeProcess, script)
	return
}

func (wallet *BCFWallet) CreateAccount(req createAccountReq.CreateAccountReq) (resp createAccountResp.CreateAccountRespResult) {
	script := fmt.Sprintf(`globalThis.bfcwalletMap.get(%s).sdk.api.basic.createAccount(%q)`, wallet.walletId, req.Secret)
	resp, _ = nodeExec[createAccountResp.CreateAccountRespResult](wallet.nodeProcess, script)
	return
}

/// transactionApis

func (wallet *BCFWallet) BroadcastCompleteTransaction(req broadcast.Params) (resp broadcastResp.BroadcastRespResult[any]) {
	reqData, err := json.Marshal(req)
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return
	}
	// w = new require('@bfmeta/wallet-bcf').BCFWalletFactory({
	//        enable: true, host: [{ip: "34.84.178.63", port: 19503}], browserPath: "https://qapmapi.pmchainbox.com/browser",
	//    });
	//sdk.api.transaction.broadcastCompleteTransaction("{\"applyBlockHeight\":114208,\"asset\":{\"transferAsset\":{\"amount\":\"185184\",\"assetType\":\"PMC\",\"sourceChainMagic\":\"XXVXQ\",\"sourceChainName\":\"paymetachain\"}},\"effectiveBlockHeight\":114258,\"fee\":\"100000\",\"fromMagic\":\"\",\"range\":[],\"rangeType\":0,\"recipientId\":\"cFqv1tiifgYE6xbhZp43XxbZVJp363BWXt\",\"remark\":{\"orderId\":\"110b45fafcb84cb7a1de7eef5a957855\"},\"senderId\":\"c6C9ycTXrPBu8wXAGhUJHau678YyQwB2Mn\",\"senderPublicKey\":\"0d3c8003248cc4c71493dd67c0c433e75b7a191758df94fb0be5db2c6a94fecd\",\"signature\":\"2d0cea07ab73be6bdab258f12e7e0aa22776a8b9dd7b130f33fdd8fce6534cb0e29bc8d4983d3564178ae4189eedba80a864bda1a4ceb8b197e530ef1774ea07\",\"storageKey\":\"assetType\",\"storageValue\":\"PMC\",\"timestamp\":31839601,\"toMagic\":\"\",\"type\":\"PMC-PAYMETACHAIN-AST-02\",\"version\":1}"))

	script := fmt.Sprintf(`globalThis.bfcwalletMap.get(%s).sdk.api.transaction.broadcastCompleteTransaction(%q)`, wallet.walletId, string(reqData))
	//fmt.Println("broadcastCompleteTransaction sc", script)
	resp, _ = nodeExec[broadcastResp.BroadcastRespResult[any]](wallet.nodeProcess, script)
	return
}
func (wallet *BCFWallet) CreateTransferAsset(req createTransferAsset.TransferAssetTransactionParams) (result createTransferAssetResp.CreateResult, err error) {
	reqData, err := json.Marshal(req)
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return result, err
	}
	script := fmt.Sprintf(`globalThis.bfcwalletMap.get(%s).sdk.api.transaction
.createTransferAsset(JSON.parse(%q))`, wallet.walletId, string(reqData))
	result, err = nodeExec[createTransferAssetResp.CreateResult](wallet.nodeProcess, script)
	return result, err
}
func (wallet *BCFWallet) PackageTransferAsset(req pkgTranscaction.PackageTransacationParams) (resp pkgTranscactionResp.PackageResult, err error) {
	reqData, err := json.Marshal(req)
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return resp, err
	}
	script := fmt.Sprintf(`globalThis.bfcwalletMap.get(%s).sdk.api.transaction
.packageTransferAsset(%q)`, wallet.walletId, string(reqData))
	resp, err = nodeExec[pkgTranscactionResp.PackageResult](wallet.nodeProcess, script)
	return resp, err
}
func (wallet *BCFWallet) BroadcastTransferAsset(req broadcastTra.BroadcastTransactionParams) (resp broadcastResultResp.BroadcastResult, err error) {
	reqData, err := json.Marshal(req)
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return resp, err
	}
	script := fmt.Sprintf(`globalThis.bfcwalletMap.get(%s).sdk.api.transaction
.broadcastTransferAsset(JSON.parse(%q))`, wallet.walletId, string(reqData))
	resp, err = nodeExec[broadcastResultResp.BroadcastResult](wallet.nodeProcess, script)
	return resp, err
}

type BCFSignUtil struct {
	nodeProcess *NodeProcess
	signUtilId  string
}

func (sdk *BCFWalletSDK) NewBCFSignUtil(prefix string) *BCFSignUtil {
	signUtilId, _ := nodeExec[int](sdk.nodeProcess, `{
		const signUtilMap = (globalThis.signUtilMap??=new Map());
		globalThis.signUtilIdAcc ??= 0
		const id = globalThis.signUtilIdAcc++
		const signUtil = new __signUtil.BFMetaSignUtil("`+prefix+`",Buffer,cryptoHelper);
		signUtilMap.set(id, signUtil)
		return id;
	}`)
	return &BCFSignUtil{nodeProcess: sdk.nodeProcess, signUtilId: strconv.Itoa(signUtilId)}
}

type KeyPair struct {
	byteSecretKey []byte `json:"byteSecretKey"`
	SecretKey     string `json:"secretKey"`
	bytePublicKey []byte `json:"bytePublicKey"`
	PublicKey     string `json:"publicKey"`
}

type ResKeyPair struct {
	SecretKey string `json:"secretKey,omitempty"`
	PublicKey string `json:"publicKey,omitempty"`
}

func (util *BCFSignUtil) CreateKeypair(secret string) (res ResKeyPair, err error) {
	var keypair KeyPair
	script := fmt.Sprintf(`{
		const keypair = await globalThis.signUtilMap.get(%s).createKeypair(%q);
		return {
			secretKey:keypair.secretKey.toString("hex"),
			publicKey:keypair.publicKey.toString("hex"),
		}
	}`, util.signUtilId, secret)
	//SRC   {"secretKey":"a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3a4465fd76c16fcc458448076372abf1912cc5b150663a64dffefe550f96feadd","publicKey":"a4465fd76c16fcc458448076372abf1912cc5b150663a64dffefe550f96feadd"}
	keypair, err = nodeExec[KeyPair](util.nodeProcess, script)
	if err != nil {
		log.Fatal("CreateKeypair err :", err)
	}
	// 将数据编码为 Base64 字符串
	//res.SecretKey = base64.StdEncoding.EncodeToString(keypair.byteSecretKey)
	//res.PublicKey = base64.StdEncoding.EncodeToString(keypair.bytePublicKey)
	res.SecretKey = keypair.SecretKey
	res.PublicKey = keypair.PublicKey
	//a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3a4465fd76c16fcc458448076372abf1912cc5b150663a64dffefe550f96feadd
	//a4465fd76c16fcc458448076372abf1912cc5b150663a64dffefe550f96feadd
	return res, err
}

// todo
func (util *BCFSignUtil) GetBinaryAddressFromPublicKey(publicKey []byte) ([]byte, error) {
	script := fmt.Sprintf(`(
	await globalThis.signUtilMap.get(%s).getBinaryAddressFromPublicKey(Buffer.from(%q,"hex")))
	.toString("hex")
`, util.signUtilId, hex.EncodeToString(publicKey))
	binaryAddress, _ := nodeExec[string](util.nodeProcess, script)
	fmt.Printf("binaryAddress %#v\n", binaryAddress)
	if binaryAddress == "" {
		return nil, errors.New("publicKey is invalid")
	}
	return hex.DecodeString(binaryAddress)
}

func (util *BCFSignUtil) GetAddressFromPublicKey(publicKey []byte, prefix string) (string, error) {
	script := fmt.Sprintf(`(
		await globalThis.signUtilMap.get(%s)
		.getAddressFromPublicKey(Buffer.from(%q,"hex"),%q)
)
		.toString("hex")
`, util.signUtilId, hex.EncodeToString(publicKey), prefix)
	address, _ := nodeExec[string](util.nodeProcess, script)
	if address == "" {
		return "", errors.New("publicKey is invalid")
	}
	return address, nil
}

func (util *BCFSignUtil) GetAddressFromPublicKeyString(publicKey, prefix string) (string, error) {
	script := fmt.Sprintf(`(
		await globalThis.signUtilMap.get(%s)
		.getAddressFromPublicKeyString(%q,%q)
)
		.toString("hex")
`, util.signUtilId, publicKey, prefix)
	address, _ := nodeExec[string](util.nodeProcess, script)
	if address == "" {
		return "", errors.New("publicKey is invalid")
	}
	return address, nil
}

func (util *BCFSignUtil) GetAddressFromSecret(secret string) (string, error) {
	script := fmt.Sprintf(`(
		await globalThis.signUtilMap.get(%s)
		.getAddressFromSecret(%q)
)
		.toString("hex")
`, util.signUtilId, secret)
	address, _ := nodeExec[string](util.nodeProcess, script)
	if address == "" {
		return "", errors.New("secret is invalid")
	}
	return address, nil
}

func (util *BCFSignUtil) GetSecondPublicKeyStringFromSecretAndSecondSecret(secret, secondSecret, encode string) (string, error) {
	var script string
	if len(encode) > 0 {
		script = fmt.Sprintf(`(
		await globalThis.signUtilMap.get(%s)
		.getSecondPublicKeyStringFromSecretAndSecondSecret(%q,%q,%q)
)
		.toString("hex")
`, util.signUtilId, secret, secondSecret, encode)
	} else {
		script = fmt.Sprintf(`(
		await globalThis.signUtilMap.get(%s)
		.getSecondPublicKeyStringFromSecretAndSecondSecret(%q,%q)
)
		.toString("hex")
`, util.signUtilId, secret, secondSecret)
	}
	got, _ := nodeExec[string](util.nodeProcess, script)
	if got == "" {
		return "", errors.New("secret or secondSecret or encode is invalid")
	}
	return got, nil
}

func (util *BCFSignUtil) CreateSecondKeypair(secret, secondSecret string) (res ResKeyPair, err error) {
	var keypair KeyPair
	script := fmt.Sprintf(`{
		const keypair = await globalThis.signUtilMap.get(%s).createSecondKeypair(%q,%q)
		return {
			secretKey:keypair.secretKey.toString("hex"),
			publicKey:keypair.publicKey.toString("hex"),
		}
	}
`, util.signUtilId, secret, secondSecret)
	keypair, err = nodeExec[KeyPair](util.nodeProcess, script)
	if err != nil {
		log.Println("CreateSecondKeypair err : ", err)
	}
	res.SecretKey = keypair.SecretKey
	res.PublicKey = keypair.PublicKey
	//if address == "" {
	//	return res, errors.New("secret is invalid")
	//}
	return res, nil
}

type ResPubKeyPair struct {
	PublicKey string `json:"publicKey,omitempty"`
}

func (util *BCFSignUtil) GetSecondPublicKeyFromSecretAndSecondSecret(secret, secondSecret string) (res ResPubKeyPair, err error) {
	var keypair KeyPair
	script := fmt.Sprintf(`{
		const got = await globalThis.signUtilMap.get(%s).getSecondPublicKeyFromSecretAndSecondSecret(%q,%q)
		return {
			publicKey:got.toString("hex")
		}
	}
`, util.signUtilId, secret, secondSecret)
	keypair, err = nodeExec[KeyPair](util.nodeProcess, script)
	if err != nil {
		log.Println("GetSecondPublicKeyFromSecretAndSecondSecret err : ", err)
	}
	res.PublicKey = keypair.PublicKey
	return res, nil
}

///
//const signature = (await bfmetaSDK.bfchainSignUtil.detachedSign(bytes, keypair.secretKey)).toString("hex");

type ResSignToString struct {
	Type string `json:"type,omitempty"`
	Data []byte `json:"data,omitempty"`
}

func (util *BCFSignUtil) DetachedSign(msg, secretKey []byte) (res ResSignToString, err error) {
	script := fmt.Sprintf(`{
		const got = await globalThis.signUtilMap.get(%s).detachedSign(Buffer.from(%q,"hex"),Buffer.from(%q,"hex"));
		console.log("DetachedSign got to str",got.toString("hex"));
		return got;
}
`, util.signUtilId, msg, secretKey)
	res, _ = nodeExec[ResSignToString](util.nodeProcess, script)
	if len(res.Data) == 0 {
		return res, errors.New("msg or secretKey is invalid")
	}
	return res, nil
}

func (util *BCFSignUtil) SignToString(msg, secretKey []byte) (string, error) {
	script := fmt.Sprintf(`(
		await globalThis.signUtilMap.get(%s)
		.signToString(Buffer.from(%q,"hex"),Buffer.from(%q,"hex"))
)
		
`, util.signUtilId, msg, secretKey)
	got, _ := nodeExec[string](util.nodeProcess, script)
	if got == "" {
		return "", errors.New("msg or secretKey is invalid")
	}
	return got, nil
}
