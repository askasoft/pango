package freshdesk

type AutomationType int
type AutomationOperator string
type AutomationPerformerType int
type AutomationMatchType string
type AutomationActionType string
type AutomationPushTo string
type AutomationResource string

const (
	AutomationTypeTicketCreation AutomationType = 1
	AutomationTypeTimeTriggers   AutomationType = 3
	AutomationTypeTicketUpdates  AutomationType = 4

	AutomationOperatorAnd AutomationOperator = "AND"
	AutomationOperatorOr  AutomationOperator = "OR"

	AutomationPerformerAgent            AutomationPerformerType = 1
	AutomationPerformerRequester        AutomationPerformerType = 2
	AutomationPerformerAgentOrRequester AutomationPerformerType = 3
	AutomationPerformerSystem           AutomationPerformerType = 4

	AutomationMatchTypeAll AutomationMatchType = "all"
	AutomationMatchTypeAny AutomationMatchType = "any"

	AutomationPushToSlack     AutomationPushTo = "Slack"
	AutomationPushToOffice365 AutomationPushTo = "Office365"

	AutomationResourceSameTicket    AutomationResource = "Same_ticket"
	AutomationResourceParentTicket  AutomationResource = "parent_ticket"
	AutomationResourceTrackerTicket AutomationResource = "tracker_ticket"
	AutomationResourceCustomObject  AutomationResource = "custom_object"
)

type AutomationRules struct {
	Rules []*AutomationRule `json:"rules,omitempty"`
}

func (ars *AutomationRules) String() string {
	return toString(ars)
}

type AutomationRule struct {
	ID         int64                  `json:"id,omitempty"`         // Id of the automation
	Name       string                 `json:"name,omitempty"`       // Name of the automation rule
	Position   int                    `json:"position,omitempty"`   // Position of the automation rule
	Active     bool                   `json:"active,omitempty"`     // Set to true if the rule is active
	Performer  *AutomationPerformer   `json:"performer,omitempty"`  // **Applicable only if automation_type_id is 4**, Any event performer (agent, customer or system) whose action triggers the rule
	Events     []*AutomationEvent     `json:"events,omitempty"`     // **Applicable only if automation_type_id is 4**, Events that are responsible for triggering the rule
	Conditions []*AutomationCondition `json:"conditions,omitempty"` // Conditions to check whether the rule can run on a ticket or not
	Operator   AutomationOperator     `json:"operator,omitempty"`   // AND/OR operator to combine multiple conditions in a rule
	Actions    []*AutomationAction    `json:"actions,omitempty"`    // sActions to be performed by the rule on matching tickets

	Summary             map[string]any `json:"summary,omitempty"`
	Outdated            bool           `json:"outdated,omitempty"`
	AffectedTicketCount int            `json:"affected_tickets_count,omitempty"`
	LastUpdatedBy       int64          `json:"last_updated_by,omitempty"`
	Meta                map[string]any `json:"meta,omitempty"`
	CreatedAt           *Time          `json:"created_at,omitempty"`
	UpdatedAt           *Time          `json:"updated_at,omitempty"`
}

func (ar *AutomationRule) String() string {
	return toString(ar)
}

type AutomationPerformer struct {
	Type    AutomationPerformerType `json:"type,omitempty"`    // Agent/Requester/AgentOrRequester/System
	Members []int64                 `json:"members,omitempty"` // IDs of the agents
}

func (ap *AutomationPerformer) String() string {
	return toString(ap)
}

type AutomationEvent struct {
	FieldName string `json:"field_name,omitempty"` // Name of the field
	From      string `json:"from,omitempty"`       // Value of the field before the event
	To        string `json:"to,omitempty"`         // Value of the field after the event
}

func (ae *AutomationEvent) String() string {
	return toString(ae)
}

type AutomationCondition struct {
	Name       string                `json:"name,omitempty"`       // 	Title of the condition
	MatchType  AutomationMatchType   `json:"match_type,omitempty"` // To check whether all conditions have to be met or atleast one. Possible values are: “all”,”any”
	Properties []*AutomationProperty `json:"properties,omitempty"` // Properties of the condition
}

func (ac *AutomationCondition) String() string {
	return toString(ac)
}

type AutomationProperty struct {
	ResourceType      string           `json:"resource_type,omitempty"`
	FieldName         string           `json:"field_name,omitempty"`
	Operator          string           `json:"operator,omitempty"`
	Value             any              `json:"value,omitempty"`
	BusinessHoursID   int64            `json:"business_hours_id,omitempty"`
	CaseSensitive     bool             `json:"case_sensitive,omitempty"`
	NestedFields      map[string]any   `json:"nested_fields,omitempty"`
	AssociatedFields  map[string]any   `json:"associated_fields,omitempty"`
	RelatedConditions []map[string]any `json:"related_conditions,omitempty"`
}

func (ap *AutomationProperty) String() string {
	return toString(ap)
}

type AutomationAction struct {
	FieldName       string             `json:"field_name,omitempty"`       //	Name of the field
	Value           any                `json:"value,omitempty"`            // Value to be set on the field
	EmailTo         int64              `json:"email_to,omitempty"`         // 	Send email to specific contact/agent/groups.
	EmailBody       string             `json:"email_body,omitempty"`       //	Content of the email
	ApiKey          string             `json:"api_key,omitempty"`          // 	API key to authenticate any HTTP requests
	AuthHeader      map[string]string  `json:"auth_header,omitempty"`      // 	Combination of user name and password to be used for HTTP requests
	CustomHeader    map[string]string  `json:"custom_headers,omitempty"`   // 	Custom header information for any HTTP request
	RequestType     string             `json:"request_type,omitempty"`     //	Type of the HTTP request
	URL             string             `json:"url,omitempty"`              //	URL for the HTTP request
	NoteBody        string             `json:"note_body,omitempty"`        // 	Content of the note added by the rule
	NotifyAgents    []int64            `json:"notify_agents,omitempty"`    // IDs of agents to be notified
	FwdTo           string             `json:"Fwd_to,omitempty"`           // Forward the ticket to an email address
	FwdCc           string             `json:"Fwd_cc,omitempty"`           // Forward the ticket to an email address
	FwdBcc          string             `json:"Fwd_bcc,omitempty"`          // Forward the ticket to an email address
	FwdNoteBody     string             `json:"fwd_note_body,omitempty"`    // Forward the ticket to an email address
	PushTo          AutomationPushTo   `json:"push_to,omitempty"`          // Channel through which the message will be sent. Possible options are:	“Slack”	“Office365”
	SlackText       string             `json:"slack_text,omitempty"`       // Content of the message sent to slack
	Office365Text   string             `json:"office365_text,omitempty"`   // 	Content of the message sent to office365
	ResourceType    AutomationResource `json:"resource_type,omitempty"`    //  	Type of the ticket. Possible values are: “Same_ticket”, ”parent_ticket”, ”tracker_ticket”, ”custom_object”
	ObjectReference string             `json:"object_reference,omitempty"` // 	Ticket’s look up field value
}

func (aa *AutomationAction) String() string {
	return toString(aa)
}
