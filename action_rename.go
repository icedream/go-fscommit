package fscommit

import "os"

type renameAction struct {
	from    string
	to      string
	tmppath string
}

func (action *renameAction) Execute() (err error) {
	// 1. Move old file/folder to backup place for Revert()
	tmppath, err := tmpPath(os.TempDir(), "fscommit")
	if err != nil {
		return
	}
	err = os.Rename(action.to, tmppath)
	if err != nil && !os.IsNotExist(err) {
		return
	} else if err == nil {
		action.tmppath = tmppath
	}

	// 2. Move new file in old file's place
	err = os.Rename(action.from, action.to)
	if err != nil {
		action.restoreBackup()
	}
	return
}

func (action *renameAction) Revert() (err error) {
	err = os.Rename(action.to, action.from)

	// If tmppath is set, that means an old file has been moved to the backup place. Restore that file/folder!
	action.restoreBackup()

	return nil
}

func (action *renameAction) restoreBackup() error {
	if len(action.tmppath) > 0 {
		return os.Rename(action.tmppath, action.to)
	}
	return nil
}

func (action *renameAction) Finalize() error {
	// If tmppath is set, that means an old file has been moved to the backup place. Delete that file/folder!
	if len(action.tmppath) > 0 {
		return os.RemoveAll(action.tmppath)
	}

	return nil
}

func Rename(from string, to string) CommitAction {
	return &renameAction{from: from, to: to}
}
