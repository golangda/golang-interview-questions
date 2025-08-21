package main

import (
	"errors"
	"fmt"
)

var errCustom = errors.New("this is my custom error")

func getErr() error{
	return  errCustom
}

func main(){
	fmt.Println(getErr())
}