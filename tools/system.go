package tools

import (
	"fmt"
	"net"

	"github.com/spf13/viper"
)

func GetEmptyPort() (port int) {
	for i := 5899; i <= 6000; i++ {
		if _, err := net.Dial("tcp", fmt.Sprintf("%s:%d", viper.GetString("external_host"), i)); err != nil {
			return i
		}
	}
	return 0
}
