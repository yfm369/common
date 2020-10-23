package gopay

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/parnurzeal/gorequest"
)

type AliPayClient struct {
	AppId            string
	PrivateKey       string
	AlipayRootCertSN string
	AppCertSN        string
	ReturnUrl        string
	NotifyUrl        string
	Charset          string
	SignType         string
	AppAuthToken     string
	AuthToken        string
	IsProd           bool
}

//初始化支付宝客户端
//    注意：如果使用支付宝公钥证书验签，请设置 支付宝根证书SN（client.SetAlipayRootCertSN()）、应用公钥证书SN（client.SetAppCertSN()）
//    appId：应用ID
//    PrivateKey：应用私钥
//    IsProd：是否是正式环境
func NewAliPayClient(appId, privateKey string, isProd bool) (client *AliPayClient) {
	return &AliPayClient{
		AppId:      appId,
		PrivateKey: privateKey,
		IsProd:     isProd,
	}
}

//alipay.trade.fastpay.refund.query(统一收单交易退款查询)
//    文档地址：https://docs.open.alipay.com/api_1/alipay.trade.fastpay.refund.query
func (a *AliPayClient) AliPayTradeFastPayRefundQuery(body BodyMap) (aliRsp *AliPayTradeFastpayRefundQueryResponse, err error) {
	var (
		p1, p2 string
		bs     []byte
	)
	p1 = body.Get("out_trade_no")
	p2 = body.Get("trade_no")
	if p1 == null && p2 == null {
		return nil, errors.New("out_trade_no and trade_no are not allowed to be null at the same time")
	}
	if bs, err = a.doAliPay(body, "alipay.trade.fastpay.refund.query"); err != nil {
		return
	}
	aliRsp = new(AliPayTradeFastpayRefundQueryResponse)
	if err = json.Unmarshal(bs, aliRsp); err != nil {
		return nil, err
	}
	if aliRsp.AliPayTradeFastpayRefundQueryResponse.Code != "10000" {
		info := aliRsp.AliPayTradeFastpayRefundQueryResponse
		return nil, fmt.Errorf(`{"code":"%v","msg":"%v","sub_code":"%v","sub_msg":"%v"}`, info.Code, info.Msg, info.SubCode, info.SubMsg)
	}
	aliRsp.SignData = getSignData(bs)
	return
}

//alipay.trade.order.settle(统一收单交易结算接口)
//    文档地址：https://docs.open.alipay.com/api_1/alipay.trade.order.settle
func (a *AliPayClient) AliPayTradeOrderSettle(body BodyMap) (aliRsp *AliPayTradeOrderSettleResponse, err error) {
	var (
		p1, p2 string
		bs     []byte
	)
	p1 = body.Get("out_request_no")
	p2 = body.Get("trade_no")
	if p1 == null || p2 == null {
		return nil, errors.New("out_request_no or trade_no are not allowed to be null")
	}
	if bs, err = a.doAliPay(body, "alipay.trade.order.settle"); err != nil {
		return
	}
	aliRsp = new(AliPayTradeOrderSettleResponse)
	if err = json.Unmarshal(bs, aliRsp); err != nil {
		return nil, err
	}
	if aliRsp.AliPayTradeOrderSettleResponse.Code != "10000" {
		info := aliRsp.AliPayTradeOrderSettleResponse
		return nil, fmt.Errorf(`{"code":"%v","msg":"%v","sub_code":"%v","sub_msg":"%v"}`, info.Code, info.Msg, info.SubCode, info.SubMsg)
	}
	aliRsp.SignData = getSignData(bs)
	return
}

//alipay.trade.create(统一收单交易创建接口)
//    文档地址：https://docs.open.alipay.com/api_1/alipay.trade.create
func (a *AliPayClient) AliPayTradeCreate(body BodyMap) (aliRsp *AliPayTradeCreateResponse, err error) {
	var (
		p1, p2 string
		bs     []byte
	)
	p1 = body.Get("out_trade_no")
	p2 = body.Get("buyer_id")
	if p1 == null && p2 == null {
		return nil, errors.New("out_trade_no and buyer_id are not allowed to be null at the same time")
	}
	if bs, err = a.doAliPay(body, "alipay.trade.create"); err != nil {
		return
	}
	aliRsp = new(AliPayTradeCreateResponse)
	if err = json.Unmarshal(bs, aliRsp); err != nil {
		return nil, err
	}
	if aliRsp.AliPayTradeCreateResponse.Code != "10000" {
		info := aliRsp.AliPayTradeCreateResponse
		return nil, fmt.Errorf(`{"code":"%v","msg":"%v","sub_code":"%v","sub_msg":"%v"}`, info.Code, info.Msg, info.SubCode, info.SubMsg)
	}
	aliRsp.SignData = getSignData(bs)
	return
}

//alipay.trade.close(统一收单交易关闭接口)
//    文档地址：https://docs.open.alipay.com/api_1/alipay.trade.close
func (a *AliPayClient) AliPayTradeClose(body BodyMap) (aliRsp *AliPayTradeCloseResponse, err error) {
	var (
		p1, p2 string
		bs     []byte
	)
	p1 = body.Get("out_trade_no")
	p2 = body.Get("trade_no")
	if p1 == null && p2 == null {
		return nil, errors.New("out_trade_no and trade_no are not allowed to be null at the same time")
	}
	if bs, err = a.doAliPay(body, "alipay.trade.close"); err != nil {
		return
	}
	aliRsp = new(AliPayTradeCloseResponse)
	if err = json.Unmarshal(bs, aliRsp); err != nil {
		return nil, err
	}
	if aliRsp.AliPayTradeCloseResponse.Code != "10000" {
		info := aliRsp.AliPayTradeCloseResponse
		return nil, fmt.Errorf(`{"code":"%v","msg":"%v","sub_code":"%v","sub_msg":"%v"}`, info.Code, info.Msg, info.SubCode, info.SubMsg)
	}
	aliRsp.SignData = getSignData(bs)
	return
}

//alipay.trade.cancel(统一收单交易撤销接口)
//    文档地址：https://docs.open.alipay.com/api_1/alipay.trade.cancel
func (a *AliPayClient) AliPayTradeCancel(body BodyMap) (aliRsp *AliPayTradeCancelResponse, err error) {
	var (
		p1, p2 string
		bs     []byte
	)
	p1 = body.Get("out_trade_no")
	p2 = body.Get("trade_no")
	if p1 == null && p2 == null {
		return nil, errors.New("out_trade_no and trade_no are not allowed to be null at the same time")
	}
	if bs, err = a.doAliPay(body, "alipay.trade.cancel"); err != nil {
		return
	}
	aliRsp = new(AliPayTradeCancelResponse)
	if err = json.Unmarshal(bs, aliRsp); err != nil {
		return nil, err
	}
	if aliRsp.AliPayTradeCancelResponse.Code != "10000" {
		info := aliRsp.AliPayTradeCancelResponse
		return nil, fmt.Errorf(`{"code":"%v","msg":"%v","sub_code":"%v","sub_msg":"%v"}`, info.Code, info.Msg, info.SubCode, info.SubMsg)
	}
	aliRsp.SignData = getSignData(bs)
	return
}

//alipay.trade.refund(统一收单交易退款接口)
//    文档地址：https://docs.open.alipay.com/api_1/alipay.trade.refund
func (a *AliPayClient) AliPayTradeRefund(body BodyMap) (aliRsp *AliPayTradeRefundResponse, err error) {
	var (
		p1, p2 string
		bs     []byte
	)
	p1 = body.Get("out_trade_no")
	p2 = body.Get("trade_no")
	if p1 == null && p2 == null {
		return nil, errors.New("out_trade_no and trade_no are not allowed to be null at the same time")
	}
	if bs, err = a.doAliPay(body, "alipay.trade.refund"); err != nil {
		return nil, err
	}
	aliRsp = new(AliPayTradeRefundResponse)
	if err = json.Unmarshal(bs, aliRsp); err != nil {
		return nil, err
	}
	if aliRsp.AlipayTradeRefundResponse.Code != "10000" {
		info := aliRsp.AlipayTradeRefundResponse
		return nil, fmt.Errorf(`{"code":"%v","msg":"%v","sub_code":"%v","sub_msg":"%v"}`, info.Code, info.Msg, info.SubCode, info.SubMsg)
	}
	aliRsp.SignData = getSignData(bs)
	return
}

//alipay.trade.refund(统一收单退款页面接口)
//    文档地址：https://docs.open.alipay.com/api_1/alipay.trade.page.refund
func (a *AliPayClient) AliPayTradePageRefund(body BodyMap) (aliRsp *AliPayTradePageRefundResponse, err error) {
	var (
		p1, p2 string
		bs     []byte
	)
	p1 = body.Get("out_trade_no")
	p2 = body.Get("trade_no")
	if p1 == null && p2 == null {
		return nil, errors.New("out_trade_no and trade_no are not allowed to be null at the same time")
	}
	if bs, err = a.doAliPay(body, "	alipay.trade.page.refund"); err != nil {
		return
	}
	aliRsp = new(AliPayTradePageRefundResponse)
	if err = json.Unmarshal(bs, aliRsp); err != nil {
		return nil, err
	}
	if aliRsp.AliPayTradePageRefundResponse.Code != "10000" {
		info := aliRsp.AliPayTradePageRefundResponse
		return nil, fmt.Errorf(`{"code":"%v","msg":"%v","sub_code":"%v","sub_msg":"%v"}`, info.Code, info.Msg, info.SubCode, info.SubMsg)
	}
	aliRsp.SignData = getSignData(bs)
	return
}

//alipay.trade.precreate(统一收单线下交易预创建)
//    文档地址：https://docs.open.alipay.com/api_1/alipay.trade.precreate
func (a *AliPayClient) AliPayTradePrecreate(body BodyMap) (aliRsp *AlipayTradePrecreateResponse, err error) {
	var bs []byte
	p := body.Get("out_trade_no")
	if p == null {
		return nil, errors.New("out_trade_no is not allowed to be null")
	}
	if bs, err = a.doAliPay(body, "alipay.trade.precreate"); err != nil {
		return
	}
	aliRsp = new(AlipayTradePrecreateResponse)
	if err = json.Unmarshal(bs, aliRsp); err != nil {
		return nil, err
	}
	if aliRsp.AlipayTradePrecreateResponse.Code != "10000" {
		info := aliRsp.AlipayTradePrecreateResponse
		return nil, fmt.Errorf(`{"code":"%v","msg":"%v","sub_code":"%v","sub_msg":"%v"}`, info.Code, info.Msg, info.SubCode, info.SubMsg)
	}
	aliRsp.SignData = getSignData(bs)
	return
}

//alipay.trade.pay(统一收单交易支付接口)
//    文档地址：https://docs.open.alipay.com/api_1/alipay.trade.pay
func (a *AliPayClient) AliPayTradePay(body BodyMap) (aliRsp *AliPayTradePayResponse, err error) {
	var bs []byte
	p := body.Get("out_trade_no")
	if p == null {
		return nil, errors.New("out_trade_no is not allowed to be null")
	}
	if bs, err = a.doAliPay(body, "alipay.trade.pay"); err != nil {
		return
	}
	aliRsp = new(AliPayTradePayResponse)
	if err = json.Unmarshal(bs, aliRsp); err != nil {
		return nil, err
	}
	if aliRsp.AliPayTradePayResponse.Code != "10000" {
		info := aliRsp.AliPayTradePayResponse
		return nil, fmt.Errorf(`{"code":"%v","msg":"%v","sub_code":"%v","sub_msg":"%v"}`, info.Code, info.Msg, info.SubCode, info.SubMsg)
	}
	aliRsp.SignData = getSignData(bs)
	return
}

//alipay.trade.query(统一收单线下交易查询)
//    文档地址：https://docs.open.alipay.com/api_1/alipay.trade.query
func (a *AliPayClient) AliPayTradeQuery(body BodyMap) (aliRsp *AliPayTradeQueryResponse, err error) {
	var (
		p1, p2 string
		bs     []byte
	)
	p1 = body.Get("out_trade_no")
	p2 = body.Get("trade_no")
	if p1 == null && p2 == null {
		return nil, errors.New("out_trade_no and trade_no are not allowed to be null at the same time")
	}
	if bs, err = a.doAliPay(body, "alipay.trade.query"); err != nil {
		return
	}
	aliRsp = new(AliPayTradeQueryResponse)
	if err = json.Unmarshal(bs, aliRsp); err != nil {
		return nil, err
	}
	if aliRsp.AliPayTradeQueryResponse.Code != "10000" {
		info := aliRsp.AliPayTradeQueryResponse
		return nil, fmt.Errorf(`{"code":"%v","msg":"%v","sub_code":"%v","sub_msg":"%v"}`, info.Code, info.Msg, info.SubCode, info.SubMsg)
	}
	aliRsp.SignData = getSignData(bs)
	return
}

//alipay.trade.app.pay(app支付接口2.0)
//    文档地址：https://docs.open.alipay.com/api_1/alipay.trade.app.pay
func (a *AliPayClient) AliPayTradeAppPay(body BodyMap) (payParam string, err error) {
	var bs []byte
	trade := body.Get("out_trade_no")
	if trade == null {
		return null, errors.New("out_trade_no is not allowed to be null")
	}
	if bs, err = a.doAliPay(body, "alipay.trade.app.pay"); err != nil {
		return null, err
	}
	payParam = string(bs)
	return
}

//alipay.trade.wap.pay(手机网站支付接口2.0)
//    文档地址：https://docs.open.alipay.com/api_1/alipay.trade.wap.pay
func (a *AliPayClient) AliPayTradeWapPay(body BodyMap) (payUrl string, err error) {
	var bs []byte
	trade := body.Get("out_trade_no")
	if trade == null {
		return null, errors.New("out_trade_no is not allowed to be null")
	}
	body.Set("product_code", "QUICK_WAP_WAY")
	if bs, err = a.doAliPay(body, "alipay.trade.wap.pay"); err != nil {
		return null, err
	}
	payUrl = string(bs)
	return
}

//alipay.trade.page.pay(统一收单下单并支付页面接口)
//    文档地址：https://docs.open.alipay.com/api_1/alipay.trade.page.pay
func (a *AliPayClient) AliPayTradePagePay(body BodyMap) (payUrl string, err error) {
	var bs []byte
	trade := body.Get("out_trade_no")
	if trade == null {
		return null, errors.New("out_trade_no is not allowed to be null")
	}
	body.Set("product_code", "FAST_INSTANT_TRADE_PAY")
	if bs, err = a.doAliPay(body, "alipay.trade.page.pay"); err != nil {
		return null, err
	}
	payUrl = string(bs)
	return
}

//alipay.fund.trans.toaccount.transfer(单笔转账到支付宝账户接口)
//    文档地址：https://docs.open.alipay.com/api_28/alipay.fund.trans.toaccount.transfer
func (a *AliPayClient) AlipayFundTransToaccountTransfer(body BodyMap) (aliRsp *AlipayFundTransToaccountTransferResponse, err error) {
	var bs []byte
	trade1 := body.Get("out_biz_no")
	if trade1 == null {
		return nil, errors.New("out_biz_no is not allowed to be null")
	}
	if bs, err = a.doAliPay(body, "alipay.fund.trans.toaccount.transfer"); err != nil {
		return
	}
	aliRsp = new(AlipayFundTransToaccountTransferResponse)
	if err = json.Unmarshal(bs, aliRsp); err != nil {
		return nil, err
	}
	if aliRsp.AlipayFundTransToaccountTransferResponse.Code != "10000" {
		info := aliRsp.AlipayFundTransToaccountTransferResponse
		return nil, fmt.Errorf(`{"code":"%v","msg":"%v","sub_code":"%v","sub_msg":"%v"}`, info.Code, info.Msg, info.SubCode, info.SubMsg)
	}
	aliRsp.SignData = getSignData(bs)
	return
}

//alipay.trade.orderinfo.sync(支付宝订单信息同步接口)
//    文档地址：https://docs.open.alipay.com/api_1/alipay.trade.orderinfo.sync
func (a *AliPayClient) AliPayTradeOrderinfoSync(body BodyMap) {

}

//alipay.system.oauth.token(换取授权访问令牌)
//    文档地址：https://docs.open.alipay.com/api_9/alipay.system.oauth.token
func (a *AliPayClient) AliPaySystemOauthToken(body BodyMap) (aliRsp *AliPaySystemOauthTokenResponse, err error) {
	var bs []byte
	grantType := body.Get("grant_type")
	if grantType == null {
		return nil, errors.New("grant_type is not allowed to be null")
	}
	code := body.Get("code")
	refreshToken := body.Get("refresh_token")
	if code == null && refreshToken == null {
		return nil, errors.New("code and refresh_token are not allowed to be null at the same time")
	}
	if bs, err = aliPaySystemOauthToken(a.AppId, a.PrivateKey, body, "alipay.system.oauth.token", a.IsProd); err != nil {
		return
	}
	aliRsp = new(AliPaySystemOauthTokenResponse)
	if err = json.Unmarshal(bs, aliRsp); err != nil {
		return nil, err
	}
	if aliRsp.AliPaySystemOauthTokenResponse.AccessToken == null {
		info := aliRsp.ErrorResponse
		return nil, fmt.Errorf(`{"code":"%v","msg":"%v","sub_code":"%v","sub_msg":"%v"}`, info.Code, info.Msg, info.SubCode, info.SubMsg)
	}
	aliRsp.SignData = getSignData(bs)
	return
}

//alipay.user.info.share(支付宝会员授权信息查询接口)
//    body：此接口无需body参数
//    文档地址：https://docs.open.alipay.com/api_2/alipay.user.info.share
func (a *AliPayClient) AlipayUserInfoShare() (aliRsp *AlipayUserInfoShareResponse, err error) {
	var bs []byte
	if bs, err = a.doAliPay(nil, "alipay.user.info.share"); err != nil {
		return nil, err
	}
	aliRsp = new(AlipayUserInfoShareResponse)
	if err = json.Unmarshal(bs, aliRsp); err != nil {
		return nil, err
	}
	if aliRsp.AlipayUserInfoShareResponse.Code != "10000" {
		info := aliRsp.AlipayUserInfoShareResponse
		return nil, fmt.Errorf(`{"code":"%v","msg":"%v","sub_code":"%v","sub_msg":"%v"}`, info.Code, info.Msg, info.SubCode, info.SubMsg)
	}
	aliRsp.SignData = getSignData(bs)
	return
}

//alipay.open.auth.token.app(换取应用授权令牌)
//    文档地址：https://docs.open.alipay.com/api_9/alipay.open.auth.token.app
func (a *AliPayClient) AlipayOpenAuthTokenApp(body BodyMap) (aliRsp *AlipayOpenAuthTokenAppResponse, err error) {
	var bs []byte
	grantType := body.Get("grant_type")
	if grantType == null {
		return nil, errors.New("grant_type is not allowed to be null")
	}
	code := body.Get("code")
	refreshToken := body.Get("refresh_token")
	if code == null && refreshToken == null {
		return nil, errors.New("code and refresh_token are not allowed to be null at the same time")
	}
	if bs, err = a.doAliPay(body, "alipay.open.auth.token.app"); err != nil {
		return
	}
	aliRsp = new(AlipayOpenAuthTokenAppResponse)
	if err = json.Unmarshal(bs, aliRsp); err != nil {
		return nil, err
	}
	if aliRsp.AlipayOpenAuthTokenAppResponse.Code != "10000" {
		info := aliRsp.AlipayOpenAuthTokenAppResponse
		return nil, fmt.Errorf(`{"code":"%v","msg":"%v","sub_code":"%v","sub_msg":"%v"}`, info.Code, info.Msg, info.SubCode, info.SubMsg)
	}
	aliRsp.SignData = getSignData(bs)
	return
}

//zhima.credit.score.get(芝麻分)
//    文档地址：https://docs.open.alipay.com/api_8/zhima.credit.score.get
func (a *AliPayClient) ZhimaCreditScoreGet(body BodyMap) (aliRsp *ZhimaCreditScoreGetResponse, err error) {
	var (
		p1, p2 string
		bs     []byte
	)
	if p1 = body.Get("product_code"); p1 == null {
		body.Set("product_code", "w1010100100000000001")
	}
	if p2 = body.Get("transaction_id"); p2 == null {
		return nil, errors.New("transaction_id is not allowed to be null")
	}
	if bs, err = a.doAliPay(body, "zhima.credit.score.get"); err != nil {
		return
	}
	aliRsp = new(ZhimaCreditScoreGetResponse)
	if err = json.Unmarshal(bs, aliRsp); err != nil {
		return nil, err
	}
	if aliRsp.ZhimaCreditScoreGetResponse.Code != "10000" {
		info := aliRsp.ZhimaCreditScoreGetResponse
		return nil, fmt.Errorf(`{"code":"%v","msg":"%v","sub_code":"%v","sub_msg":"%v"}`, info.Code, info.Msg, info.SubCode, info.SubMsg)
	}
	aliRsp.SignData = getSignData(bs)
	return
}

//向支付宝发送请求
func (a *AliPayClient) doAliPay(body BodyMap, method string) (bytes []byte, err error) {
	var (
		bodyStr, sign, url, urlParam string
		bodyBs                       []byte
		res                          gorequest.Response
		errs                         []error
	)
	if body != nil {
		if bodyBs, err = json.Marshal(body); err != nil {
			return nil, fmt.Errorf("json.Marshal：%v", err.Error())
		}
		bodyStr = string(bodyBs)
	}
	pubBody := make(BodyMap)
	pubBody.Set("app_id", a.AppId)
	pubBody.Set("method", method)
	pubBody.Set("format", "JSON")
	if a.AppCertSN != null {
		pubBody.Set("app_cert_sn", a.AppCertSN)
	}
	if a.AlipayRootCertSN != null {
		pubBody.Set("alipay_root_cert_sn", a.AlipayRootCertSN)
	}
	if a.ReturnUrl != null {
		pubBody.Set("return_url", a.ReturnUrl)
	}
	if a.Charset == null {
		pubBody.Set("charset", "utf-8")
	} else {
		pubBody.Set("charset", a.Charset)
	}
	if a.SignType == null {
		pubBody.Set("sign_type", "RSA2")
	} else {
		pubBody.Set("sign_type", a.SignType)
	}
	pubBody.Set("timestamp", time.Now().Format(TimeLayout))
	pubBody.Set("version", "1.0")
	if a.NotifyUrl != null {
		pubBody.Set("notify_url", a.NotifyUrl)
	}
	if a.AppAuthToken != null {
		pubBody.Set("app_auth_token", a.AppAuthToken)
	}
	if a.AuthToken != null {
		pubBody.Set("auth_token", a.AuthToken)
	}
	if bodyStr != null {
		pubBody.Set("biz_content", bodyStr)
	}
	if sign, err = getRsaSign(pubBody, pubBody.Get("sign_type"), FormatPrivateKey(a.PrivateKey)); err != nil {
		return
	}
	pubBody.Set("sign", sign)
	urlParam = FormatAliPayURLParam(pubBody)
	if method == "alipay.trade.app.pay" {
		return []byte(urlParam), nil
	}
	if method == "alipay.trade.page.pay" {
		if !a.IsProd {
			return []byte(zfbSandboxBaseUrl + "?" + urlParam), nil
		} else {
			return []byte(zfbBaseUrl + "?" + urlParam), nil
		}
	}
	agent := HttpAgent()
	if !a.IsProd {
		url = zfbSandboxBaseUrlUtf8
	} else {
		url = zfbBaseUrlUtf8
	}
	if res, bytes, errs = agent.Post(url).Type("form-data").SendString(urlParam).EndBytes(); len(errs) > 0 {
		return nil, errs[0]
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP Request Error, StatusCode = %v", res.StatusCode)
	}
	if method == "alipay.trade.wap.pay" {
		if res.Request.URL.String() == zfbSandboxBaseUrl || res.Request.URL.String() == zfbBaseUrl {
			return nil, errors.New("alipay.trade.wap.pay error,please check the parameters")
		}
		return []byte(res.Request.URL.String()), nil
	}
	return
}

func getSignData(bs []byte) (signData string) {
	str := string(bs)
	indexStart := strings.Index(str, `":`)
	indexEnd := strings.Index(str, `,"sign"`)
	signData = str[indexStart+2 : indexEnd]
	return
}
