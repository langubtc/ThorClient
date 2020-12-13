package monitor

import "fmt"

func RebootCommand(password string) (Command string) {
	command := fmt.Sprintf("echo \"%s\" | sudo -S reboot", password)
	return command
}

func PowerOffCommand(password string) (Command string) {
	command := fmt.Sprintf("echo \"%s\" | sudo -S poweroff", password)
	return command
}
