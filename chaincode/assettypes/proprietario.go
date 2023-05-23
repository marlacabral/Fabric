package assettypes

import (
	"fmt"

	"github.com/goledgerdev/cc-tools/assets"
)

var Proprietario = assets.AssetType{
	Tag:         "proprietario",
	Label:       "Proprietario",
	Description: "Proprietario da rede",

	Props: []assets.AssetProp{
		{
			// Primary key
			Required: true,
			IsKey:    true,
			Tag:      "id",
			Label:    "ID",
			DataType: "string",
			Writers:  []string{`org1MSP`}, // This means only org1 can create the asset (others can edit)
		},
		{
			// Mandatory property
			Required: true,
			Tag:      "name",
			Label:    "Nome do proprietario",
			DataType: "string",
			// Validate funcion
			Validate: func(name interface{}) error {
				nameStr := name.(string)
				if nameStr == "" {
					return fmt.Errorf("name must be non-empty")
				}
				return nil
			},
		},
	},
}
