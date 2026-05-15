package event

import "sync"

// Bus is a simple event bus using Go channels
type Bus struct {
	mu       sync.RWMutex
	handlers map[string][]chan interface{}
}

func NewBus() *Bus {
	return &Bus{
		handlers: make(map[string][]chan interface{}),
	}
}

// Subscribe registers a handler for an event type
func (b *Bus) Subscribe(eventType string) chan interface{} {
	ch := make(chan interface{}, 100)
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventType] = append(b.handlers[eventType], ch)
	return ch
}

// Publish sends an event to all subscribers
func (b *Bus) Publish(eventType string, data interface{}) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, ch := range b.handlers[eventType] {
		select {
		case ch <- data:
		default:
			// Channel full, skip
		}
	}
}

// Unsubscribe removes a subscriber
func (b *Bus) Unsubscribe(eventType string, ch chan interface{}) {
	b.mu.Lock()
	defer b.mu.Unlock()
	handlers := b.handlers[eventType]
	for i, h := range handlers {
		if h == ch {
			b.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
			close(ch)
			break
		}
	}
}

// Event types
const (
	EventDeviceOnline       = "device:online"
	EventDeviceOffline      = "device:offline"
	EventDeviceRegister     = "device:register"
	EventStreamStart        = "stream:start"
	EventStreamStop         = "stream:stop"
	EventStreamChange       = "stream:change"
	EventAlarmReceived      = "alarm:received"
	EventCatalogUpdated     = "catalog:updated"
	EventMediaServerOnline  = "media:server:online"
	EventMediaServerOffline = "media:server:offline"
)

// StreamEvent is published when a stream state changes
type StreamEvent struct {
	Stream        string
	DeviceID      string
	ChannelID     string
	App           string
	MediaServerID string
	SSRC          string
	Online        bool
}

// DeviceEvent is published for device-related events
type DeviceEvent struct {
	DeviceID string
	IP       string
	Port     int
}

// AlarmEvent is published when an alarm is received
type AlarmEvent struct {
	DeviceID  string
	ChannelID string
	Type      string
}
