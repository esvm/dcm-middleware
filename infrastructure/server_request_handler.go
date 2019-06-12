package infrastructure

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"

	log "github.com/sirupsen/logrus"
)

type ServerRequestHandler struct {
	Port int
}

func NewServerRequestHandler(port int) *ServerRequestHandler {
	return &ServerRequestHandler{
		Port: port,
	}
}

func (server *ServerRequestHandler) Listen(callback func(conn net.Conn, host string, bytes []byte) error) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", server.Port))
	defer listener.Close()
	if err != nil {
		return err
	}

	log.Infof("Listening on port: %d\n", server.Port)

	for {
		conn, err := listener.Accept()
		host := conn.RemoteAddr().(*net.TCPAddr).IP.String()
		log.Infof("Receive connection from: %s", host)
		if err != nil {
			log.Debug(err.Error())
			continue
		}

		go func() {
			defer conn.Close()
			for {
				var size = make([]byte, 4)
				_, err = io.ReadFull(conn, size)

				if err != nil {
					log.Infof("%s disconnected", host)
					log.Debug(err)
					return
				}
				log.Debugf("Read size: %d", binary.BigEndian.Uint32(size))

				var bytes = make([]byte, binary.BigEndian.Uint32(size))
				_, err = io.ReadFull(conn, bytes)
				if err != nil {
					log.Infof("%s disconnected", host)
					log.Debug(err)
					return
				}

				err = callback(conn, host, bytes)
				if err != nil {
					log.Infof("%s disconnected", host)
					log.Debug(err)
				}
			}
		}()
	}
}

func Send(conn net.Conn, msg []byte) error {
	_, err := conn.Write(msg)
	return err
}
