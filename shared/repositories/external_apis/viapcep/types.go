package viacep

type ViaCepResponse struct {
	Cep        string `json:"cep"`
	Street     string `json:"logradouro"`
	Complement string `json:"complemento"`
	District   string `json:"bairro"`
	City       string `json:"localidade"`
	State      string `json:"uf"`
	IbgeCode   string `json:"ibge"`
	GiaCode    string `json:"gia"`
	SiafiCode  string `json:"siafi"`
}
