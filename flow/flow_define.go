package flow

import (
	"forjj/forjfile"
	"github.com/forj-oss/forjj-modules/trace"
	"fmt"
	"forjj/utils"
)

type FlowDefine struct { // Yaml structure
	Name   string
	Title  string // Flow title
	Define map[string]FlowPluginTypeDef
	OnRepo map[string]FlowTaskDef `yaml:"on-repo-do"`
	OnForj map[string]FlowTaskDef `yaml:"on-forjfile-do"`
}

func (fd *FlowDefine)apply(repo *forjfile.RepoStruct, Forjfile *forjfile.Forge) error {
	bInError := false

	var tasks map[string]FlowTaskDef
	if repo == nil {
		tasks = fd.OnForj
	} else {
		tasks = fd.OnRepo
	}

	for _, flowTask := range tasks {
		onWhat := "Forjfile"
		if repo != nil {
			onWhat = fmt.Sprintf("repository '%s'", repo.GetString("name"))
		} else {

		}
		gotrace.Trace("flow '%s': %s on %s", fd.Name, flowTask.Description, onWhat)

		task_to_set, err := flowTask.if_section(repo, Forjfile)
		if err != nil {
			gotrace.Error("Flow '%s' - if section: Unable to apply flow task '%s'.", fd.Name, err)
			bInError = true
			continue
		}

		if ! task_to_set {
			gotrace.Trace("Flow task not applied to %s. if condition fails.", onWhat)
			continue
		}

		gotrace.Trace("Flow task '%s' applying to %s.", fd.Name, onWhat)

		tmpl_data := New_FlowTaskModel(repo, Forjfile.Forjfile())

		if flowTask.List == nil {
			if err := flowTask.Set.apply(tmpl_data, Forjfile); err != nil {
				gotrace.Error("Unable to apply flow task '%s' on %s. %s", fd.Name, onWhat, err)
				continue
			}
			gotrace.Trace("flow task '%s' applied on %s.", fd.Name, repo.GetString("name"))
			continue
		}

		// Load list
		max := make([]int, len(flowTask.List))

		for index, taskList := range flowTask.List {
			taskList.list = taskList.Get(repo, Forjfile)
			max[index] = len(taskList.list)
		}

		// Loop on list and set CurrentList
		looplist := utils.NewMLoop(max...)
		tmpl_data.List = make(map[string]interface{})
		for !looplist.Eol() {
			for index, pos := range looplist.Cur() {
				flowTaskList := flowTask.List[index]
				tmpl_data.List[flowTaskList.Name] = flowTaskList.list[pos]
			}

			if err := flowTask.Set.apply(tmpl_data, Forjfile); err != nil {
				gotrace.Error("Unable to apply flow task '%s' on %s. %s", fd.Name, onWhat, err)
			} else {
				gotrace.Trace("flow task '%s' applied on %s.", fd.Name, onWhat)
			}

			looplist.Increment()
		}
	}
	if bInError {
		return fmt.Errorf("Failed to apply '%s'. Errors detected.", fd.Name)
	}

	return nil
}

func (ftd *FlowTaskDef)if_section(repo *forjfile.RepoStruct, Forjfile *forjfile.Forge) (task_to_set bool, _ error) {
	task_to_set = true
	if ftd.If != nil {
		for _, ftif := range ftd.If {
			if v, err := ftif.IfEvaluate(repo, Forjfile); err != nil {
				return false, err
			} else if !v {
				task_to_set = false
				break
			}
		}
	}
	return
}
