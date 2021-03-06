package triggers

import (
	"fmt"
	"time"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

type baseTriggerEnvelope struct {
	Flow        *flows.FlowReference `json:"flow" validate:"required"`
	TriggeredOn time.Time            `json:"triggered_on" validate:"required"`
}

func ReadTrigger(session flows.Session, envelope *utils.TypedEnvelope) (flows.Trigger, error) {
	switch envelope.Type {

	case TypeManual:
		return ReadManualTrigger(session, envelope)
	case TypeFlowAction:
		return ReadFlowActionTrigger(session, envelope)

	default:
		return nil, fmt.Errorf("unknown trigger type: %s", envelope.Type)
	}
}

func readBaseTrigger(session flows.Session, base *baseTrigger, envelope *baseTriggerEnvelope) error {
	var err error

	base.triggeredOn = envelope.TriggeredOn

	if base.flow, err = session.Assets().GetFlow(envelope.Flow.UUID); err != nil {
		return err
	}

	return nil
}
