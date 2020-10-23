/*!
 * <邮件变量>
 *
 * Copyright (c) 2020.4 by <yfm/ ALADING Co.>
 */

package mailcfg

import (
	"common/config"
	"fmt"
)

const (
	MAILID_CASHOUTFAIL = 1 //提现失败返现

	MAIL_TYPE_PLAER = 1 //玩家邮件
	MAIL_TYPE_SYS   = 2 //系统邮件
)

func GetTitle(id int32, param ...interface{}) (title string) {
	pCfg := config.GetMailCfg(id)
	if pCfg == nil {
		return
	}

	title = fmt.Sprintf(pCfg.Title, param...)

	return
}

func GetContext(id int32, param ...interface{}) (context string) {
	pCfg := config.GetMailCfg(id)
	if pCfg == nil {
		return
	}

	context = fmt.Sprintf(pCfg.Context, param...)

	return
}
