package distribution

import (
	"net"

	"github.com/esvm/dcm-middleware/distribution/models"
	ifs "github.com/esvm/dcm-middleware/infrastructure"
	log "github.com/sirupsen/logrus"
)

type Invoker struct {
	srh *ifs.ServerRequestHandler
}

func NewInvoker(port int) *Invoker {
	return &Invoker{srh: ifs.NewServerRequestHandler(port)}
}

func (invoker *Invoker) Invoke(callback func(string, models.Request, chan models.Response) error) error {
	return invoker.srh.Listen(func(conn net.Conn, host string, bytes []byte) error {
		//ack(conn)
		var request models.Request
		err := ifs.Unmarshall(bytes, &request)
		if err != nil {
			return err
		}
		log.Infof("Received request: %s", request.Operation())

		chRes := make(chan models.Response)
		go func(err *error) {
			*err = callback(host, request, chRes)
			if *err != nil {
				close(chRes)
			}
		}(&err)

		for response := range chRes {
			bytes, err := ifs.Marshall(&response)
			if err != nil {
				log.Debug("porra: ", err)
			}

			var re models.Response
			err = ifs.Unmarshall(bytes[4:], &re)
			ifs.Send(conn, bytes)
		}

		log.Debug("encerrou")
		log.Debug(err)
		return err
	})
}

func ack(conn net.Conn) {
	bytes, _ := ifs.Ack()
	ifs.Send(conn, bytes)
}
