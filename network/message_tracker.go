package network

import (
	"container/list"
	"errors"
	"sync"
)

// MessageTracker tracks a configurable fixed amount of messages.
// Messages are stored first-in-first-out.  Duplicate messages should not be stored in the queue.
type MessageTracker interface {
	// Add will add a message to the tracker, deleting the oldest message if necessary
	Add(message *Message) (err error)

	// Delete will delete message from tracker
	Delete(id string) (err error)

	// Message returns a message for a given ID.  Message is retained in tracker.
	Message(id string) (message *Message, err error)

	// Messages returns messages in FIFO order.
	Messages() (messages []*Message)
}

// ErrMessageNotFound is an error returned by MessageTracker when a message with specified id is not found
var ErrMessageNotFound = errors.New("message not found")

type messageTracker struct {
	sync.RWMutex

	capacity   int
	linkedList *list.List
	mapping    map[string]*list.Element
}

// Add adds a message to the tracker, deleting the oldest message if the capacity is overloaded. It is concurrently safe.
func (m *messageTracker) Add(message *Message) error {
	m.Lock()
	defer m.Unlock()

	element, exists := m.mapping[message.ID]
	if exists {
		m.linkedList.MoveToFront(element)
		return nil
	}

	if m.linkedList.Len() >= m.capacity {
		oldestElement := m.linkedList.Back()
		oldestMessage := oldestElement.Value.(*Message)
		m.linkedList.Remove(oldestElement)
		delete(m.mapping, oldestMessage.ID)
	}

	element = m.linkedList.PushFront(message)
	m.mapping[message.ID] = element
	return nil
}

// Delete removes a message from the tracker by ID.
func (m *messageTracker) Delete(id string) error {
	m.Lock()
	defer m.Unlock()
	element, exist := m.mapping[id]
	if !exist {
		return ErrMessageNotFound
	}

	m.linkedList.Remove(element)
	delete(m.mapping, id)
	return nil
}

// Message retrieves a message by ID
func (m *messageTracker) Message(id string) (*Message, error) {
	m.RLock()
	defer m.RUnlock()
	element, exists := m.mapping[id]
	if !exists {
		return nil, ErrMessageNotFound
	}

	return element.Value.(*Message), nil
}

// Messages returns all messages in FIFO order
func (m *messageTracker) Messages() []*Message {
	m.RLock()
	defer m.RUnlock()
	messages := make([]*Message, m.linkedList.Len())
	i := 0
	for e := m.linkedList.Back(); e != nil; e = e.Prev() {
		messages[i] = e.Value.(*Message)
		i++
	}
	return messages
}

// NewMessageTracker creates a new message tracker with a specified capacity
func NewMessageTracker(capacity int) MessageTracker {
	return &messageTracker{
		capacity:   capacity,
		linkedList: list.New(),
		mapping:    make(map[string]*list.Element),
	}
}
