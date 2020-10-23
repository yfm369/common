package smscode

import (
	"fmt"
	pub "public_yfm"
	"runtime"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

const (
	NOTICE_PHONE   = "18851196486"
	ALI_APP_ID     = "LTAI4Fo6xomXzjxAU1U8RWMo"
	ALI_APP_SECRET = "1GPSzTSmm8dJAAXKBhWhTOzMetf1J2"
	ALI_SIGNNAME   = "阿拉丁寻宝"
	ALI_TEMPLATE   = "SMS_186966122"
	DUMP_SIGNNAME  = "云牧网耕"
	DUMP_TEMPLATE  = "SMS_189713630"
	add_test 	   = "123"
)

func SendBatchSms(phones string) {
	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", ALI_APP_ID, ALI_APP_SECRET)
	if err != nil {
		pub.PrintLog("SendSmsCode error :", err.Error())
		return
	}

	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = phones
	request.SignName = ALI_SIGNNAME
	request.TemplateCode = "SMS_181862284"

	response, err := client.SendSms(request)
	if err != nil {
		pub.PrintLog("sendSms send error", err.Error())
		return
	}

	pub.PrintLog(fmt.Sprintf("response is %#v", response))
}

func SendDumpSms(phone, svrname, msg string) bool {
	if runtime.GOOS == `windows` {
		return false
	}

	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", ALI_APP_ID, ALI_APP_SECRET)
	if err != nil {
		pub.PrintLog("SendSmsCode error :", err.Error())
		return false
	}

	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = phone
	request.SignName = DUMP_SIGNNAME
	request.TemplateCode = DUMP_TEMPLATE
	request.TemplateParam = fmt.Sprintf("{\"name\":\"%s\",\"msg\":\"%s\"}", svrname, msg)

	response, err := client.SendSms(request)
	if err != nil {
		pub.PrintLog("sendSms send error", err.Error())
		return false
	}

	pub.PrintLog(fmt.Sprintf("response is %#v", response))

	return true
}
