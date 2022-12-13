package freshdesk

type ListTicketsOption struct {
	Filter           string // The various filters available are new_and_my_open, watching, spam, deleted.
	RequestID        int64
	Email            string
	UniqueExternalID string
	CompanyID        int64
	UpdatedSince     Time
	Include          string // stats, requester, description
	OrderBy          string // created_at, due_by, updated_at, status
	OrderType        string // asc, desc (default)
	Page             int
	PerPage          int
}

func (lto *ListTicketsOption) Values() Values {
	q := Values{}
	q.SetString("filter", lto.Filter)
	q.SetInt64("request_id", lto.RequestID)
	q.SetString("email", lto.Email)
	q.SetString("unique_external_id", lto.UniqueExternalID)
	q.SetInt64("company_id", lto.CompanyID)
	q.SetTime("updated_since", lto.UpdatedSince)
	q.SetString("include", lto.Include)
	q.SetString("order_by", lto.OrderBy)
	q.SetString("order_type", lto.OrderType)
	q.SetInt("page", lto.Page)
	q.SetInt("per_page", lto.PerPage)
	return q
}

type ListContactsOption struct {
	Email        string
	Mobile       string
	Phone        string
	CompanyID    int64
	UpdatedSince Time
	State        string // [blocked/deleted/unverified/verified]
	Page         int
	PerPage      int
}

func (lco *ListContactsOption) Values() Values {
	q := Values{}
	q.SetString("email", lco.Email)
	q.SetString("mobile", lco.Mobile)
	q.SetString("phone", lco.Phone)
	q.SetInt64("company_id", lco.CompanyID)
	q.SetString("state", lco.State)
	q.SetTime("updated_since", lco.UpdatedSince)
	q.SetInt("page", lco.Page)
	q.SetInt("per_page", lco.PerPage)
	return q
}
