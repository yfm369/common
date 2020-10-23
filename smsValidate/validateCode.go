package smscode

import (
	"db"
	"fmt"
	"math/rand"
	pub "public_yfm"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/gomodule/redigo/redis"
)

const (
	CODE_TIMEOUT = 60
)

type ValidateInfo struct {
	Phone string
	Time  int64
	Code  string
}

var gValidate []*ValidateInfo

func GenValidateCode(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}

func SendSmsCode(phone string) bool {
	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", ALI_APP_ID, ALI_APP_SECRET)
	if err != nil {
		pub.PrintLog("SendSmsCode error :", err.Error())
		return false
	}

	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = phone
	request.SignName = ALI_SIGNNAME
	request.TemplateCode = ALI_TEMPLATE
	code := GenValidateCode(6)
	request.TemplateParam = fmt.Sprintf("{\"code\":\"%s\"}", code)

	redisc := db.RedisPool.Get()
	defer redisc.Close()

	_, errd := redisc.Do("SET", phone, code, "EX", "60")
	if errd != nil {
		pub.PrintLog("sendSms redis error", errd.Error())
		return false
	}

	pub.PrintLog("sendSmsv set code phone :", phone, " code:", code)

	go func() {
		response, err := client.SendSms(request)
		if err != nil {
			pub.PrintLog("sendSms send error", err.Error())
			return
		}

		pub.PrintLog(fmt.Sprintf("response is %#v", response))
	}()

	return true
}

func FindValidate(phone string) string {
	redisc := db.RedisPool.Get()
	defer redisc.Close()

	code, err := redis.String(redisc.Do("GET", phone))
	if err != nil {
		return ""
	}

	return code
}
