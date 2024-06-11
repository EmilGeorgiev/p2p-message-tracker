package network_test

import (
	"fmt"
	"testing"

	"github.com/ChainSafe/gossamer-go-interview/network"
	"github.com/stretchr/testify/assert"
)

func generateMessage(n int) *network.Message {
	return &network.Message{
		ID:     fmt.Sprintf("someID%d", n),
		PeerID: fmt.Sprintf("somePeerID%d", n),
		Data:   []byte{0, 1, 1},
	}
}

func TestMessageTracker_Add(t *testing.T) {
	t.Run("add, get, then all messages", func(t *testing.T) {
		length := 5
		mt := network.NewMessageTracker(length)

		for i := 0; i < length; i++ {
			err := mt.Add(generateMessage(i))
			assert.NoError(t, err)

			msg, err := mt.Message(generateMessage(i).ID)
			assert.NoError(t, err)
			assert.NotNil(t, msg)
		}

		msgs := mt.Messages()
		assert.Equal(t, []*network.Message{
			generateMessage(0),
			generateMessage(1),
			generateMessage(2),
			generateMessage(3),
			generateMessage(4),
		}, msgs)
	})

	t.Run("add, get, then all messages, delete some", func(t *testing.T) {
		length := 5
		mt := network.NewMessageTracker(length)

		for i := 0; i < length; i++ {
			err := mt.Add(generateMessage(i))
			assert.NoError(t, err)

			msg, err := mt.Message(generateMessage(i).ID)
			assert.NoError(t, err)
			assert.NotNil(t, msg)
		}

		msgs := mt.Messages()
		assert.Equal(t, []*network.Message{
			generateMessage(0),
			generateMessage(1),
			generateMessage(2),
			generateMessage(3),
			generateMessage(4),
		}, msgs)

		for i := 0; i < length-2; i++ {
			err := mt.Delete(generateMessage(i).ID)
			assert.NoError(t, err)
		}

		msgs = mt.Messages()
		assert.Equal(t, []*network.Message{
			generateMessage(3),
			generateMessage(4),
		}, msgs)

	})

	t.Run("not full, with duplicates", func(t *testing.T) {
		length := 5
		mt := network.NewMessageTracker(length)

		for i := 0; i < length-1; i++ {
			err := mt.Add(generateMessage(i))
			assert.NoError(t, err)
		}
		for i := 0; i < length-1; i++ {
			err := mt.Add(generateMessage(length - 2))
			assert.NoError(t, err)
		}

		msgs := mt.Messages()
		assert.Equal(t, []*network.Message{
			generateMessage(0),
			generateMessage(1),
			generateMessage(2),
			generateMessage(3),
		}, msgs)
	})

	t.Run("not full, with duplicates from other peers", func(t *testing.T) {
		length := 5
		mt := network.NewMessageTracker(length)

		for i := 0; i < length-1; i++ {
			err := mt.Add(generateMessage(i))
			assert.NoError(t, err)
		}
		for i := 0; i < length-1; i++ {
			msg := generateMessage(length - 2)
			msg.PeerID = "somePeerID0"
			err := mt.Add(msg)
			assert.NoError(t, err)
		}

		msgs := mt.Messages()
		assert.Equal(t, []*network.Message{
			generateMessage(0),
			generateMessage(1),
			generateMessage(2),
			generateMessage(3),
		}, msgs)
	})
}

func TestMessageTracker_Cleanup(t *testing.T) {
	t.Run("overflow and cleanup", func(t *testing.T) {
		length := 5
		mt := network.NewMessageTracker(length)

		for i := 0; i < length*2; i++ {
			err := mt.Add(generateMessage(i))
			assert.NoError(t, err)
		}

		msgs := mt.Messages()
		assert.Equal(t, []*network.Message{
			generateMessage(5),
			generateMessage(6),
			generateMessage(7),
			generateMessage(8),
			generateMessage(9),
		}, msgs)
	})

	t.Run("overflow and cleanup with duplicate", func(t *testing.T) {
		length := 5
		mt := network.NewMessageTracker(length)

		for i := 0; i < length*2; i++ {
			err := mt.Add(generateMessage(i))
			assert.NoError(t, err)
		}

		for i := length; i < length*2; i++ {
			err := mt.Add(generateMessage(i))
			assert.NoError(t, err)
		}

		msgs := mt.Messages()
		assert.Equal(t, []*network.Message{
			generateMessage(5),
			generateMessage(6),
			generateMessage(7),
			generateMessage(8),
			generateMessage(9),
		}, msgs)
	})

	t.Run("duplicated messages are moved to the front", func(t *testing.T) {
		length := 5
		mt := network.NewMessageTracker(length)

		assert.NoError(t, mt.Add(generateMessage(0)))
		assert.NoError(t, mt.Add(generateMessage(1)))
		assert.NoError(t, mt.Add(generateMessage(2)))
		assert.NoError(t, mt.Add(generateMessage(3)))
		assert.NoError(t, mt.Add(generateMessage(4)))
		assert.NoError(t, mt.Add(generateMessage(0)))
		assert.NoError(t, mt.Add(generateMessage(1)))

		msgs := mt.Messages()
		assert.Equal(t, []*network.Message{
			generateMessage(2),
			generateMessage(3),
			generateMessage(4),
			generateMessage(0),
			generateMessage(1),
		}, msgs)
	})
}

func TestMessageTracker_Delete(t *testing.T) {
	t.Run("empty tracker", func(t *testing.T) {
		length := 5
		mt := network.NewMessageTracker(length)
		err := mt.Delete("bleh")
		assert.ErrorIs(t, err, network.ErrMessageNotFound)
	})
}

func TestMessageTracker_Message(t *testing.T) {
	t.Run("empty tracker", func(t *testing.T) {
		length := 5
		mt := network.NewMessageTracker(length)
		msg, err := mt.Message("bleh")
		assert.ErrorIs(t, err, network.ErrMessageNotFound)
		assert.Nil(t, msg)
	})
}

// createTestMessages generates a list of test messages
func createTestMessages(n int) []*network.Message {
	messages := make([]*network.Message, n)
	for i := 0; i < n; i++ {
		messages[i] = generateMessage(i)
	}
	return messages
}

func BenchmarkAdd(b *testing.B) {
	tracker := network.NewMessageTracker(1000)
	messages := createTestMessages(b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tracker.Add(messages[i])
	}
}

func BenchmarkDelete(b *testing.B) {
	tracker := network.NewMessageTracker(1000)
	messages := createTestMessages(b.N)

	for i := 0; i < b.N; i++ {
		tracker.Add(messages[i])
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tracker.Delete(messages[i].ID)
	}
}

func BenchmarkMessage(b *testing.B) {
	tracker := network.NewMessageTracker(1000)
	messages := createTestMessages(b.N)

	for i := 0; i < b.N; i++ {
		tracker.Add(messages[i])
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tracker.Message(messages[i].ID)
	}
}

func BenchmarkMessages(b *testing.B) {
	tracker := network.NewMessageTracker(1000)
	messages := createTestMessages(1000)

	for i := 0; i < 1000; i++ {
		tracker.Add(messages[i])
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tracker.Messages()
	}
}
