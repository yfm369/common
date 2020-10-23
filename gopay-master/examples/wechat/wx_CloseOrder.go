package wechat

import (
	"fmt"

	"github.com/iGoogle-ink/gopay"
)

func CloseOrder() {
	//初始化微信客户端
	//    appId：应用ID
	//    MchID：商户ID
	//    ApiKey：Key值
	//    isProd：是否是正式环境
	client := gopay.NewWeChatClient("wxdaa2ab9ef87b5497", "1368139502", "GFDS8j98rewnmgl45wHTt980jg543abc", false)

	//初始化参数结构体
	body := make(gopay.BodyMap)
	body.Set("out_trade_no", "MfZC2segKxh0bnJSELbvKNeH3d9oWvvQ")
	body.Set("nonce_str", gopay.GetRandomString(32))
	body.Set("sign_type", gopay.SignType_MD5)

	//请求关闭订单，成功后得到结果
	wxRsp, err := client.CloseOrder(body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("wxRsp：", *wxRsp)
}
