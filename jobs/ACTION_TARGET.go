package jobs

func SendTo_ActionTarget(job Job) {
	if job.Spec[job.State] != ACTION_TARGET {
		PassToNext(job)
	}

	job.State++

	switch job.Spec[job.State] {
	case SOURCE_VLAN:

	}
}
