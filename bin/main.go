package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	colorcandy "github.com/improvemedia/colorcandy.git"
)

const (
	defaultThriftAddr = ":3033"
)

var (
	url        = flag.String("url", "", "Path or URL to image (gif, jpeg, png)")
	thrift     = flag.String("thrift", defaultThriftAddr, "Thrift address to listen")
	configFile = flag.String("config", "etc/candy.json", "Path to candy config")
)

func main() {
	flag.Parse()
	var config colorcandy.Config

	r, err := os.Open(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	d := json.NewDecoder(r)
	d.Decode(&config)

	c := colorcandy.New(config)

	if *url != "" {
		fmt.Printf("URL: %s\n", *url)

		res, err := c.Candify(*url, []string{})
		if err != nil {
			log.Fatal(err)
		}
		enc := json.NewEncoder(os.Stdout)
		fmt.Println(enc.Encode(res))

		return
	}

	errc := make(chan error)
	go func() {
		errc <- interrupt()
	}()

	s := colorcandy.NewThriftService(c)
	go func() {
		errc <- s.Start(*thrift)
	}()

	log.Fatal(<-errc)
}

func interrupt() error {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	return fmt.Errorf("%s", <-c)
}
