package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang/glog"
	"go.etcd.io/etcd/clientv3"
)

func toString(obj interface{}) (string, error) {
	bytes, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// 写入ETCD
func PUT(cli *clientv3.Client, key, val string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	resp, err := cli.Put(ctx, key, val)
	cancel()
	if err != nil {
		return "", err
	}
	return toString(resp)
}

// 查询ETCD
func GET(cli *clientv3.Client, key string, opts ...clientv3.OpOption) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	resp, err := cli.Get(ctx, key, opts...)
	cancel()
	if err != nil {
		return "", err
	}
	return toString(resp)
}

// 删除ETCD
func DELETE(cli *clientv3.Client, key string, opts ...clientv3.OpOption) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	resp, err := cli.Delete(ctx, key, opts...)
	cancel()
	if err != nil {
		return "", err
	}
	return toString(resp)
}

func main() {

	config := clientv3.Config{
		Endpoints:   []string{"10.0.6.247:2379", "10.0.6.247:22379", "10.0.6.247:32379"},
		DialTimeout: 5 * time.Second,
	}

	cli, err := clientv3.New(config)
	if err != nil {
		glog.Fatal(err.Error())
	}

	defer func() {
		if err := cli.Close(); err != nil {
			glog.Error(err.Error())
		}
	}()

	if v, err := PUT(cli, "sample_key", "{'ds1':'dss'}"); err != nil {
		glog.Errorf(err.Error())
	} else {
		fmt.Printf("PUT RESULT: %s\n", v)
	}

	if v, err := GET(cli, "sample_key"); err != nil {
		glog.Errorf(err.Error())
	} else {
		fmt.Printf("GET RESULT: %s\n", v)
	}
	if v, err := DELETE(cli, "sample1_key2"); err != nil {
		glog.Errorf(err.Error())
	} else {
		fmt.Printf("DELETE RESULT: %s\n", v)
	}
}
