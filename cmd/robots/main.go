package main

import (
	"flag"
	"fmt"
	"mrdmitry/blog/pkg/monke"
	"os"
	"text/template"
)

func main() {
	protocol := "https"
	flag.Func("protocol", "webserver protocol: [http, https] (default \"https\")", func(value string) error {
		v, err := monke.ProtocolValidator(value)
		protocol = v
		return err
	})

	hostname := flag.String("hostname", "aidenhale.dev", "hostname")

	flag.Parse()

	urlPrefix := monke.SanitizeUrl(fmt.Sprintf("%s://%s/", protocol, *hostname))

	tmpl, err := template.ParseFiles("./web/templates/robots.txt")
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(os.Stdout, struct {
		UrlPrefix string
	}{
		UrlPrefix: urlPrefix,
	})

	if err != nil {
		panic(err)
	}
}
