package server

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/revolutionchain/charon/pkg/revo"
)

var ErrNoRevoConnections = errors.New("revod has no connections")
var ErrCannotGetConnectedChain = errors.New("Cannot detect chain revod is connected to")
var ErrBlockSyncingSeemsStalled = errors.New("Block syncing seems stalled")
var ErrLostLotsOfBlocks = errors.New("Lost a lot of blocks, expected block height to be higher")
var ErrLostFewBlocks = errors.New("Lost a few blocks, expected block height to be higher")

func (s *Server) testConnectionToRevod() error {
	networkInfo, err := s.revoRPCClient.GetNetworkInfo(s.revoRPCClient.GetContext())
	if err == nil {
		// chain can theoretically block forever if revod isn't up
		// but then GetNetworkInfo would be erroring
		chainChan := make(chan string)
		getChainTimeout := time.NewTimer(10 * time.Second)
		go func(ch chan string) {
			chain := s.revoRPCClient.Chain()
			chainChan <- chain
		}(chainChan)

		select {
		case chain := <-chainChan:
			if chain == revo.ChainRegTest {
				// ignore how many connections there are
				return nil
			}
			if networkInfo.Connections == 0 {
				s.logger.Log("liveness", "Revod has no network connections")
				return ErrNoRevoConnections
			}
			break
		case <-getChainTimeout.C:
			s.logger.Log("liveness", "Revod getnetworkinfo request timed out")
			return ErrCannotGetConnectedChain
		}
	} else {
		s.logger.Log("liveness", "Revod getnetworkinfo errored", "err", err)
	}
	return err
}

func (s *Server) testLogEvents() error {
	_, err := s.revoRPCClient.GetTransactionReceipt(s.revoRPCClient.GetContext(), "0000000000000000000000000000000000000000000000000000000000000000")
	if err == revo.ErrInternalError {
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

	blockChainInfo, err := s.revoRPCClient.GetBlockChainInfo(s.revoRPCClient.GetContext())
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

func (s *Server) testRevodErrorRate() error {
	minimumSuccessRate := float32(*s.healthCheckPercent / 100)
	revoSuccessRate := s.revoRequestAnalytics.GetSuccessRate()

	if revoSuccessRate < minimumSuccessRate {
		s.logger.Log("liveness", "revod request success rate is low", "rate", revoSuccessRate)
		return errors.New(fmt.Sprintf("revod request success rate is %f<%f", revoSuccessRate, minimumSuccessRate))
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
