package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bobgo0912/b0b-common/pkg/config"
	"github.com/bobgo0912/b0b-common/pkg/constant"
	"github.com/bobgo0912/b0b-common/pkg/log"
	"github.com/bobgo0912/b0b-common/pkg/server/balancer"
	"github.com/bobgo0912/b0b-common/pkg/server/common"
	"github.com/bobgo0912/b0b-common/pkg/server/model"
	"github.com/bobgo0912/b0b-common/pkg/util"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	clientv3 "go.etcd.io/etcd/client/v3"
	"strings"
	"sync"
	"time"
)

type Servers struct {
	Ctx        context.Context
	EtcdClient *clientv3.Client
	Lock       sync.Mutex
	Services   map[string]*model.EtcdReg
	Balancer   *balancer.P2cBalancer
	Error      error
}

type Discover struct {
	EtcdClient *clientv3.Client
	Server     Server
	ServerCtx  context.Context
}

func NewDiscover(ctx context.Context, etcdClient *clientv3.Client, server Server) *Discover {
	return &Discover{
		Server:     server,
		EtcdClient: etcdClient,
		ServerCtx:  ctx,
	}
}

func (d *Discover) Start() {
	key := fmt.Sprintf(constant.EtcdServers, config.Cfg.ENV, d.Server.Type, config.Cfg.ServiceName, config.Cfg.NodeId)
	lease := clientv3.NewLease(d.EtcdClient)
	var curLeaseId clientv3.LeaseID = 0
	for {
		select {
		case <-d.ServerCtx.Done():
			log.Infof("discover %s server stop ", d.Server.Type)
			return
		default:
			if curLeaseId == 0 {
				leaseResp, err := lease.Grant(d.ServerCtx, 10)
				if err != nil {
					log.Error("grant fail err=", err)
					time.Sleep(10 * time.Second)
					continue
				}
				toJson, err := d.Server.ToJson()
				if err != nil {
					log.Error("ToJson fail err=", err)
					time.Sleep(10 * time.Second)
					continue
				}
				if _, err := d.EtcdClient.Put(d.ServerCtx, key, toJson, clientv3.WithLease(leaseResp.ID)); err != nil {
					log.Error("etcd Put fail err=", err)
					time.Sleep(10 * time.Second)
					continue
				}
				curLeaseId = leaseResp.ID
			} else {
				if _, err := lease.KeepAliveOnce(d.ServerCtx, curLeaseId); err == rpctypes.ErrLeaseNotFound {
					curLeaseId = 0
					continue
				}
			}
			time.Sleep(2 * time.Second)
		}

	}
}

func NewServices(ctx context.Context, etcdClient *clientv3.Client) *Servers {
	return &Servers{Ctx: ctx, Lock: sync.Mutex{},
		EtcdClient: etcdClient,
		Services:   make(map[string]*model.EtcdReg, 0),
		Balancer:   balancer.NewP2cBalancer(),
	}
}

func (s *Servers) Start() {
	go func() {
		err := s.initService()
		if err != nil {
			s.Error = err
		}
		s.watchServiceUpdate()
	}()
}
func (s *Servers) initService() error {
	key := fmt.Sprintf(constant.EtcdServersPre, config.Cfg.ENV)
	rangeResp, err := s.EtcdClient.Get(s.Ctx, key, clientv3.WithPrefix())
	if err != nil {
		log.Error(err)
		return err
	}
	s.Lock.Lock()
	for _, kv := range rangeResp.Kvs {
		k := strings.TrimPrefix(string(kv.Key), key)
		split := strings.Split(k, "/")
		if len(split) != 3 {
			continue
		}
		mKey := util.GetStrings(split[0], split[1], split[2])
		var sd *model.EtcdReg
		err := json.Unmarshal(kv.Value, &sd)
		if err != nil {
			log.Warn("InitService etcd Unmarshal %s fail err=%v", string(kv.Value), err)
			continue
		}
		rKey := util.GetStrings(split[0], split[1])
		sd.Key = mKey
		s.Services[mKey] = sd
		s.Balancer.Add(rKey, sd)
	}
	s.Lock.Unlock()
	return nil
}

func (s *Servers) watchServiceUpdate() {
	key := fmt.Sprintf(constant.EtcdServersPre, config.Cfg.ENV)
	watchChan := s.EtcdClient.Watch(s.Ctx, key, clientv3.WithPrefix())
	for watchResp := range watchChan {
		for _, event := range watchResp.Events {
			k := strings.TrimPrefix(string(event.Kv.Key), key)
			split := strings.Split(k, "/")
			if len(split) != 3 {
				continue
			}
			mKey := util.GetStrings(split[0], split[1], split[2])
			s.Lock.Lock()
			switch event.Type {
			case mvccpb.PUT:
				var sd *model.EtcdReg
				err := json.Unmarshal(event.Kv.Value, &sd)
				if err != nil {
					log.Warn("InitService etcd Unmarshal %s fail err=%v", string(event.Kv.Value), err)
					continue
				}
				rKey := util.GetStrings(split[0], split[1])
				sd.Key = mKey
				s.Services[mKey] = sd
				s.Balancer.Add(rKey, sd)
				break
			case mvccpb.DELETE: //DELETE事件，目录中有key被删掉(Lease过期，key 也会被删掉)
				delete(s.Services, mKey)
				s.Balancer.Remove(mKey)
				break
			}
			s.Lock.Unlock()
		}
	}
}

func GetNode(serviceName string, serverType common.Type) *model.EtcdReg {
	if MainServers == nil {
		return nil
	}
	key := util.GetStrings(string(serverType), serviceName)
	for ke, reg := range MainServers.DiscoverServers.Services {
		if strings.Contains(ke, key) {
			return reg
		}
	}
	return nil
}
func GetNodeN(serviceName string, serverType common.Type) (*model.EtcdReg, error) {
	if MainServers == nil {
		return nil, errors.New("main server not init")
	}
	key := util.GetStrings(string(serverType), serviceName)
	balance, err := MainServers.DiscoverServers.Balancer.Balance(key)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func GetRpcNodeAddress(serviceName string, serverType common.Type) string {
	node := GetNode(serviceName, serverType)
	if node == nil {
		return ""
	}
	return fmt.Sprintf("%s:%d", node.Host, node.Port)
}

func GetNodeList(serviceName string, serverType common.Type) []*model.EtcdReg {
	if MainServers == nil {
		return nil
	}
	regs := make([]*model.EtcdReg, 0)
	key := util.GetStrings(string(serverType), serviceName)
	for ke, reg := range MainServers.DiscoverServers.Services {
		if strings.Contains(ke, key) {
			regs = append(regs, reg)
		}
	}
	return regs
}
