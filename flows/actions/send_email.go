package actions

import (
	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

// TypeSendEmail is our type for the email action
const TypeSendEmail string = "send_email"

// EmailAction can be used to send an email to one or more recipients. The subject, body and emails can all contain templates.
//
// A `send_email` event will be created for each email address.
//
// ```
//   {
//     "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "type": "send_email",
//     "subject": "Here is your activation token",
//     "body": "Your activation token is @contact.fields.activation_token",
//     "emails": ["@contact.urns.email"]
//   }
// ```
//
// @action send_email
type EmailAction struct {
	BaseAction
	Emails  []string `json:"emails"   validate:"required,min=1"`
	Subject string   `json:"subject"  validate:"required"`
	Body    string   `json:"body"     validate:"required"`
}

// Type returns the type of this action
func (a *EmailAction) Type() string { return TypeSendEmail }

// Validate validates the fields on this action
func (a *EmailAction) Validate() error {
	return utils.ValidateAll(a)
}

// Execute creates the email events
func (a *EmailAction) Execute(run flows.FlowRun, step flows.Step) error {
	for _, email := range a.Emails {
		email, err := excellent.EvaluateTemplateAsString(run.Environment(), run.Context(), email)
		if err != nil {
			run.AddError(step, err)
		}

		subject, err := excellent.EvaluateTemplateAsString(run.Environment(), run.Context(), a.Subject)
		if err != nil {
			run.AddError(step, err)
		}

		body, err := excellent.EvaluateTemplateAsString(run.Environment(), run.Context(), a.Body)
		if err != nil {
			run.AddError(step, err)
		}
		run.AddEvent(step, events.NewSendEmailEvent(email, subject, body))
	}
	return nil
}
