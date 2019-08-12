package bridge

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"

	"github.com/philborlin/committed/syncable"
	"github.com/philborlin/committed/topic"
	"github.com/pkg/errors"
)

// Bridge manages the interactions between a topic and a syncable
type Bridge struct {
	Name      string
	Syncable  syncable.Syncable
	topics    map[string]topic.Topic
	lastIndex uint64
}

// Snapshot is the snapshot struct
type Snapshot struct {
	LastIndex uint64
}

// NewBridge creates a wrapper that will
func NewBridge(name string, s syncable.Syncable, topics map[string]topic.Topic) (*Bridge, error) {
	if len(s.Topics()) == 0 {
		return nil, fmt.Errorf("[%s.bridge] No topics so there is nothing to sync", name)
	}

	if len(s.Topics()) > 1 {
		// There is going to be some serious syncronization work that needs to happen to support multiple
		// topics. Deferring this until later.
		return nil, fmt.Errorf("[%s.bridge] We don't support more than one topic in a syncable yet", name)
	}

	// Create a map that only has entries for the topics we are listening to
	var tmap = make(map[string]topic.Topic)
	for _, item := range s.Topics() {
		t, ok := topics[item]
		if !ok {
			return nil, fmt.Errorf("syncable %s is trying to listen to topic %s which does not exist", name, item)
		}
		tmap[item] = t
	}

	return &Bridge{Name: name, Syncable: s, topics: tmap}, nil
}

// GetSnapshot implements Snapshotter
func (b *Bridge) GetSnapshot() ([]byte, error) {
	s := &Snapshot{LastIndex: b.lastIndex}
	var buf bytes.Buffer
	_ = gob.NewEncoder(&buf).Encode(s)

	return buf.Bytes(), nil
}

// ApplySnapshot implements Snapshotter
func (b *Bridge) ApplySnapshot(snap []byte) error {
	var s Snapshot
	dec := gob.NewDecoder(bytes.NewBuffer(snap))
	if err := dec.Decode(&s); err != nil {
		return errors.Wrap(err, "Could not decode snapshot")
	}

	b.lastIndex = s.LastIndex

	return nil
}

// Init implements Syncable
// To close the syncable send a message to the ctx.Done() channel
// It is the caller's responsibility to listen to any errors on the errorC channel passed in
func (b *Bridge) Init(ctx context.Context, errorC chan<- error) error {
	err := b.Syncable.Init(ctx)
	if err != nil {
		return errors.Wrapf(err, "[%s.bridge] Init of internal syncable failed", b.Name)
	}

	for _, t := range b.topics {
		reader, err := t.NewReader(0)
		if err != nil {
			return errors.Wrapf(err, "[%s.bridge] Could not create reader", b.Name)
		}

		go func(t topic.Topic) {
			for {
				select {
				// TODO In order for done to work,
				// walTopicReader has to implement it also since calls to next just block
				// case <-ctx.Done():
				// err := b.Syncable.Close()
				// if err != nil {
				// 	errorC <- errors.Wrapf(err, "[%s.bridge] Problem closing wrapped syncable", b.Name)
				// }
				// return
				default:
					ap, err := reader.Next(ctx)
					if err != nil {
						errorC <- errors.Wrapf(err,
							"[%s.bridge] Problem getting the next accepted proposal from topic %s", b.Name, t.Name())
						continue
					}
					if err := b.Syncable.Sync(ctx, ap); err != nil {
						errorC <- errors.Wrapf(err, "[%s.bridge] Problem syncing", b.Name)
						continue
					}
					b.lastIndex = ap.Index
				}
			}
		}(t)
	}

	return nil
}