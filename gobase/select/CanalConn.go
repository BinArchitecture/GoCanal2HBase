package _select

import (
	"time"
	"net"
	"github.com/golang/glog"
	"strings"
	"sync"
)

type CanalConn struct {
	addr string
	soTimeout time.Duration
	conn net.Conn
	locker sync.RWMutex
}
func NewCanalConn(addr string,soTimeout time.Duration) (con *CanalConn,err error){
	con=new(CanalConn)
	con.addr=addr
	con.soTimeout=soTimeout
	con.conn,err=net.Dial("tcp",addr)
	if err != nil {
		glog.Error(err)
		return nil,err
	}
	con.conn.SetDeadline(time.Now().Add(soTimeout))
	return con,nil
}

func (self *CanalConn) ReConnect() (err error){
	self.locker.Lock()
	defer self.locker.Unlock()
	self.conn.Close()
	self.conn,err=net.Dial("tcp",self.addr)
	if err != nil {
		glog.Error(err)
		return err
	}
	self.conn.SetDeadline(time.Now().Add(self.soTimeout))
	return nil
}

func (self *CanalConn) Read(b []byte) (n int, err error){
	n,err= self.conn.Read(b)
	if err!=nil {
		netErr,ok:=err.(*net.OpError)
		if ok{
			cf:=netErr.Err.Error()
			if strings.Contains(cf,"timeout"){
				glog.Warning(netErr.Err)
				return 0,nil
			}
		}
	}else if err!=nil{
		return -1,err
	}
	return n,err
}

func (self *CanalConn) Write(b []byte) (n int, err error){
	return self.conn.Write(b)
}

func (self *CanalConn) Close() error{
	return self.conn.Close()
}

func (self *CanalConn) LocalAddr() net.Addr{
	return self.conn.LocalAddr()
}

func (self *CanalConn) RemoteAddr() net.Addr{
	return self.conn.RemoteAddr()
}

func (self *CanalConn) SetDeadline(t time.Time) error{
	return self.conn.SetDeadline(t)
}

func (self *CanalConn) SetReadDeadline(t time.Time) error{
	return self.conn.SetReadDeadline(t)
}

func (self *CanalConn) SetWriteDeadline(t time.Time) error{
	return self.conn.SetWriteDeadline(t)
}