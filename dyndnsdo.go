package main

import (
	"flag"
	"fmt"
	"github.com/coreos/go-systemd/activation"
	"github.com/digitalocean/godo"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"os"
)

var (
	certFile = flag.String("cert", "", "certificate file")
	keyFile  = flag.String("key", "", "key file")
)

func update(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			log.Printf("%v: %v\n", name, h)
		}
	}

	// Obtain data from the query string which we'll need for our API update
	domain := req.FormValue("domain")
	ip := req.FormValue("ip")

	ctx := context.TODO()
	client := godo.NewFromToken(os.Getenv("DO_API_TOKEN"))

	// Specifically we will be updating an A record, for IPv4 naked domain address.
	records, _, _ := client.Domains.RecordsByTypeAndName(ctx, domain, "A", domain, &godo.ListOptions{})

	for _, record := range records {
		log.Printf("DigialOcean: %s: %s\n", domain, record.Data)
		log.Printf("Router:      %s: %s\n", domain, ip)
		if record.Data != ip {
			log.Printf("Updating:    %s -> %s\n", record.Data, ip)
			_, _, err := client.Domains.EditRecord(ctx, domain, record.ID, &godo.DomainRecordEditRequest{
				Type:     record.Type,
				Name:     record.Name,
				Data:     ip,
				Priority: record.Priority,
				Port:     record.Port,
				TTL:      record.TTL,
				Weight:   record.Weight,
				Flags:    record.Flags,
				Tag:      record.Tag,
			})
			if err != nil {
				fmt.Fprintln(w, "<ErrCount>1</ErrCount>")
			} else {
				fmt.Fprintln(w, "<ErrCount>0</ErrCount>")
			}
		} else {
			fmt.Fprintln(w, "<ErrCount>0</ErrCount>")
		}
	}
}

func main() {
	listeners, err := activation.Listeners()
	if err != nil {
		panic(err)
	}
	if len(listeners) != 1 {
		panic("Unexpected number of socket activation fds")
	}
	flag.Parse()
	log.Printf("cert: %s\n", *certFile)
	log.Printf("key:  %s\n", *keyFile)
	http.HandleFunc("/update", update)
	http.ServeTLS(listeners[0], nil, *certFile, *keyFile)
}
