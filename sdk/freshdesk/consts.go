package freshdesk

const (
	SourceEmail          = 1
	SourceProtal         = 2
	SourcePhone          = 3
	SourceChat           = 7
	SourceFeedbackWidget = 9
	SourceOutboundEmail  = 10

	StatusOpen     = 2
	StatusPending  = 3
	StatusResolved = 4
	StatusClosed   = 5

	PriorityLow    = 1
	PriorityMedium = 2
	PriorityHigh   = 3
	PriorityUrgent = 4

	OrderByCreatedAt = "created_at"
	OrderByDueBy     = "due_by"
	OrderByUpdatedAt = "updated_at"
	OrderByStatus    = "status"

	OrderAsc  = "asc"
	OrderDesc = "desc"

	IncludeDescription   = "description"
	IncludeCompany       = "company"
	IncludeConversations = "conversations"
	IncludeRequester     = "requester"
	IncludeStats         = "stats"
)
