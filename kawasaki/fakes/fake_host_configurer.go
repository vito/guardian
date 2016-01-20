// This file was generated by counterfeiter
package fakes

import (
	"os"
	"sync"

	"github.com/cloudfoundry-incubator/guardian/kawasaki"
	"github.com/pivotal-golang/lager"
)

type FakeHostConfigurer struct {
	ApplyStub        func(logger lager.Logger, cfg kawasaki.NetworkConfig, netnsFD *os.File) error
	applyMutex       sync.RWMutex
	applyArgsForCall []struct {
		logger  lager.Logger
		cfg     kawasaki.NetworkConfig
		netnsFD *os.File
	}
	applyReturns struct {
		result1 error
	}
	DestroyStub        func(cfg kawasaki.NetworkConfig) error
	destroyMutex       sync.RWMutex
	destroyArgsForCall []struct {
		cfg kawasaki.NetworkConfig
	}
	destroyReturns struct {
		result1 error
	}
}

func (fake *FakeHostConfigurer) Apply(logger lager.Logger, cfg kawasaki.NetworkConfig, netnsFD *os.File) error {
	fake.applyMutex.Lock()
	fake.applyArgsForCall = append(fake.applyArgsForCall, struct {
		logger  lager.Logger
		cfg     kawasaki.NetworkConfig
		netnsFD *os.File
	}{logger, cfg, netnsFD})
	fake.applyMutex.Unlock()
	if fake.ApplyStub != nil {
		return fake.ApplyStub(logger, cfg, netnsFD)
	} else {
		return fake.applyReturns.result1
	}
}

func (fake *FakeHostConfigurer) ApplyCallCount() int {
	fake.applyMutex.RLock()
	defer fake.applyMutex.RUnlock()
	return len(fake.applyArgsForCall)
}

func (fake *FakeHostConfigurer) ApplyArgsForCall(i int) (lager.Logger, kawasaki.NetworkConfig, *os.File) {
	fake.applyMutex.RLock()
	defer fake.applyMutex.RUnlock()
	return fake.applyArgsForCall[i].logger, fake.applyArgsForCall[i].cfg, fake.applyArgsForCall[i].netnsFD
}

func (fake *FakeHostConfigurer) ApplyReturns(result1 error) {
	fake.ApplyStub = nil
	fake.applyReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeHostConfigurer) Destroy(cfg kawasaki.NetworkConfig) error {
	fake.destroyMutex.Lock()
	fake.destroyArgsForCall = append(fake.destroyArgsForCall, struct {
		cfg kawasaki.NetworkConfig
	}{cfg})
	fake.destroyMutex.Unlock()
	if fake.DestroyStub != nil {
		return fake.DestroyStub(cfg)
	} else {
		return fake.destroyReturns.result1
	}
}

func (fake *FakeHostConfigurer) DestroyCallCount() int {
	fake.destroyMutex.RLock()
	defer fake.destroyMutex.RUnlock()
	return len(fake.destroyArgsForCall)
}

func (fake *FakeHostConfigurer) DestroyArgsForCall(i int) kawasaki.NetworkConfig {
	fake.destroyMutex.RLock()
	defer fake.destroyMutex.RUnlock()
	return fake.destroyArgsForCall[i].cfg
}

func (fake *FakeHostConfigurer) DestroyReturns(result1 error) {
	fake.DestroyStub = nil
	fake.destroyReturns = struct {
		result1 error
	}{result1}
}

var _ kawasaki.HostConfigurer = new(FakeHostConfigurer)
