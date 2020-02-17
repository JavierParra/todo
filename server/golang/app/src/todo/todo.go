package todo

type Todo struct {
	Id        string `json:"id"`
	Complete  bool   `json:"complete"`
	Name      string `json:"name"`
	Created   int64  `json:"created"`
	Completed int64  `json:"completed"`
	Notes     string `json:"notes"`
}
