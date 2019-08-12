// Code generated by counterfeiter. DO NOT EDIT.
package topicfakes

import (
	"sync"

	"github.com/philborlin/committed/topic"
)

type FakeTopic struct {
	AppendStub        func(topic.Data) error
	appendMutex       sync.RWMutex
	appendArgsForCall []struct {
		arg1 topic.Data
	}
	appendReturns struct {
		result1 error
	}
	appendReturnsOnCall map[int]struct {
		result1 error
	}
	CloseStub        func() error
	closeMutex       sync.RWMutex
	closeArgsForCall []struct {
	}
	closeReturns struct {
		result1 error
	}
	closeReturnsOnCall map[int]struct {
		result1 error
	}
	NameStub        func() string
	nameMutex       sync.RWMutex
	nameArgsForCall []struct {
	}
	nameReturns struct {
		result1 string
	}
	nameReturnsOnCall map[int]struct {
		result1 string
	}
	NewReaderStub        func(uint64) (topic.Reader, error)
	newReaderMutex       sync.RWMutex
	newReaderArgsForCall []struct {
		arg1 uint64
	}
	newReaderReturns struct {
		result1 topic.Reader
		result2 error
	}
	newReaderReturnsOnCall map[int]struct {
		result1 topic.Reader
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeTopic) Append(arg1 topic.Data) error {
	fake.appendMutex.Lock()
	ret, specificReturn := fake.appendReturnsOnCall[len(fake.appendArgsForCall)]
	fake.appendArgsForCall = append(fake.appendArgsForCall, struct {
		arg1 topic.Data
	}{arg1})
	fake.recordInvocation("Append", []interface{}{arg1})
	fake.appendMutex.Unlock()
	if fake.AppendStub != nil {
		return fake.AppendStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.appendReturns
	return fakeReturns.result1
}

func (fake *FakeTopic) AppendCallCount() int {
	fake.appendMutex.RLock()
	defer fake.appendMutex.RUnlock()
	return len(fake.appendArgsForCall)
}

func (fake *FakeTopic) AppendCalls(stub func(topic.Data) error) {
	fake.appendMutex.Lock()
	defer fake.appendMutex.Unlock()
	fake.AppendStub = stub
}

func (fake *FakeTopic) AppendArgsForCall(i int) topic.Data {
	fake.appendMutex.RLock()
	defer fake.appendMutex.RUnlock()
	argsForCall := fake.appendArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeTopic) AppendReturns(result1 error) {
	fake.appendMutex.Lock()
	defer fake.appendMutex.Unlock()
	fake.AppendStub = nil
	fake.appendReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeTopic) AppendReturnsOnCall(i int, result1 error) {
	fake.appendMutex.Lock()
	defer fake.appendMutex.Unlock()
	fake.AppendStub = nil
	if fake.appendReturnsOnCall == nil {
		fake.appendReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.appendReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeTopic) Close() error {
	fake.closeMutex.Lock()
	ret, specificReturn := fake.closeReturnsOnCall[len(fake.closeArgsForCall)]
	fake.closeArgsForCall = append(fake.closeArgsForCall, struct {
	}{})
	fake.recordInvocation("Close", []interface{}{})
	fake.closeMutex.Unlock()
	if fake.CloseStub != nil {
		return fake.CloseStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.closeReturns
	return fakeReturns.result1
}

func (fake *FakeTopic) CloseCallCount() int {
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	return len(fake.closeArgsForCall)
}

func (fake *FakeTopic) CloseCalls(stub func() error) {
	fake.closeMutex.Lock()
	defer fake.closeMutex.Unlock()
	fake.CloseStub = stub
}

func (fake *FakeTopic) CloseReturns(result1 error) {
	fake.closeMutex.Lock()
	defer fake.closeMutex.Unlock()
	fake.CloseStub = nil
	fake.closeReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeTopic) CloseReturnsOnCall(i int, result1 error) {
	fake.closeMutex.Lock()
	defer fake.closeMutex.Unlock()
	fake.CloseStub = nil
	if fake.closeReturnsOnCall == nil {
		fake.closeReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.closeReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeTopic) Name() string {
	fake.nameMutex.Lock()
	ret, specificReturn := fake.nameReturnsOnCall[len(fake.nameArgsForCall)]
	fake.nameArgsForCall = append(fake.nameArgsForCall, struct {
	}{})
	fake.recordInvocation("Name", []interface{}{})
	fake.nameMutex.Unlock()
	if fake.NameStub != nil {
		return fake.NameStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.nameReturns
	return fakeReturns.result1
}

func (fake *FakeTopic) NameCallCount() int {
	fake.nameMutex.RLock()
	defer fake.nameMutex.RUnlock()
	return len(fake.nameArgsForCall)
}

func (fake *FakeTopic) NameCalls(stub func() string) {
	fake.nameMutex.Lock()
	defer fake.nameMutex.Unlock()
	fake.NameStub = stub
}

func (fake *FakeTopic) NameReturns(result1 string) {
	fake.nameMutex.Lock()
	defer fake.nameMutex.Unlock()
	fake.NameStub = nil
	fake.nameReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeTopic) NameReturnsOnCall(i int, result1 string) {
	fake.nameMutex.Lock()
	defer fake.nameMutex.Unlock()
	fake.NameStub = nil
	if fake.nameReturnsOnCall == nil {
		fake.nameReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.nameReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *FakeTopic) NewReader(arg1 uint64) (topic.Reader, error) {
	fake.newReaderMutex.Lock()
	ret, specificReturn := fake.newReaderReturnsOnCall[len(fake.newReaderArgsForCall)]
	fake.newReaderArgsForCall = append(fake.newReaderArgsForCall, struct {
		arg1 uint64
	}{arg1})
	fake.recordInvocation("NewReader", []interface{}{arg1})
	fake.newReaderMutex.Unlock()
	if fake.NewReaderStub != nil {
		return fake.NewReaderStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.newReaderReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeTopic) NewReaderCallCount() int {
	fake.newReaderMutex.RLock()
	defer fake.newReaderMutex.RUnlock()
	return len(fake.newReaderArgsForCall)
}

func (fake *FakeTopic) NewReaderCalls(stub func(uint64) (topic.Reader, error)) {
	fake.newReaderMutex.Lock()
	defer fake.newReaderMutex.Unlock()
	fake.NewReaderStub = stub
}

func (fake *FakeTopic) NewReaderArgsForCall(i int) uint64 {
	fake.newReaderMutex.RLock()
	defer fake.newReaderMutex.RUnlock()
	argsForCall := fake.newReaderArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeTopic) NewReaderReturns(result1 topic.Reader, result2 error) {
	fake.newReaderMutex.Lock()
	defer fake.newReaderMutex.Unlock()
	fake.NewReaderStub = nil
	fake.newReaderReturns = struct {
		result1 topic.Reader
		result2 error
	}{result1, result2}
}

func (fake *FakeTopic) NewReaderReturnsOnCall(i int, result1 topic.Reader, result2 error) {
	fake.newReaderMutex.Lock()
	defer fake.newReaderMutex.Unlock()
	fake.NewReaderStub = nil
	if fake.newReaderReturnsOnCall == nil {
		fake.newReaderReturnsOnCall = make(map[int]struct {
			result1 topic.Reader
			result2 error
		})
	}
	fake.newReaderReturnsOnCall[i] = struct {
		result1 topic.Reader
		result2 error
	}{result1, result2}
}

func (fake *FakeTopic) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.appendMutex.RLock()
	defer fake.appendMutex.RUnlock()
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	fake.nameMutex.RLock()
	defer fake.nameMutex.RUnlock()
	fake.newReaderMutex.RLock()
	defer fake.newReaderMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeTopic) recordInvocation(key string, args []interface{}) {
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

var _ topic.Topic = new(FakeTopic)
