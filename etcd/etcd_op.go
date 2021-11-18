package etcd

import (
	"cmit.com/crd/domain-config/api/v1alpha1"
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/go-logr/logr"
	"go.etcd.io/etcd/clientv3"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

var (
	cli *clientv3.Client
)

type CliEmtpyError struct {
	cli *clientv3.Client
}

func (e *CliEmtpyError) Error() string {
	return "client not conn!"
}

type CliConnError struct {
	cli   *clientv3.Client
	state string
}

func (e *CliConnError) Error() string {
	return "client conn err: " + e.state
}

func ConnCheck() (bool, error) {
	if cli == nil {
		return false, &CliEmtpyError{cli}
	}
	state := cli.ActiveConnection().GetState().String()
	if state != "READY" {
		return false, &CliConnError{cli, state}
	}
	log.Println("conn check ok!")
	return true, nil
}
func InitConn(addr []string, timeout time.Duration, configDir string) (err error) {
	var etcdCertPath = configDir + "healthcheck-client.crt"
	var etcdCertKeyPath = configDir + "healthcheck-client.key"
	var etcdCaPath = configDir + "ca.crt"

	//load cert
	cert, err := tls.LoadX509KeyPair(etcdCertPath, etcdCertKeyPath)
	if err != nil {
		log.Println("etcd load cert failed", err.Error())
		return
	}

	//load root ca
	caData, err := ioutil.ReadFile(etcdCaPath)
	if err != nil {
		log.Println("etcd load root ca failed", err.Error())
		return
	}

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caData)

	_tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      pool,
	}

	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   addr,
		TLS:         _tlsConfig,
		DialTimeout: timeout,
	})
	if err != nil {
		log.Println("etcd connect to etcd failed", err.Error())
		return
	}
	return
}

func PutValue(key string, value string) {
	//设置超时及设置值
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := cli.Put(ctx, key, value)
	cancelFunc()
	if err != nil {
		log.Println("etcd.Put", err.Error())
	}
}

func GetValue(key string) (value string) {
	//取值
	result := "no value"
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	res, err := cli.Get(ctx, key)
	cancelFunc()
	if err != nil {
		log.Println("etcd.Get", err.Error())
		return "get value err"
	}

	for _, v := range res.Kvs {
		result = string(v.Value)
	}
	return result
}

func DelValue(key string) {
	//删除
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := cli.Delete(ctx, key)
	cancelFunc()
	if err != nil {
		log.Println("etcd.Del", err.Error())
		return
	}
}

func Reverse(s string) string {
	sa := strings.Split(s, ".")
	var res []string
	for i := len(sa) - 1; i >= 0; i-- {
		res = append(res, sa[i])
	}
	return strings.Join(res, "/")
}

func Option(ipHosts []v1alpha1.IPHosts, reqLogger logr.Logger, ctx context.Context) {
	_, err := ConnCheck()
	if err != nil {
		log.Println("etcd option", err.Error())
		return
	}
	ips := ipHosts
	reqLogger.Info("etcd  ----------操作开始-----------")

	for i := 0; i < len(ips); i++ {
		opt := ips[i].Option
		if opt == "put" {
			confHost := "/skydns/" + Reverse(ips[i].ConfHost) + "/"
			confIp := "{\"host\":\"" + ips[i].ConfIp + "\"}"
			PutValue(confHost, confIp)
			v := GetValue(confHost)
			reqLogger.Info("etcd put", "key: ", confHost, "value: ", v)
		} else if opt == "del" {
			confHost := "/skydns/" + Reverse(ips[i].ConfHost) + "/"
			v := GetValue(confHost)
			reqLogger.Info("etcd del", "key: ", confHost, "value: ", v)
			DelValue(confHost)
		} else {
			confHost := "/skydns/" + Reverse(ips[i].ConfHost) + "/"
			v := GetValue(confHost)
			reqLogger.Info("etcd get", "key: ", confHost, "value:", v)
		}
	}
	reqLogger.Info("etcd ", "操作完成！  ipHosts: ", ips)
}
