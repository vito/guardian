// This file was generated by counterfeiter
package fakes

import (
	"sync"

	"github.com/cloudfoundry-incubator/guardian/rundmc"
)

type FakeEventStore struct {
	OnEventStub        func(id string, event string)
	onEventMutex       sync.RWMutex
	onEventArgsForCall []struct {
		id    string
		event string
	}
	EventsStub        func(id string) []string
	eventsMutex       sync.RWMutex
	eventsArgsForCall []struct {
		id string
	}
	eventsReturns struct {
		result1 []string
	}
}

func (fake *FakeEventStore) OnEvent(id string, event string) {
	fake.onEventMutex.Lock()
	fake.onEventArgsForCall = append(fake.onEventArgsForCall, struct {
		id    string
		event string
	}{id, event})
	fake.onEventMutex.Unlock()
	if fake.OnEventStub != nil {
		fake.OnEventStub(id, event)
	}
}

func (fake *FakeEventStore) OnEventCallCount() int {
	fake.onEventMutex.RLock()
	defer fake.onEventMutex.RUnlock()
	return len(fake.onEventArgsForCall)
}

func (fake *FakeEventStore) OnEventArgsForCall(i int) (string, string) {
	fake.onEventMutex.RLock()
	defer fake.onEventMutex.RUnlock()
	return fake.onEventArgsForCall[i].id, fake.onEventArgsForCall[i].event
}

func (fake *FakeEventStore) Events(id string) []string {
	fake.eventsMutex.Lock()
	fake.eventsArgsForCall = append(fake.eventsArgsForCall, struct {
		id string
	}{id})
	fake.eventsMutex.Unlock()
	if fake.EventsStub != nil {
		return fake.EventsStub(id)
	} else {
		return fake.eventsReturns.result1
	}
}

func (fake *FakeEventStore) EventsCallCount() int {
	fake.eventsMutex.RLock()
	defer fake.eventsMutex.RUnlock()
	return len(fake.eventsArgsForCall)
}

func (fake *FakeEventStore) EventsArgsForCall(i int) string {
	fake.eventsMutex.RLock()
	defer fake.eventsMutex.RUnlock()
	return fake.eventsArgsForCall[i].id
}

func (fake *FakeEventStore) EventsReturns(result1 []string) {
	fake.EventsStub = nil
	fake.eventsReturns = struct {
		result1 []string
	}{result1}
}

var _ rundmc.EventStore = new(FakeEventStore)
