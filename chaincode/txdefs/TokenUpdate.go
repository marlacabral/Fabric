package txdefs

import (
	"encoding/json"

	"github.com/goledgerdev/cc-tools/assets"
	"github.com/goledgerdev/cc-tools/errors"
	sw "github.com/goledgerdev/cc-tools/stubwrapper"
	tx "github.com/goledgerdev/cc-tools/transactions"
)

// Transfer Token
// PUT Method
var TransferToken = tx.Transaction{
	Tag:         "transferToken",
	Label:       "Transfer Token",
	Description: "Transfer tokens between owners",
	Method:      "PUT",
	Callers:     []string{`$org\dMSP`, "org1MSP"}, // Any orgs can call this transaction

	Args: []tx.Argument{
		{
			Tag:         "tokenOrigem",
			Label:       "Token Origem",
			Description: "Token to transfer from",
			DataType:    "->token",
			Required:    true,
		},
		{
			Tag:         "tokenDestino",
			Label:       "Token Destino",
			Description: "Token to transfer to",
			DataType:    "->token",
			Required:    true,
		},
		{
			Tag:         "quantidadeTransferida",
			Label:       "Quantidade Transferida",
			Description: "Quantity of tokens to transfer",
			DataType:    "number",
			Required:    true,
		},
		{
			Tag:         "idNovoToken",
			Label:       "ID Novo Token",
			Description: "ID of the new token",
			DataType:    "string",
			Required:    true,
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		tokenOrigemKey, ok := req["tokenOrigem"].(assets.Key)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter tokenOrigem must be an asset")
		}
		tokenDestinoKey, ok := req["tokenDestino"].(assets.Key)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter tokenDestino must be an asset")
		}

		// Returns tokenOrigem from channel
		tokenOrigemAsset, err := tokenOrigemKey.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "failed to get tokenOrigem from the ledger")
		}
		tokenOrigemMap := (map[string]interface{})(*tokenOrigemAsset)

		// Returns tokenDestino from channel
		tokenDestinoAsset, err := tokenDestinoKey.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "failed to get tokenDestino from the ledger")
		}
		tokenDestinoMap := (map[string]interface{})(*tokenDestinoAsset)

		quantidadeTransferida, ok := req["quantidadeTransferida"].(float64)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter quantidadeTransferida must be a number")
		}

		idNovoToken, ok := req["idNovoToken"].(string)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter idNovoToken must be a string")
		}

		// Mark tokenOrigem as burned
		tokenOrigemMap["burned"] = true

		// Create tokenDestino with the transferred quantity
		tokenDestinoMap["quantidade"] = quantidadeTransferida
		tokenDestinoMap["burned"] = false

		// Create new tokenOrigem with the remaining quantity
		novoTokenOrigem := make(map[string]interface{})
		novoTokenOrigem["@assetType"] = "token"
		novoTokenOrigem["@key"] = idNovoToken
		novoTokenOrigem["quantidade"] = tokenOrigemMap["quantidade"].(float64) - quantidadeTransferida
		novoTokenOrigem["burned"] = false

		// Update tokenOrigem data
		tokenOrigemMap["@key"] = idNovoToken

		// Update tokenDestino and novoTokenOrigem assets
		_, err = tokenDestinoAsset.Update(stub, tokenDestinoMap)
		if err != nil {
			return nil, errors.WrapError(err, "failed to update tokenDestino asset")
		}

		_, err = tokenOrigemAsset.Update(stub, tokenOrigemMap)
		if err != nil {
			return nil, errors.WrapError(err, "failed to update tokenOrigem asset")
		}

		// Marshal novoTokenOrigem to JSON format
		novoTokenOrigemJSON, nerr := json.Marshal(novoTokenOrigem)
		if nerr != nil {
			return nil, errors.WrapError(nerr, "failed to marshal novoTokenOrigem")
		}

		// Put novoTokenOrigem into the ledger
		err = stub.PutState(idNovoToken, novoTokenOrigemJSON)
		if err != nil {
			return nil, errors.WrapError(err, "failed to put novoTokenOrigem into the ledger")
		}

		// Marshal tokenOrigem back to JSON format
		tokenOrigemJSON, jerr := json.Marshal(tokenOrigemMap)
		if jerr != nil {
			return nil, errors.WrapError(jerr, "failed to marshal response")
		}

		return tokenOrigemJSON, nil
	},
}
