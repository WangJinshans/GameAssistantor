package service

import (
	"errors"
	"game_assistantor/network"

	"github.com/rs/zerolog/log"
)

var exitChan chan bool
var nativeServer *network.NaiveServer

func StartDeviceService() {

	exitChan = make(chan bool, 1)

	serverConfig := network.ServerConfig{
		Address:       "0.0.0.0:12000",
		Timeout:       60,
		MaxConnection: 100,
	}
	nativeServer = network.NewNativeServer(&serverConfig)
	nativeServer.RegisterCallbacks(connectionMade, connectionLost, messageReceived)

	nativeServer.Listen()
	<-exitChan
}

func connectionMade(c *network.Connection, vin string) {
	log.Info().Msgf("Receive new connection from %v, vin: %s", c.RemoteAddr(), vin)
	c.SetID(vin) // 设置当前连接Id
}

// 不会很频繁
func messageReceived(c *network.Connection, segment []byte) {
	log.Info().Msgf("receive message: %s", segment)
	connectionMade(c, "deviceId01")
}

func connectionLost(c *network.Connection, err error) {
	log.Info().Msgf("Connection lost with client %v, vin: %s, err: %v", c.RemoteAddr(), c.GetID(), err)
	// vin := c.GetID()
}

func GetDeviceList() (deviceList []string, err error) {
	deviceMap := nativeServer.GetIDSet()
	if len(deviceMap) == 0 {
		err = errors.New("empty device")
		return
	}

	for deviceId, _ := range deviceMap {
		if deviceId != "" {
			deviceList = append(deviceList, deviceId)
		}
	}

	return
}

func StopService() {
	nativeServer.Stop()
	exitChan <- true
}

func SendCommand(deviceId string, command []byte) (err error) {
	conn := nativeServer.GetConn(deviceId)
	if conn != nil {
		err = conn.Send(command)
	}

	return
}
