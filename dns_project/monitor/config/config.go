package config

import (
	"fmt"
	"github.com/astaxie/beego/config"
)

func Initconfig(conftype, filename string) (conf config.Configer, err error) {
	conf, err = config.NewConfig(conftype, filename)
	if err != nil {
		fmt.Println("load conf failed,", err)
		return
	}
	return

}
