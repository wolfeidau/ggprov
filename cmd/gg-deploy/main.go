package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/wolfeidau/ggprov"
)

var (
	hostname = kingpin.Arg("hostname", "Name of hostname of the device to deploy.").Required().String()
	ggName   = kingpin.Arg("name", "Name of greengrass core.").Required().String()

	username = kingpin.Flag("username", "Username used to access board.").Default("linaro").String()
	port     = kingpin.Flag("port", "Port used to access board.").Default("22").String()
)

func main() {
	kingpin.Parse()

	keyPath := os.Getenv("HOME") + "/.ssh/id_rsa"

	ss, err := ggprov.NewSSHSession(*hostname, *port, *username, keyPath)
	if err != nil {
		log.Fatalf("%+v\n", err)
	}

	defer ggprov.DoClose(ss.Client)

	sourcePath := fmt.Sprintf("%s.yaml", *ggName)

	log.Println("Transferring config file for", sourcePath)

	err = ss.CopyPath(sourcePath, sourcePath)
	if err != nil {
		log.Fatalf("%+v\n", err)
	}

	greengrassPath := "greengrass-linux-aarch64-1.0.0.tar.gz"

	log.Println("Transferring greengrass file for", greengrassPath)

	err = ss.CopyPath(greengrassPath, greengrassPath)
	if err != nil {
		log.Fatalf("%+v\n", err)
	}

}
