package sdk

import (
	"bufio"
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/jbase"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/req/account"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/req/accountAsset"
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
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/resp/accountAssetResp"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/resp/accountResp"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/resp/assetDetailsResp"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/resp/assetsResp"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/resp/asymmetricDecryptResp"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/resp/asymmetricEncryptResp"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/resp/blockResp"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/resp/broadcastResp"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/resp/broadcastResultResp"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/resp/createAccountResp"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/resp/createTransferAssetResp"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/resp/generateSecretResp"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/resp/lastBlockResp"
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/resp/pkgTranscactionResp"
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
	OnClose    OnNodeProcessClose
	IsClosed   bool
	execLock   sync.Mutex // 互斥锁
}

type OnNodeProcessClose func()

func newNodeProcess(cmd string, args []string, debug bool, onclose OnNodeProcessClose) *NodeProcess {
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
		log.Fatalf("Failed to start process: %v", err)
		panic(err)
	}
	// 在一个goroutine中读取stdout
	var readLines = func(name string, reader io.Reader, resultCode int) {
		scanner := bufio.NewReader(reader)
		for {
			line, err := scanner.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Println(name+" read line error:", err)
				if strings.Contains(err.Error(), "file already closed") {
					break
				}
				continue
			}
			if debug {
				log.Printf("line: %v\n", line)
			}
			parts := strings.SplitN(line, "Result ", 2)
			if len(parts) == 2 {
				// fmt.Printf("node: %v\n", parts[1])
				// 解析行并提取req_id和json_result
				result := parts[1]
				idIndex := strings.Index(result, " ")
				if idIndex == -1 {
					log.Println(name+" Invalid line format:", line)
					continue
				}
				reqIDStr := result[0:idIndex]
				resultData := result[idIndex+1:]
				// 将req_id转换为整数
				reqID, err := strconv.Atoi(reqIDStr)
				if err != nil {
					log.Println(name+" Invalid req_id:", reqIDStr)
					continue
				}
				successIndex := strings.Index(resultData, " ")
				if successIndex == -1 {
					log.Println(name+" Invalid line format:", line)
					continue
				}
				successStr := resultData[0:successIndex]
				messageStr := resultData[successIndex+1:]
				var Code int
				if successStr == "true" {
					Code = 1
				} else {
					Code = 0
				}
				// 从map中获取并移除通道
				if ch, ok := nodeProcess.ChannelMap.LoadAndDelete(reqID); ok {
					channel := ch.(chan Result)
					// 向通道发送结果
					channel <- Result{Code: Code, Message: messageStr}
				} else {
					log.Println(name+" Channel not found for req_id:", reqID)
				}
			}
		}
		log.Printf("node-process stdio %d end", resultCode)
		if resultCode == 1 {
			nodeProcess.IsClosed = true
			nodeProcess.OnClose()
		}
	}
	go readLines("stdout", nodeProcess.Stdout, 1)
	go readLines("stderr", nodeProcess.Stderr, 0)
	return nodeProcess
}

type BCFWalletSDK struct {
	debug       bool
	nodeProcess *NodeProcess
}

func (sdk *BCFWalletSDK) SetDebug(debug bool) {
	sdk.debug = debug
}
func (sdk *BCFWalletSDK) SetOnClose(onclose OnNodeProcessClose) {
	sdk.nodeProcess.OnClose = onclose
}
func (sdk *BCFWalletSDK) GetIsClosed() bool {
	return sdk.nodeProcess.IsClosed
}
func (sdk *BCFWalletSDK) Close() error {
	err := sdk.nodeProcess.Cmd.Process.Kill()
	sdk.nodeProcess.Stdout.Close()
	sdk.nodeProcess.Stderr.Close()
	sdk.nodeProcess.Stdin.Close()
	return err
}
func NewLocalBCFWalletSDK(debug bool) BCFWalletSDK {
	// 启动Node.js进程
	// onclose ON= func() {}
	nodeProcess := newNodeProcess("node", []string{"--no-warnings", "./sdk.js"}, debug, func() {})
	return BCFWalletSDK{nodeProcess: nodeProcess}
}
func NewBCFWalletSDK() BCFWalletSDK {
	var try = 0
	var repl string

	packageJson := readPackageJson()
	name := packageJson["name"].(string)
	version := packageJson["version"].(string)
	log.Printf("require %s@%s", name, version)

	for try < 3 {
		var err error
		repl, err = exec.LookPath(name)
		if errors.Is(err, exec.ErrDot) {
			err = nil
		}
		if err == nil {
			cmd := exec.Command(name, "--version")
			var out bytes.Buffer
			cmd.Stdout = &out
			err = cmd.Run()
			if err == nil {
				if !strings.Contains(out.String(), (name + " " + version)) {
					err = fmt.Errorf("locale %s need upgrade to v%s", out.String(), version)
				}
			}
		}
		if err != nil {
			fmt.Printf("installing %s@%s...", name, version)
			install := exec.Command("npm", "i", "-g", name+"@"+version)
			err = install.Run()
			if err != nil {
				log.Fatal(err)
			}
			continue
		}
		break
	}
	nodeProcess := newNodeProcess(repl, []string{}, false, func() {})
	return BCFWalletSDK{nodeProcess: nodeProcess}
}

//go:embed package.json
var packageJsonData []byte

func readPackageJson() (result map[string]interface{}) {
	// 映射JSON数据
	err := json.Unmarshal(packageJsonData, &result)
	if err != nil {
		log.Fatal("fail to parse package.json")
		panic(err)
	}
	return
}

var reqIdAcc = 0

func nodeExec[T any](nodeProcess *NodeProcess, jsCode string) (T, error) {
	nodeProcess.execLock.Lock()
	defer nodeProcess.execLock.Unlock()
	var res T
	reqIdAcc += 1
	req_id := reqIdAcc
	channel := make(chan Result)
	nodeProcess.ChannelMap.Store(req_id, channel)
	var evalCode = fmt.Sprintf("await returnToGo(%d, async()=>%v)\r\n\n", req_id, jsCode)
	// fmt.Println("evalCode", evalCode)
	_, err := nodeProcess.Stdin.Write([]byte(evalCode))
	if err != nil {
		return res, err
	}
	result := <-channel
	// fmt.Println("evalResult", result)
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

// / baseApi
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
	script := fmt.Sprintf(`globalThis.bfcwalletMap.get(%s).sdk.api.basic.generateSecret(%q)`, wallet.walletId, req.Lang)
	resp, _ = nodeExec[generateSecretResp.GenerateSecretRespResult](wallet.nodeProcess, script)
	return
}
func (wallet *BCFWallet) CreateAccount(req createAccountReq.CreateAccountReq) (resp createAccountResp.CreateAccountRespResult) {
	script := fmt.Sprintf(`globalThis.bfcwalletMap.get(%s).sdk.api.basic.createAccount(%q)`, wallet.walletId, req.Secret)
	resp, _ = nodeExec[createAccountResp.CreateAccountRespResult](wallet.nodeProcess, script)
	return
}

// / transactionApis
func (wallet *BCFWallet) BroadcastCompleteTransaction(req broadcast.Params) (resp broadcastResp.BroadcastRespResult[any], err error) {
	reqData, err := json.Marshal(req)
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return
	}
	// w = new require('@bfmeta/wallet-bcf').BCFWalletFactory({
	//        enable: true, host: [{ip: "34.84.178.63", port: 19503}], browserPath: "https://qapmapi.pmchainbox.com/browser",
	//    });
	//sdk.api.transaction.broadcastCompleteTransaction("{\"applyBlockHeight\":114208,\"asset\":{\"transferAsset\":{\"amount\":\"185184\",\"assetType\":\"PMC\",\"sourceChainMagic\":\"XXVXQ\",\"sourceChainName\":\"paymetachain\"}},\"effectiveBlockHeight\":114258,\"fee\":\"100000\",\"fromMagic\":\"\",\"range\":[],\"rangeType\":0,\"recipientId\":\"cFqv1tiifgYE6xbhZp43XxbZVJp363BWXt\",\"remark\":{\"orderId\":\"110b45fafcb84cb7a1de7eef5a957855\"},\"senderId\":\"c6C9ycTXrPBu8wXAGhUJHau678YyQwB2Mn\",\"senderPublicKey\":\"0d3c8003248cc4c71493dd67c0c433e75b7a191758df94fb0be5db2c6a94fecd\",\"signature\":\"2d0cea07ab73be6bdab258f12e7e0aa22776a8b9dd7b130f33fdd8fce6534cb0e29bc8d4983d3564178ae4189eedba80a864bda1a4ceb8b197e530ef1774ea07\",\"storageKey\":\"assetType\",\"storageValue\":\"PMC\",\"timestamp\":31839601,\"toMagic\":\"\",\"type\":\"PMC-PAYMETACHAIN-AST-02\",\"version\":1}"))
	script := fmt.Sprintf(`globalThis.bfcwalletMap.get(%s).sdk.api.transaction.broadcastCompleteTransaction(JSON.parse(%q))`, wallet.walletId, string(reqData))
	resp, err = nodeExec[broadcastResp.BroadcastRespResult[any]](wallet.nodeProcess, script)
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
	Prefix      string
	signUtilId  string
}

func (util *BCFSignUtil) AsJsSignUtil() string {
	return fmt.Sprintf("globalThis.signUtilMap.get(%s)", util.signUtilId)
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
	return &BCFSignUtil{nodeProcess: sdk.nodeProcess, Prefix: prefix, signUtilId: strconv.Itoa(signUtilId)}
}

type ResKeyPair struct {
	SecretKey jbase.HexStringBuffer `json:"secretKey,omitempty"`
	PublicKey jbase.HexStringBuffer `json:"publicKey,omitempty"`
}

// ts
// Se 03ac674216f3e15c761ee1a5e255f067953623c8b388b4459e13f978d7c846f4caf0f4c00cf9240771975e42b6672c88a832f98f01825dda6e001e2aab0bc0cc
// Pu caf0f4c00cf9240771975e42b6672c88a832f98f01825dda6e001e2aab0bc0cc
// go
// Se 03ac674216f3e15c761ee1a5e255f067953623c8b388b4459e13f978d7c846f4caf0f4c00cf9240771975e42b6672c88a832f98f01825dda6e001e2aab0bc0cc
// Pu caf0f4c00cf9240771975e42b6672c88a832f98f01825dda6e001e2aab0bc0cc
func (util *BCFSignUtil) CreateKeypair(secret string) (keypair ResKeyPair, err error) {
	script := fmt.Sprintf(`{
		const keypair = await globalThis.signUtilMap.get(%s).createKeypair(%q);
		return {
			SecretKey:keypair.secretKey.toString("hex"),
			PublicKey:keypair.publicKey.toString("hex"),
		}
	}`, util.signUtilId, secret)

	keypair, err = nodeExec[ResKeyPair](util.nodeProcess, script)
	if err != nil {
		log.Fatal("CreateKeypair err :", err)
	}
	return
}
func (util *BCFSignUtil) CreateKeypairBySecretKey(secret jbase.StringBuffer) (keypair ResKeyPair, err error) {
	script := fmt.Sprintf(`{
		const keypair = await globalThis.signUtilMap.get(%s).createKeypairBySecretKey(%s);
		return {
			secretKey:keypair.secretKey.toString("hex"),
			publicKey:keypair.publicKey.toString("hex"),
		}
	}`, util.signUtilId, secret.AsJsBuffer())
	//SRC   {"secretKey":"a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3a4465fd76c16fcc458448076372abf1912cc5b150663a64dffefe550f96feadd","publicKey":"a4465fd76c16fcc458448076372abf1912cc5b150663a64dffefe550f96feadd"}
	keypair, err = nodeExec[ResKeyPair](util.nodeProcess, script)
	if err != nil {
		log.Fatal("CreateKeypair err :", err)
	}
	return keypair, err
}
func (util *BCFSignUtil) GetPublicKeyFromSecret(secret string) (res string, err error) {
	script := fmt.Sprintf(`(
	await globalThis.signUtilMap.get(%s).getPublicKeyFromSecret(%q)
)`, util.signUtilId, secret)
	res, err = nodeExec[string](util.nodeProcess, script)
	if err != nil {
		log.Fatal("GetPublicKeyFromSecret err :", err)
	}
	return res, err
}

func (util *BCFSignUtil) GetAddressFromPublicKey(publicKey jbase.StringBuffer, prefix string) (string, error) {
	script := fmt.Sprintf(`(
		await globalThis.signUtilMap.get(%s)
		.getAddressFromPublicKey(%s,%q))
		`, util.signUtilId, publicKey.AsJsBuffer(), prefix)
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
		`, util.signUtilId, secret)
	address, _ := nodeExec[string](util.nodeProcess, script)
	if address == "" {
		return "", errors.New("secret is invalid")
	}
	return address, nil
}
func (util *BCFSignUtil) GetSecondPublicKeyStringFromSecretAndSecondSecret(secret, secondSecret string) (publicKey jbase.HexStringBuffer, err error) {
	script := fmt.Sprintf(`(
	await globalThis.signUtilMap.get(%s)
	.getSecondPublicKeyStringFromSecretAndSecondSecret(%q,%q)
	)
	`, util.signUtilId, secret, secondSecret)
	publicKey, err = nodeExec[jbase.HexStringBuffer](util.nodeProcess, script)
	if publicKey.Value == "" {
		err = errors.New("secret or secondSecret or encode is invalid")
	}
	return
}

// 根据私钥获取公钥String
func (util *BCFSignUtil) GetSecondPublicKeyFromSecretAndSecondSecretV2(secret, secondSecret string) (publicKey jbase.HexStringBuffer, err error) {
	script := fmt.Sprintf(`(
	await globalThis.signUtilMap.get(%s)
	.getSecondPublicKeyStringFromSecretAndSecondSecretV2(%q,%q)
	)
	`, util.signUtilId, secret, secondSecret)
	publicKey, err = nodeExec[jbase.HexStringBuffer](util.nodeProcess, script)
	if publicKey.Value == "" {
		err = errors.New("secret or secondSecret or encode is invalid")
	}
	return
}
func (util *BCFSignUtil) CreateSecondKeypair(secret, secondSecret string) (keypair ResKeyPair, err error) {
	script := fmt.Sprintf(`{
		const keypair = await globalThis.signUtilMap.get(%s).createSecondKeypair(%q,%q)
		return {
			secretKey:keypair.secretKey.toString("hex"),
			publicKey:keypair.publicKey.toString("hex"),
		}
	}
	`, util.signUtilId, secret, secondSecret)
	keypair, err = nodeExec[ResKeyPair](util.nodeProcess, script)
	if err != nil {
		log.Println("CreateSecondKeypair err : ", err)
	}
	return keypair, nil
}

type ResPubKeyPair struct {
	PublicKey jbase.HexStringBuffer `json:"publicKey,omitempty"`
}

func (util *BCFSignUtil) GetSecondPublicKeyFromSecretAndSecondSecret(secret, secondSecret string) (keypair ResPubKeyPair, err error) {
	script := fmt.Sprintf(`{
		const got = await globalThis.signUtilMap.get(%s).getSecondPublicKeyFromSecretAndSecondSecret(%q,%q)
		return {
			publicKey:got.toString("hex")
		}
	}
	`, util.signUtilId, secret, secondSecret)
	keypair, err = nodeExec[ResPubKeyPair](util.nodeProcess, script)
	if err != nil {
		log.Println("GetSecondPublicKeyFromSecretAndSecondSecret err : ", err)
	}
	return keypair, nil
}

// /
// const signature = (await bfmetaSDK.bfchainSignUtil.detachedSign(bytes, keypair.secretKey)).toString("hex");
type ResSignToString struct {
	Type string `json:"type,omitempty"`
	Data []byte `json:"data,omitempty"`
}

func (util *BCFSignUtil) DetachedSign(msg, secretKey jbase.StringBuffer) (signature jbase.HexStringBuffer, err error) {
	script := fmt.Sprintf(`{
		const got = await globalThis.signUtilMap.get(%s).detachedSign(%s,%s);
		return got.toString("hex");
	}
	`, util.signUtilId, msg.AsJsBuffer(), secretKey.AsJsBuffer())
	signature, err = nodeExec[jbase.HexStringBuffer](util.nodeProcess, script)
	return
}

// /**
//   - 验证签名
//     *
//   - @param message
//   - @param signatureBuffer
//   - @param publicKeyBuffer
//   - @returns
//     */
//
// detachedVeriy(message: Uint8Array, signatureBuffer: Uint8Array, publicKeyBuffer: Uint8Array): Promise<boolean>;
func (util *BCFSignUtil) DetachedVerify(message, signature, publicKey jbase.StringBuffer) (verified bool, err error) {
	script := fmt.Sprintf(`(
		await globalThis.signUtilMap.get(%s)
		.detachedVeriy(%s,%s,%s)
	)
	`, util.signUtilId, message.AsJsBuffer(), signature.AsJsBuffer(), publicKey.AsJsBuffer())
	// log.Printf("script=%s", script)
	verified, err = nodeExec[bool](util.nodeProcess, script)
	return
}

/**
 * 非对称加密
 *
 * @param msg
 * @param decryptPK
 * @param encryptSK
 */
//asymmetricEncrypt(msg: Uint8Array, decryptPK: Uint8Array, encryptSK: Uint8Array): {
//encryptedMessage: Uint8Array;
//nonce: Uint8Array;
//};
func (util *BCFSignUtil) AsymmetricEncrypt(msg, decryptPK, encryptSK []byte) (res asymmetricEncryptResp.ResAsymmetricEncrypt, err error) {
	script := fmt.Sprintf(`{
		const got = await globalThis.signUtilMap.get(%s).asymmetricEncrypt((Buffer.from(%q,"hex")),(Buffer.from(%q,"hex")),(Buffer.from(%q,"hex")));
		return {
			encryptedMessage:got.encryptedMessage.toString("hex"),
			nonce:got.nonce.toString("hex"),
		}
	}
`, util.signUtilId, msg, decryptPK, encryptSK)
	res, err = nodeExec[asymmetricEncryptResp.ResAsymmetricEncrypt](util.nodeProcess, script)
	if err != nil {
		return res, err
	}
	return res, nil
}

// 非对称解密
//
//	asymmetricDecrypt(encryptedMessage: Uint8Array, encryptPK: Uint8Array, decryptSK: Uint8Array, nonce?: Uint8Array): false | Uint8Array;
func (util *BCFSignUtil) AsymmetricDecrypt(req asymmetricDecrypt.Req) (res asymmetricDecryptResp.ResAsymmetricDecrypt, err error) {
	script := fmt.Sprintf(`{
		const got = await globalThis.signUtilMap.get(%s).asymmetricDecrypt(
(Buffer.from(%q,"hex")),(Buffer.from(%q,"hex")),(Buffer.from(%q,"hex")),(Buffer.from(%q,"hex"))
);
		return {
			encryptedMessage:got.encryptedMessage.toString("hex"),
			nonce:got.nonce.toString("hex"),
		}
	}
`, util.signUtilId, req.EncryptedMessage, req.EncryptPK, req.DecryptSK, req.Nonce)
	res, err = nodeExec[asymmetricDecryptResp.ResAsymmetricDecrypt](util.nodeProcess, script)
	if err != nil {
		return res, err
	}
	return res, nil
}

// checkSecondSecret(secret: string, secondSecret: string, secondPublicKey: string): Promise<boolean>;
func (util *BCFSignUtil) CheckSecondSecret(secret, secondSecret, secondPublicKey string) (res bool, err error) {
	script := fmt.Sprintf(`(
		await globalThis.signUtilMap.get(%s)
		.checkSecondSecret(%q,%q,%q)
)
`, util.signUtilId, secret, secondSecret, secondPublicKey)
	res, err = nodeExec[bool](util.nodeProcess, script)
	if res == false {
		return res, errors.New("secret or secondSecret or secondPublicKey is invalid")
	}
	return res, nil
}

// checkSecondSecret(secret: string, secondSecret: string, secondPublicKey: string): Promise<boolean>;
func (util *BCFSignUtil) CheckSecondSecretV2(secret, secondSecret, secondPublicKey string) (res bool, err error) {
	script := fmt.Sprintf(`(
		await globalThis.signUtilMap.get(%s)
		.checkSecondSecretV2(%q,%q,%q)
)
`, util.signUtilId, secret, secondSecret, secondPublicKey)
	res, err = nodeExec[bool](util.nodeProcess, script)
	if res == false {
		return res, errors.New("secret or secondSecret or secondPublicKey is invalid")
	}
	return res, nil
}

// 根据安全密码的公私钥对
// createSecondKeypairV2
func (util *BCFSignUtil) CreateSecondKeypairV2(secret, secondSecret string) (keypair ResKeyPair, err error) {
	script := fmt.Sprintf(`{
		const keypair = await globalThis.signUtilMap.get(%s).createSecondKeypairV2(%q,%q);
		return {
			secretKey:keypair.secretKey.toString("hex"),
			publicKey:keypair.publicKey.toString("hex"),
		}
	}`, util.signUtilId, secret, secondSecret)
	//SRC   {"secretKey":"a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3a4465fd76c16fcc458448076372abf1912cc5b150663a64dffefe550f96feadd","publicKey":"a4465fd76c16fcc458448076372abf1912cc5b150663a64dffefe550f96feadd"}
	keypair, err = nodeExec[ResKeyPair](util.nodeProcess, script)
	if err != nil {
		log.Fatal("CreateSecondKeypairV2 err :", err)
	}
	return keypair, err
}
