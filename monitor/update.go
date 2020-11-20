package monitor

import "fmt"

func UpdateConfig(oldPool, newPool string) (command string) {
	commandString := fmt.Sprintf("sed -i \"s/%s/%s/g\"  桌面/钱包配置.sh", oldPool, newPool)
	fmt.Println(commandString)
	return commandString
}
