package main

import (
    // "encoding/json"
    "fmt"
    // "log"
)

func main () {
	a := "啊"
	sa  := fmt.Sprintf("%02X", a)
	fmt.Println(sa)
}
