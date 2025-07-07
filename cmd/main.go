package main

import (
	"fmt"
	"go_novel/conf"
	"go_novel/router"
)

func main() {
	fmt.Println("Hello go_novel")
	conf.Init()

	r := router.NewRouter()
	_ = r.Run(conf.HttpPort)
}
