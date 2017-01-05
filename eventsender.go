package eventsender

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type EventSender struct {
	buff *Buffer
	err  error
}

// New creates a new http client event sender
func New(method, url string) *EventSender {
	sender := &EventSender{NewBuffer(), nil}
	go func() {
		client := http.Client{}
		req, err := http.NewRequest(method, url, sender.buff)
		if err != nil {
			sender.err = err
			return
		}
		resp, err := client.Do(req)
		if err != nil {
			sender.err = err
		}
		if resp.StatusCode != http.StatusOK {
			sender.err = errors.New(fmt.Sprintf("http status %v", resp.StatusCode))
		}
	}()
	return sender
}

// SendEvent sends a event to the server
func (sender *EventSender) SendEvent(value interface{}) error {
	bs, err := json.Marshal(value)
	if err != nil {
		return err
	}
	bs = append(bs, '\n')
	_, err = sender.buff.Write(bs)
	return err
}

// Close closes the event request
func (sender *EventSender) Close() error {
	return sender.buff.Close()
}

// Error returns any occured http error (in request creation or submission)
func (sender *EventSender) Error() error {
	return sender.err
}
