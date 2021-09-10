package etcdfucn

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"gitlab.tools.vipshop.com/terry.zeng/golang-etcd3/g"
	"go.etcd.io/etcd/clientv3"
	"strconv"

	"golang.org/x/net/context"
	"io/ioutil"
	"log"
	"time"
)

func EtcdControl() (cli *clientv3.Client) {


	endPoints := g.Config().EtcdConfig.Host

	etcdCaFile := g.Config().EtcdSSL.CaFile
	etcdCertFile := g.Config().EtcdSSL.CertFile
	etcdCertKeyFile := g.Config().EtcdSSL.CertKeyFile
	// method := g.Config().EtcdSetting.EtcdMethod

	cert, err := tls.LoadX509KeyPair(etcdCertFile, etcdCertKeyFile)
	if err != nil {
		log.Println("load etcd cert file faile.")
		return
	}

	caData, err := ioutil.ReadFile(etcdCaFile)
	if err != nil {
		log.Println("load etcd ca file faile.")
		return
	}

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caData)

	_tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      pool,
	}

	cfg := clientv3.Config{
		Endpoints: endPoints,
		TLS:       _tlsConfig,
		DialTimeout:  5 * time.Second,
	}

	cli, err = clientv3.New(cfg)
	if err != nil {
		log.Print (" error 001: ", err)
	}

	return cli

}

func EtcdWrite(cli *clientv3.Client, hostname string) {

	etcdpath  := g.Config().EtcdSetting.EtcdPath
	debug := g.Config().Debug

	//requestTimeOut :=  5 * time.Second
	//ctx, cancel := context.WithTimeout(context.Background(), requestTimeOut)

	etcdKey, etcdValue := etcdpath + hostname, strconv.FormatInt(time.Now().Unix(), 10)

	if debug == true {
		fmt.Println( "key: ",  etcdKey , " value: ", etcdValue )
	}

	_, err := cli.Put(context.TODO(), etcdKey, etcdValue)
	// cancel()
	if err != nil {
		log.Println("Put failed. ", err)
		return
	}

}

func EtcdSingalRead(cli *clientv3.Client, hostname string) {

	etcdpath  := g.Config().EtcdSetting.EtcdPath
	// debug := g.Config().Debug

	requestTimeOut :=  5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeOut)

	etcdKey := etcdpath + hostname

	resp, err := cli.Get(ctx, etcdKey)
	cancel()
	if err != nil {
		log.Println("Get failed. ", err)
		return
	}

	for _, kv := range resp.Kvs {
		log.Printf("Get: {%s:%s} \n", kv.Key, kv.Value)
	}

}

func EtcdRangeRead(cli *clientv3.Client) {
	defer cli.Close()
	etcdpath  := g.Config().EtcdSetting.EtcdPath

	requestTimeOut :=  50 * time.Second
	ctx, _ := context.WithTimeout(context.Background(), requestTimeOut)

	resp, err := cli.Get(ctx, etcdpath, clientv3.WithPrefix())
	if err != nil {
		log.Println("Get failed. ", err)
		return
	}

	for _, kv := range resp.Kvs {
		log.Printf("Get: {%s:%s} \n", kv.Key, kv.Value)
	}
}