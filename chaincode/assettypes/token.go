import (
	"asset"
	"assettypes"
	"fmt"
	"github.com/goledgerdev/cc-tools/assets"
	"github.com/goledgerdev/cc-tools/transactions"
)

// Definindo o tipo de ativo "Token"
var Token = assets.AssetType{
	Tag:         "token",
	Label:       "Token",
	Description: "Token",

	Props: []assets.AssetProp{
		{
			// Primary key
			Required:     true,
			IsKey:        true,
			Tag:          "id",
			Label:        "ID",
			DefaultValue: 0,
			DataType:     "string",
			Writers:      []string{"org2MSP"}, // Isso significa que apenas org2 pode criar o ativo (outros podem editar)
		},
		{
			// Composite Key
			Tag:      "currentTenant",
			Label:    "currentTenant",
			DataType: "->proprietario",
		},
		{
			Tag:      "burned",
			Label:    "Burned",
			DataType: "boolean",
		},
	},
}

// Método para verificar se um token foi queimado (utilizado)
func IsTokenBurned(ctx assets.TransactionContextInterface, tokenID string) (bool, error) {
	token, err := ctx.GetAsset("token", tokenID)
	if err != nil {
		return false, fmt.Errorf("erro ao obter o token: %s", err.Error())
	}
	if token == nil {
		return false, fmt.Errorf("token não encontrado")
	}

	burned, err := token.GetBoolPropertyValue("burned")
	if err != nil {
		return false, fmt.Errorf("erro ao obter o valor de 'burned': %s", err.Error())
	}

	return burned, nil
}
