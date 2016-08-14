package fscommit

type Commit []CommitAction

type CommitAction interface {
	Execute() error
	Revert() error
	Finalize() error
}

type CommitResult struct {
	CommitError    error
	RevertErrors   []error
	FinalizeErrors []error
}

// Executes the given commit.
//
// Execution will continue until the end of the list of actions or until an error occurs in one of the given actions,
// leading to the reverse of all actions executed so far.
func (commit Commit) Execute() (result CommitResult) {
	finishedActions := []CommitAction{}
	result = CommitResult{
		RevertErrors:   []error{},
		FinalizeErrors: []error{},
	}
	for _, action := range commit {
		result.CommitError = action.Execute()
		if result.CommitError != nil {
			break
		}
		finishedActions = append(finishedActions, action)
	}

	if result.CommitError != nil {
		// Error, reverse all actions done so far!
		for i := len(finishedActions) - 1; i >= 0; i-- {
			action := finishedActions[i]
			if err := action.Revert(); err != nil {
				result.RevertErrors = append(result.RevertErrors, err)
			}
		}
	} else {
		// No error, finalize all actions to actually complete what they were intended to do instead of being safe
		for _, action := range finishedActions {
			if err := action.Finalize(); err != nil {
				result.FinalizeErrors = append(result.FinalizeErrors, err)
			}
		}
	}

	return
}
