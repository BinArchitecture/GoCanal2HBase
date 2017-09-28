package main

import (
	"fmt"
	"github.com/BinArchitecture/GoCanal2HBase/gobase"
	"github.com/golang/glog"
	"os"
	"time"
)

func main() {
	conn:=gobase.NewCanalConnector("10.6.30.109:11111","","","example")
	conn.Connect()
	conn.Subscribe("")
	for {
		msg,err:=conn.GetWithoutAck(1024)
		if err!=nil{
			glog.Error(err)
			os.Exit(1)
		}
		id:=msg.GetId()
		entries:=msg.GetEntries()
		fmt.Println(id)
		fmt.Printf("entry:%v\n",entries)
		time.Sleep(10*time.Second)
	}
	//fmt.Printf("succ")
}
