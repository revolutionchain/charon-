package server

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/revolutionchain/charon/pkg/qtum"
)

var ErrNoQtumConnections = errors.New("revod has no connections")
var ErrCannotGetConnectedChain = errors.New("Cannot detect chain revod is connected to")
var ErrBlockSyncingSeemsStalled = errors.New("Block syncing seems stalled")
var ErrLostLotsOfBlocks = errors.New("Lost a lot of blocks, expected block height to be higher")
var ErrLostFewBlocks = errors.New("Lost a few blocks, expected block height to be higher")

func (s *Server) testConnectionToQtumd() error {
	networkInfo, err := s.qtumRPCClient.GetNetworkInfo(s.qtumRPCClient.GetContext())
	if err == nil {
		// chain can theoretically block forever if revod isn't up
		// but then GetNetworkInfo would be erroring
		chainChan := make(chan string)
		getChainTimeout := time.NewTimer(10 * time.Second)
		go func(ch chan string) {
			chain := s.qtumRPCClient.Chain()
			chainChan <- chain
		}(chainChan)

		select {
		case chain := <-chainChan:
			if chain == qtum.ChainRegTest {
				// ignore how many connections there are
				return nil
			}
			if networkInfo.Connections == 0 {
				s.logger.Log("liveness", "Qtumd has no network connections")
				return ErrNoQtumConnections
			}
			break
		case <-getChainTimeout.C:
			s.logger.Log("liveness", "Qtumd getnetworkinfo request timed out")
			return ErrCannotGetConnectedChain
		}
	} else {
		s.logger.Log("liveness", "Qtumd getnetworkinfo errored", "err", err)
	}
	return err
}

func (s *Server) testLogEvents() error {
	_, err := s.qtumRPCClient.GetTransactionReceipt(s.qtumRPCClient.GetContext(), "0000000000000000000000000000000000000000000000000000000000000000")
	if err == qtum.ErrInternalError {
		s.logger.Log("liveness", "-logevents might not be enabled")
		return errors.Wrap(err, "-logevents might not be enabled")
	}
	return nil
}

func (s *Server) testBlocksSyncing() error {
	s.blocksMutex.RLock()
	nextBlockCheck := s.nextBlockCheck
	lastBlockStatus := s.lastBlockStatus
	s.blocksMutex.RUnlock()
	now := time.Now()
	if nextBlockCheck == nil {
		nextBlockCheckTime := time.Now().Add(-30 * time.Minute)
		nextBlockCheck = &nextBlockCheckTime
	}
	if nextBlockCheck.After(now) {
		if lastBlockStatus != nil {
			s.logger.Log("liveness", "blocks syncing", "err", lastBlockStatus)
		}
		return lastBlockStatus
	}
	s.blocksMutex.Lock()
	if s.nextBlockCheck != nil && nextBlockCheck != s.nextBlockCheck {
		// multiple threads were waiting on write lock
		s.blocksMutex.Unlock()
		return s.testBlocksSyncing()
	}
	defer s.blocksMutex.Unlock()

	blockChainInfo, err := s.qtumRPCClient.GetBlockChainInfo(s.qtumRPCClient.GetContext())
	if err != nil {
		s.logger.Log("liveness", "getblockchainfo request failed", "err", err)
		return err
	}

	nextBlockCheckTime := time.Now().Add(5 * time.Minute)
	s.nextBlockCheck = &nextBlockCheckTime

	if blockChainInfo.Blocks == s.lastBlock {
		// stalled
		nextBlockCheckTime = time.Now().Add(15 * time.Second)
		s.nextBlockCheck = &nextBlockCheckTime
		s.lastBlockStatus = ErrBlockSyncingSeemsStalled
	} else if blockChainInfo.Blocks < s.lastBlock {
		// lost some blocks...?
		if s.lastBlock-blockChainInfo.Blocks > 10 {
			// lost a lot of blocks
			// probably a real problem
			s.lastBlock = 0
			nextBlockCheckTime = time.Now().Add(60 * time.Second)
			s.nextBlockCheck = &nextBlockCheckTime
			s.logger.Log("liveness", "Lost lots of blocks")
			s.lastBlockStatus = ErrLostLotsOfBlocks
		} else {
			// lost a few blocks
			// could be revod nodes out of sync behind a load balancer
			nextBlockCheckTime = time.Now().Add(10 * time.Second)
			s.nextBlockCheck = &nextBlockCheckTime
			s.logger.Log("liveness", "Lost a few blocks")
			s.lastBlockStatus = ErrLostFewBlocks
		}
	} else {
		// got a higher block height than last time
		s.lastBlock = blockChainInfo.Blocks
		nextBlockCheckTime = time.Now().Add(90 * time.Second)
		s.nextBlockCheck = &nextBlockCheckTime
		s.lastBlockStatus = nil
	}

	return s.lastBlockStatus
}

func (s *Server) testQtumdErrorRate() error {
	minimumSuccessRate := float32(*s.healthCheckPercent / 100)
	qtumSuccessRate := s.qtumRequestAnalytics.GetSuccessRate()

	if qtumSuccessRate < minimumSuccessRate {
		s.logger.Log("liveness", "revod request success rate is low", "rate", qtumSuccessRate)
		return errors.New(fmt.Sprintf("revod request success rate is %f<%f", qtumSuccessRate, minimumSuccessRate))
	} else {
		return nil
	}
}

func (s *Server) testCharonErrorRate() error {
	minimumSuccessRate := float32(*s.healthCheckPercent / 100)
	ethSuccessRate := s.ethRequestAnalytics.GetSuccessRate()

	if ethSuccessRate < minimumSuccessRate {
		s.logger.Log("liveness", "client eth success rate is low", "rate", ethSuccessRate)
		return errors.New(fmt.Sprintf("client eth request success rate is %f<%f", ethSuccessRate, minimumSuccessRate))
	} else {
		return nil
	}
}
