package monitor

import (
	"strings"
)

type ResultDecode struct {
	Stratum  string
	Miner    string
	Worker   string
	PoolIP   string
	PoolPort string
	PoolStr  string
}

// 扫描矿机启动程序进行分段拆分

func DecodeMinerInfo(resultPool string) ResultDecode {
	//fmt.Println(resultPool)
	// 以@作为分隔符切出矿池接入地址和端口
	countSplit := strings.Split(resultPool, "@")
	poolResult := strings.Split(countSplit[1], ":")

	//以//作为分割符号切分出矿工地址和矿工号
	minerResult := strings.Split(countSplit[0], "//")
	workerResult := strings.Split(minerResult[1], ":")
	worker := strings.Split(workerResult[0], ".")

	return ResultDecode{
		Stratum:  "stratum2+tcp://" + minerResult[1],
		Miner:    worker[0],
		Worker:   worker[1],
		PoolStr:  strings.Replace(countSplit[1], "\n", "", -1),
		PoolIP:   poolResult[0],
		PoolPort: poolResult[1],
	}
}
