// This file was generated by counterfeiter
package fakes

import (
	"sync"

	"github.com/cloudfoundry-incubator/guardian/kawasaki"
)

type FakePortForwarder struct {
	ForwardStub        func(spec kawasaki.PortForwarderSpec) error
	forwardMutex       sync.RWMutex
	forwardArgsForCall []struct {
		spec kawasaki.PortForwarderSpec
	}
	forwardReturns struct {
		result1 error
	}
}

func (fake *FakePortForwarder) Forward(spec kawasaki.PortForwarderSpec) error {
	fake.forwardMutex.Lock()
	fake.forwardArgsForCall = append(fake.forwardArgsForCall, struct {
		spec kawasaki.PortForwarderSpec
	}{spec})
	fake.forwardMutex.Unlock()
	if fake.ForwardStub != nil {
		return fake.ForwardStub(spec)
	} else {
		return fake.forwardReturns.result1
	}
}

func (fake *FakePortForwarder) ForwardCallCount() int {
	fake.forwardMutex.RLock()
	defer fake.forwardMutex.RUnlock()
	return len(fake.forwardArgsForCall)
}

func (fake *FakePortForwarder) ForwardArgsForCall(i int) kawasaki.PortForwarderSpec {
	fake.forwardMutex.RLock()
	defer fake.forwardMutex.RUnlock()
	return fake.forwardArgsForCall[i].spec
}

func (fake *FakePortForwarder) ForwardReturns(result1 error) {
	fake.ForwardStub = nil
	fake.forwardReturns = struct {
		result1 error
	}{result1}
}

var _ kawasaki.PortForwarder = new(FakePortForwarder)
