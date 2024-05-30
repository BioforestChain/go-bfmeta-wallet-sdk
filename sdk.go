package sdk

import (
	"bufio"
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
				//if ch, ok := channelMap.LoadAndDelete(reqID); ok {
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
	p.Stdin.Close()
	return p.Cmd.Wait()
}

type BCFWalletSDK struct {
	debug       bool
	nodeProcess *NodeProcess
}

func (sdk *BCFWalletSDK) Close() {
	sdk.nodeProcess.CloseProcess()
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

//func nodeExecCommonResult[T any](nodeProcess *NodeProcess, jsCode string) (T, error) {
//	var res T
//	reqIdAcc += 1
//	req_id := reqIdAcc
//	channel := make(chan Result)
//	nodeProcess.ChannelMap.Store(req_id, channel)
//
//	var evalCode = fmt.Sprintf("await returnToGo(%d, async()=>%v)\r\n\n", req_id, jsCode)
//	_, err := nodeProcess.Stdin.Write([]byte(evalCode))
//	if err != nil {
//		return res, err
//	}
//
//	result := <-channel
//	if result.Code == 1 {
//		var commonResult resp.CommonResult
//		msgBytes := []byte(result.Message)
//		err := json.Unmarshal(msgBytes, &commonResult)
//		if err != nil {
//			return res, err
//		}
//		if commonResult.Success {
//			err := json.Unmarshal(msgBytes, &res)
//		} else {
//
//		}
//		return res, err
//	} else {
//		return res, errors.New(result.Message)
//	}
//}

type BCFWallet struct {
	nodeProcess *NodeProcess
	walletId    string
}

func (sdk BCFWalletSDK) NewBCFWallet(ip string, port int, browserPath string) *BCFWallet {
	//TODO
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
// todo

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

// todo resp 同上
//func (wallet *BCFWallet) GetLastBlock() (resp lastBlockResp.GetBlockResult) {
//	script := fmt.Sprintf(`globalThis.bfcwalletMap.get(%s).getLastBlock()`, wallet.walletId)
//	resp, _ = nodeExec[lastBlockResp.GetBlockResult](wallet.nodeProcess, script)
//	return
//}

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
	script := fmt.Sprintf(`globalThis.bfcwalletMap.get(%s).getTransactions(%v)`, wallet.walletId, string(reqData))
	resp, _ = nodeExec[transactionsResp.GetTransactionsResult](wallet.nodeProcess, script)
	return resp
}

func (wallet *BCFWallet) GenerateSecret(req generateSecretReq.GenerateSecretParams) (resp generateSecretResp.GenerateSecretRespResult) {
	//reqData, err := json.Marshal(req)
	//if err != nil {
	//	fmt.Println("Error marshalling to JSON:", err)
	//	return
	//}
	script := fmt.Sprintf(`globalThis.bfcwalletMap.get(%s).generateSecret(%q)`, wallet.walletId, req.Lang)
	resp, _ = nodeExec[generateSecretResp.GenerateSecretRespResult](wallet.nodeProcess, script)
	return
}

func (wallet *BCFWallet) CreateAccount(req createAccountReq.CreateAccountReq) (resp createAccountResp.CreateAccountRespResult) {
	script := fmt.Sprintf(`globalThis.bfcwalletMap.get(%s).createAccount(%q)`, wallet.walletId, req.Secret)
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
	script := fmt.Sprintf(`globalThis.bfcwalletMap.get(%s).broadcastCompleteTransaction(%q)`, wallet.walletId, string(reqData))
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
.createTransferAsset(%q)`, wallet.walletId, string(reqData))
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
.broadcastTransferAsset(%q)`, wallet.walletId, string(reqData))
	resp, err = nodeExec[broadcastResultResp.BroadcastResult](wallet.nodeProcess, script)
	return resp, err
}
