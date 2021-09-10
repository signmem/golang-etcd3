package main

import (
	"flag"
	"fmt"
	"gitlab.tools.vipshop.com/terry.zeng/golang-etcd3/etcdfucn"
	"gitlab.tools.vipshop.com/terry.zeng/golang-etcd3/g"
	"log"
)

func main() {

	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "sow version")
	flag.Parse()

	if *version {
		fmt.Sprintf(g.Version)
	}

	g.ParseConfig(*cfg)

	if g.Config().Debug {
		g.InitLog("debug")
	} else {
		g.InitLog("info")
	}

	filePath := g.Config().EtcdSetting.LoadFile

	allHostList, err := etcdfucn.ReadHost(filePath)

	if err != nil {
		log.Println("read file error.")
	}

	cli := etcdfucn.EtcdControl()

	method := g.Config().EtcdSetting.EtcdMethod

	if method == "write" {
		for _, hostname := range allHostList {
			etcdfucn.EtcdWrite(cli, hostname)
		}
		defer cli.Close()
	}

	if method == "read" {
		for _, hostname := range allHostList {
			etcdfucn.EtcdSingalRead(cli, hostname)
		}
		defer cli.Close()
	}

	if method == "range" {
		etcdfucn.EtcdRangeRead(cli)
	}

}