package main

import (
	"io"

	"github.com/hairyhenderson/gomplate/v3"
)

// the gomplate config file/structure
type gconfig struct {
	Verbose  bool
	ExecPipe bool     `mapstructure:"exec-pipe"`
	Includes []string `mapstructure:"include"`

	Input       string   `mapstructure:"in"`
	InputFiles  []string `mapstructure:"file"`
	InputDir    string   `mapstructure:"input-dir"`
	ExcludeGlob []string `mapstructure:"exclude"`
	OutputFiles []string `mapstructure:"out"`
	OutputDir   string   `mapstructure:"output-dir"`
	OutputMap   string   `mapstructure:"output-map"`
	OutMode     string   `mapstructure:"chmod"`
	out         io.Writer

	DataSources       []string `mapstructure:"datasource"`
	DataSourceHeaders []string `mapstructure:"datasource-header"`
	Contexts          []string `mapstructure:"context"`

	Plugins []string `mapstructure:"plugin"`

	LDelim string `mapstructure:"left-delim"`
	RDelim string `mapstructure:"right-delim"`

	Templates []string `mapstructure:"template"`
}

func (c gconfig) toConfig() *gomplate.Config {
	g := &gomplate.Config{
		Input:       c.Input,
		InputFiles:  c.InputFiles,
		InputDir:    c.InputDir,
		ExcludeGlob: c.ExcludeGlob,
		OutputFiles: c.OutputFiles,
		OutputDir:   c.OutputDir,
		OutputMap:   c.OutputMap,
		OutMode:     c.OutMode,
		Out:         c.out,

		DataSources:       c.DataSources,
		DataSourceHeaders: c.DataSourceHeaders,
		Contexts:          c.Contexts,

		Plugins: c.Plugins,

		LDelim: c.LDelim,
		RDelim: c.RDelim,

		Templates: c.Templates,
	}

	return g
}

func (c gconfig) String() string {
	return c.toConfig().String()
}
