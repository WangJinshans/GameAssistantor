package main

import (
	"fmt"

	"github.com/wenlng/go-captcha/captcha"
)

func main() {
	// Captcha Single Instances
	capt := captcha.GetCaptcha()

	// Generate Captcha
	dots, b64, tb64, key, err := capt.Generate()
	if err != nil {
		panic(err)
	}

	// Main image base64 code
	fmt.Println(b64)

	// Thumb image base64 code
	fmt.Println(len(tb64))

	// Only key
	fmt.Println(key)

	// Dot data For verification
	fmt.Println(dots) // 与前端直接提交点击的位置信息进行比对
}
