package updater

//UpdateOptions defines options for UpdateModules process
type UpdateOptions struct {
	IsIndirect bool //Updater all direct and indirect modules
	IsPrompt   bool //Confirm in prompt updates with changes
	IsSelect   bool //Select direct modules to update
	IsTests    bool //Run local tests after updating each module and rollback in case of errors
	IsVerbose  bool //Print output
}
