package filterprocessor

import (
	"context"
	_ "embed" // used to embed config
	"fmt"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"

	"github.com/algorand/indexer/conduit"
	"github.com/algorand/indexer/conduit/data"
	"github.com/algorand/indexer/conduit/plugins"
	"github.com/algorand/indexer/conduit/plugins/processors"
	"github.com/algorand/indexer/conduit/plugins/processors/filterprocessor/expression"
	"github.com/algorand/indexer/conduit/plugins/processors/filterprocessor/fields"

	sdk "github.com/algorand/go-algorand-sdk/v2/types"
)

// PluginName to use when configuring.
const PluginName = "filter_processor"

// package-wide init function
func init() {
	processors.Register(PluginName, processors.ProcessorConstructorFunc(func() processors.Processor {
		return &FilterProcessor{}
	}))
}

// FilterProcessor filters transactions by a variety of means
type FilterProcessor struct {
	FieldFilters []fields.Filter

	logger *log.Logger
	cfg    Config
	ctx    context.Context
}

//go:embed sample.yaml
var sampleConfig string

// Metadata returns metadata
func (a *FilterProcessor) Metadata() conduit.Metadata {
	return conduit.Metadata{
		Name:         PluginName,
		Description:  "Filter transactions out of the results according to a configurable pattern.",
		Deprecated:   false,
		SampleConfig: sampleConfig,
	}
}

// Config returns the config
func (a *FilterProcessor) Config() string {
	s, _ := yaml.Marshal(a.cfg)
	return string(s)
}

// Init initializes the filter processor
func (a *FilterProcessor) Init(ctx context.Context, _ data.InitProvider, cfg plugins.PluginConfig, logger *log.Logger) error {
	a.logger = logger
	a.ctx = ctx

	err := cfg.UnmarshalConfig(&a.cfg)
	if err != nil {
		return fmt.Errorf("filter processor init error: %w", err)
	}

	// configMaps is the "- any: ...." portion of the filter config
	for _, configMaps := range a.cfg.Filters {

		// We only want one key in the map (i.e. either "any" or "all").  The reason we use a list is that want
		// to maintain ordering of the filters and a straight-up map doesn't do that.
		if len(configMaps) != 1 {
			return fmt.Errorf("filter processor Init(): illegal filter tag formation.  tag length was: %d", len(configMaps))
		}

		for key, subConfigs := range configMaps {

			if !fields.ValidFieldOperation(key) {
				return fmt.Errorf("filter processor Init(): filter key was not a valid value: %s", key)
			}

			var searcherList []*fields.Searcher

			for _, subConfig := range subConfigs {

				t, err := fields.LookupFieldByTag(subConfig.FilterTag, &sdk.SignedTxnWithAD{})
				if err != nil {
					return err
				}

				exp, err := expression.MakeExpression(subConfig.ExpressionType, subConfig.Expression, t)
				if err != nil {
					return fmt.Errorf("filter processor Init(): could not make expression: %w", err)
				}

				searcher, err := fields.MakeFieldSearcher(exp, subConfig.ExpressionType, subConfig.FilterTag, a.cfg.SearchInner)
				if err != nil {
					return fmt.Errorf("filter processor Init(): error making field searcher - %w", err)
				}

				searcherList = append(searcherList, searcher)
			}

			ff := fields.Filter{
				Op:        fields.Operation(key),
				Searchers: searcherList,
			}

			a.FieldFilters = append(a.FieldFilters, ff)

		}
	}

	return nil

}

// Close a no-op for this processor
func (a *FilterProcessor) Close() error {
	return nil
}

// Process processes the input data
func (a *FilterProcessor) Process(input data.BlockData) (data.BlockData, error) {
	var err error
	payset := input.Payset
	for _, searcher := range a.FieldFilters {
		payset, err = searcher.SearchAndFilter(payset)
		if err != nil {
			return data.BlockData{}, err
		}
	}
	input.Payset = payset
	return input, err
}
