package monitor

import "strconv"

var logfile = "/var/log/thor.log"

// 重新组装命令，可以根据用户选择的行数进行展示
func ResponseCommand(filerow int) (parsCommand string) {
	cmdTail := "tail -n " + strconv.Itoa(filerow) + " " + logfile
	return cmdTail
}
