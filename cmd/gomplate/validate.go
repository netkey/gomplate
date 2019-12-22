package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func notTogether(cmd *cobra.Command, flags ...string) error {
	found := ""
	for _, flag := range flags {
		f := cmd.Flag(flag)
		if f != nil && f.Changed {
			if found != "" {
				a := make([]string, len(flags))
				for i := range a {
					a[i] = "--" + flags[i]
				}
				return fmt.Errorf("only one of these flags is supported at a time: %s", strings.Join(a, ", "))
			}
			found = flag
		}
	}
	return nil
}

func mustTogether(cmd *cobra.Command, left, right string) error {
	l := cmd.Flag(left)
	if l != nil && l.Changed {
		r := cmd.Flag(right)
		if r != nil && !r.Changed {
			return fmt.Errorf("--%s must be set when --%s is set", right, left)
		}
	}

	return nil
}

func validateOpts(cmd *cobra.Command, args []string) (err error) {
	err = notTogether(cmd, "in", "file", "input-dir")
	if err == nil {
		err = notTogether(cmd, "out", "output-dir", "output-map", "exec-pipe")
	}

	if err == nil && len(conf.InputFiles) != len(conf.OutputFiles) {
		err = fmt.Errorf("must provide same number of --out (%d) as --file (%d) options", len(conf.OutputFiles), len(conf.InputFiles))
	}

	if err == nil && cmd.Flag("exec-pipe").Changed && len(args) == 0 {
		err = fmt.Errorf("--exec-pipe may only be used with a post-exec command after --")
	}

	if err == nil {
		err = mustTogether(cmd, "output-dir", "input-dir")
	}

	if err == nil {
		err = mustTogether(cmd, "output-map", "input-dir")
	}

	return err
}
