package main

import "fmt"

type UserError struct {
	code string
	msg  string
}

func (e *UserError) Error() string {
	return fmt.Sprintf("code：%d，%s", e.code, e.msg)
}

func main() {

}
