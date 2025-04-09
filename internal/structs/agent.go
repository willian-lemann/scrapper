package structs

type AgentDTO struct {
	Id              int      `json:"id"`
	Creci           int      `json:"creci"`
	Nome            string   `json:"nome"`
	Tipo            int16    `json:"tipo"`
	Secundaria      bool     `json:"secundaria"`
	Regular         bool     `json:"regular"`
	Situacao        int64    `json:"situacao"`
	Foto            bool     `json:"foto"`
	Cpf             string   `json:"cpf"`
	Termo           int      `json:"termo"`
	Telefones       []string `json:"telefones"`
	DataFalecimento *string  `json:"datafalecimento"`
}

type Agent struct {
	Id        int      `json:"id"`
	AgentId   int      `json:"agentId"`
	Name      string   `json:"name"`
	Type      int16    `json:"type"`
	IsRegular bool     `json:"isRegular"`
	HasPhoto  bool     `json:"hasPhoto"`
	Cpf       string   `json:"cpf"`
	IsActive  bool     `json:"isActive"`
	Phones    []string `json:"phones"`
}
