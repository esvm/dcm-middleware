package infrastructure

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"

	log "github.com/sirupsen/logrus"
)

type ClientRequestHandler struct {
	Host string
	Port int
	conn net.Conn
}

func NewClientRequestHandler(host string, port int) *ClientRequestHandler {
	return &ClientRequestHandler{
		Host: host,
		Port: port,
	}
}

func (client *ClientRequestHandler) Connect() error {
	var err error
	client.conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", client.Host, client.Port))
	return err
}

func (client *ClientRequestHandler) Close() error {
	if client.conn == nil {
		return errors.New("Connection is nil")
	}

	return client.conn.Close()
}

func (client *ClientRequestHandler) Send(msg []byte) error {
	_, err := client.conn.Write(msg)
	return err
}

func (client *ClientRequestHandler) Receive() ([]byte, error) {
	log.Debug("conn", client.conn)
	var size = make([]byte, 4)
	_, err := io.ReadFull(client.conn, size)
	if err != nil {
		return nil, err
	}

	log.Debug("conn 3", client.conn)

	var bytes = make([]byte, binary.BigEndian.Uint32(size))
	_, err = io.ReadFull(client.conn, bytes)

	return bytes, err
}
