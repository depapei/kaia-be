package responses

type Success struct {
	Data    interface{} `json:"data"`
	Success bool        `json:"success"`
}

type Fail struct {
	Message string `json:"message"`
}
