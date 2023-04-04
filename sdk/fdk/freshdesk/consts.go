package freshdesk

const (
	//JobStatusInProgress  = "IN PROGRESS"
	JobStatusInProgress = "in_progress"
	JobStatusCompleted  = "completed"

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

	TicketFilterNewAndMyOpen = "new_and_my_open"
	TicketFilterWatching     = "watching"
	TicketFilterSpam         = "spam"
	TicketFilterDeleted      = "deleted"

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
	IncludeCompany       = "company"       // Will return the company's id and name.
	IncludeConversations = "conversations" // Will return up to 10 conversations sorted by "created_at" in ascending order
	IncludeRequester     = "requester"     // Will return the requester's email, id, mobile, name, and phone
	IncludeStats         = "stats"         // Will return the ticket’s closed_at, resolved_at and first_responded_at time

	ContactStateBlocked    = "blocked"
	ContactStateDeleted    = "deleted"
	ContactStateUnverified = "unverified"
	ContactStateVerified   = "verified"

	AgentStateFulltime   = "fulltime"
	AgentStateOccasional = "occasional"

	AgentTicketScopeGlobal     = 1
	AgentTicketScopeGroup      = 2
	AgentTicketScopeRestricted = 3

	AssignDisabled             = 0
	AssignRoundRobin           = 1
	AssignSkillBasedRoundRobin = 2
	AssignLoadBasedRoundRobin  = 3
	AssignOmniroute            = 4

	AutomationTypeTicketCreation = 1
	AutomationTypeTimeTriggers   = 3
	AutomationTypeTicketUpdates  = 4

	AutomationOperatorAnd = "AND"
	AutomationOperatorOr  = "OR"

	AutomationPerformerAgent            = 1
	AutomationPerformerRequester        = 2
	AutomationPerformerAgentOrRequester = 3
	AutomationPerformerSystem           = 4

	AutomationMatchTypeAll = "all"
	AutomationMatchTypeAny = "any"

	AutomationActionPushToSlack     = "Slack"
	AutomationActionPushToOffice365 = "Office365"

	AutomationActionSameTicket    = "Same_ticket"
	AutomationActionParentTicket  = "parent_ticket"
	AutomationActionTrackerTicket = "tracker_ticket"
	AutomationActionCustomObject  = "custom_object"
)
