package main

import (
	"encoding/json"
	"fmt"
	"imooc.com/resk/services"
)

func main() {
	data, _ := json.Marshal(&services.RedEnvelopeSendingDTO{})
	fmt.Println(string(data))
}
