package bridge

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/philborlin/committed/internal/node/syncable"
	"github.com/philborlin/committed/internal/node/syncable/syncablefakes"
	"github.com/philborlin/committed/internal/node/topic"
	"github.com/philborlin/committed/internal/node/topic/topicfakes"
	"github.com/philborlin/committed/internal/node/types"
	"github.com/philborlin/committed/internal/node/types/typesfakes"
)

var _ = Describe("Topic", func() {
	const (
		bridgeName string = "qux"
	)

	var (
		fooFakeTopic *topicfakes.FakeTopic
		barFakeTopic *topicfakes.FakeTopic
		fakeReader   *topicfakes.FakeReader
		fakeLeader   *typesfakes.FakeLeader
		fakeProposer *typesfakes.FakeProposer

		topics        map[string]topic.Topic
		fakeSyncable  *syncablefakes.FakeSyncable
		bridgeFactory Factory

		zeroIndex types.Index
		oneIndex  types.Index
		twoIndex  types.Index
	)

	JustBeforeEach(func() {
		fooFakeTopic = &topicfakes.FakeTopic{}
		fooFakeTopic.NameReturns("foo")
		barFakeTopic = &topicfakes.FakeTopic{}
		barFakeTopic.NameReturns("bar")

		fakeReader = &topicfakes.FakeReader{}
		fooFakeTopic.NewReaderReturns(fakeReader, nil)
		fakeLeader = &typesfakes.FakeLeader{}
		fakeLeader.IsLeaderReturns(true)
		fakeProposer = &typesfakes.FakeProposer{}

		topics = map[string]topic.Topic{"foo": fooFakeTopic, "bar": barFakeTopic}
		fakeSyncable = &syncablefakes.FakeSyncable{}
		bridgeFactory = &TopicSyncableBridgeFactory{}

		zeroIndex = types.Index{Index: 0, Term: 0}
		oneIndex = types.Index{Index: 1, Term: 1}
		twoIndex = types.Index{Index: 2, Term: 2}
	})

	new := func(name string, s syncable.Syncable, topics map[string]topic.Topic,
		snap *Snapshot) (*TopicSyncableBridge, error) {
		b, err := bridgeFactory.New(name, fakeSyncable, topics, fakeLeader, fakeProposer, snap)
		if err != nil {
			return nil, err
		}
		return b.(*TopicSyncableBridge), err
	}

	Describe("NewBridge()", func() {
		It("should create a new bridge", func() {
			fakeSyncable.TopicsReturns([]string{"foo"})
			b, err := new(bridgeName, fakeSyncable, topics, nil)
			Expect(err).To(BeNil())
			Expect(len(b.topics)).To(Equal(1))
			Expect(b.topics["foo"]).To(Equal(fooFakeTopic))
			Expect(b.Name).To(Equal(bridgeName))
			Expect(b.Syncable).To(Equal(fakeSyncable))
			Expect(b.lastIndex).To(Equal(zeroIndex))
		})

		It("should apply the snapshot if there is one", func() {
			snap := &Snapshot{LastIndex: oneIndex}

			fakeSyncable.TopicsReturns([]string{"foo"})
			b, err := new(bridgeName, fakeSyncable, topics, snap)
			Expect(err).To(BeNil())
			Expect(len(b.topics)).To(Equal(1))
			Expect(b.topics["foo"]).To(Equal(fooFakeTopic))
			Expect(b.Name).To(Equal(bridgeName))
			Expect(b.Syncable).To(Equal(fakeSyncable))
			Expect(b.lastIndex).To(Equal(oneIndex))
		})

		It("should error if there are no topics", func() {
			fakeSyncable.TopicsReturns([]string{})
			b, err := new(bridgeName, fakeSyncable, topics, nil)
			Expect(b).To(BeNil())
			Expect(err.Error()).To(ContainSubstring("No topics so there is nothing to sync"))
		})

		It("should error if there are too many topics", func() {
			fakeSyncable.TopicsReturns([]string{"foo", "bar"})
			b, err := new(bridgeName, fakeSyncable, topics, nil)
			Expect(b).To(BeNil())
			Expect(err.Error()).To(ContainSubstring("We don't support more than one topic in a syncable yet"))
		})

		It("should error if topic does not exist", func() {
			fakeSyncable.TopicsReturns([]string{"baz"})
			b, err := new(bridgeName, fakeSyncable, topics, nil)
			Expect(b).To(BeNil())
			Expect(err.Error()).To(ContainSubstring("is trying to listen to topic"))
		})
	})

	Describe("GetSnapshot()/ApplySnapshot()", func() {
		var (
			b   *TopicSyncableBridge
			err error
		)

		JustBeforeEach(func() {
			fakeSyncable.TopicsReturns([]string{"foo"})
			b, err = new(bridgeName, fakeSyncable, topics, nil)
			Expect(err).To(BeNil())
		})

		It("should create and restore a snapshot", func() {
			b.appliedIndex = oneIndex
			s := b.GetSnapshot()

			b.lastIndex = zeroIndex
			b.ApplySnapshot(*s)
			Expect(b.lastIndex).To(Equal(oneIndex))
			Expect(b.appliedIndex).To(Equal(oneIndex))
		})
	})

	Describe("Init()", func() {
		var (
			errorC    chan error
			fakeError error
			err       error
			b         *TopicSyncableBridge
			ctx       context.Context
			tick      time.Duration
		)

		JustBeforeEach(func() {
			errorC = make(chan error)
			fakeError = fmt.Errorf("fake error")
			fakeSyncable.TopicsReturns([]string{"foo"})
			b, err = new(bridgeName, fakeSyncable, topics, nil)
			Expect(err).To(BeNil())
			ctx = context.Background()
			tick = 2 * time.Millisecond
		})

		It("should error if syncable init fails", func(done Done) {
			fakeSyncable.InitReturns(fakeError)
			err = b.Init(context.Background(), errorC, tick)
			Expect(err).To(BeNil())

			err = <-errorC
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(ContainSubstring("Init of internal syncable failed"))
			close(done)
		}, 0.2)

		It("should error if topic can't create a new reader", func(done Done) {
			fooFakeTopic.NewReaderReturns(nil, fakeError)
			err = b.Init(context.Background(), errorC, tick)
			Expect(err).To(BeNil())

			err = <-errorC
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(ContainSubstring("Could not create reader"))
			close(done)
		}, 0.2)

		// TODO Cancelling isn't ready yet
		XIt("should error if ctx is done and the syncable has a close error", func(done Done) {
			fakeSyncable.CloseReturns(fakeError)
			ctx, cancel := context.WithCancel(ctx)

			go func() {
				_ = b.Init(ctx, errorC, tick)
			}()

			cancel()

			err = <-errorC
			Expect(err.Error()).To(ContainSubstring("Problem closing wrapped syncable"))
			close(done)
		}, 0.2)

		It("should err if next call on reader has an error", func(done Done) {
			fakeReader.NextReturns(nil, fakeError)

			go func() {
				_ = b.Init(ctx, errorC, tick)
			}()

			err = <-errorC
			Expect(err.Error()).To(ContainSubstring("Problem getting the next accepted proposal from topic"))
			// Once we read from errorC the bridge loops through and calls the next Next() before we do
			// the following check. This is why NextCallCount() equals 2 even though it seems like it should be 1
			Expect(fakeReader.NextCallCount()).To(Equal(2))

			close(done)
		}, 0.2)

		It("should err if the call to sync has an error", func(done Done) {
			fakeSyncable.SyncReturns(fakeError)

			go func() {
				_ = b.Init(ctx, errorC, tick)
			}()

			err = <-errorC
			Expect(err.Error()).To(ContainSubstring("Problem syncing"))
			Expect(fakeReader.NextCallCount()).To(Equal(2))

			close(done)
		}, 0.2)

		It("should update lastIndex if a successful sync occurs", func(done Done) {
			ap := &types.AcceptedProposal{Index: 1, Term: 1}
			fakeReader.NextReturnsOnCall(0, ap, nil)
			fakeReader.NextReturns(ap, fakeError)

			Expect(b.lastIndex).To(Equal(zeroIndex))

			go func() {
				_ = b.Init(ctx, errorC, tick)
			}()

			err = <-errorC
			Expect(err.Error()).To(ContainSubstring("Problem getting the next accepted proposal from topic"))
			Expect(b.lastIndex).To(Equal(oneIndex))
			Expect(fakeReader.NextCallCount()).To(Equal(3))

			close(done)
		}, 0.2)

		It("should stop syncing when node is no longer the leader", func(done Done) {
			ap := &types.AcceptedProposal{Index: 2}
			fakeReader.NextReturns(ap, nil)
			fakeLeader.IsLeaderReturnsOnCall(0, true)
			fakeLeader.IsLeaderReturnsOnCall(1, true)
			fakeLeader.IsLeaderReturns(false)
			fakeSyncable.CloseReturns(fakeError)

			go func() {
				_ = b.Init(ctx, errorC, tick)
			}()

			err = <-errorC
			Expect(err.Error()).To(ContainSubstring("Problem closing syncable"))
			Expect(fakeReader.NextCallCount()).To(Equal(1))
			close(done)
		}, 0.2)

		It("should restart syncing when node becomes the leader again", func(done Done) {
			fakeReader.NextReturns(nil, fakeError)
			fakeLeader.IsLeaderReturnsOnCall(0, true)
			fakeLeader.IsLeaderReturnsOnCall(1, false)
			fakeLeader.IsLeaderReturns(true)
			fakeSyncable.CloseReturns(fakeError)

			go func() {
				_ = b.Init(ctx, errorC, tick)
			}()

			err = <-errorC
			Expect(err.Error()).To(ContainSubstring("Problem closing syncable"))

			err = <-errorC
			Expect(err.Error()).To(ContainSubstring("Problem getting the next accepted proposal from topic"))
			Expect(fakeReader.NextCallCount()).To(Equal(2))
			close(done)
		}, 0.2)

		It("should update the applied index but not the lastIndex if this is the leader", func() {
			b.lastIndex = twoIndex
			Expect(b.appliedIndex).To(Equal(zeroIndex))
			b.UpdateIndex(oneIndex)
			Expect(b.appliedIndex).To(Equal(oneIndex))
			Expect(b.lastIndex).To(Equal(twoIndex))
		})

		It("should update the applied index and the lastIndex if this is not the leader", func() {
			fakeLeader.IsLeaderReturns(false)
			b.lastIndex = twoIndex
			Expect(b.appliedIndex).To(Equal(zeroIndex))
			b.UpdateIndex(oneIndex)
			Expect(b.appliedIndex).To(Equal(oneIndex))
			Expect(b.lastIndex).To(Equal(oneIndex))
			fakeLeader.IsLeaderReturns(true)
		})
	})
})
