package core

type NuxDevice struct {
	ID               string    `json:"id"`        // Nux id
	Name             string    `json:"name"`      // Nux name
	BootTime         int64     `json:"boot_time"` // Last boot time
	Address          string    `json:"address"`
	Disks            []NuxDisk `json:"disks"`      // Disk list
	DiskTotalBytes   int64     `json:"disk_total"` // Disk total bytes
	DiskFreeBytes    int64     `json:"disk_free"`  // Disk free bytes
	MemoryTotalBytes int64     `json:"mem_total"`  // Memory total bytes
	MemoryFreeBytes  int64     `json:"mem_free"`   // Memory free bytes
	CPURate          int32     `json:"cpu_rate"`   // CPU use rate
	CPUTemp          int32     `json:"cpu_temp"`   // CPU temperature
}

type NuxDisk struct {
	ID             string `json:"id"`         // Disk id
	Name           string `json:"name"`       // Disk name
	UnixPath       string `json:"path"`       // Disk relative path: /{DiskId}
	UnixAbsolute   string `json:"absolute"`   // Disk absolute path: /c/uki, Unix Format
	DiskTotalBytes int64  `json:"disk_total"` // Disk total bytes
	DiskFreeBytes  int64  `json:"disk_free"`  // Disk free bytes
}

type NuxFile struct {
	Name         string `json:"name"`        // File name
	UnixPath     string `json:"path"`        // File relative path: /{DiskId}/a/b/c, Unix Format
	UnixAbsolute string `json:"absolute"`    // File absolute path: /c/uki/a/b/c, Unix Format
	Size         int64  `json:"size"`        // File size
	UpdateTime   int64  `json:"update_time"` // Last update time
	IsDir        bool   `json:"is_dir"`      // Is Directory
	MD5          string `json:"md5"`         // File MD5
	Thumbnail    string `json:"thumbnail"`   // Thumbnail url for image or video
	IsCollected  bool   `json:"is_col"`      // Is collected
	GhostUrl     string `json:"ghost_url"`   // Only for recent delete
}
