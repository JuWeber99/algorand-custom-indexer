package initialize

import (
	_ "embed"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algorand/indexer/conduit/pipeline"
	"github.com/algorand/indexer/conduit/plugins/exporters/filewriter"
	noopExporter "github.com/algorand/indexer/conduit/plugins/exporters/noop"
	algodimporter "github.com/algorand/indexer/conduit/plugins/importers/algod"
	fileimporter "github.com/algorand/indexer/conduit/plugins/importers/filereader"
	"github.com/algorand/indexer/conduit/plugins/processors/filterprocessor"
	noopProcessor "github.com/algorand/indexer/conduit/plugins/processors/noop"
)

//go:embed conduit.test.init.default.yml
var defaultYml string

// TestInitDataDirectory tests the initialization of the data directory
func TestInitDataDirectory(t *testing.T) {
	verifyFile := func(file string, importer string, exporter string, processors []string) {
		require.FileExists(t, file)
		data, err := os.ReadFile(file)
		require.NoError(t, err)
		var cfg pipeline.Config
		require.NoError(t, yaml.Unmarshal(data, &cfg))
		assert.Equal(t, importer, cfg.Importer.Name)
		assert.Equal(t, exporter, cfg.Exporter.Name)
		require.Equal(t, len(processors), len(cfg.Processors))
		for i := range processors {
			assert.Equal(t, processors[i], cfg.Processors[i].Name)
		}
	}

	// Defaults
	dataDirectory := t.TempDir()
	runConduitInit(dataDirectory, "", []string{}, "")
	verifyFile(fmt.Sprintf("%s/conduit.yml", dataDirectory), algodimporter.PluginName, filewriter.PluginName, nil)

	// Explicit defaults
	dataDirectory = t.TempDir()
	runConduitInit(dataDirectory, algodimporter.PluginName, []string{noopProcessor.PluginName}, filewriter.PluginName)
	verifyFile(fmt.Sprintf("%s/conduit.yml", dataDirectory), algodimporter.PluginName, filewriter.PluginName, []string{noopProcessor.PluginName})

	// Different
	dataDirectory = t.TempDir()
	runConduitInit(dataDirectory, fileimporter.PluginName, []string{noopProcessor.PluginName, filterprocessor.PluginName}, noopExporter.PluginName)
	verifyFile(fmt.Sprintf("%s/conduit.yml", dataDirectory), fileimporter.PluginName, noopExporter.PluginName, []string{noopProcessor.PluginName, filterprocessor.PluginName})
}
