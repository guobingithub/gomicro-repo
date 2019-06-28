package main

import (
	"fmt"
	"gomicro-repo/api/logger"
	"gomicro-repo/api/server"
	"gomicro-repo/api/zkmgr"
)

func main() {
	logger.Info(fmt.Sprintf("gomicro-repo api start ..."))

	zkmgr.StartZkService()

	server.StartHttpServer()
}
