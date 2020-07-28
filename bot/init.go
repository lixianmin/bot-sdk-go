package bot

import (
	"github.com/lixianmin/logo"
	"github.com/lixianmin/bot-sdk-go/bot/logger"
)

/********************************************************************
created:    2020-07-20
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type InitArgs struct {
	Logger logo.ILogger // 自定义日志对象，默认只输出到控制台
}

func Init(args InitArgs) {
	logger.Init(args.Logger)
}
