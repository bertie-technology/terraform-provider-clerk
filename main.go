package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// version can be set via ldflags at build time
var version string = "dev"

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/bertie-technology/clerk",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), New(version), opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
