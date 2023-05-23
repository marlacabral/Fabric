package txdefs

import (
	"encoding/json"

	"github.com/goledgerdev/cc-tools/assets"
	"github.com/goledgerdev/cc-tools/errors"
	sw "github.com/goledgerdev/cc-tools/stubwrapper"
	tx "github.com/goledgerdev/cc-tools/transactions"
)

// Método GET
var GetContabilidadeToken = tx.Transaction{
	Tag:     "getContabilidadeToken",
	Label:   "Obter Token de Contabilidade",
	Method:  "GET",
	Callers: []string{"$org1MSP", "$org2MSP"}, // Apenas org1 e org2 podem chamar essa transação

	Args: []tx.Argument{
		{
			Required: true,
			Tag:      "proprietario",
			Label:    "Nome do proprietário",
			DataType: "string",
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		proprietario, _ := req["proprietario"].(string)
		limite, hasLimite := req["limite"].(float64)

		if hasLimite && limite <= 0 {
			return nil, errors.NewCCError("O limite deve ser maior que zero", 400)
		}

		// Preparar consulta para o CouchDB
		consulta := map[string]interface{}{
			"selector": map[string]interface{}{
				"@assetType":   "token",
				"proprietario": proprietario,
			},
		}

		if hasLimite {
			consulta["limit"] = limite
		}

		var err error
		resposta, err := assets.Search(stub, consulta, "", true)
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "Erro ao buscar o token de contabilidade", 500)
		}

		respostaJSON, err := json.Marshal(resposta)
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "Erro ao converter a resposta para JSON", 500)
		}

		return respostaJSON, nil
	},
}
