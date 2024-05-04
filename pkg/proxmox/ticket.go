package proxmox

type Ticket struct {
	Data struct {
		CSRFPreventionToken string `json:"CSRFPreventionToken"`
		Username            string `json:"username"`
		Ticket              string `json:"ticket"`
		Cap                 struct {
			SoftwareDefinedNetwork struct {
				SDNAudit          int `json:"SDN.Audit"`
				PermissionsModify int `json:"Permissions.Modify"`
				SDNAllocate       int `json:"SDN.Allocate"`
				SDNUse            int `json:"SDN.Use"`
			} `json:"sdn"`
			VirtualMachines struct {
				VMBackup           int `json:"VM.Backup"`
				VMConfigMemory     int `json:"VM.Config.Memory"`
				VMConfigOptions    int `json:"VM.Config.Options"`
				VMMigrate          int `json:"VM.Migrate"`
				VMConfigHWType     int `json:"VM.Config.HWType"`
				VMConfigCloudinit  int `json:"VM.Config.Cloudinit"`
				VMConsole          int `json:"VM.Console"`
				VMSnapshotRollback int `json:"VM.Snapshot.Rollback"`
				VMConfigCPU        int `json:"VM.Config.CPU"`
				PermissionsModify  int `json:"Permissions.Modify"`
				VMConfigDisk       int `json:"VM.Config.Disk"`
				VMConfigNetwork    int `json:"VM.Config.Network"`
				VMAudit            int `json:"VM.Audit"`
				VMClone            int `json:"VM.Clone"`
				VMMonitor          int `json:"VM.Monitor"`
				VMAllocate         int `json:"VM.Allocate"`
				VMSnapshot         int `json:"VM.Snapshot"`
				VMPowerMgmt        int `json:"VM.PowerMgmt"`
				VMConfigCDROM      int `json:"VM.Config.CDROM"`
			} `json:"vms"`
			Storage struct {
				DatastoreAudit            int `json:"Datastore.Audit"`
				DatastoreAllocateTemplate int `json:"Datastore.AllocateTemplate"`
				PermissionsModify         int `json:"Permissions.Modify"`
				DatastoreAllocate         int `json:"Datastore.Allocate"`
				DatastoreAllocateSpace    int `json:"Datastore.AllocateSpace"`
			} `json:"storage"`
			DataCenter struct {
				SDNAllocate int `json:"SDN.Allocate"`
				SysAudit    int `json:"Sys.Audit"`
				SDNUse      int `json:"SDN.Use"`
				SDNAudit    int `json:"SDN.Audit"`
			} `json:"dc"`
			Nodes struct {
				PermissionsModify int `json:"Permissions.Modify"`
				SysSyslog         int `json:"Sys.Syslog"`
				SysPowerMgmt      int `json:"Sys.PowerMgmt"`
				SysIncoming       int `json:"Sys.Incoming"`
				SysConsole        int `json:"Sys.Console"`
				SysAudit          int `json:"Sys.Audit"`
				SysModify         int `json:"Sys.Modify"`
			} `json:"nodes"`
			Access struct {
				GroupAllocate     int `json:"Group.Allocate"`
				PermissionsModify int `json:"Permissions.Modify"`
				UserModify        int `json:"User.Modify"`
			} `json:"access"`
		} `json:"cap"`
	} `json:"data"`
}
