package freshservice

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
