package txdefs

import (
	"encoding/json"

	"github.com/goledgerdev/cc-tools/assets"
	"github.com/goledgerdev/cc-tools/errors"
	sw "github.com/goledgerdev/cc-tools/stubwrapper"
	tx "github.com/goledgerdev/cc-tools/transactions"
)

// Create a new Library on channel
// POST Method
var CriarToken = tx.Transaction{
	Tag:         "criarToken",
	Label:       "Criar Novo token",
	Description: "Criar Novo token",
	Method:      "POST",
	Callers:     []string{"$org1MSP"}, // Only org1 can call this transaction

	Args: []tx.Argument{
		{
			Tag:         "token",
			Label:       "Token",
			Description: "Token",
			DataType:    "string",
			Required:    true,
		},
		{
			Tag:         "quantidade",
			Label:       "Quantidade",
			Description: "Quantidade do token",
			DataType:    "number",
			Required:    true,
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		token, _ := req["token"].(string)
		quantidade, _ := req["quantidade"].(float64)

		if quantidade <= 0 {
			return nil, errors.WrapError(nil, "A quantidade deve ser maior que 0")
		}

		novoTokenMap := make(map[string]interface{})
		novoTokenMap["@assetType"] = "novoToken"
		novoTokenMap["token"] = token
		novoTokenMap["quantidade"] = quantidade

		novoTokenAsset, err := assets.NewAsset(novoTokenMap)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to create a new asset")
		}

		// Save the new token on channel
		_, err = novoTokenAsset.PutNew(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Error saving asset on blockchain")
		}

		// Marshal asset back to JSON format
		novoTokenJSON, nerr := json.Marshal(novoTokenAsset)
		if nerr != nil {
			return nil, errors.WrapError(nil, "failed to encode asset to JSON format")
		}

		return novoTokenJSON, nil
	},
}
