package module

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

const surveyPageSize = 25

// SelectCLI allows interactively select modules to update
func (mds *Modules) SelectCLI() error {
	prompt := &survey.MultiSelect{
		Message:  "Choose which modules to update",
		Options:  mds.surveyOptions(),
		PageSize: surveyPageSize,
	}

	var chs []int
	if err := survey.AskOne(prompt, &chs); err != nil {
		return err
	}

	var sltMds Modules
	for _, n := range chs {
		sltMds = append(sltMds, (*mds)[n])
	}
	*mds = sltMds

	return nil
}

func (mds Modules) surveyOptions() []string {
	var ops []string
	for _, md := range mds {
		ops = append(ops, fmt.Sprintf("%s %s -> %s", md.Path, md.Version, md.newVersion()))
	}
	return ops
}
