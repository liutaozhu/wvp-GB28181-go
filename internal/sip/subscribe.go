package sip

import (
	"sync"
	"time"

	"wvp-pro-go/internal/config"

	"go.uber.org/zap"
)

// EventResult contains the result of a SIP subscription event
type EventResult struct {
	StatusCode int
	Msg        string
	Source     string // source IP
	Raw        string // raw SIP message
}

// SipEvent represents a pending SIP response subscription
type SipEvent struct {
	Key       string
	OkChan    chan *EventResult
	ErrorChan chan *EventResult
	Timeout   time.Duration
	Timer     *time.Timer
}

// Subscribe manages SIP response subscriptions with timeout
type Subscribe struct {
	mu         sync.RWMutex
	events     map[string]*SipEvent
	log        *zap.Logger
	cfg        config.UserSettingConfig
}

func NewSubscribe(log *zap.Logger, cfg config.UserSettingConfig) *Subscribe {
	return &Subscribe{
		events: make(map[string]*SipEvent),
		log:    log,
		cfg:    cfg,
	}
}

// AddSubscribe adds a new subscription with callback channels
func (s *Subscribe) AddSubscribe(key string, timeout time.Duration) *SipEvent {
	event := &SipEvent{
		Key:       key,
		OkChan:    make(chan *EventResult, 1),
		ErrorChan: make(chan *EventResult, 1),
		Timeout:   timeout,
	}

	s.mu.Lock()
	s.events[key] = event
	s.mu.Unlock()

	// Start timeout timer
	event.Timer = time.AfterFunc(timeout, func() {
		s.mu.Lock()
		delete(s.events, key)
		s.mu.Unlock()

		select {
		case event.ErrorChan <- &EventResult{
			StatusCode: 408,
			Msg:        "请求超时",
		}:
		default:
		}
		s.log.Debug("sip event timeout", zap.String("key", key))
	})

	return event
}

// RemoveSubscribe removes a subscription
func (s *Subscribe) RemoveSubscribe(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if event, ok := s.events[key]; ok {
		if event.Timer != nil {
			event.Timer.Stop()
		}
		delete(s.events, key)
	}
}

// Notify sends a notification to a subscription
func (s *Subscribe) Notify(key string, result *EventResult) {
	s.mu.RLock()
	event, ok := s.events[key]
	s.mu.RUnlock()

	if !ok {
		return
	}

	s.RemoveSubscribe(key)

	if result.StatusCode >= 200 && result.StatusCode < 300 {
		select {
		case event.OkChan <- result:
		default:
		}
	} else {
		select {
		case event.ErrorChan <- result:
		default:
		}
	}
}

// WaitForResult waits for a SIP event result with timeout (replaces Java DeferredResult)
func (s *Subscribe) WaitForResult(event *SipEvent) (*EventResult, error) {
	select {
	case result := <-event.OkChan:
		return result, nil
	case result := <-event.ErrorChan:
		return result, nil
	case <-time.After(event.Timeout):
		return &EventResult{StatusCode: 408, Msg: "请求超时"}, nil
	}
}

// Count returns the number of active subscriptions
func (s *Subscribe) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.events)
}
