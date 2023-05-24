package qtum

import (
	"context"
	"errors"
	"sync"
	"time"
)

var errorHandlers map[error]errorHandler
var ErrErrorHandlerFailed = errors.New("failed to recover from error")
var ErrErrorHandlerRunning = errors.New("error recovery routine already running")

type errorHandler func(ctx context.Context, qtum *Qtum, state *errorState, method *Method) error

func newErrorState() *errorState {
	return &errorState{
		state: make(map[string]interface{}),
	}
}

type errorState struct {
	mutex sync.RWMutex
	state map[string]interface{}
}

func (e *errorState) Get(variable string) interface{} {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	return e.state[variable]
}

func (e *errorState) Put(variable string, value interface{}) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.state[variable] = value
}

func errWalletNotFoundHandler(ctx context.Context, qtum *Qtum, state *errorState, method *Method) error {
	if state.Get("createwallet") != nil {
		qtum.GetLogger().Log("msg", "createwallet call failed")
		return ErrErrorHandlerFailed
	}

	if state.Get("createwalletRunning") != nil {
		qtum.GetLogger().Log("msg", "Wallet error handler already running")
		return ErrErrorHandlerRunning
	}

	state.Put("createwalletRunning", true)
	defer func() {
		state.Put("createwalletRunning", nil)
	}()

	listWalletsReq := ListWalletsRequest([]string{})
	wallets, err := method.ListWallets(ctx, &listWalletsReq)

	if err != nil {
		qtum.GetLogger().Log("msg", "Error listing wallets", "err", err)
		return err
	}

	if len(*wallets) == 1 {
		// should be fixed, another node probably fixed the wallets, if it becomes a blocking problem the healthcheck will pick it up
		qtum.GetLogger().Log("msg", "Only one wallet loaded in Qtumd, will not try to fix wallet issues, another node might have already fixed things")
		return nil
	}

	if len(*wallets) == 0 {
		err := loadWallets(ctx, qtum, state, method)
		if err != nil {
			qtum.GetLogger().Log("msg", "Error loading wallets", "err", err)
			return err
		}
	} else {
		// multiple wallets loaded, unload them all
		for _, wallet := range *wallets {
			unloadWalletReq := UnloadWalletRequest([]string{wallet})
			_, err := method.UnloadWallet(ctx, &unloadWalletReq)
			if err != nil {
				qtum.GetLogger().Log("msg", "Error unloading wallet", "wallet", wallet, "err", err)
			}
		}

		wallets, err = method.ListWallets(ctx, &listWalletsReq)

		if err != nil {
			qtum.GetLogger().Log("msg", "Error listing wallets after unloading all", "err", err)
			return err
		}

		if len(*wallets) == 1 {
			qtum.GetLogger().Log("msg", "Unloaded all wallets but there is one loaded still, assuming another node fixed things")
			return nil
		} else if len(*wallets) != 0 {
			qtum.GetLogger().Log("msg", "Failed to unload wallets, multiple still are loaded")
			return errors.New("Failed to unload all wallets")
		}

		err := loadWallets(ctx, qtum, state, method)
		if err != nil {
			qtum.GetLogger().Log("msg", "Error loading wallets after unloading all", "err", err)
			return err
		}
	}

	return nil
}

func loadWallets(ctx context.Context, qtum *Qtum, state *errorState, method *Method) error {
	// no wallets loaded, we need to listwalletdir and load them each until requests go through
	// some wallets might have a password so we need to take that into account
	listWalletDirReq := ListWalletDirRequest([]string{})
	listedWallets, err := method.ListWalletDir(ctx, &listWalletDirReq)

	if err != nil {
		return err
	}

	if len(listedWallets.Wallets) == 0 {
		createWalletReq := CreateWalletRequest([]string{"wallet"})
		_, err = method.CreateWallet(ctx, &createWalletReq)

		if err == nil {
			state.Put("createwallet", true)
			go func() {
				select {
				case <-time.After(1 * time.Minute):
					// expire after a little bit - in case the qtum node changes
					state.Put("createwallet", false)
					/*
						case <-ctx.Done():
							return
					*/
				}
			}()
		}
	}

	loadWallet := state.Get("loadwallet")
	loadWalletIndex := -1
	if loadWallet != nil {
		loadWalletIndex = loadWallet.(int)
	}

	if loadWalletIndex > len(listedWallets.Wallets) {
		loadWalletIndex = -1
	}

	qtum.GetLogger().Log("msg", "loadwallet", "index", loadWalletIndex, "available", len(listedWallets.Wallets))

	for index, wallet := range listedWallets.Wallets {
		if loadWalletIndex > index {
			qtum.GetLogger().Log("msg", "continue")
			continue
		}

		loadWalletReq := LoadWalletRequest([]string{wallet.Name})
		_, err := method.LoadWallet(ctx, &loadWalletReq)
		if err == nil {
			qtum.GetLogger().Log("msg", "break")
			break
		}

		state.Put("loadwallet", index)

		go func() {
			time.Sleep(30 * time.Second)
			loadWallet := state.Get("loadwallet")
			if loadWallet != nil {
				loadWalletIndex := loadWallet.(int)
				if loadWalletIndex == index {
					state.Put("loadwallet", nil)
				}
			}
		}()

		break
	}

	return nil
}

func init() {
	errorHandlers = make(map[error]errorHandler)
	errorHandlers[ErrWalletNotFound] = errWalletNotFoundHandler
	errorHandlers[ErrWalletNotSpecified] = errWalletNotFoundHandler
}
