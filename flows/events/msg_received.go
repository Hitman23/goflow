package events

import (
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/inputs"
)

// TypeMsgReceived is a constant for incoming messages
const TypeMsgReceived string = "msg_received"

// MsgReceivedEvent events are used for starting flows or resuming flows which are waiting for a message.
// They represent an MO message for a contact.
//
// ```
//   {
//     "type": "msg_received",
//     "created_on": "2006-01-02T15:04:05Z",
//     "msg_uuid": "2d611e17-fb22-457f-b802-b8f7ec5cda5b",
//     "urn": "tel:+12065551212",
//     "channel": {"uuid": "61602f3e-f603-4c70-8a8f-c477505bf4bf", "name": "Twilio"},
//     "contact": {"uuid": "0e06f977-cbb7-475f-9d0b-a0c4aaec7f6a", "name": "Bob"},
//     "text": "hi there",
//     "attachments": ["https://s3.amazon.com/mybucket/attachment.jpg"]
//   }
// ```
//
// @event msg_received
type MsgReceivedEvent struct {
	BaseEvent
	MsgUUID     flows.InputUUID         `json:"msg_uuid" validate:"required,uuid4"`
	Channel     *flows.ChannelReference `json:"channel,omitempty"`
	URN         urns.URN                `json:"urn" validate:"urn"`
	Contact     *flows.ContactReference `json:"contact"`
	Text        string                  `json:"text"`
	Attachments []flows.Attachment      `json:"attachments,omitempty"`
}

// NewMsgReceivedEvent creates a new incoming msg event for the passed in channel, contact and string
func NewMsgReceivedEvent(uuid flows.InputUUID, channel *flows.ChannelReference, contact *flows.ContactReference, urn urns.URN, text string, attachments []flows.Attachment) *MsgReceivedEvent {
	return &MsgReceivedEvent{
		BaseEvent:   NewBaseEvent(),
		MsgUUID:     uuid,
		Channel:     channel,
		Contact:     contact,
		URN:         urn,
		Text:        text,
		Attachments: attachments,
	}
}

// Type returns the type of this event
func (e *MsgReceivedEvent) Type() string { return TypeMsgReceived }

// Apply applies this event to the given run
func (e *MsgReceivedEvent) Apply(run flows.FlowRun) error {
	var channel flows.Channel
	var err error

	if e.Channel != nil {
		channel, err = run.Session().Assets().GetChannel(e.Channel.UUID)
		if err != nil {
			return err
		}
	}

	// update this run's input
	run.SetInput(inputs.NewMsgInput(e.MsgUUID, channel, e.CreatedOn(), e.URN, e.Text, e.Attachments))

	run.ResetExpiration(nil)
	return nil
}
