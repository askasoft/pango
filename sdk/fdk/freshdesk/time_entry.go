package freshdesk

type TimeEntry struct {
	// Id of the time entry
	ID int64 `json:"id,omitempty"`

	// Boolean	Set as true if the time-entry is billable
	Billable bool `json:"billable"`

	// Set to true if timer is currently running
	TimerRunning bool `json:"timer_running"`

	// The total amount of time spent by the timer in hh::mm format.
	// This field cannot be set if timer_running is true.
	// Mandatory if timer_running is false
	TimeSpent TimeSpent `json:"time_spent,omitempty"`

	// Description of the time-entry
	Note string `json:"note,omitempty"`

	// The user/agent to whom this time-entry is assigned
	AgentID int64 `json:"agent_id,omitempty"`

	// The ID of the ticket to which this time entry is associated
	TicketID int64 `json:"ticket_id,omitempty"`

	// Key value pairs containing the names and values of custom fields.
	CustomFields map[string]any `json:"custom_fields,omitempty"`

	// Time at which this time-entry was added/created
	ExecutedAt Time `json:"executed_at,omitempty"`

	// The time at which the time-entry is added. If a timer, which is in stopped state, is started again, this holds date_time at which the timer is started again
	StartTime Time `json:"start_time,omitempty"`

	// Ticket creation timestamp
	CreatedAt Time `json:"created_at,omitempty"`

	// Ticket updated timestamp
	UpdatedAt Time `json:"updated_at,omitempty"`
}

func (te *TimeEntry) String() string {
	return toString(te)
}

type TimeEntryCreate struct {
	// Boolean	Set as true if the time-entry is billable
	Billable bool `json:"billable"`

	// Set to true if timer is currently running
	TimerRunning bool `json:"timer_running"`

	// The total amount of time spent by the timer in hh::mm format.
	// This field cannot be set if timer_running is true.
	// Mandatory if timer_running is false
	TimeSpent TimeSpent `json:"time_spent,omitempty"`

	// Description of the time-entry
	Note string `json:"note,omitempty"`

	// The user/agent to whom this time-entry is assigned
	AgentID int64 `json:"agent_id,omitempty"`

	// Time at which this time-entry was added/created
	ExecutedAt *Time `json:"executed_at,omitempty"`

	// The time at which the time-entry is added. If a timer, which is in stopped state, is started again, this holds date_time at which the timer is started again
	StartTime *Time `json:"start_time,omitempty"`
}

func (te *TimeEntryCreate) String() string {
	return toString(te)
}

type TimeEntryUpdate = TimeEntryCreate
