package distribution

import (
	"github.com/esvm/dcm-middleware/distribution/models"
	ifs "github.com/esvm/dcm-middleware/infrastructure"
	log "github.com/sirupsen/logrus"
)

type Requestor struct {
	client *ifs.ClientRequestHandler
}

func NewRequestor(host string, port int) (*Requestor, error) {
	req := new(Requestor)
	req.client = ifs.NewClientRequestHandler(host, port)
	err := req.client.Connect()
	return req, err
}

func (req *Requestor) Close() error {
	return req.client.Close()
}

func (req *Requestor) Invoke(request models.Request) (chan models.Response, error) {
	bytes, err := ifs.Marshall(&request)
	if err != nil {
		log.Debugf("infrastructure.Marshall: %s", err.Error())
		return nil, err
	}

	ch := make(chan models.Response)

	err = req.client.Send(bytes)
	if err != nil {
		log.Debugf("clientRequest.Send: %s", err.Error())
		return nil, err
	}

	go func() {
		for {
			responseBytes, err := req.client.Receive()
			if err != nil {
				log.Debugf("clientRequest.Receive: %s", err.Error())
				break
			}

			var response models.Response
			err = ifs.Unmarshall(responseBytes, &response)
			log.Debug(response)
			if err != nil {
				log.Debugf("infrastructure.Unmarshall: %s", err.Error())
				break
			}

			ch <- response
			if !response.KeepAlive() {
				close(ch)
				break
			}
		}
	}()

	return ch, nil
}
