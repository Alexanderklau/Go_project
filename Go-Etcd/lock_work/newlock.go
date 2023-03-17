package mian

import (
	"context"
	"fmt"
	"time"

	//"reflect"

	"go.etcd.io/etcd/clientv3"
)

var (
	lease                  clientv3.Lease
	ctx                    context.Context
	cancelFunc             context.CancelFunc
	leaseId                clientv3.LeaseID
	leaseGrantResponse     *clientv3.LeaseGrantResponse
	leaseKeepAliveChan     <-chan *clientv3.LeaseKeepAliveResponse
	leaseKeepAliveResponse *clientv3.LeaseKeepAliveResponse
	txn                    clientv3.Txn
	txnResponse            *clientv3.TxnResponse
	kv                     clientv3.KV
)

type ETCD struct {
	client *clientv3.Client
	cfg    clientv3.Config
	err    error
}

func New(endpoints ...string) (*ETCD, error) {
	cfg := clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: time.Second * 5,
	}

	client, err := clientv3.New(cfg)
	if err != nil {
		fmt.Println("连接ETCD失败")
		return nil, err
	}

	etcd := &ETCD{
		cfg:    cfg,
		client: client,
	}

	fmt.Println("连接ETCD成功")
	return etcd, nil
}

func (etcd *ETCD) Newleases_lock(ip string) error {
	lease := clientv3.NewLease(etcd.client)
	leaseGrantResponse, err := lease.Grant(context.TODO(), 5)
	if err != nil {
		fmt.Println(err)
		return err
	}
	leaseId := leaseGrantResponse.ID
	ctx, cancelFunc := context.WithCancel(context.TODO())
	defer cancelFunc()
	defer lease.Revoke(context.TODO(), leaseId)
	leaseKeepAliveChan, err := lease.KeepAlive(ctx, leaseId)
	if err != nil {
		fmt.Println(err)
		return err
	}

	kv := clientv3.NewKV(etcd.client)
	txn := kv.Txn(context.TODO())
	txn.If(clientv3.Compare(clientv3.CreateRevision("/dev/lock"), "=", 0)).Then(
		clientv3.OpPut("/dev/lock", ip, clientv3.WithLease(leaseId))).Else(
		clientv3.OpGet("/dev/lock"))
	txnResponse, err := txn.Commit()
	if err != nil {
		fmt.Println(err)
		return err
	}
	if txnResponse.Succeeded {
		fmt.Println("抢到锁了")
		fmt.Println("选定主节点 %s", ip)
		for {
			select {
			case leaseKeepAliveResponse = <-leaseKeepAliveChan:
				if leaseKeepAliveResponse != nil {
					fmt.Println("续租成功,leaseID :", leaseKeepAliveResponse.ID)
				} else {
					fmt.Println("续租失败")
				}

			}
		}
	} else {
		fmt.Println("没抢到锁", txnResponse.Responses[0].GetResponseRange().Kvs[0].Value)
		fmt.Println("继续抢")
		time.Sleep(time.Second * 1)
	}
	return nil
}

func main() {
	etcd, err := New("xxxxxxxxx:2379")
	if err != nil {
		fmt.Println(err)
	}
	for {
		etcd.Newleases_lock("node1")
	}
}
