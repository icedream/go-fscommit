package fscommit

import "os"

type removeAllAction struct {
	path    string
	tmppath string
}

func (action *removeAllAction) Execute() (err error) {
	// Move file/folder out of the way
	tmppath, err := tmpPath(os.TempDir(), "fscommit")
	if err != nil {
		return
	}
	err = os.Rename(action.path, tmppath)
	if os.IsNotExist(err) {
		// File/folder does not exist, so nothing to do here!
		err = nil
		return
	} else if err != nil {
		return
	}
	action.tmppath = tmppath
	return
}

func (action *removeAllAction) Revert() (err error) {
	// If tmppath is set, that means an old file/folder has been moved to the backup place. Restore that file/folder!
	if len(action.tmppath) > 0 {
		return os.Rename(action.tmppath, action.path)
	}

	return nil
}

func (action *removeAllAction) Finalize() error {
	// If tmppath is set, that means an old file/folder has been moved to the backup place. Delete that file/folder!
	if len(action.tmppath) > 0 {
		return os.RemoveAll(action.tmppath)
	}

	return nil
}

func RemoveAll(path string) CommitAction {
	return &removeAllAction{path: path}
}
