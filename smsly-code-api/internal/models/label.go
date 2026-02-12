package models

type Label struct {
	ID          int64  `json:"id"`
	RepoID      int64  `json:"repo_id"`
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description"`
}
