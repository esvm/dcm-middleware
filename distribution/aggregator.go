package distribution

import (
	"errors"

	"github.com/esvm/dcm-middleware/distribution/models"
	"github.com/montanaflynn/stats"
)

type Aggregator struct {
	topicsMessage map[string][]*models.Message
	bufferSize    int
}

func NewAggregator(bufferSize int) *Aggregator {
	return &Aggregator{
		topicsMessage: map[string][]*models.Message{},
		bufferSize:    bufferSize,
	}
}

func (agg *Aggregator) AppendMetric(topicID string, message *models.Message) (bool, error) {
	size := len(agg.topicsMessage[topicID])
	if size >= agg.bufferSize {
		return true, errors.New("Buffer is fully")
	}

	agg.topicsMessage[topicID] = append(agg.topicsMessage[topicID], message)
	return size+1 == agg.bufferSize, nil
}

func (agg *Aggregator) Aggregate(topicID string) *models.Metric {
	payloads := []float64{}
	for _, message := range agg.topicsMessage[topicID] {
		payloads = append(payloads, message.Payload)
	}

	// calcula as métricas (média, mediana, variança, desvio padrão, moda)
	average, _ := stats.Mean(payloads)
	median, _ := stats.Median(payloads)
	variance, _ := stats.Variance(payloads)
	standardDeviation, _ := stats.StandardDeviation(payloads)
	mode, _ := stats.Mode(payloads)

	// se não tem moda, setamos 0
	var modeValue float64
	if len(mode) != 0 {
		modeValue = mode[0]
	}

	// zera o buffer
	agg.topicsMessage[topicID] = []*models.Message{}

	return &models.Metric{
		Average:           average,
		Median:            median,
		Variance:          variance,
		StandardDeviation: standardDeviation,
		Mode:              modeValue,
	}
}
