// Code generated by counterfeiter. DO NOT EDIT.
package bridgefakes

import (
	"sync"

	"github.com/philborlin/committed/bridge"
	"github.com/philborlin/committed/syncable"
	"github.com/philborlin/committed/topic"
)

type FakeFactory struct {
	NewStub        func(string, syncable.Syncable, map[string]topic.Topic) (bridge.Bridge, error)
	newMutex       sync.RWMutex
	newArgsForCall []struct {
		arg1 string
		arg2 syncable.Syncable
		arg3 map[string]topic.Topic
	}
	newReturns struct {
		result1 bridge.Bridge
		result2 error
	}
	newReturnsOnCall map[int]struct {
		result1 bridge.Bridge
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeFactory) New(arg1 string, arg2 syncable.Syncable, arg3 map[string]topic.Topic) (bridge.Bridge, error) {
	fake.newMutex.Lock()
	ret, specificReturn := fake.newReturnsOnCall[len(fake.newArgsForCall)]
	fake.newArgsForCall = append(fake.newArgsForCall, struct {
		arg1 string
		arg2 syncable.Syncable
		arg3 map[string]topic.Topic
	}{arg1, arg2, arg3})
	fake.recordInvocation("New", []interface{}{arg1, arg2, arg3})
	fake.newMutex.Unlock()
	if fake.NewStub != nil {
		return fake.NewStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.newReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeFactory) NewCallCount() int {
	fake.newMutex.RLock()
	defer fake.newMutex.RUnlock()
	return len(fake.newArgsForCall)
}

func (fake *FakeFactory) NewCalls(stub func(string, syncable.Syncable, map[string]topic.Topic) (bridge.Bridge, error)) {
	fake.newMutex.Lock()
	defer fake.newMutex.Unlock()
	fake.NewStub = stub
}

func (fake *FakeFactory) NewArgsForCall(i int) (string, syncable.Syncable, map[string]topic.Topic) {
	fake.newMutex.RLock()
	defer fake.newMutex.RUnlock()
	argsForCall := fake.newArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeFactory) NewReturns(result1 bridge.Bridge, result2 error) {
	fake.newMutex.Lock()
	defer fake.newMutex.Unlock()
	fake.NewStub = nil
	fake.newReturns = struct {
		result1 bridge.Bridge
		result2 error
	}{result1, result2}
}

func (fake *FakeFactory) NewReturnsOnCall(i int, result1 bridge.Bridge, result2 error) {
	fake.newMutex.Lock()
	defer fake.newMutex.Unlock()
	fake.NewStub = nil
	if fake.newReturnsOnCall == nil {
		fake.newReturnsOnCall = make(map[int]struct {
			result1 bridge.Bridge
			result2 error
		})
	}
	fake.newReturnsOnCall[i] = struct {
		result1 bridge.Bridge
		result2 error
	}{result1, result2}
}

func (fake *FakeFactory) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.newMutex.RLock()
	defer fake.newMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeFactory) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ bridge.Factory = new(FakeFactory)
