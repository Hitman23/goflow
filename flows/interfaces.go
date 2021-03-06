package flows

import (
	"time"

	"github.com/nyaruka/goflow/utils"
)

type UUID string

type NodeUUID UUID

func (u NodeUUID) String() string { return string(u) }

type ExitUUID UUID

func (u ExitUUID) String() string { return string(u) }

type FlowUUID UUID

func (u FlowUUID) String() string { return string(u) }

type ActionUUID UUID

func (u ActionUUID) String() string { return string(u) }

type ContactUUID UUID

func (u ContactUUID) String() string { return string(u) }

type ChannelUUID UUID

func (u ChannelUUID) String() string { return string(u) }

type RunUUID UUID

func (u RunUUID) String() string { return string(u) }

type StepUUID UUID

func (u StepUUID) String() string { return string(u) }

type LabelUUID UUID

func (u LabelUUID) String() string { return string(u) }

type GroupUUID UUID

func (u GroupUUID) String() string { return string(u) }

type InputUUID UUID

func (u InputUUID) String() string { return string(u) }

// RunStatus represents the current status of the engine session
type SessionStatus string

const (
	// SessionStatusActive represents a session that is still active
	SessionStatusActive SessionStatus = "active"

	// SessionStatusCompleted represents a session that has run to completion
	SessionStatusCompleted SessionStatus = "completed"

	// SessionStatusWaiting represents a session which is waiting for something from the caller
	SessionStatusWaiting SessionStatus = "waiting"

	// SessionStatusErrored represents a session that encountered an error
	SessionStatusErrored SessionStatus = "errored"
)

func (r SessionStatus) String() string { return string(r) }

// RunStatus represents the current status of the flow run
type RunStatus string

const (
	// RunStatusActive represents a run that is still active
	RunStatusActive RunStatus = "active"

	// RunStatusCompleted represents a run that has run to completion
	RunStatusCompleted RunStatus = "completed"

	// RunStatusWaiting represents a run which is waiting for something from the caller
	RunStatusWaiting RunStatus = "waiting"

	// RunStatusErrored represents a run that encountered an error
	RunStatusErrored RunStatus = "errored"

	// RunStatusExpired represents a run that expired due to inactivity
	RunStatusExpired RunStatus = "expired"

	// RunStatusInterrupted represents a run that was interrupted by another flow
	RunStatusInterrupted RunStatus = "interrupted"
)

func (r RunStatus) String() string { return string(r) }

type SessionAssets interface {
	GetChannel(ChannelUUID) (Channel, error)

	GetField(FieldKey) (*Field, error)
	GetFieldSet() (*FieldSet, error)

	GetFlow(FlowUUID) (Flow, error)

	GetGroup(GroupUUID) (*Group, error)
	GetGroupSet() (*GroupSet, error)

	GetLabel(LabelUUID) (*Label, error)
	GetLabelSet() (*LabelSet, error)

	HasLocations() bool
	GetLocationHierarchy() (*utils.LocationHierarchy, error)
}

type Flow interface {
	UUID() FlowUUID
	Name() string
	Language() utils.Language
	ExpireAfterMinutes() int
	Translations() FlowTranslations

	Validate(SessionAssets) error
	Nodes() []Node
	GetNode(uuid NodeUUID) Node

	Reference() *FlowReference
}

type Node interface {
	UUID() NodeUUID

	Router() Router
	Actions() []Action
	Exits() []Exit
	Wait() Wait
}

type EventLog interface {
	Add(Event)
	Events() []Event
}

type Action interface {
	UUID() ActionUUID

	Execute(FlowRun, Step, EventLog) error
	Validate(SessionAssets) error
	utils.Typed
}

type Router interface {
	PickRoute(FlowRun, []Exit, Step) (interface{}, Route, error)
	Validate([]Exit) error
	ResultName() string
	utils.Typed
}

type Route struct {
	exit  ExitUUID
	match string
}

func (r Route) Exit() ExitUUID { return r.exit }
func (r Route) Match() string  { return r.match }

var NoRoute = Route{}

func NewRoute(exit ExitUUID, match string) Route {
	return Route{exit, match}
}

type Exit interface {
	UUID() ExitUUID
	DestinationNodeUUID() NodeUUID
	Name() string
}

type Wait interface {
	utils.Typed

	Begin(FlowRun, Step)
	CanResume(FlowRun, Step) bool
	HasTimedOut() bool

	Resume(FlowRun)
	ResumeByTimeOut(FlowRun)
}

// FlowTranslations provide a way to get the Translations for a flow for a specific language
type FlowTranslations interface {
	GetLanguageTranslations(utils.Language) Translations
}

// Translations provide a way to get the translation for a specific language for a uuid/key pair
type Translations interface {
	GetTextArray(uuid UUID, key string) []string
}

type Trigger interface {
	utils.Typed

	TriggeredOn() time.Time
	Flow() Flow
}

type Event interface {
	CreatedOn() time.Time
	SetCreatedOn(time.Time)

	FromCaller() bool
	SetFromCaller(bool)

	Apply(FlowRun) error

	utils.Typed
}

type Input interface {
	utils.VariableResolver
	utils.Typed

	UUID() InputUUID
	CreatedOn() time.Time
	Channel() Channel
}

type Step interface {
	utils.VariableResolver

	UUID() StepUUID
	NodeUUID() NodeUUID
	ExitUUID() ExitUUID

	ArrivedOn() time.Time
	LeftOn() *time.Time

	Leave(ExitUUID)

	Events() []Event
}

// LogEntry is a container for a new event generated by the engine, i.e. not from the caller
type LogEntry interface {
	Step() Step
	Action() Action
	Event() Event
}

// Session represents the session of a flow run which may contain many runs
type Session interface {
	Assets() SessionAssets

	Environment() utils.Environment
	SetEnvironment(utils.Environment)

	Contact() *Contact
	SetContact(*Contact)

	Status() SessionStatus
	Trigger() Trigger
	PushFlow(Flow, FlowRun)
	Wait() Wait
	FlowOnStack(FlowUUID) bool

	Start(Trigger, []Event) error
	Resume([]Event) error
	Runs() []FlowRun
	GetRun(RunUUID) (FlowRun, error)
	GetCurrentChild(FlowRun) FlowRun

	Log() []LogEntry
	LogEvent(Step, Action, Event)
	ClearLog()
}

// RunSummary represents the minimum information available about all runs (current or related) and is the
// representation of runs made accessible to router tests.
type RunSummary interface {
	UUID() RunUUID
	Contact() *Contact
	Flow() Flow
	Status() RunStatus
	Results() Results
}

// FlowRun represents a run in the current session
type FlowRun interface {
	RunSummary

	Environment() utils.Environment
	Session() Session
	Context() utils.VariableResolver
	Input() Input
	Extra() utils.JSONFragment
	Webhook() *utils.RequestResponse

	SetContact(*Contact)
	SetInput(Input)
	SetStatus(RunStatus)
	SetWebhook(*utils.RequestResponse)
	SetExtra(utils.JSONFragment)

	ApplyEvent(Step, Action, Event) error
	AddError(Step, Action, error)
	AddFatalError(Step, Action, error)

	CreateStep(Node) Step
	Path() []Step
	PathLocation() (Step, Node, error)

	GetText(uuid UUID, key string, native string) string
	GetTextArray(uuid UUID, key string, native []string) []string

	Snapshot() RunSummary
	Parent() RunSummary
	SessionParent() FlowRun
	Ancestors() []FlowRun

	CreatedOn() time.Time
	ExpiresOn() *time.Time
	ResetExpiration(*time.Time)
	ExitedOn() *time.Time
	Exit(RunStatus)
}

// ChannelType represents the type of a Channel
type ChannelType string

func (ct ChannelType) String() string { return string(ct) }

// Channel represents a channel for sending and receiving messages
type Channel interface {
	UUID() ChannelUUID
	Name() string
	Address() string
	Type() ChannelType
}

// MsgDirection is the direction of a Msg (either in or out)
type MsgDirection string

const (
	// MsgOut represents an outgoing message
	MsgOut MsgDirection = "O"

	// MsgIn represents an incoming message
	MsgIn MsgDirection = "I"
)
