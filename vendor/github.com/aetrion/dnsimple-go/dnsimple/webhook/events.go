package webhook

import (
	"github.com/aetrion/dnsimple-go/dnsimple"
)

func switchEvent(name string, payload []byte) (Event, error) {
	var event Event

	switch name {
	case "domain.create":
		event = &DomainEvent{}
	case "domain.delete":
		event = &DomainEvent{}
	case "domain.token_reset":
		event = &DomainEvent{}
	case "domain.auto_renew_enable":
		event = &DomainEvent{}
	case "domain.auto_renew_disable":
		event = &DomainEvent{}
	case "webhook.create":
		event = &WebhookEvent{}
	default:
		event = &GenericEvent{}
	}

	return event, event.parse(payload)
}

//
// GenericEvent represents a generic event, where the data is a simple map of strings.
//
type GenericEvent struct {
	Event_Header
	Data interface{} `json:"data"`
}

func (e *GenericEvent) parse(payload []byte) error {
	e.payload = payload
	return unmashalEvent(payload, e)
}

// ParseGenericEvent unpacks the data into a GenericEvent.
func ParseGenericEvent(e *GenericEvent, payload []byte) error {
	return e.parse(payload)
}

//
// DomainEvent represents the base event sent for a domain action.
//
type DomainEvent struct {
	Event_Header
	Data   *DomainEvent     `json:"data"`
	Domain *dnsimple.Domain `json:"domain"`
}

// ParseDomainEvent unpacks the payload into a DomainEvent.
func ParseDomainEvent(e *DomainEvent, payload []byte) error {
	return e.parse(payload)
}

func (e *DomainEvent) parse(payload []byte) error {
	e.payload, e.Data = payload, e
	return unmashalEvent(payload, e)
}

//
// Webhook represents a generic event, where the data is a simple map of strings.
//
type WebhookEvent struct {
	Event_Header
	Data *WebhookEvent `json:"data"`
}

// ParseWebhookEvent unpacks the data into a WebhookEvent.
func ParseWebhookEvent(e *WebhookEvent, payload []byte) error {
	return e.parse(payload)
}

func (e *WebhookEvent) parse(payload []byte) error {
	e.payload, e.Data = payload, e
	return unmashalEvent(payload, e)
}
