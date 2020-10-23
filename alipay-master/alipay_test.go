package alipay_test

import (
	pub "public_yfm"

	"common/alipay-master"
)

var (
	appID = "2019082266405426"

	// RSA2(SHA256)
	aliPublicKey = "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA2MhEVUp+rRRyAD9HZfiSg8LLxRAX18XOMJE8/MNnlSSTWCCoHnM+FIU+AfB+8FE+gGIJYXJlpTIyWn4VUMtewh/4C8uwzBWod/3ilw9Uy7lFblXDBd8En8a59AxC6c9YL1nWD7/sh1szqej31VRI2OXQSYgvhWNGjzw2/KS1GdrWmdsVP2hOiKVy6TNtH7XnCSRfBBCQ+LgqO1tE0NHDDswRwBLAFmIlfZ//qZ+a8FvMc//sUm+CV78pQba4nnzsmh10fzVVFIWiKw3VDsxXPRrAtOJCwNsBwbvMuI/ictvxxjUl4nBZDw4lXt5eWWqBrnTSzogFNOk06aNmEBTUhwIDAQAB"

	privateKey = "MIIEogIBAAKCAQEAg6OetFjLGIBg/mR5FA5cGEk99my08rxq2meeewukvQRy6PyORj/kIfQaK40p4Y4b/jos66M6FOHGRuszIx+3X3/TYzdem05BzcYscn7eVrJPJ735qRAC+iS7KNGPsAi0CyZ8sGcy0ff4xqpFIem0c1NMpcio6PRceh+UcG2+UEbB1C+QSV5N4zjltgy4WbUuO9req0VHFxWN+eGaLQVe+hZJx4pVYDajeLpStwhuC51JBOY80ka3ZvJyGVUixWFcjMqG+NAyWMlk9gMqaehMtyQu0hPd3CiX0WzFzodl2h4VAuzJxv/2KGwtde6yMo+GO+WyvliavX6ZKvRdc6aGjQIDAQABAoIBAG+K+TN/+bftMELfB+mCsW6ywRpJypnUJgaivpsspo6zclsRhczo3noWyQYWO2Kwc+/bg6y3RsPi+4ukSFR+z3bQbWIozLV0fjXKsmbiMavssz3Nr/sdYihrb1uLFuHmvgikuAsRpvJb2CUeqi2uRgVilBP2D4o1ZRbmI3WQyC5kTZwZgJOkpFAd1p2wM/EIIbun6DwTDszSEgRRCp0/METXVfd3yJ8Mh7oYo1EtXfcuWbAGGWNquxdtsacdlsXixnzMyR/LFVn1lAFIM9Hs5kmHNdkpm2PqxhkvdAy2ijB8ZoACwqUnMfl+jgd03daEtRernHiyWSiWlwvi2FfAtIECgYEAy5WxcLTpEa+oFbDY9+dEjBwMrxrF0OXlP5XK+mjBPBSaT01Wz9TBJErK8dTcySjKJshNtlVEq3AsbavBcpwKXFEF0C8AOIHAv+yXvqrv2/UKoTL2h/6yb6JO8teP8hJixIYl4ph7hJcJ0kSjWCZHguHZN5jkha36TVt3pXY2ODkCgYEApYf6vI2E2piRr2JPbCJt5DapyT3yMw2msfkydGdzQvHXPhoHj4tl+eJirRY9yvY0GtAdhnZ/iVz7VxvDTznMqeZzYrM5AP1kvhVUrfQcV4C+AJxLujnaf0aHJYmWjltBuhchgQtAfzXPQ5SPnChVwKc3VaNnIQO3aboWmRKoePUCgYAzfR2OcsLLjVCGg96r/BqzENkIZE4JgktTpI/ceyf8CP3p9pZxI87hXeUr+nkIiz9tRZWZ+sDOVyV1a04WrW5VYMyGlYyJvg9Auxa5y0O0rqnMkTYWuQzp/PPYqTonsAy4xXDJeWUr6IM8Yc2qGqxVZsdoL0wEnzbB39NHzrjxyQKBgDpGRN9ccwkB7UfxNES9WjKdi+htBncytxywvjJ8uPc4bK5QO5ktWhk+ub51tgtd4boOylYsIXoaYeGoxHl/v62Qk86LieXvTygcGlOjPNcRW9KbM428EE/+ZFWyum4jcmAxBHqJm4stRmpkQqqXCJlqRPDBNe1JgaiW+p2pE+aBAoGAYswlFfb9MlPJYM695EIpaYV7lvic1MKZlsNZNMjqLtYTqwibYxE43lX0+Fu/cE4cdPuTJsz1weh5wBKI9CrpeXFJ6zWrVqG8aSA5Ikr3doMP+WMlTwGtmkIOkIbPyqKI1RqUGDW39rTYc+Gj2lNdNQXQ5vH5/kJMMhqkxh8sRYA="
)

var client *alipay.Client

func init() {
	return
	var err error
	client, err = alipay.New(appID, "", privateKey, false)
	if err != nil {
		pub.PrintLog("初始化支付宝失败, 错误信息为", err.Error())
		panic("初始化支付宝失败")
	}

	if client.IsProduction() {
		// https://docs.open.alipay.com/291/105971#Krqvg
		pub.PrintLog("加载证书", client.LoadAppPublicCertFromFile("appCertPublicKey_2019082266405426.crt"))
		pub.PrintLog("加载证书", client.LoadAliPayRootCertFromFile("alipayRootCert.crt"))
		pub.PrintLog("加载证书", client.LoadAliPayPublicCertFromFile("alipayCertPublicKey_RSA2.crt"))
	}
}
