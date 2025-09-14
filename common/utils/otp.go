package utils

import (
	"fmt"

	"github.com/pquerna/otp/totp"
)

const (
	OTP_APP   = "V8PAY"
	OTP_EMAIL = "@v8pay.com"
)

// 生成密钥和地址
func GenOtpKey(name string) string {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      OTP_APP,          // 应用或公司名称
		AccountName: name + OTP_EMAIL, // 用户名/邮箱
	})
	if err != nil {
		return ""
	}

	return key.Secret()
}

// 生成密钥和地址
func GenOtpKeyAndUrl(name string) (string, string) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      name,
		AccountName: name + "@v8pay.com",
	})
	if err != nil {
		return "", ""
	}

	return key.Secret(), key.URL()
}

// 生成otp链接
func GenOtpUrl(name, secret string) string {
	if name == "" || secret == "" {
		return ""
	}
	// otpauth://totp/<Issuer>:<AccountName>?secret=<Secret>&issuer=<Issuer>&period=30
	otpUrl := fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s", OTP_APP, name+OTP_EMAIL, secret, OTP_APP)
	return otpUrl
}
