package payload

type Task struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Status  string `json:"status"`
}
