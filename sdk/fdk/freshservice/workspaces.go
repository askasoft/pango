package freshservice

// ---------------------------------------------------
// Workspace

type ListWorkspacesOption = PageOption

func (fs *Freshservice) GetWorkspace(id int64) (*Workspace, error) {
	url := fs.endpoint("/workspaces/%d", id)
	result := &workspaceResult{}
	err := fs.doGet(url, result)
	return result.Workspace, err
}

func (fs *Freshservice) ListWorkspaces(lwo *ListWorkspacesOption) ([]*Workspace, bool, error) {
	url := fs.endpoint("/workspaces")
	result := &workspacesResult{}
	next, err := fs.doList(url, lwo, result)
	return result.Workspaces, next, err
}

func (fs *Freshservice) IterWorkspaces(lwo *ListWorkspacesOption, iwf func(*Workspace) error) error {
	if lwo == nil {
		lwo = &ListWorkspacesOption{}
	}
	if lwo.Page < 1 {
		lwo.Page = 1
	}
	if lwo.PerPage < 1 {
		lwo.PerPage = 100
	}

	for {
		ws, next, err := fs.ListWorkspaces(lwo)
		if err != nil {
			return err
		}
		for _, ag := range ws {
			if err = iwf(ag); err != nil {
				return err
			}
		}
		if !next {
			break
		}
		lwo.Page++
	}
	return nil
}
