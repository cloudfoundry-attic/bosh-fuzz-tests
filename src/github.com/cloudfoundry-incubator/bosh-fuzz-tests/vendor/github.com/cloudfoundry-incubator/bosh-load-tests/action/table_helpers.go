package action

type Table struct {
	Rows []map[string]interface{} `json:"Rows"`
}

type Output struct {
	Tables []Table `json:"Tables"`
}
