package main

import "fmt"
type Payment interface{
	Pay(amount float64) error
}
type Alipay struct{}
func (a Alipay) Pay(amount float64) error{
	fmt.Printf("使用支付宝支付 %.2f 元\n", amount)
	return nil
}
type WechatPay struct{}
func (w WechatPay) Pay(amount float64) error{
	fmt.Printf("使用微信支付 %.2f元",amount)
	return nil
}

func Checkout(p Payment,amount float64) {
	err := p.Pay(amount)
	if err != nil{
		fmt.Println("支付失败",err)
	}
	fmt.Println("支付成功")
}

func main() {
    Checkout(Alipay{}, 100.5545666666666666666666666666666666666666666666666660)
    Checkout(WechatPay{}, 88.88)
}