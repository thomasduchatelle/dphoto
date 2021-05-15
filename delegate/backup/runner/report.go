package runner

import "sync"

// Report is returned with the backup is finished ; if Errors isn't empty, backup has failed or has been interupted.
type Report struct {
	errorsMutex sync.Mutex
	Errors      []error
}

func (r *Report) AppendError(err error) {
	if err == nil {
		return
	}

	r.errorsMutex.Lock()
	defer r.errorsMutex.Unlock()

	r.Errors = append(r.Errors, err)
}
