package main

import "tcp_pkg"

func main() {
	var client = new(tcp_pkg.TcpClient)
	client.Start()
}