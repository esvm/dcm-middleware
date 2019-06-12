package broker

import (
	"errors"
	"os"
	"strconv"

	dst "github.com/dwsb/projetomiddleware/distribution"
	"github.com/dwsb/projetomiddleware/distribution/models"
	log "github.com/sirupsen/logrus"
)

type Exchange struct {
	invoker     *dst.Invoker
	topicMgmt   *TopicMgmt
	sessionMgmt *SessionMgmt
}

const consumeBatch = 10
const defaultPort = 8426

func NewExchange(topicMgmt *TopicMgmt, sessionMgmt *SessionMgmt) *Exchange {
	port, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		log.Warn("Environment not setted. Using default port: ", defaultPort)
		port = defaultPort
	}

	return &Exchange{
		invoker:     dst.NewInvoker(port),
		topicMgmt:   topicMgmt,
		sessionMgmt: sessionMgmt,
	}
}

func (exchange *Exchange) Exchange() {
	chCon := map[string][]chan models.Response{}

	exchange.invoker.Invoke(func(host string, request models.Request, chRes chan models.Response) error {
		op := request.Operation()

		switch op {
		case models.CreateTopicOperation:
			response := exchange.topicMgmt.CreateTopic(request.(*models.CreateTopicRequest))
			chRes <- response
			close(chRes)
		case models.ConsumeOperation:
			consumeRequest := request.(*models.ConsumeRequest)

			stop := true
			EOF := &models.ErrEOF{}
			for stop {
				resList := exchange.topicMgmt.Consume(consumeRequest, consumeBatch)
				for _, res := range resList {
					switch res.Error() {
					case "":
						chRes <- res
					case EOF.Error():
						stop = false
						break
					default:
						chRes <- &models.ConsumeResponse{
							Err:   res.Error(),
							Alive: false,
						}
						close(chRes)
						return nil
					}
				}
			}

			chCon[consumeRequest.TopicID] = append(chCon[consumeRequest.TopicID], chRes)
		case models.PublishOperation:
			publishRequest := request.(*models.PublishRequest)
			res := exchange.topicMgmt.Publish(publishRequest)

			chRes <- res
			close(chRes)

			response := &models.ConsumeResponse{
				MessageResult: publishRequest.Metrics,
				Err:           res.Error(),
				Alive:         res.Error() == "",
			}

			for _, ch := range chCon[publishRequest.TopicID] {
				ch <- response
			}
		default:
			return errors.New("Invalid request")
		}

		return nil
	})
}
