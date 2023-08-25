package network

import (
	"bufio"
	"bytes"
	"net"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

var (
	mapMutex = sync.RWMutex{}
)

// Connection holds info about connection
type Connection struct {
	conn           *net.TCPConn
	Server         Server
	ResidueBytes   []byte
	id             string
	peerAddress    string
	MessageChan    chan []byte
	ExitChan       chan struct{}
	IsFirstMessage bool // 是否是第一次发送消息
}

// ServerConfig involve server's configurations
type ServerConfig struct {
	Address       string
	Timeout       int
	MaxQps        int64
	MaxConnection int
}

// 处理离线删除key后又更新key的问题
type WarpedConnectionInfo struct {
	Data map[string]interface{}
	Conn *Connection
}

// NaiveServer represent a server
type NaiveServer struct {
	Config *ServerConfig

	funcOnConnectionMade   func(c *Connection, vin string)
	funcOnConnectionClosed func(c *Connection, err error)
	funcOnMessageReceived  func(c *Connection, message []byte)

	started bool

	deviceMap map[string]*Connection
}

func NewNativeServer(config *ServerConfig) (server *NaiveServer) {
	server = &NaiveServer{
		Config:    config,
		deviceMap: make(map[string]*Connection),
	}
	return
}

func (c *Connection) String() string {
	return c.peerAddress
}

// SetID set id for connection
func (c *Connection) SetID(id string) {
	c.id = id
	c.Server.addConn(id, c)
}

// GetID return a connection's id
func (c *Connection) GetID() string {
	return c.id
}

func (c *Connection) listen() {
	reader := bufio.NewReader(c.conn)
	buffer := make([]byte, 4096, 4096)
	var read int
	var err error
	for {
		err = c.conn.SetReadDeadline(time.Now().Add(c.Server.GetTimeout()))
		if err != nil {
			log.Error().Msg(err.Error())
			return
		}
		var buffers bytes.Buffer
		read, err = reader.Read(buffer)
		if err != nil {
			log.Error().Msgf("close connection: %s, connection is: %v", c.id, c)
			c.ExitChan <- struct{}{}
			c.Server.OnConnectionClosed(c, err)
			c.conn.Close()
			c.Server.removeConn(c.id)
			c.conn = nil
			return
		}
		// go transfer byte slice argument as point, so we copy
		// bytes here to prevent race conditions
		buffers.Write(buffer[:read])
		c.MessageChan <- buffers.Bytes()
		buffers.Reset()
	}
}

func (c *Connection) dispatchMessage() {
	for {
		select {
		case <-c.ExitChan:
			log.Info().Msgf("stop working, connection:%s", c.id)
			return
		case message := <-c.MessageChan:
			c.Server.OnMessageReceived(c, message)
		}
	}
}

func (c *Connection) Send(message []byte) error {
	_, err := c.conn.Write(message)
	return err
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

// RegisteCallbacks registe three callbacks
func (s *NaiveServer) RegisterCallbacks(onConnectionMade func(c *Connection, vin string), onConnectionClosed func(c *Connection, err error), onMessageReceived func(c *Connection, message []byte)) {
	s.funcOnConnectionMade = onConnectionMade
	s.funcOnConnectionClosed = onConnectionClosed
	s.funcOnMessageReceived = onMessageReceived
}

func (s *NaiveServer) OnConnectionClosed(c *Connection, err error) {
	s.funcOnConnectionClosed(c, err)
}

func (s *NaiveServer) OnConnectionMade(c *Connection, vin string) {
	s.funcOnConnectionMade(c, vin)
}

func (s *NaiveServer) OnMessageReceived(c *Connection, message []byte) {
	s.funcOnMessageReceived(c, message)
}

func (s *NaiveServer) GetTimeout() time.Duration {
	return time.Duration(s.Config.Timeout) * time.Second
}

func (s *NaiveServer) addConn(id string, conn *Connection) {
	if id == "" {
		return
	}

	mapMutex.Lock()
	defer mapMutex.Unlock()

	s.deviceMap[id] = conn
}

func (s *NaiveServer) removeConn(id string) {
	if id == "" {
		return
	}

	mapMutex.Lock()
	defer mapMutex.Unlock()

	delete(s.deviceMap, id)
}

func (s *NaiveServer) GetConn(id string) *Connection {
	mapMutex.RLock()
	defer mapMutex.RUnlock()

	return s.deviceMap[id]
}

func (s *NaiveServer) GetIDSet() map[string]bool {
	mapMutex.RLock()
	defer mapMutex.RUnlock()

	log.Debug().Msgf("device map: %v", s.deviceMap)
	m := make(map[string]bool, len(s.deviceMap))
	for id := range s.deviceMap {
		m[id] = true
	}
	return m
}

func (s *NaiveServer) Listen() {
	s.started = true
	defer func() { s.started = false }()

	addr, err := net.ResolveTCPAddr("tcp", s.Config.Address)
	if err != nil {
		log.Panic().Err(err)
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Panic().Err(err)
	}
	defer listener.Close()
	//netutil.LimitListener(listener, s.MaxConnection) // 最大连接数

	for {
		if !s.started {
			log.Info().Msg("server is going down")
			return
		}
		err = listener.SetDeadline(time.Now().Add(3 * time.Second))
		if err != nil {
			log.Error().Msg(err.Error())
			return
		}
		var conn *net.TCPConn
		conn, err = listener.AcceptTCP()
		if err != nil {
			continue
		}

		c := Connection{
			conn:           conn,
			Server:         s,
			MessageChan:    make(chan []byte, 20),
			ExitChan:       make(chan struct{}),
			IsFirstMessage: true,
		}
		// set peer address at start to avoid frequently system calls
		c.peerAddress = c.RemoteAddr().String()

		go c.dispatchMessage()
		go c.listen()
	}
}

func (s *NaiveServer) Stop() {
	s.started = false
}
