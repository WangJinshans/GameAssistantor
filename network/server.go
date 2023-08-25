package network

import "time"

// Server is a general purpose interface for tcp server
type Server interface {
	RegisterCallbacks(
		onConnectionMade func(c *Connection, vin string),
		onConnectionClosed func(c *Connection, err error),
		onMessageReceived func(c *Connection, message []byte),
	)

	OnConnectionMade(c *Connection, vin string)
	OnConnectionClosed(c *Connection, err error)
	OnMessageReceived(c *Connection, message []byte)

	GetTimeout() time.Duration

	// addConn add a id-conn pair to connections
	addConn(id string, conn *Connection)
	// removeConn remove a named connection to connections
	removeConn(id string)

	// GetConn get connection object via id
	GetConn(id string) *Connection
	// GetIDSet get id set (for connection reporting)
	GetIDSet() map[string]bool

	// Listen start server's running loop
	Listen()
	// Stop stop server's running loop
	Stop()
}
