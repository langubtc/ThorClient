package thormonitor

import (
	"ThorGui/config"
	"ThorGui/monitor"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hsw409328/ip-range-lib"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

type CommandResult struct {
	Ip     string
	Target string
	Cmdout string
	Status int32
	Optype string
}

type Foo struct {
	Index      int
	Worker     string
	MinerIp    string
	Wallet     string
	ServerIp   string
	ServerStr  string
	ServerPort string
	ScanTime   string
}

func MinerIPFunc(ipString string) ([]string, error) {
	t := ip_range_lib.NewIpRangeLib()

	ipRange := ipString
	result, err := t.IpRangeToIpList(ipRange)
	if err != nil {
		return result, err
	}
	return result, nil

}

// 执行命令主要函数方法
func remoteExec(user, ip, password, opType string, port int, command string) (result CommandResult) {
	//fmt.Println(command)
	var commandStr string

	// 拼接目标地址 ip:port
	sshHost := ip + ":" + strconv.Itoa(port)

	client, err := ssh.Dial("tcp", sshHost, &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: 3 * time.Second,
	})

	if err != nil {
		commandStr = "SSH Connect Error"
		return CommandResult{
			Ip:     ip,
			Target: ip,
			Cmdout: commandStr,
			Status: 1,
			Optype: opType,
		}
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:  0,
		ssh.IGNCR: 1,
	}

	// 建立新会话
	session, err := client.NewSession()

	if err := session.RequestPty("vt220", 80, 40, modes); err != nil {
		log.Fatal("error:", err)
	}

	if err != nil {
		commandStr = "SSH session error"
		return CommandResult{
			Ip:     ip,
			Target: ip,
			Cmdout: commandStr,
			Status: 1,
			Optype: opType,
		}
	}

	// 关闭会话
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b

	if err := session.Run(command); err != nil {

		commandStr = "Failed to run: " + err.Error()
	}

	//fmt.Println(b.String())

	commandStr = b.String()

	// 关闭连接
	defer client.Close()

	return CommandResult{
		Ip:     ip,
		Target: ip,
		Cmdout: commandStr,
		Status: 0,
		Optype: opType,
	}

}

func RunMonitor(opType string, ipList []string, oldstring, newstring string) (total, active, inactive int, result string) {

	//var opType = "update"
	var cmdCommand string
	var ServerIp []string
	var countResult = 0
	var ActiveResult = 0
	var InactiveResult = 0
	var statsResult string

	conf := config.LoadConfig() //加载配置文件

	switch opType {

	case "scan":
		cmdCommand = "grep 'miner' 桌面/钱包配置.sh"
		ServerIp = ipList
	case "reboot":
		cmdCommand = monitor.RebootCommand(conf.Password)
		ServerIp = ipList
	case "stats":
		cmdCommand = monitor.ResponseCommand(conf.DefaultLogRow)
		ServerIp = ipList
	case "update":
		cmdCommand = monitor.UpdateConfig(oldstring, newstring)
		ServerIp = ipList
	default:
		cmdCommand = "grep 'miner' 桌面/钱包配置.sh"
		ServerIp, _ = MinerIPFunc(conf.IpRange)
		ServerIp = ipList
	}

	outchan := make(chan CommandResult)
	var wg_command sync.WaitGroup
	var wg_processing sync.WaitGroup

	var saveMiner []Foo

	// 开启执行线程
	for _, t := range ServerIp {
		wg_command.Add(1)

		target := t + " (" + conf.User + "@" + t + ")"
		go func(dst, user, ip, command string, out chan CommandResult) {
			defer wg_command.Done()
			result := remoteExec(conf.User, ip, conf.Password, opType, conf.Port, cmdCommand)
			out <- CommandResult{
				ip,
				dst,
				result.Cmdout,
				result.Status,
				result.Optype,
			}
		}(target, conf.User, t, cmdCommand, outchan)
	}

	// 开启读取线程
	wg_processing.Add(1)
	go func() {
		defer wg_processing.Done()
		for o := range outchan {

			scanTime := time.Now().Format("2006-01-02 15:04:05")

			if o.Status == 0 {
				ActiveResult += 1 //记录成功数
				if opType == "stats" {
					statsResult = statsResult + o.Cmdout
				}
				if opType == "scan" {
					returnDecodeResult := monitor.DecodeMinerInfo(o.Cmdout)
					//fmt.Println("解析结果", o.Target,
					//	returnDecodeResult.Stratum, returnDecodeResult.Miner, returnDecodeResult.Worker, returnDecodeResult.PoolIP, returnDecodeResult.PoolPort)
					saveMiner = append(saveMiner, Foo{countResult,
						returnDecodeResult.Worker,
						o.Ip,
						returnDecodeResult.Miner,
						returnDecodeResult.PoolIP,
						returnDecodeResult.PoolStr,
						returnDecodeResult.PoolPort,
						scanTime,
					})
				}
			} else {
				InactiveResult += 1 //记录失败
				if opType == "stats" {
					statsResult = statsResult + o.Cmdout
				}
				if opType == "scan" {
					saveMiner = append(saveMiner, Foo{countResult, "", o.Ip, o.Cmdout, "", "", "", scanTime})
				}
			}
			countResult += 1 //记录扫描总数
		}
	}()

	// wait untill all goroutines to finish and close the channel
	wg_command.Wait()
	close(outchan)

	wg_processing.Wait()

	fmt.Printf("miners count:%d", countResult)
	//保存结果成json文件
	if opType == "scan" {
		file, _ := json.MarshalIndent(saveMiner, "", "")
		_ = ioutil.WriteFile("./config/result.json", file, 0644)
	}
	return countResult, ActiveResult, InactiveResult, statsResult
}
