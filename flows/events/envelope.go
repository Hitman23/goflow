package events

import (
	"fmt"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

func ReadEvents(envelopes []*utils.TypedEnvelope) ([]flows.Event, error) {
	events := make([]flows.Event, len(envelopes))
	for e, envelope := range envelopes {
		event, err := EventFromEnvelope(envelope)
		if err != nil {
			return nil, err
		}
		event.SetFromCaller(true)
		events[e] = event
	}
	return events, nil
}

func EventFromEnvelope(envelope *utils.TypedEnvelope) (flows.Event, error) {
	var event flows.Event

	switch envelope.Type {
	case TypeAddLabel:
		event = &AddLabelEvent{}
	case TypeAddToGroup:
		event = &AddToGroupEvent{}
	case TypeAddURN:
		event = &AddURNEvent{}
	case TypeSendEmail:
		event = &SendEmailEvent{}
	case TypeError:
		event = &ErrorEvent{}
	case TypeFlowTriggered:
		event = &FlowTriggeredEvent{}
	case TypeSessionTriggered:
		event = &SessionTriggeredEvent{}
	case TypeRunExpired:
		event = &RunExpiredEvent{}
	case TypeMsgReceived:
		event = &MsgReceivedEvent{}
	case TypeSendMsg:
		event = &SendMsgEvent{}
	case TypeMsgWait:
		event = &MsgWaitEvent{}
	case TypeNothingWait:
		event = &NothingWaitEvent{}
	case TypeRemoveFromGroup:
		event = &RemoveFromGroupEvent{}
	case TypeSaveFlowResult:
		event = &SaveFlowResultEvent{}
	case TypeSaveContactField:
		event = &SaveContactFieldEvent{}
	case TypePreferredChannel:
		event = &PreferredChannelEvent{}
	case TypeSetEnvironment:
		event = &SetEnvironmentEvent{}
	case TypeSetExtra:
		event = &SetExtraEvent{}
	case TypeSetContact:
		event = &SetContactEvent{}
	case TypeUpdateContact:
		event = &UpdateContactEvent{}
	case TypeWebhookCalled:
		event = &WebhookCalledEvent{}
	default:
		return nil, fmt.Errorf("Unknown event type: %s", envelope.Type)
	}

	return event, utils.UnmarshalAndValidate(envelope.Data, event, fmt.Sprintf("event[type=%s]", envelope.Type))
}
