package main

import (
	"log"

	dnstap "github.com/dnstap/golang-dnstap"

	"github.com/golang/protobuf/proto"
)

type inverseDnsFiller struct {
	data    chan []byte
	inverse *inverseDnsMap
}

func newInverseDnsFiller(ipss *inverseDnsMap) *inverseDnsFiller {
	return &inverseDnsFiller{
		data:    make(chan []byte, 8),
		inverse: ipss,
	}
}

func (this *inverseDnsFiller) RunOutputLoop() {

	dt := &dnstap.Dnstap{}

	for b := range this.data {

		if err := proto.Unmarshal(b, dt); err != nil {
			log.Fatalf("proto.Unmarshal() failed: %s\n", err)
			break
		}

		if dt.Message.GetType() == dnstap.Message_FORWARDER_RESPONSE {
			this.inverse.Add(dt.Message.ResponseMessage)
		}

	}
}

func (this *inverseDnsFiller) GetOutputChannel() chan []byte {
	return this.data
}

func (_ *inverseDnsFiller) Close() {
}
