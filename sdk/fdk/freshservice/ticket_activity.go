package freshservice

type Actor struct {
	ID   int64  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type TicketActivity struct {
	Actor       *Actor   `json:"actor,omitempty"`
	Content     string   `json:"content,omitempty"`
	SubContents []string `json:"sub_contents,omitempty"`
	CreatedAt   *Time    `json:"created_at,omitempty"`
}

func (ta *TicketActivity) String() string {
	return toString(ta)
}

type ticketActivityResult struct {
	TicketActivity  *TicketActivity   `json:"ticket_activity,omitempty"`
	TicketActivitys []*TicketActivity `json:"ticket_activities,omitempty"`
}
