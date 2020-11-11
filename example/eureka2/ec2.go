package main

import (
	"fmt"
	"github.com/tietang/go-eureka-client/eureka"
	"github.com/tietang/props/ini"
)

func main() {
	fmt.Println("main")
	conf := ini.NewIniFileConfigSource("ec.ini")
	client := eureka.NewClient(conf)
	client.Start()
	c := make(chan int, 1)
	<-c
}
