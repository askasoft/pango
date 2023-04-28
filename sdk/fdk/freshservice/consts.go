package freshservice

const (
	OrderAsc  = "asc"
	OrderDesc = "desc"

	IncludeConversations  = "conversations"
	IncludeRequester      = "requester"
	IncludeRequestedFor   = "requested_for"
	IncludeStats          = "stats"
	IncludeProblem        = "problem"
	IncludeAssets         = "assets"
	IncludeChange         = "change"
	IncludeRelatedTickets = "related_tickets"

	AgentStateFulltime   = "fulltime"
	AgentStateOccasional = "occasional"

	AgentGroupUnassignedFor30m = "30m"
	AgentGroupUnassignedFor1h  = "1h"
	AgentGroupUnassignedFor2h  = "2h"
	AgentGroupUnassignedFor4h  = "4h"
	AgentGroupUnassignedFor8h  = "8h"
	AgentGroupUnassignedFor12h = "12h"
	AgentGroupUnassignedFor1d  = "1d"
	AgentGroupUnassignedFor2d  = "2d"
	AgentGroupUnassignedFor3d  = "3d"

	RequesterGroupTypeManual    = "manual"
	RequesterGroupTypeRuleBased = "rule_based"

	TicketSourceEmail          = 1
	TicketSourceProtal         = 2
	TicketSourcePhone          = 3
	TicketSourceChat           = 4
	TicketSourceFeedbackWidget = 5
	TicketSourceYammer         = 6
	TicketSourceAWSCloudwatch  = 7
	TicketSourcePagerduty      = 8
	TicketSourceWalkup         = 9
	TicketSourceSlack          = 10

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

	ConversationSourceEmail          = 0
	ConversationSourceForm           = 1
	ConversationSourceNote           = 2
	ConversationSourceStatus         = 3
	ConversationSourceMeta           = 4
	ConversationSourceFeedback       = 5
	ConversationSourceForwardedEmail = 6

	FolderVisibilityAllUsers      = 1
	FolderVisibilityLoggedInUsers = 2
	FolderVisibilityAgentsOnly    = 3
	FolderVisibilityDepartments   = 4
	FolderVisibilityAgentGroups   = 5
	FolderVisibilityContactGroups = 6

	ArticleTypePermanent  = 1
	ArticleTypeWorkaround = 1

	ArticleStatusDraft     = 1
	ArticleStatusPublished = 2
)
