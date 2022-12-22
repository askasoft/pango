package freshdesk

type Job struct {
	ID string `json:"id,omitempty"`

	Name string `json:"name,omitempty"`

	Status string `json:"status,omitempty"`

	DownloadURL string `json:"download_url,omitempty"`

	CreatedAt *Time `json:"created_at,omitempty"`

	UpdatedAt *Time `json:"updated_at,omitempty"`

	StatusUpdatedAt *Time `json:"status_updated_at,omitempty"`

	Progress int `json:"progress,omitempty"`
}

func (job *Job) IsCompleted() bool {
	return job.Status == JobStatusCompleted
}

func (job *Job) IsInProgress() bool {
	return job.Status == JobStatusInProgress
}

func (job *Job) String() string {
	return toString(job)
}
