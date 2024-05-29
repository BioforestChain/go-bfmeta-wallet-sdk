package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
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
	"io"
	"os/exec"
	"strconv"
	"strings"
	"sync"
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

func newNodeProcess(cmd string, args ...string) *NodeProcess {
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
			//fmt.Printf("line: %v\n", line)
			parts := strings.SplitN(line, "Result ", 2)
			if len(parts) == 2 {
				//fmt.Printf("node: %v\n", parts[1])
				// 解析行并提取req_id和json_result
				parts := strings.SplitN(parts[1], " ", 3)
				if len(parts) < 2 {
					fmt.Println(name+" Invalid line format:", line)
					continue
				}
				reqIDStr := parts[0]
				resultData := parts[1]

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
	nodeProcess *NodeProcess
}

func (sdk *BCFWalletSDK) Close() {
	sdk.nodeProcess.CloseProcess()
}

func newBCFWalletSDK() BCFWalletSDK {
	//init 结构体
	// 启动Node.js进程
	nodeProcess := newNodeProcess("node", "--no-warnings", "./sdk.js")
	return BCFWalletSDK{nodeProcess: nodeProcess}
}

var reqIdAcc = 0

// func NodeExec[T any](jsCode string) (T, error) {
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

type BCFWallet struct {
	nodeProcess *NodeProcess
	walletId    string
}

func (sdk BCFWalletSDK) newBCFWallet(ip string, port int, browserPath string) *BCFWallet {
	//TODO
	bfcWalletId, _ := nodeExec[int](sdk.nodeProcess, `{
		const bfcwalletMap = (globalThis.bfcwalletMap??=new Map());
		globalThis.bfcwalletIdAcc ??= 0
		const id = globalThis.bfcwalletIdAcc++
		const bfcwallet = new require('@bfmeta/wallet-bcf').BCFWalletFactory({
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

// func (wallet *BCFWallet) getAddressBalance(address, magic, assetType string) (res BalanceResult) {
func (wallet *BCFWallet) getAddressBalance(req address.Params) (res BalanceResult) {
	script := fmt.Sprintf(`globalThis.bfcwalletMap.get(%s).getAddressBalance(%q, %q, %q)`, wallet.walletId, req.Address, req.Magic, req.AssetType)
	res, _ = nodeExec[BalanceResult](wallet.nodeProcess, script)
	return res
}

/// baseApi

func (wallet *BCFWallet) getTransactionsByBrowser(req transactions.GetTransactionsParams) (resp transactionsResp.GetTransactionsByBrowserResult, err error) {
	reqData, err := json.Marshal(req)
	if err != nil {
		return resp, fmt.Errorf("failed to marshal request: %v", err)
	}
	script := fmt.Sprintf(`globalThis.bfcwalletMap.get(%s).getTransactionsByBrowser(%v)`, wallet.walletId, string(reqData))
	resp, _ = nodeExec[transactionsResp.GetTransactionsByBrowserResult](wallet.nodeProcess, script)
	return resp, nil
}

func (wallet *BCFWallet) getAccountInfo(req account.GetAccountInfoParams) (resp accountResp.GetAccountInfoRespResult) {
	resp, _ = nodeExec[accountResp.GetAccountInfoRespResult](wallet.nodeProcess, `
		globalThis.bfcwalletMap.get(`+wallet.walletId+`).getAccountInfo("`+req.Address+`")
	`)
	return resp
}

func (wallet *BCFWallet) getAccountAsset(req accountAsset.GetAccountAssetParams) (resp accountAssetResp.GetAccountAssetRespResult) {
	resp, _ = nodeExec[accountAssetResp.GetAccountAssetRespResult](wallet.nodeProcess, `
		globalThis.bfcwalletMap.get(`+wallet.walletId+`).getAccountAsset("`+req.Address+`")
	`)
	return resp
}

func (wallet *BCFWallet) getAssets(req assets.PaginationOptions) (resp assetsResp.GetAssetsRespResult) {
	script := fmt.Sprintf(`
        globalThis.bfcwalletMap.get(%q).getAssets(%d, %d, %q)
    `, wallet.walletId, req.Page, req.PageSize, req.AssetType)
	resp, _ = nodeExec[assetsResp.GetAssetsRespResult](wallet.nodeProcess, script)
	return resp
}

func (wallet *BCFWallet) getAssetDetails(req assetDetails.Req) (resp assetDetailsResp.GetAssetDetailsRespResult) {
	script := fmt.Sprintf(`globalThis.bfcwalletMap.get(%s).getAssetDetails(%q)`, wallet.walletId, req.AssetType)
	resp, _ = nodeExec[assetDetailsResp.GetAssetDetailsRespResult](wallet.nodeProcess, script)
	return
}
func (wallet *BCFWallet) getAllAccountAsset(req accountAsset.GetAllAccountAssetReq) (resp accountAssetResp.GetAllAccountAssetRespResult) {
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
//func getBlock(req block.GetBlockParams) (resp blockResp.GetAllAccountAssetResp) {
//	return
//}

func (wallet *BCFWallet) getLastBlock() (resp lastBlockResp.GetLastBlockInfoRespResult) {
	script := fmt.Sprintf(`globalThis.bfcwalletMap.get(%s).getLastBlock()`, wallet.walletId)
	resp, _ = nodeExec[lastBlockResp.GetLastBlockInfoRespResult](wallet.nodeProcess, script)
	return
}

func (wallet *BCFWallet) getTransactions(req transactions.GetTransactionsParams) (resp transactionsResp.GetTransactionsResult) {
	reqData, err := json.Marshal(req)
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return
	}
	script := fmt.Sprintf(`globalThis.bfcwalletMap.get(%s).getTransactions(%v)`, wallet.walletId, string(reqData))
	resp, _ = nodeExec[transactionsResp.GetTransactionsResult](wallet.nodeProcess, script)
	return resp
}

func (wallet *BCFWallet) generateSecret(req generateSecretReq.GenerateSecretParams) (resp generateSecretResp.GenerateSecretRespResult) {
	//reqData, err := json.Marshal(req)
	//if err != nil {
	//	fmt.Println("Error marshalling to JSON:", err)
	//	return
	//}
	script := fmt.Sprintf(`globalThis.bfcwalletMap.get(%s).generateSecret(%q)`, wallet.walletId, req.Lang)
	resp, _ = nodeExec[generateSecretResp.GenerateSecretRespResult](wallet.nodeProcess, script)
	return
}

func (wallet *BCFWallet) createAccount(req createAccountReq.CreateAccountReq) (resp createAccountResp.CreateAccountRespResult) {
	script := fmt.Sprintf(`globalThis.bfcwalletMap.get(%s).createAccount(%q)`, wallet.walletId, req.Secret)
	resp, _ = nodeExec[createAccountResp.CreateAccountRespResult](wallet.nodeProcess, script)
	return
}

/// transactionApis

func (wallet *BCFWallet) broadcastCompleteTransaction(req broadcast.Params) (resp broadcastResp.BroadcastRespResult[any]) {
	reqData, err := json.Marshal(req)
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return
	}
	script := fmt.Sprintf(`globalThis.bfcwalletMap.get(%s).broadcastCompleteTransaction(%q)`, wallet.walletId, string(reqData))
	resp, _ = nodeExec[broadcastResp.BroadcastRespResult[any]](wallet.nodeProcess, script)
	return
}
func (wallet *BCFWallet) createTransferAsset(req createTransferAsset.TransferAssetTransactionParams) (resp interface{}) {
	//reqData, err := json.Marshal(req)
	//if err != nil {
	//	fmt.Println("Error marshalling to JSON:", err)
	//	return
	//}
	//script := fmt.Sprintf(`globalThis.bfcwalletMap.get(%s).createTransferAsset(%q)`, wallet.walletId, string(reqData))
	//resp, _ = nodeExec[createTransferAssetResp.CreateResult](wallet.nodeProcess, script)
	//if resp.Success {
	//	return resp.SuccessCreateResult
	//}
	return
}
func (wallet *BCFWallet) packageTransferAsset(req createAccountReq.CreateAccountReq) (resp createAccountResp.CreateAccountResp) {
	return
}
func (wallet *BCFWallet) broadcastTransferAsset(req createAccountReq.CreateAccountReq) (resp createAccountResp.CreateAccountResp) {
	return
}
