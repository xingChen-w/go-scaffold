// HTTP 服务入口

package main

import (
	"fmt"
	"go-scaffold/internal/constant"
	"go-scaffold/internal/http"
	"go-scaffold/pkg/conf"
	"go-scaffold/pkg/log"
	"go-scaffold/pkg/xgorm"
	"go-scaffold/pkg/xredis"
	"os"
)

func main() {
	fmt.Printf("processID: %d\n", os.Getppid())
	conf.Init()                                                                         // 加载配置文件
	log.RawConfig("app.log", conf.GetHandler()).Init()                                  // 加载日志
	defer log.Sync()                                                                    // 日志落盘
	xgorm.RawConfig("app.mysql.db1", conf.GetHandler()).Append(constant.DBStoreNameDB1) // 加载gorm
	xgorm.RawConfig("app.mysql.db2", conf.GetHandler()).Append(constant.DBStoreNameDB2) // 加载gorm
	defer xgorm.CloseAll()
	xredis.RawConfig("app.redis", conf.GetHandler()).Init() // 加载redis
	defer xredis.Cli.Close()
	http.Serve() // 启动服务
}
