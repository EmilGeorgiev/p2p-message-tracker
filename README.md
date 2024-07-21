## Description

This is a peer-to-peer (p2p) message tracker. There is a `Message` type that is found in `network/message.go`.  
```go 
// Message is received from peers in a p2p network.
type Message struct {
	ID     string
	PeerID string
	Data   []byte
}
```
Each message is uniquely identified by the `Message.ID`. Messages with the same ID may be received by multiple peers.  Peers are uniquely identified by their own ID stored in `Message.PeerID`. 

The interface for the message tracker is defined in `network/message_tracker.go`.  
```go 
// MessageTracker tracks a configurable fixed amount of messages.
// Messages are stored first-in-first-out.  Duplicate messages should not be stored in the queue.
type MessageTracker interface {
	// Add will add a message to the tracker
	Add(message *Message) (err error)
	// Delete will delete message from tracker
	Delete(id string) (err error)
	// Get returns a message for a given ID.  Message is retained in tracker
	Message(id string) (message *Message, err error)
	// All returns messages in the order in which they were received
	Messages() (messages []*Message)
}
```

There is an exported constructor `network.NewMessageTracker(length int)` which accepts a length parameter.  This parameter is used to constrain the number of messages the tracker can track.

There are a few key points:
- Duplicate messages based on `Message.ID` are returned by `MessageTracker.All()` once.
- The tracker holds only a configurable maximum amount of messages so it does not grow in size indefinitely.
