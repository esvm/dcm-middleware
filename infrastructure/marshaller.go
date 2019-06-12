package infrastructure

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"

	log "github.com/sirupsen/logrus"
)

func Ack() ([]byte, error) {
	return Marshall("ok")
}

func Marshall(value interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(value)
	if err != nil {
		return nil, err
	}

	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, uint32(buf.Len()))
	log.Debugf("Encoded bytes: %d", buf.Len())

	return append(bs, buf.Bytes()...), nil
}

func Unmarshall(value []byte, ref interface{}) error {
	buf := bytes.NewBuffer(value)
	dec := gob.NewDecoder(buf)
	return dec.Decode(ref)
}
