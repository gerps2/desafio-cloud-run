package valueObjects

import (
	"errors"
	"regexp"
	"strings"
)

type Cep string

func NewCep(codigo string) (Cep, error) {
	codigo = strings.TrimSpace(codigo)

	if !isValidCep(codigo) {
		return "", errors.New("CEP inválido")
	}

	if len(codigo) == 8 && !strings.Contains(codigo, "-") {
		codigo = codigo[:5] + "-" + codigo[5:]
	}

	return Cep(codigo), nil
}

func isValidCep(codigo string) bool {
	reDash := regexp.MustCompile(`^\d{5}-\d{3}$`)
	reNoDash := regexp.MustCompile(`^\d{8}$`)
	return reDash.MatchString(codigo) || reNoDash.MatchString(codigo)
}

func (c Cep) String() string {
	return string(c)
}
