package generator

import (
	"encoding/json"
	"os"

	"github.com/topoface/snippet-challenge/model"
)

// GenerateDefaultConfig writes default config to outputFile.
func GenerateDefaultConfig(outputFile *os.File) error {
	defaultCfg := &model.Config{}
	defaultCfg.SetDefaults()
	if data, err := json.MarshalIndent(defaultCfg, "", "  "); err != nil {
		return err
	} else if _, err := outputFile.Write(data); err != nil {
		return err
	}
	return nil
}
