// Code generated by counterfeiter. DO NOT EDIT.
package bridgefakes

import (
	"context"
	"sync"
	"time"

	"github.com/philborlin/committed/bridge"
)

type FakeBridge struct {
	InitStub        func(context.Context, chan<- error, time.Duration) error
	initMutex       sync.RWMutex
	initArgsForCall []struct {
		arg1 context.Context
		arg2 chan<- error
		arg3 time.Duration
	}
	initReturns struct {
		result1 error
	}
	initReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeBridge) Init(arg1 context.Context, arg2 chan<- error, arg3 time.Duration) error {
	fake.initMutex.Lock()
	ret, specificReturn := fake.initReturnsOnCall[len(fake.initArgsForCall)]
	fake.initArgsForCall = append(fake.initArgsForCall, struct {
		arg1 context.Context
		arg2 chan<- error
		arg3 time.Duration
	}{arg1, arg2, arg3})
	fake.recordInvocation("Init", []interface{}{arg1, arg2, arg3})
	fake.initMutex.Unlock()
	if fake.InitStub != nil {
		return fake.InitStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.initReturns
	return fakeReturns.result1
}

func (fake *FakeBridge) InitCallCount() int {
	fake.initMutex.RLock()
	defer fake.initMutex.RUnlock()
	return len(fake.initArgsForCall)
}

func (fake *FakeBridge) InitCalls(stub func(context.Context, chan<- error, time.Duration) error) {
	fake.initMutex.Lock()
	defer fake.initMutex.Unlock()
	fake.InitStub = stub
}

func (fake *FakeBridge) InitArgsForCall(i int) (context.Context, chan<- error, time.Duration) {
	fake.initMutex.RLock()
	defer fake.initMutex.RUnlock()
	argsForCall := fake.initArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeBridge) InitReturns(result1 error) {
	fake.initMutex.Lock()
	defer fake.initMutex.Unlock()
	fake.InitStub = nil
	fake.initReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeBridge) InitReturnsOnCall(i int, result1 error) {
	fake.initMutex.Lock()
	defer fake.initMutex.Unlock()
	fake.InitStub = nil
	if fake.initReturnsOnCall == nil {
		fake.initReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.initReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeBridge) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.initMutex.RLock()
	defer fake.initMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeBridge) recordInvocation(key string, args []interface{}) {
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

var _ bridge.Bridge = new(FakeBridge)
