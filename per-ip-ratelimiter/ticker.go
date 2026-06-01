package main

import "time"

type CleanManager struct {
	ticker      *time.Ticker
	stopChan    chan struct{}
	doneChan    chan struct{}
	maxIdleTime time.Duration
}

func newCleanupManager(cleanupInterval, maxIdleTime time.Duration) *CleanManager {
	return &CleanManager{
		ticker:      time.NewTicker(cleanupInterval),
		stopChan:    make(chan struct{}),
		doneChan:    make(chan struct{}),
		maxIdleTime: maxIdleTime,
	}
}
func (cm *CleanManager) Start() {
	go func() {
		for {
			select {
			case <-cm.ticker.C:
				cm.doCleanUp()
			case <-cm.stopChan:
				cm.ticker.Stop()
				close(cm.doneChan)
				return
			}
		}
	}()
}
func (cm *CleanManager) doCleanUp() {
	mu.Lock()
	defer mu.Unlock()

	now := time.Now()
	for ip, client := range clients {
		if now.Sub(client.lastSeen) > cm.maxIdleTime {
			delete(clients, ip)
		}
	}
}
func (cm *CleanManager) Stop() {
	close(cm.stopChan)
	<-cm.doneChan
}
