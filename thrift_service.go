package colorcandy

import (
	"log"
	"time"

	"git.apache.org/thrift.git/lib/go/thrift"

	"github.com/improvemedia/colorcandy.git/candy"
)

type ThriftService struct {
	*ColorCandy
}

func NewThriftService(c *ColorCandy) *ThriftService {
	return &ThriftService{c}
}

func (s *ThriftService) Start(addr string) {
	log.Printf("[thrift] starting service at %s", addr)

	transport, err := thrift.NewTServerSocket(addr)
	if err != nil {
		log.Fatal(err)
	}

	server := thrift.NewTSimpleServer4(
		candy.NewCandyProcessor(s),
		transport,
		thrift.NewTTransportFactory(),
		thrift.NewTBinaryProtocolFactoryDefault(),
	)

	server.Serve()
}

func (s *ThriftService) Candify(url string) (res map[int64]*candy.ColorMeta, err error) {
	log.Printf("processing: %s", url)

	start := time.Now()
	res, err = s.ColorCandy.Candify(url)
	if err != nil {
		log.Printf("Error: %s", err)
	}
	log.Printf("Finished in: %s", time.Since(start))

	return
}
