package jobs

import (
	"fmt"

	"github.com/philcantcode/localmapper/capability"
	"github.com/philcantcode/localmapper/cmdb"
	"github.com/philcantcode/localmapper/utils"
)

func PreProcess(job JobSpec) {
	switch job.JobID {
	case 1:
		job = job_1_PingSweepVLANs(job)
	}

	runCapabilities(job)
}

// These packer functions should put the IPs into the jobspec for processing
func job_1_PingSweepVLANs(job JobSpec) JobSpec {
	utils.Log(fmt.Sprintf("job_1_PingSweepVLANs recieved JobSpec %+v\n", job), true)
	vlans := cmdb.SelectAllVlans()

	for _, vlan := range vlans {
		// Convert the highIP and lowIP to a list of CIDR ranges
		cidr, err := utils.IPv4RangeToCIDRRange(vlan.HighIP, vlan.LowIP)
		utils.ErrorLog(fmt.Sprintf("Couldn't convert IPs to CIDR (%s > %s)", vlan.HighIP, vlan.LowIP), err, true)

		// For each CIDR range (e.g., 10.0.0.0/8) add to the job targets
		job.Targets = append(job.Targets, cidr...)
	}

	return job
}

func runCapabilities(job JobSpec) {
	for _, capID := range job.Capabilities {
		for _, target := range job.Targets {
			runCapability := capability.SELECT_Capability_ByID(capID)

			for i, param := range runCapability.Command.Params {
				if job.DataType == param.MetaType {
					runCapability.Command.Params[i].Value = target
				}
			}

			utils.Log(fmt.Sprintf("Running Job Capability: %+v\n", runCapability), true)
			result := capability.ProcessCapability(runCapability)
			utils.Log(fmt.Sprintf("Result From Job Capability: %s\n", string(result)), true)
		}
	}
}
