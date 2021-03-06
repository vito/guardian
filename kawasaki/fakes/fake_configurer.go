// This file was generated by counterfeiter
package fakes

import (
	"sync"

	"github.com/cloudfoundry-incubator/guardian/kawasaki"
	"github.com/pivotal-golang/lager"
)

type FakeConfigurer struct {
	ApplyStub        func(log lager.Logger, cfg kawasaki.NetworkConfig, nsPath string) error
	applyMutex       sync.RWMutex
	applyArgsForCall []struct {
		log    lager.Logger
		cfg    kawasaki.NetworkConfig
		nsPath string
	}
	applyReturns struct {
		result1 error
	}
	DestroyStub        func(log lager.Logger, cfg kawasaki.NetworkConfig) error
	destroyMutex       sync.RWMutex
	destroyArgsForCall []struct {
		log lager.Logger
		cfg kawasaki.NetworkConfig
	}
	destroyReturns struct {
		result1 error
	}
}

func (fake *FakeConfigurer) Apply(log lager.Logger, cfg kawasaki.NetworkConfig, nsPath string) error {
	fake.applyMutex.Lock()
	fake.applyArgsForCall = append(fake.applyArgsForCall, struct {
		log    lager.Logger
		cfg    kawasaki.NetworkConfig
		nsPath string
	}{log, cfg, nsPath})
	fake.applyMutex.Unlock()
	if fake.ApplyStub != nil {
		return fake.ApplyStub(log, cfg, nsPath)
	} else {
		return fake.applyReturns.result1
	}
}

func (fake *FakeConfigurer) ApplyCallCount() int {
	fake.applyMutex.RLock()
	defer fake.applyMutex.RUnlock()
	return len(fake.applyArgsForCall)
}

func (fake *FakeConfigurer) ApplyArgsForCall(i int) (lager.Logger, kawasaki.NetworkConfig, string) {
	fake.applyMutex.RLock()
	defer fake.applyMutex.RUnlock()
	return fake.applyArgsForCall[i].log, fake.applyArgsForCall[i].cfg, fake.applyArgsForCall[i].nsPath
}

func (fake *FakeConfigurer) ApplyReturns(result1 error) {
	fake.ApplyStub = nil
	fake.applyReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeConfigurer) Destroy(log lager.Logger, cfg kawasaki.NetworkConfig) error {
	fake.destroyMutex.Lock()
	fake.destroyArgsForCall = append(fake.destroyArgsForCall, struct {
		log lager.Logger
		cfg kawasaki.NetworkConfig
	}{log, cfg})
	fake.destroyMutex.Unlock()
	if fake.DestroyStub != nil {
		return fake.DestroyStub(log, cfg)
	} else {
		return fake.destroyReturns.result1
	}
}

func (fake *FakeConfigurer) DestroyCallCount() int {
	fake.destroyMutex.RLock()
	defer fake.destroyMutex.RUnlock()
	return len(fake.destroyArgsForCall)
}

func (fake *FakeConfigurer) DestroyArgsForCall(i int) (lager.Logger, kawasaki.NetworkConfig) {
	fake.destroyMutex.RLock()
	defer fake.destroyMutex.RUnlock()
	return fake.destroyArgsForCall[i].log, fake.destroyArgsForCall[i].cfg
}

func (fake *FakeConfigurer) DestroyReturns(result1 error) {
	fake.DestroyStub = nil
	fake.destroyReturns = struct {
		result1 error
	}{result1}
}

var _ kawasaki.Configurer = new(FakeConfigurer)
