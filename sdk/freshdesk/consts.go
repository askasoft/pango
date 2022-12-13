package freshdesk

const (
	TicketSourceEmail          = 1
	TicketSourceProtal         = 2
	TicketSourcePhone          = 3
	TicketSourceChat           = 7
	TicketSourceFeedbackWidget = 9
	TicketSourceOutboundEmail  = 10

	TicketStatusOpen     = 2
	TicketStatusPending  = 3
	TicketStatusResolved = 4
	TicketStatusClosed   = 5

	TicketPriorityLow    = 1
	TicketPriorityMedium = 2
	TicketPriorityHigh   = 3
	TicketPriorityUrgent = 4

	ConversationSourceReply     = 0
	ConversationSourceNote      = 2
	ConversationSourceTweets    = 5
	ConversationSourceSurvey    = 6
	ConversationSourceFacebook  = 7
	ConversationSourceForwarded = 8
	ConversationSourcePhone     = 9
	ConversationSourceECommerce = 11

	FolderVisibilityAllUsers                = 1
	FolderVisibilityLoggedInUsers           = 2
	FolderVisibilityAgents                  = 3
	FolderVisibilitySelectedCompanies       = 4
	FolderVisibilityBots                    = 5
	FolderVisibilitySelectedContactSegments = 6
	FolderVisibilitySelectedCompanySegments = 7

	ArticleStatusDraft     = 1
	ArticleStatusPublished = 2

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

	StateBlocked    = "blocked"
	StateDeleted    = "deleted"
	StateUnverified = "unverified"
	StateVerified   = "verified"
)
