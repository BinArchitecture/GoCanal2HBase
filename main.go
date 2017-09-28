package main

import (
	"fmt"
	"github.com/BinArchitecture/GoCanal2HBase/gobase"
)

func main() {
	conn:=gobase.NewCanalConnector("10.6.30.109:11111","","","example")
	conn.Connect()
	conn.Subscribe("")
	fmt.Printf("succ")
}
