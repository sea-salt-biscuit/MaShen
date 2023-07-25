package main

import (
	"awesomeProject/XueLang/MaShen/config"
	"awesomeProject/XueLang/MaShen/net"
)

func main() {
	//config.A()
	//fmt.Println(config.File)
	//host := config.File.MustValue("login_server", "host", "127.0.0.1")
	//fmt.Println(host)
	host := config.File.MustValue("login_server", "host", "127.0.0.1")
	port := config.File.MustValue("login_server", "port", "8003")
	s := net.NewServer(host + ":" + port)
	s.Start()
}
