package monitor

import "fmt"

func RebootCommand(password string) (Command string) {
	command := fmt.Sprintf("echo \"%s\" | sudo -S reboot", password)
	return command
}
