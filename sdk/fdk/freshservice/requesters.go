package freshservice

import (
	"github.com/askasoft/pango/str"
)

// ---------------------------------------------------
// Requester

type ListMembersOption = PageOption

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

func (lro *ListRequestersOption) IsNil() bool {
	return lro == nil
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

func (fs *Freshservice) CreateRequesterGroup(rg *RequesterGroup) (*RequesterGroup, error) {
	url := fs.endpoint("/requester_groups")
	result := &requesterGroupResult{}
	err := fs.doPost(url, rg, result)
	return result.RequesterGroup, err
}

func (fs *Freshservice) GetRequesterGroup(id int64) (*RequesterGroup, error) {
	url := fs.endpoint("/requester_groups/%d", id)
	result := &requesterGroupResult{}
	err := fs.doGet(url, result)
	return result.RequesterGroup, err
}

func (fs *Freshservice) ListRequesterGroups() ([]*RequesterGroup, error) {
	url := fs.endpoint("/requester_groups")
	result := &requesterGroupResult{}
	err := fs.doGet(url, result)
	return result.RequesterGroups, err
}

// Note:
// Only groups of type “manual” can be updated through this API.
func (fs *Freshservice) UpdateRequesterGroup(id int64, rg *RequesterGroup) (*RequesterGroup, error) {
	url := fs.endpoint("/requester_groups/%d", id)
	result := &requesterGroupResult{}
	err := fs.doPut(url, rg, result)
	return result.RequesterGroup, err
}

// Delete Requester Group
// Note:
// 1. Deleting a Requester Group will only disband the requester group and will not delete its members.
// 2. Deleted requester groups cannot be restored.
func (fs *Freshservice) DeleteRequesterGroup(id int64) error {
	url := fs.endpoint("/requester_groups/%d", id)
	return fs.doDelete(url)
}

// Add Requester to Requester Group
// Note:
// 1.Requesters can be added only to manual requester groups.
// 2.Requester can be added one at a time.
func (fs *Freshservice) AddRequesterToRequesterGroup(rgid, rid int64) error {
	url := fs.endpoint("/requester_groups/%d/members/%d", rgid, rid)
	return fs.doPost(url, nil, nil)
}

// Delete Requester from Requester Group
// Note:
// 1.Requesters can be removed only from manual requester groups.
// 2.Requester can be removed one at a time.
func (fs *Freshservice) DeleteRequesterFromRequesterGroup(rgid, rid int64) error {
	url := fs.endpoint("/requester_groups/%d/members/%d", rgid, rid)
	return fs.doDelete(url)
}

func (fs *Freshservice) ListRequesterGroupMembers(rgid int64, lmo *ListMembersOption) ([]*Requester, bool, error) {
	url := fs.endpoint("/requester_groups/%d/members", rgid)
	result := &requesterResult{}
	next, err := fs.doList(url, lmo, result)
	return result.Requesters, next, err
}

func (fs *Freshservice) CreateRequester(requester *Requester) (*Requester, error) {
	url := fs.endpoint("/requesters")
	result := &requesterResult{}
	err := fs.doPost(url, requester, result)
	return result.Requester, err
}

func (fs *Freshservice) GetRequester(id int64) (*Requester, error) {
	url := fs.endpoint("/requesters/%d", id)
	result := &requesterResult{}
	err := fs.doGet(url, result)
	return result.Requester, err
}

// List Requesters
// Use Requester attributes to filter your list.
// Note:
// 1. Filtered results cannot be sorted. By default it is sorted by created_at in descending order.
// 2. Adding "include_agents=true" to the query string will include agents in the response. The default response includes only requesters and not agents. Only users who also have the "Manage Agents" permission will be able to use this modifier.
// 3. The query must be URL encoded (see example).
// 4. Query can be framed using the name of the requester fields, which can be obtained from the Supported Requester Fields Section.
// 5. Query string must be enclosed between a pair of double quotes and can have up to 512 characters.
// 6. Logical operators AND, OR along with parenthesis( ) can be used to group conditions.
// 7. Relational operators greater than or equal to :> and less than or equal to :< can be used along with date fields and numeric fields.
// 8. Input for date field should be in UTC Format.
// 9. The number of objects returned per page is 30.
// 10. To scroll through the pages add the page parameter to the url. The page number starts with 1 and should not exceed 40.
// 11. To filter for fields with no values assigned, use the null keyword.
// 12. The "~" query operator can be used for "starts with" text searches. "Starts with" search is supported for one or more of the following attributes: first_name, last_name, name, primary_email, mobile_phone_number, work_phone_number. The query format is https://domain.freshservice.com/api/v2/requesters?query="~[attribute_1|attribute_2]:'somestring'". The query needs to be URL encoded. This would return a list of users for whom attribute_1 OR attribute_2 starts with "somestring". Refer to examples 13, 14, and 15.
// 13. Please note that any update made to requester either in Freshservice application or through API may take a few minutes to get indexed, after which the updated results will be available through API.
// == Custom Fields Supported	Type
// Single line text	string
// Number	integer
// Dropdown	string
// Date	date
// Phone number	string
func (fs *Freshservice) ListRequesters(lro *ListRequestersOption) ([]*Requester, bool, error) {
	url := fs.endpoint("/requesters")
	result := &requesterResult{}
	next, err := fs.doList(url, lro, result)
	return result.Requesters, next, err
}

func (fs *Freshservice) GetRequesterFields() ([]*RequesterField, error) {
	url := fs.endpoint("/requester_fields")
	result := &requesterFieldResult{}
	err := fs.doGet(url, result)
	return result.RequesterFields, err
}

// Update a Requester
// This operation allows you to modify the profile of a particular requester.
// Note:
// can_see_all_tickets_from_associated_departments will automatically be set to false unless it is explicitly set to true in the payload, irrespective of the previous value of the field.
func (fs *Freshservice) UpdateRequester(id int64, requester *Requester) (*Requester, error) {
	url := fs.endpoint("/requesters/%d", id)
	result := &requesterResult{}
	err := fs.doPut(url, requester, result)
	return result.Requester, err
}

// Deactivate a Requester
// This operation allows you to deactivate a requester.
func (fs *Freshservice) DeactivateRequester(id int64) error {
	url := fs.endpoint("/requesters/%d", id)
	return fs.doDelete(url)
}

// Forget a Requester
// This operation allows you to permanently delete a requester and the tickets that they requested.
func (fs *Freshservice) ForgetRequester(id int64) error {
	url := fs.endpoint("/requesters/%d/forget", id)
	return fs.doDelete(url)
}

// Convert a requester to an occasional agent with SD Agent role and no group memberships.
func (fs *Freshservice) ConvertRequesterToAgent(id int64) (*Agent, error) {
	url := fs.endpoint("/requesters/%d/convert_to_agent", id)
	result := &agentResult{}
	err := fs.doPut(url, nil, result)
	return result.Agent, err
}

// Merge secondary requesters into a primary requester.
func (fs *Freshservice) MergeRequesters(id int64, ids ...int64) (*Requester, error) {
	url := fs.endpoint("/requesters/%d/merge?secondary_requesters=%s", id, str.JoinInt64s(ids, ","))
	result := &requesterResult{}
	err := fs.doPut(url, nil, result)
	return result.Requester, err
}

// Reactivate a Requester
// This operation allows you to reactivate a particular deactivated requester.
func (fs *Freshservice) ReactivateRequester(id int64) (*Requester, error) {
	url := fs.endpoint("/requesters/%d/reactivate", id)
	result := &requesterResult{}
	err := fs.doPut(url, nil, result)
	return result.Requester, err
}
