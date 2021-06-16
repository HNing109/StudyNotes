package main

import "tcp_pkg"

func main() {
	var server = new(tcp_pkg.TcpServer)
	server.Start()
}
