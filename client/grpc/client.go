package grpc

import (
	"context"
	"github.com/Gitforxuyang/walle/config"
	"github.com/Gitforxuyang/walle/util/utils"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"sync"
)

var (
	invokers map[string]*Proxy   = make(map[string]*Proxy)
	services map[string]*Service = make(map[string]*Service)
	initOnce sync.Once
	etcd     *clientv3.Client
)

type Service struct {
	Package string            `json:"package"`
	Name    string            `json:"name"`
	Methods map[string]Method `json:"methods"`
}
type Method struct {
	Req  Message `json:"req"`
	Resp Message `json:"resp"`
}
type Message map[string]string

const (
	ETCD_WALLE_SERVICE_PREFIX = "/eva/walle/service/"
)

func Init() {
	initOnce.Do(func() {
		etcd := config.GetEtcdClient()
		res, err := etcd.Get(context.TODO(), ETCD_WALLE_SERVICE_PREFIX, clientv3.WithPrefix())
		utils.Must(err)
		for _, v := range res.Kvs {
			service := Service{}
			err := utils.JsonToStruct(string(v.Value), &service)
			utils.Must(err)
			proxy := NewProxy(service.Name)
			invokers[service.Name] = proxy
			services[service.Name] = &service
		}
		go watch()
	})
}

func watch() {
	for {
		etcd := config.GetEtcdClient()
		chs := etcd.Watch(context.TODO(), ETCD_WALLE_SERVICE_PREFIX, clientv3.WithPrefix())
		for ch := range chs {
			for _, event := range ch.Events {
				switch event.Type {
				case mvccpb.PUT:
					service := Service{}
					err := utils.JsonToStruct(string(event.Kv.Value), &service)
					utils.Must(err)
					if invokers[service.Name] == nil {
						proxy := NewProxy(service.Name)
						invokers[service.Name] = proxy
					}
					services[service.Name] = &service
				case mvccpb.DELETE:

				}
			}
		}

	}
}