//+build integration

package integration

import (
	. "gopkg.in/check.v1"

	"gotest.tools/v3/fs"
	"gotest.tools/v3/icmd"
)

type ConfigSuite struct {
	tmpDir *fs.Dir
}

var _ = Suite(&ConfigSuite{})

func (s *ConfigSuite) SetUpTest(c *C) {
	s.tmpDir = fs.NewDir(c, "gomplate-inttests",
		fs.WithFile(".gomplate.yaml", "in: hello world\n"),
	)
}

func (s *ConfigSuite) TearDownTest(c *C) {
	s.tmpDir.Remove()
}

func (s *ConfigSuite) TestReadsFromConfigFile(c *C) {
	result := icmd.RunCmd(icmd.Command(GomplateBin), func(cmd *icmd.Cmd) {
		cmd.Dir = s.tmpDir.Path()
	})
	result.Assert(c, icmd.Expected{ExitCode: 0, Out: "hello world"})
}
