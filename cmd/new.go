package cmd

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/b4b4r07/gist/pkg/gist"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

type newCmd struct {
	meta

	private bool
}

// newNewCmd creates a new new command
func newNewCmd() *cobra.Command {
	c := &newCmd{}

	newCmd := &cobra.Command{
		Use:                   "new",
		Short:                 "Create gist file",
		Aliases:               []string{},
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		SilenceErrors:         true,
		Args:                  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.meta.init(args); err != nil {
				return err
			}
			return c.run(args)
		},
	}

	f := newCmd.Flags()
	f.BoolVarP(&c.private, "private", "p", false, "make private")

	return newCmd
}

func (c *newCmd) run(args []string) error {
	validate := func(input string) error {
		if len(input) < 3 {
			return errors.New("Filename must have more than 3 characters")
		}
		return nil
	}

	files := make(map[string]string)
	for _, arg := range args {
		f, err := os.Open(arg)
		if err != nil {
			return err
		}
		defer f.Close()
		b, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}
		prompt := promptui.Prompt{
			Label:     "Filename",
			Validate:  validate,
			AllowEdit: true,
			Default:   filepath.Base(arg),
		}
		name, err := prompt.Run()
		if err != nil {
			return err
		}
		files[name] = string(b)
	}

	prompt := promptui.Prompt{Label: "Description"}
	desc, err := prompt.Run()
	if err != nil {
		return err
	}

	return gist.Create(gist.Page{
		Files:       files,
		Description: desc,
		Public:      !c.private,
	})
}
