package gobase

import (
	"bytes"
	"encoding/binary"
	"github.com/golang/glog"
	"sync"
	"time"
	otterproto "github.com/BinArchitecture/GoCanal2HBase/com_alibaba_otter_canal_protocol/canal"
	"github.com/golang/protobuf/proto"
	"strconv"
	canal_entry "github.com/BinArchitecture/GoCanal2HBase/com_alibaba_otter_canal_protocol/entry"
)

type CanalConnector struct {
	addr                  string
	username              string
	password              string
	soTimeout             time.Duration
	filter                string
	supportedCompressions []otterproto.Compression
	clientIdentity        *ClientIdentity
	connected             bool
	rollbackOnConnect     bool
	rollbackOnDisConnect  bool
	readDataLock          sync.RWMutex
	writeDataLock         sync.RWMutex
	readHeader            *bytes.Buffer
	writeHeader           *bytes.Buffer
	conn                  *CanalConn
}

type Message struct {
	id      int64
	entries []*canal_entry.Entry
}

type ClientIdentity struct {
	destination string
	clientId    int16
	filter      string
}

func NewCanalConnector(addr string, username string, password string, destination string) *CanalConnector {
	conn := new(CanalConnector)
	conn.addr = addr
	conn.username = username
	conn.password = password
	conn.soTimeout = 60 * time.Second
	conn.supportedCompressions = make([]otterproto.Compression, 1, 100)
	conn.connected = false
	conn.rollbackOnConnect = true
	conn.rollbackOnDisConnect = false
	conn.readHeader = bytes.NewBuffer(make([]byte, 4))
	conn.writeHeader = bytes.NewBuffer(make([]byte, 4))
	conn.clientIdentity = new(ClientIdentity)
	conn.clientIdentity.clientId = 1001
	conn.clientIdentity.destination = destination
	conn.clientIdentity.filter = ""
	return conn
}

func (self *CanalConnector) Connect() (err error) {
	if self.connected {
		return nil
	}
	self.conn, err = NewCanalConn(self.addr, self.soTimeout)
	if err != nil {
		glog.Error(err)
		return err
	}
	p, err := self.parsePacket()
	if err != nil {
		glog.Error(err)
		return err
	}
	if p.GetVersion() != 1 {
		panic("unsupported version at this client.")
	}
	if p.GetType() != otterproto.PacketType_HANDSHAKE {
		panic("expect handshake but found other type.")
	}
	handshake := new(otterproto.Handshake)
	proto.Unmarshal(p.Body, handshake)
	self.supportedCompressions = append(self.supportedCompressions, handshake.SupportedCompressions...)

	ca := &otterproto.ClientAuth{
		Username:        &self.username,
		Password:        []byte(self.password),
		NetReadTimeout:  proto.Int32(30000),
		NetWriteTimeout: proto.Int32(30000),
	}
	bq, err := proto.Marshal(ca)
	p = new(otterproto.Packet)
	ptype := otterproto.PacketType_CLIENTAUTHENTICATION
	p.Type = &ptype
	p.Body = bq
	bqq, err := proto.Marshal(p)
	self.writeWithHeader(bqq)
	p, err = self.parsePacket()
	if err != nil {
		glog.Error(err)
		return err
	}
	if p.GetType() != otterproto.PacketType_ACK {
		panic("unexpected packet type when ack is expected")
	}
	ack := new(otterproto.Ack)
	proto.Unmarshal(p.Body, ack)
	if ack.ErrorCode != nil {
		code := int(*ack.ErrorCode)
		if code > 0 {
			str := "something goes wrong when doing authentication:"
			str += ack.GetErrorMessage()
			panic(str)
		}
	}
	self.connected = true
	return nil
}

func (self *CanalConnector) Subscribe(filter string) error {
	pty := otterproto.PacketType_SUBSCRIPTION
	return self.subunsub(filter, &pty)
}

func (self *CanalConnector) subunsub(filter string, sutype *otterproto.PacketType) error {
	sub := new(otterproto.Sub)
	sub.Destination = proto.String(self.clientIdentity.destination)
	sub.ClientId = proto.String(strconv.Itoa(int(self.clientIdentity.clientId)))
	if filter != "" {
		sub.Filter = proto.String(filter)
	}
	p := new(otterproto.Packet)
	p.Type = sutype
	p.Body, _ = proto.Marshal(sub)
	by, _ := proto.Marshal(p)
	err := self.writeWithHeader(by)
	if err != nil {
		glog.Error(err)
		return err
	}
	p, err = self.parsePacket()
	if err != nil {
		glog.Error(err)
		return err
	}
	ack := new(otterproto.Ack)
	proto.Unmarshal(p.GetBody(), ack)
	if ack.ErrorCode != nil {
		code := int(*ack.ErrorCode)
		if code > 0 {
			str := "failed to subscribe with reason: "
			str += ack.GetErrorMessage()
			panic(str)
		}
	}
	if filter != "" {
		self.clientIdentity.filter = filter
	}
	return nil
}

func (self *CanalConnector) Unsubscribe() error {
	pty := otterproto.PacketType_UNSUBSCRIPTION
	return self.subunsub("", &pty)
}

func (self *CanalConnector) Disconnect() error {
	if self.rollbackOnDisConnect && self.connected {
		err := self.Rollback()
		if err != nil {
			glog.Error(err)
			return err
		}
	}
	self.connected = false
	if self.conn.conn != nil {
		self.conn.conn.Close()
	}
	return nil
}

func (self *CanalConnector) Rollback() error {
	ca := new(otterproto.ClientRollback)
	ca.ClientId = proto.String(strconv.Itoa(int(self.clientIdentity.clientId)))
	ca.BatchId = proto.Int64(0)
	ca.Destination = proto.String(self.clientIdentity.destination)
	p := new(otterproto.Packet)
	pty := otterproto.PacketType_CLIENTROLLBACK
	p.Type = &pty
	p.Body, _ = proto.Marshal(ca)
	by, _ := proto.Marshal(p)
	err := self.writeWithHeader(by)
	if err != nil {
		glog.Error(err)
		return err
	}
	return nil
}

func (self *CanalConnector) GetWithoutAck(batchSize int) (*Message, error) {
	if batchSize <= 0 {
		batchSize = 1000
	}
	get := new(otterproto.Get)
	get.ClientId = proto.String(strconv.Itoa(int(self.clientIdentity.clientId)))
	get.Destination = proto.String(self.clientIdentity.destination)
	get.AutoAck = proto.Bool(false)
	get.FetchSize = proto.Int32(int32(batchSize))
	get.Timeout = proto.Int64(-1)
	get.Unit = proto.Int32(2) //java TimeUnit.MILLISECONDS.ordinal()
	p := new(otterproto.Packet)
	pType := otterproto.PacketType_GET
	p.Type = &pType
	p.Body, _ = proto.Marshal(get)
	pp, _ := proto.Marshal(p)
	err := self.writeWithHeader(pp)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	msg, err := self.receiveMessages()
	return msg, err
}

func (self *CanalConnector) Get(batchId int) (*Message, error) {
	msg, err := self.GetWithoutAck(batchId)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	self.Ack(int(msg.id))
	return msg, nil
}

func (self *CanalConnector) Ack(batchId int) error {
	ca := new(otterproto.ClientAck)
	ca.ClientId = proto.String(strconv.Itoa(int(self.clientIdentity.clientId)))
	ca.Destination = proto.String(self.clientIdentity.destination)
	ca.BatchId = proto.Int64(int64(batchId))
	p := new(otterproto.Packet)
	pType := otterproto.PacketType_CLIENTACK
	p.Type = &pType
	p.Body, _ = proto.Marshal(ca)
	pp, _ := proto.Marshal(p)
	err := self.writeWithHeader(pp)
	if err != nil {
		glog.Error(err)
		return err
	}
	return nil
}

func (self *CanalConnector) receiveMessages() (*Message, error) {
	p, err := self.parsePacket()
	if err != nil {
		return nil, err
	}
	switch p.GetType() {
	case otterproto.PacketType_MESSAGES:
		if p.GetCompression() != otterproto.Compression_NONE {
			panic("compression is not supported in this connector")
		}
		msgs := new(otterproto.Messages)
		proto.Unmarshal(p.GetBody(), msgs)
		msg := new(Message)
		msg.id = *msgs.BatchId
		msg.entries = make([]*canal_entry.Entry, len(msgs.GetMessages()))
		for _, bb := range msgs.GetMessages() {
			entry := new(canal_entry.Entry)
			proto.Unmarshal(bb, entry)
			msg.entries = append(msg.entries, entry)
		}
		return msg, nil
	case otterproto.PacketType_ACK:
		ack := new(otterproto.Ack)
		proto.Unmarshal(p.GetBody(), ack)
		str := "something goes wrong with reason: "
		str += ack.GetErrorMessage()
		panic(str)
	default:
		str := "unexpected packet type: "
		str += p.GetType().String()
		panic(str)
	}
}

func (self *CanalConnector) parsePacket() (*otterproto.Packet, error) {
	bb, err := self.readNextPacket()
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	p := new(otterproto.Packet)
	proto.Unmarshal(bb, p)
	return p, nil
}

func (self *CanalConnector) writeWithHeader(body []byte) error {
	self.writeDataLock.Lock()
	defer self.writeDataLock.Unlock()
	self.writeHeader.Reset()
	err := binary.Write(self.writeHeader, binary.BigEndian, int32(len(body)))
	if err != nil {
		glog.Error(err, self.addr)
		return err
	}
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, self.writeHeader.Bytes())
	binary.Write(buf, binary.BigEndian, body)
	self.conn.Write(buf.Bytes())
	//self.conn.Write(body)
	return nil
}

func (self *CanalConnector) readNextPacket() ([]byte, error) {
	self.readDataLock.Lock()
	defer self.readDataLock.Unlock()
	self.readHeader.Reset()
	size := 1024
	b := make([]byte, size)
	for {
		self.conn.SetReadDeadline(time.Now().Add(300 * time.Second))
		n, err := self.conn.Read(b)
		if err != nil {
			glog.Error(err, self.addr)
			self.conn.ReConnect()
			return nil, err
		}
		self.readHeader.Write(b[:n])
		if n == size {
			continue
		}
		break
	}
	var bodyLen int32
	err := binary.Read(self.readHeader, binary.BigEndian, &bodyLen)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	bodyBuf := make([]byte, int(bodyLen))
	self.readHeader.Read(bodyBuf)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	return bodyBuf, nil
}

//func (self *DefalutRemotingClient) handlerConn(conn net.Conn, addr string) {
//	b := make([]byte, 1024)
//	var length, headerLength, bodyLength int32
//	var buf = bytes.NewBuffer([]byte{})
//	var header, body []byte
//	var flag int = 0
//	//var cf int32 =0
//	for {
//		n, err := conn.Read(b)
//		if err != nil {
//			glog.Error(err, addr)
//			self.releaseConn(addr, conn)
//			return
//		}
//
//		_, err = buf.Write(b[:n])
//		if err != nil {
//			glog.Error(err, addr)
//			self.releaseConn(addr, conn)
//			return
//		}
//
//		for {
//			if flag == 0 {
//				if buf.Len() >= 4 {
//					err = binary.Read(buf, binary.BigEndian, &length)
//					if err != nil {
//						glog.Error(err)
//						return
//					}
//					flag = 1
//				} else {
//					break
//				}
//			}
//
//			if flag == 1 {
//				if buf.Len() >= 4 {
//					err = binary.Read(buf, binary.BigEndian, &headerLength)
//					if err != nil {
//						glog.Error(err)
//						return
//					}
//					flag = 2
//				} else {
//					break
//				}
//
//			}
//
//			if flag == 2 {
//				if (buf.Len() > 0) && (buf.Len() >= int(headerLength)) {
//					header = make([]byte, headerLength)
//					_, err = buf.Read(header)
//					if err != nil {
//						glog.Error(err)
//						return
//					}
//					flag = 3
//				} else {
//					break
//				}
//			}
//
//			if flag == 3 {
//				bodyLength = length - 4 - headerLength
//				if bodyLength == 0 {
//					flag = 0
//				} else {
//
//					if buf.Len() >= int(bodyLength) {
//						body = make([]byte, int(bodyLength))
//						_, err = buf.Read(body)
//						if err != nil {
//							glog.Error(err)
//							return
//						}
//						flag = 0
//					} else {
//						break
//					}
//				}
//			}
//
//			if flag == 0 {
//				headerCopy := make([]byte, len(header))
//				bodyCopy := make([]byte, len(body))
//				copy(headerCopy, header)
//				copy(bodyCopy, body)
//				go func() {
//					cmd := decodeRemoteCommand(headerCopy, bodyCopy)
//					resp, ok := self.responseTable.Load(cmd.Opaque)
//					self.responseTable.Delete(cmd.Opaque)
//					if ok {
//						response:=resp.(*ResponseFuture)
//						response.ResponseCommand = cmd
//						if response.invokeCallback != nil {
//							response.SendRequestOK=true
//							//cf=atomic.AddInt32(&cf,1)
//							response.invokeCallback(response)
//							//fmt.Printf("invoke is %d\n",int(cf))
//						}
//
//						if response.done != nil {
//							response.done <- true
//						}
//					} else {
//						//if cmd.Code == NOTIFY_CONSUMER_IDS_CHANGED {
//						//	return
//						//}
//						jsonCmd, err := json.Marshal(cmd)
//
//						if err != nil {
//							glog.Error(err)
//						}
//						glog.Errorf("consume cmd error:%s\n",string(jsonCmd))
//					}
//				}()
//			}
//		}
//
//	}
//}
//
//func (self *DefalutRemotingClient) sendRequest(header, body []byte, conn net.Conn, addr string) error {
//	buf := bytes.NewBuffer([]byte{})
//	binary.Write(buf, binary.BigEndian, int32(len(header) + len(body) + 4))
//	binary.Write(buf, binary.BigEndian, int32(len(header)))
//	binary.Write(buf, binary.BigEndian, header)
//	if body != nil && len(body) > 0 {
//		binary.Write(buf, binary.BigEndian, body)
//	}
//	_, err := conn.Write(buf.Bytes())
//	if err != nil {
//		glog.Error(err)
//		self.releaseConn(addr, conn)
//		return err
//	}
//	return nil
//}
//
//func (self *DefalutRemotingClient) releaseConn(addr string, conn net.Conn) {
//	self.connTableLock.Lock()
//	delete(self.connTable, addr)
//	conn.Close()
//	self.connTableLock.Unlock()
//}
