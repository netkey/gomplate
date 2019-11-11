/*
The gomplate command

*/
package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/signal"

	"github.com/hairyhenderson/gomplate/v3"
	"github.com/hairyhenderson/gomplate/v3/env"
	"github.com/hairyhenderson/gomplate/v3/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile      string
	conf         gconfig
	printVer     bool
	postRunInput *bytes.Buffer
)

const (
	defaultConfigName = ".gomplate"
	defaultConfigFile = ".gomplate.yaml"
)

func printVersion(name string) {
	fmt.Printf("%s version %s\n", name, version.Version)
}

// postRunExec - if templating succeeds, the command following a '--' will be executed
func postRunExec(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		name := args[0]
		args = args[1:]
		// nolint: gosec
		c := exec.Command(name, args...)
		if conf.ExecPipe {
			c.Stdin = postRunInput
		} else {
			c.Stdin = os.Stdin
		}
		c.Stderr = os.Stderr
		c.Stdout = os.Stdout

		// make sure all signals are propagated
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs)
		go func() {
			// Pass signals to the sub-process
			sig := <-sigs
			if c.Process != nil {
				// nolint: gosec
				_ = c.Process.Signal(sig)
			}
		}()

		return c.Run()
	}
	return nil
}

// optionalExecArgs - implements cobra.PositionalArgs. Allows extra args following
// a '--', but not otherwise.
func optionalExecArgs(cmd *cobra.Command, args []string) error {
	if cmd.ArgsLenAtDash() == 0 {
		return nil
	}
	return cobra.NoArgs(cmd, args)
}

// process --include flags - these are analogous to specifying --exclude '*',
// then the inverse of the --include options.
func processIncludes(includes, excludes []string) []string {
	out := []string{}
	// if any --includes are set, we start by excluding everything
	if len(includes) > 0 {
		out = make([]string, 1+len(includes))
		out[0] = "*"
	}
	for i, include := range includes {
		// includes are just the opposite of an exclude
		out[i+1] = "!" + include
	}
	out = append(out, excludes...)
	return out
}

func newGomplateCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "gomplate",
		Short: "Process text files with Go templates",
		// TODO: cmd.SetVersionTemplate(s string)
		// Version: version.Version,
		PreRunE: validateOpts,
		RunE: func(cmd *cobra.Command, args []string) error {
			// fmt.Fprintf(os.Stderr, "viper settings: %#v\nds: %#v\n", viper.AllSettings(), viper.GetStringSlice("datasource"))
			conf := gconfig{}
			if err := viper.Unmarshal(&conf); err != nil {
				return fmt.Errorf("unmarshaling error on %#v: %w", conf, err)
			}
			if printVer {
				printVersion(cmd.Name())
				return nil
			}
			if conf.Verbose {
				// nolint: errcheck
				fmt.Fprintf(os.Stderr, "%s version %s, build %s\nconfig is:\n%s\n\n",
					cmd.Name(), version.Version, version.GitCommit,
					conf)
			}

			// support --include
			conf.ExcludeGlob = processIncludes(conf.Includes, conf.ExcludeGlob)

			if conf.ExecPipe {
				postRunInput = &bytes.Buffer{}
				conf.out = postRunInput
			}
			err := gomplate.RunTemplates(conf.toConfig())
			cmd.SilenceErrors = true
			cmd.SilenceUsage = true
			if conf.Verbose {
				// nolint: errcheck
				fmt.Fprintf(os.Stderr, "rendered %d template(s) with %d error(s) in %v\n",
					gomplate.Metrics.TemplatesProcessed, gomplate.Metrics.Errors, gomplate.Metrics.TotalRenderDuration)
			}
			return err
		},
		PostRunE: postRunExec,
		Args:     optionalExecArgs,
	}
	return rootCmd
}

func initFlags(command *cobra.Command) {
	cobra.OnInitialize(initConfig)

	command.Flags().SortFlags = false

	command.Flags().StringSliceVarP(&conf.DataSources, "datasource", "d", nil, "`datasource` in alias=URL form. Specify multiple times to add multiple sources.")
	command.Flags().StringSliceVarP(&conf.DataSourceHeaders, "datasource-header", "H", nil, "HTTP `header` field in 'alias=Name: value' form to be provided on HTTP-based data sources. Multiples can be set.")

	command.Flags().StringSliceVarP(&conf.Contexts, "context", "c", nil, "pre-load a `datasource` into the context, in alias=URL form. Use the special alias `.` to set the root context.")

	command.Flags().StringSliceVar(&conf.Plugins, "plugin", nil, "plug in an external command as a function in name=path form. Can be specified multiple times")

	command.Flags().StringSliceVarP(&conf.InputFiles, "file", "f", []string{"-"}, "Template `file` to process. Omit to use standard input, or use --in or --input-dir")
	command.Flags().StringVarP(&conf.Input, "in", "i", "", "Template `string` to process (alternative to --file and --input-dir)")
	command.Flags().StringVar(&conf.InputDir, "input-dir", "", "`directory` which is examined recursively for templates (alternative to --file and --in)")

	command.Flags().StringSliceVar(&conf.ExcludeGlob, "exclude", []string{}, "glob of files to not parse")
	command.Flags().StringSliceVar(&conf.Includes, "include", []string{}, "glob of files to parse")

	command.Flags().StringSliceVarP(&conf.OutputFiles, "out", "o", []string{"-"}, "output `file` name. Omit to use standard output.")
	command.Flags().StringSliceVarP(&conf.Templates, "template", "t", []string{}, "Additional template file(s)")
	command.Flags().StringVar(&conf.OutputDir, "output-dir", ".", "`directory` to store the processed templates. Only used for --input-dir")
	command.Flags().StringVar(&conf.OutputMap, "output-map", "", "Template `string` to map the input file to an output path")
	command.Flags().StringVar(&conf.OutMode, "chmod", "", "set the mode for output file(s). Omit to inherit from input file(s)")

	command.Flags().BoolVar(&conf.ExecPipe, "exec-pipe", false, "pipe the output to the post-run exec command")

	ldDefault := env.Getenv("GOMPLATE_LEFT_DELIM", "{{")
	rdDefault := env.Getenv("GOMPLATE_RIGHT_DELIM", "}}")
	command.Flags().StringVar(&conf.LDelim, "left-delim", ldDefault, "override the default left-`delimiter` [$GOMPLATE_LEFT_DELIM]")
	command.Flags().StringVar(&conf.RDelim, "right-delim", rdDefault, "override the default right-`delimiter` [$GOMPLATE_RIGHT_DELIM]")

	command.Flags().BoolVarP(&conf.Verbose, "verbose", "V", false, "output extra information about what gomplate is doing")

	command.Flags().BoolVarP(&printVer, "version", "v", false, "print the version")

	command.Flags().StringVar(&cfgFile, "config", defaultConfigFile, "config file (overridden by commandline flags)")

	viper.BindPFlags(command.Flags())
}

func initConfig() {
	configRequired := false
	if cfgFile != defaultConfigFile {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
		configRequired = true
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName(defaultConfigName)
	}

	err := viper.ReadInConfig()
	if err != nil && configRequired {
		panic(err)
	}
	if err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func main() {
	command := newGomplateCmd()
	initFlags(command)
	if err := command.Execute(); err != nil {
		// nolint: errcheck
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
