package main

import (
	"cloudc/controller"
)

func main() {
	server := &controller.Server{}
	server.Initialize("mysql", "root", "xiaohei", "3306", "127.0.0.1", "ccdb")
	server.Run("0.0.0.0:8009")
}
