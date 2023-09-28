package service

import (
	"context"
	"errors"
	"fmt"
	"game_assistantor/network"
	"game_assistantor/utils"
	"strings"

	"github.com/go-redis/redis"
	"github.com/rs/zerolog/log"
)

type GroupMessage struct {
	SenderId   string // 发送id
	ReceiverId string // 接受id
	DeviceType string
	Message    string
}

var (
	env string
	app string

	host        string
	hostAddress string

	commandPort    int // grpc 命令下行服务端口
	serverType     string
	connectionType string // 连接类型 是否断开连接
	nativeServer   *network.NaiveServer
	redisClient    *redis.Client
	logLevel       string
	platForm       string
	protocol       string // protocol

	redisHost        string
	redisPort        int
	redisPassword    string
	redisDB          int
	redisReadTimeout int

	port            int
	socketTimeout   int
	sendCommandPort int
	messageChan     chan GroupMessage
)

func init() {
	messageChan = make(chan GroupMessage, 100)
}

func HandMessage(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("start to quit handle")
			return
		case msg := <-messageChan:
			SendGroupMessage(msg)
		}
	}
}

func SendGroupMessage(msg GroupMessage) {
	connectionMap := nativeServer.GetConnectionMap()
	fromId := msg.SenderId
	receiverId := msg.ReceiverId
	log.Info().Msgf("from id is: %s, receive id is: %s", fromId, receiverId)
	log.Info().Msgf("connectionMap is: %v", connectionMap)

	if receiverId == "all" {
		for _, conn := range connectionMap {
			log.Info().Msgf("device id: %s, connection type is: %s", conn.GetID(), conn.DeviceType)
			if conn.GetID() == fromId {
				// 不转发自己
				continue
			}
			conn.SendMessageChan <- fmt.Sprintf("%s\n", msg.Message)
		}
	} else {
		conn, ok := connectionMap[receiverId]
		if !ok {
			log.Error().Msgf("attemp to sync data to empty conn")
			return
		}
		conn.SendMessageChan <- fmt.Sprintf("%s\n", msg.Message)
	}

}

func StartDeviceService(ctx context.Context) {

	//redisClient = util.GetRedisClientWithTimeOut(redisHost, redisPort, redisPassword, redisDB, redisReadTimeout)

	serverConfig := network.ServerConfig{
		Address: "0.0.0.0:12000",
		Timeout: socketTimeout,
	}
	nativeServer = network.NewNativeServer(&serverConfig)
	nativeServer.RegisterCallbacks(connectionMade, connectionLost, messageReceived)

	go HandMessage(ctx)
	go nativeServer.Listen()

	select {
	case <-ctx.Done():
		nativeServer.Stop()
		break
	}
}

func connectionMade(c *network.Connection, deviceId string) {
	c.MarkConnection(deviceId)
	log.Info().Msgf("Receive new connection from device id: %s, device type is: %s", deviceId, c.DeviceType)

	//dataSet := make(map[string]interface{})
	//dataSet["deviceId"] = deviceId
	//dataSet["host"] = host
	//dataSet["address"] = hostAddress
	//dataSet["last_updated"] = time.Now().Unix()
	//vinKey := fmt.Sprintf("%s_%s", common.ConnectionKey, deviceId)
	//
	//.HMSet(vinKey, dataSet)
}

func messageReceived(c *network.Connection, data []byte) {

	segment := string(data)
	log.Info().Msgf("Receive segment: %s", segment)
	if len(c.Left) > 0 {
		segment = c.Left + segment
	}

	messages, _, _ := Split(segment)
	for _, message := range messages {
		info := strings.Split(message, "@")
		deviceId := info[0][1:]
		deviceType := info[1]
		receiverId := info[2]

		if deviceType != "" {
			c.DeviceType = deviceType
		}
		if c.IsFirstMessage {
			connectionMade(c, deviceId)
			c.IsFirstMessage = false
		}
		var msg GroupMessage
		msg.SenderId = deviceId
		msg.DeviceType = deviceType
		msg.Message = message
		msg.ReceiverId = receiverId
		if receiverId != "init" {
			messageChan <- msg
		}
	}
}

func connectionLost(c *network.Connection, err error) {
	deviceId := c.GetID()
	log.Info().Msgf("Connection lost with client, deviceId: %s, err: %v", deviceId, err)
}

// #device_server@publisher@32@up@end_string#g#
// #device_server@publisher@88@up@end_string#ring#
// #device_server@publisher@88@down@end_strin
func Split(segment string) (messages []string, left string, invalidMessage [][]byte) {

	if len(segment) < 12 {
		left = segment
		return
	}
	startFlag := "#"
	var indexList []int
	for index := 0; index < len(segment); index++ {
		sf := string(segment[index])
		if sf == startFlag {
			indexList = append(indexList, index)
		}
	}
	if len(indexList) == 1 {
		left = segment
	}
	var i int
	log.Info().Msgf("segment is: %s, index list is: %v", segment, indexList)
	for i < len(indexList) {
		message := segment[indexList[i] : indexList[i+1]+1]
		if len(message) == 0 {
			break
		}
		if message == "##" {
			break
		}
		// log.Info().Msgf("message is: %s", message)
		ok := utils.CheckPackage(message)

		if !ok {
			log.Error().Msgf("package error: %s", message)
			continue
		}
		messages = append(messages, message)
		i += 2
	}

	return
}

func GetDeviceList() (deviceList []string, err error) {
	deviceMap := nativeServer.GetConnectionMap()
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
}

func SendCommand(deviceId string, command []byte) (err error) {
	conn := nativeServer.GetConn(deviceId)
	if conn != nil {
		err = conn.Send(command)
	}
	return
}
