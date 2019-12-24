package hue

import (
	"context"
	"github.com/collinux/gohue"
	"github.com/gissleh/lucifer"
	"net"
	"sync"
	"time"
)

func New() lucifer.Driver {
	return &driver{
		bridgeList: make([]*bridge, 0, 64),
		bridgeMap:  make(map[string]*bridge, 64),
	}
}

type driver struct {
	mutex      sync.Mutex
	bridgeList []*bridge
	bridgeMap  map[string]*bridge
}

func (driver *driver) SetupBridge(ctx context.Context, ip string) (lucifer.Bridge, string, error) {
	ghBridge, err := hue.NewBridge(ip)
	if err != nil {
		return nil, "", err
	}

	var key string
	for {
		newKey, err := ghBridge.CreateUser("github.com/gissleh/lucifer")
		if err == nil {
			if _, ok := err.(net.Error); ok {
				return nil, "", err
			}

			key = newKey
			break
		}

		select {
		case <-time.After(time.Second):
		case <-ctx.Done():
			return nil, "", ctx.Err()
		}
	}

	err = ghBridge.Login(key)
	if err != nil {
		return nil, key, err
	}

	bridge := &bridge{gh: ghBridge}

	driver.mutex.Lock()
	driver.bridgeList = append(driver.bridgeList, bridge)
	driver.bridgeMap[bridge.ID()] = bridge
	driver.mutex.Unlock()

	return bridge, key, nil
}

func (driver *driver) AddBridge(ctx context.Context, ip, key string) (lucifer.Bridge, error) {
	ghBridge, err := hue.NewBridge(ip)
	if err != nil {
		return nil, err
	}

	err = ghBridge.Login(key)
	if err != nil {
		return nil, err
	}

	bridge := &bridge{gh: ghBridge}

	driver.mutex.Lock()
	driver.bridgeList = append(driver.bridgeList, bridge)
	driver.bridgeMap[bridge.ID()] = bridge
	driver.mutex.Unlock()

	return bridge, nil
}

func (driver *driver) Bridge(id string) lucifer.Bridge {
	driver.mutex.Lock()
	bridge := driver.bridgeMap[id]
	driver.mutex.Unlock()

	return bridge
}

func (driver *driver) Bridges() []lucifer.Bridge {
	driver.mutex.Lock()
	list := make([]lucifer.Bridge, 0, len(driver.bridgeList))
	for _, bridge := range driver.bridgeList {
		list = append(list, bridge)
	}
	driver.mutex.Unlock()

	return list
}
