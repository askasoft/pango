package freshservice

type ListMembersOption = PageOption

type ListTicketsOption struct {
	Filter       string // The various filters available are new_and_my_open, watching, spam, deleted.
	WorkspaceID  int64
	RequestID    int64
	Type         string
	UpdatedSince Time
	Include      string // stats, requester, requester_for
	//OrderBy      string // created_at, due_by, updated_at, status
	OrderType string // asc, desc (default)
	Page      int
	PerPage   int
}

func (lto *ListTicketsOption) Values() Values {
	q := Values{}
	q.SetString("filter", lto.Filter)
	q.SetInt64("workspace_id", lto.WorkspaceID)
	q.SetInt64("request_id", lto.RequestID)
	q.SetString("type", lto.Type)
	q.SetTime("updated_since", lto.UpdatedSince)
	q.SetString("include", lto.Include)
	//q.SetString("order_by", lto.OrderBy)
	q.SetString("order_type", lto.OrderType)
	q.SetInt("page", lto.Page)
	q.SetInt("per_page", lto.PerPage)
	return q
}

type FilterTicketsOption struct {
	Query   string
	Page    int
	PerPage int
}

func (fto *FilterTicketsOption) Values() Values {
	q := Values{}
	q.SetString("query", fto.Query)
	q.SetInt("page", fto.Page)
	q.SetInt("per_page", fto.PerPage)
	return q
}

// PerPage: 1 ~ 100, default: 30
type ListConversationsOption = PageOption

type ListAgentsOption struct {
	FirstName          string
	LastName           string
	Name               string // Concatenation of first_name and last_name with single space in-between fields.
	JobTitle           string
	Email              string
	MobilePhoneNumber  string
	WorkPhoneNumber    string
	Active             bool
	State              string // fulltime or occasional
	DepartmentID       int64
	ReportingManagerID int64
	TimeZone           string
	Language           string
	LocationID         int64
	CreatedAt          Date // Date (YYYY-MM-DD) when the requester is created.
	UpdatedAt          Date // Date (YYYY-MM-DD) when the requester is updated.
	Page               int
	PerPage            int
}

func (lao *ListAgentsOption) Values() Values {
	q := Values{}
	q.SetString("first_name", lao.FirstName)
	q.SetString("last_name", lao.LastName)
	q.SetString("name", lao.Name)
	q.SetString("job_title", lao.JobTitle)
	q.SetString("email", lao.Email)
	q.SetString("mobile_phone_number", lao.MobilePhoneNumber)
	q.SetString("work_phone_number", lao.WorkPhoneNumber)
	q.SetBool("active", lao.Active)
	q.SetString("state", lao.State)
	q.SetInt64("department_id", lao.DepartmentID)
	q.SetInt64("reporting_manager_id", lao.ReportingManagerID)
	q.SetString("time_zone", lao.TimeZone)
	q.SetString("language", lao.Language)
	q.SetInt64("location_id", lao.LocationID)
	q.SetDate("created_at", lao.CreatedAt)
	q.SetDate("updated_at", lao.UpdatedAt)
	q.SetInt("page", lao.Page)
	q.SetInt("per_page", lao.PerPage)
	return q
}

type ListRequestersOption struct {
	FirstName         string
	LastName          string
	Name              string // Concatenation of first_name and last_name with single space in-between fields.
	JobTitle          string
	PrimaryEmail      string
	Email             string
	MobilePhoneNumber string
	WorkPhoneNumber   string
	DepartmentID      int64
	TimeZone          string
	Language          string
	LocationID        int64
	CreatedAt         Date // Date (YYYY-MM-DD) when the requester is created.
	UpdatedAt         Date // Date (YYYY-MM-DD) when the requester is updated.
	IncludeAgents     bool
	Page              int
	PerPage           int
}

func (lro *ListRequestersOption) Values() Values {
	q := Values{}
	q.SetString("first_name", lro.FirstName)
	q.SetString("last_name", lro.LastName)
	q.SetString("name", lro.Name)
	q.SetString("job_title", lro.JobTitle)
	q.SetString("primary_email", lro.PrimaryEmail)
	q.SetString("email", lro.Email)
	q.SetString("mobile_phone_number", lro.MobilePhoneNumber)
	q.SetString("work_phone_number", lro.WorkPhoneNumber)
	q.SetInt64("department_id", lro.DepartmentID)
	q.SetString("time_zone", lro.TimeZone)
	q.SetString("language", lro.Language)
	q.SetInt64("location_id", lro.LocationID)
	q.SetDate("created_at", lro.CreatedAt)
	q.SetDate("updated_at", lro.UpdatedAt)
	q.SetBool("include_agents", lro.IncludeAgents)
	q.SetInt("page", lro.Page)
	q.SetInt("per_page", lro.PerPage)
	return q
}

type SearchArticlesOption struct {
	SearchTerm string // The keywords for which the solution articles have to be searched.
	UserEmail  string // By default, the API will search the articles for the user whose API key is provided. If you want to search articles for a different user, please provide their user_email.
	Page       int
	PerPage    int
}

func (sao *SearchArticlesOption) Values() Values {
	q := Values{}
	q.SetString("search_term", sao.SearchTerm)
	q.SetString("user_email", sao.UserEmail)
	q.SetInt("page", sao.Page)
	q.SetInt("per_page", sao.PerPage)
	return q
}
