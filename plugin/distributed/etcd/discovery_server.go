package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/fankane/go-utils/str"
	"github.com/fankane/go-utils/utime"
	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	registeredServers    = make(map[string]map[string]*ServerInfo)
	defaultCheckInterval = int64(3) //
)

type serverHandle struct {
	l          *sync.RWMutex
	e          *etcd
	serverChan chan map[string]*ServerInfo
}

func (h *serverHandle) PutHandle(key, value []byte, version int64) {
	if !strings.HasPrefix(string(key), h.e.conf.SInfo.ServerName) {
		return
	}
	sInfo := &ServerInfo{}
	if err := json.Unmarshal(value, sInfo); err != nil {
		log.Println("unmarshal err:", err)
		return
	}
	rsMap := make(map[string]*ServerInfo)
	rsMap[h.e.serverIDFromKey(string(key))] = sInfo
	h.l.Lock()
	registeredServers[h.e.confName] = rsMap
	h.l.Unlock()
	if h.serverChan != nil {
		h.serverChan <- rsMap
	}
}
func (h *serverHandle) DelHandle(key []byte) {
	if !strings.HasPrefix(string(key), h.e.conf.SInfo.ServerName) {
		return
	}
	h.l.Lock()
	delete(registeredServers[h.e.confName], h.e.serverIDFromKey(string(key)))
	h.l.Unlock()
	if h.serverChan != nil {
		h.serverChan <- registeredServers[h.e.confName]
	}
}

func (e *etcd) RegisterServer() error {
	if !e.conf.OpenDiscovery {
		return nil
	}
	if e.conf.SInfo.CheckInterval <= 0 {
		e.conf.SInfo.CheckInterval = defaultCheckInterval
	}
	lease := clientv3.NewLease(e.cli)
	var leaseID clientv3.LeaseID
	registerFunc := func() error {
		if e.conf.SInfo.stop {
			return fmt.Errorf("server stoped") //返回错误，退出定期续约
		}
		ctx := context.Background()
		if leaseID == 0 { //第一次创建 服务信息
			leaseResp, err := lease.Grant(ctx, e.conf.SInfo.CheckInterval+1) //租约时间比检测时间多1秒
			if err != nil {
				return fmt.Errorf("grant err:%s", err)
			}
			if _, err = e.Put(ctx, e.getServerKey(), str.ToJSON(e.conf.SInfo),
				clientv3.WithLease(leaseResp.ID)); err != nil {
				return fmt.Errorf("register server:%s err:%s", e.conf.SInfo.ServerID, err)
			}
			leaseID = leaseResp.ID
			rsMap := make(map[string]*ServerInfo)
			rsMap[e.conf.SInfo.ServerID] = e.conf.SInfo
			registeredServers[e.confName] = rsMap
			return nil
		}
		// 续约租约，如果租约已经过期将curLeaseId复位到0重新走创建租约的逻辑\
		_, err := lease.KeepAliveOnce(ctx, leaseID)
		if err != nil {
			if err == rpctypes.ErrLeaseNotFound {
				leaseID = 0
				return nil
			}
			log.Println("keep alive once err:", err)
			return err
		}
		return nil
	}
	if err := registerFunc(); err != nil {
		return err
	}
	go func() {
		e.WatchServers(make(chan map[string]*ServerInfo))
	}()
	go func() {
		if err := utime.TickerDo(time.Second*time.Duration(e.conf.SInfo.CheckInterval),
			registerFunc, utime.WithReturn(true)); err != nil {
			log.Printf(fmt.Sprintf("ticker do err:%s", err))
			return
		}
	}()
	return nil
}

func (e *etcd) UnRegisterServer() error {
	if !e.conf.OpenDiscovery { //没有注册服务的，直接返回
		return nil
	}
	e.conf.SInfo.stop = true
	if _, err := e.Delete(context.Background(), e.getServerKey()); err != nil {
		return err
	}
	return nil
}

func (e *etcd) GetServers() map[string]*ServerInfo {
	return registeredServers[e.confName]
}

// WatchServers 监听服务信息，有变动时会写入channel servers
func (e *etcd) WatchServers(servers chan map[string]*ServerInfo) {
	e.Watch(context.Background(), e.conf.SInfo.ServerName, &serverHandle{
		l:          &sync.RWMutex{},
		e:          e,
		serverChan: servers,
	}, clientv3.WithPrefix())
}

func (e *etcd) getServerKey() string {
	return fmt.Sprintf("%s-%s", e.conf.SInfo.ServerName, e.conf.SInfo.ServerID)
}

func (e *etcd) serverIDFromKey(key string) string {
	return strings.Replace(key, fmt.Sprintf("%s-", e.conf.SInfo.ServerName), "", 1)
}
