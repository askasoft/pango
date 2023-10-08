package freshservice

type TimeEntry struct {

	// Id of the time entry
	ID int64 `json:"id,omitempty"`

	// ID of the workspace to which the time entry belongs(inherited from the ticket's workspace). The attribute is applicable only for accounts with the 'Workspaces' feature enabled.
	WorkspaceID int64 `json:"workspace_id,omitempty"`

	// Set to true if timer is currently running. Default value is false.
	// At a time, only one timer can be running for an agent across the account
	TimerRunning bool `json:"timer_running"`

	// Boolean	Set as true if the time-entry is billable. Default value is true
	Billable bool `json:"billable"`

	// The total amount of time spent by the timer in hh::mm format.
	// This field cannot be set if timer_running is true.
	// Mandatory if timer_running is false
	TimeSpent TimeSpent `json:"time_spent,omitempty"`

	// Time at which the timer is executed.
	// Default value (unless given in request) is the time at which timer is added.
	// Should be less than or equal to current date_time
	ExecutedAt *Time `json:"executed_at,omitempty"`

	// Id of the task assigned to the time-entry. Task should be valid on the given ticket and assigned to agent_id
	TaskID int64 `json:"task_id,omitempty"`

	// Description of the time-entry
	Note string `json:"note,omitempty"`

	// The user/agent to whom this time-entry is assigned
	AgentID int64 `json:"agent_id,omitempty"`

	// Key value pairs containing the names and values of custom fields.
	CustomFields map[string]any `json:"custom_fields,omitempty"`

	// The time at which the time-entry is added. If a timer, which is in stopped state, is started again, this holds date_time at which the timer is started again
	StartTime *Time `json:"start_time,omitempty"`

	// Ticket creation timestamp
	CreatedAt *Time `json:"created_at,omitempty"`

	// Ticket updated timestamp
	UpdatedAt *Time `json:"updated_at,omitempty"`
}

func (te *TimeEntry) String() string {
	return toString(te)
}

type timeEntryResult struct {
	TimeEntry   *TimeEntry   `json:"time_entry,omitempty"`
	TimeEntries []*TimeEntry `json:"time_entries,omitempty"`
}
