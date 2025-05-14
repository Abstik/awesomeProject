package main

import (
	"fmt"

	"awesomeProject/utils"
)

func main() {
	pas := utils.EncryptPassword("123456")
	fmt.Println(pas)
}
