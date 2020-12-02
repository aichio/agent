package docker

import "time"

type Inspect struct {
	AppArmorProfile string          `json:"AppArmorProfile"`
	Args            []string        `json:"Args"`
	Config          Config          `json:"Config"`
	Created         time.Time       `json:"Created"`
	Driver          string          `json:"Driver"`
	ExecIDs         []string        `json:"ExecIDs"`
	HostConfig      HostConfig      `json:"HostConfig"`
	HostnamePath    string          `json:"HostnamePath"`
	HostsPath       string          `json:"HostsPath"`
	LogPath         string          `json:"LogPath"`
	ID              string          `json:"Id"`
	Image           string          `json:"Image"`
	MountLabel      string          `json:"MountLabel"`
	Name            string          `json:"Name"`
	NetworkSettings NetworkSettings `json:"NetworkSettings"`
	Path            string          `json:"Path"`
	ProcessLabel    string          `json:"ProcessLabel"`
	ResolvConfPath  string          `json:"ResolvConfPath"`
	RestartCount    int             `json:"RestartCount"`
	State           State           `json:"State"`
	Mounts          []Mounts        `json:"Mounts"`
}
type Healthcheck struct {
	Test []string `json:"Test"`
}
type Labels struct {
	ComExampleVendor  string `json:"com.example.vendor"`
	ComExampleLicense string `json:"com.example.license"`
	ComExampleVersion string `json:"com.example.version"`
}
type VolumesData struct {
}
type Volumes struct {
	VolumesData VolumesData `json:"/volumes/data"`
}
type Config struct {
	AttachStderr    bool        `json:"AttachStderr"`
	AttachStdin     bool        `json:"AttachStdin"`
	AttachStdout    bool        `json:"AttachStdout"`
	Cmd             []string    `json:"Cmd"`
	Domainname      string      `json:"Domainname"`
	Env             []string    `json:"Env"`
	Healthcheck     Healthcheck `json:"Healthcheck"`
	Hostname        string      `json:"Hostname"`
	Image           string      `json:"Image"`
	Labels          Labels      `json:"Labels"`
	MacAddress      string      `json:"MacAddress"`
	NetworkDisabled bool        `json:"NetworkDisabled"`
	OpenStdin       bool        `json:"OpenStdin"`
	StdinOnce       bool        `json:"StdinOnce"`
	Tty             bool        `json:"Tty"`
	User            string      `json:"User"`
	Volumes         Volumes     `json:"Volumes"`
	WorkingDir      string      `json:"WorkingDir"`
	StopSignal      string      `json:"StopSignal"`
	StopTimeout     int         `json:"StopTimeout"`
}
type BlkioWeightDevice struct {
}
type BlkioDeviceReadBps struct {
}
type BlkioDeviceWriteBps struct {
}
type BlkioDeviceReadIOps struct {
}
type BlkioDeviceWriteIOps struct {
}
type Options struct {
	Property1 string `json:"property1"`
	Property2 string `json:"property2"`
}
type DeviceRequests struct {
	Driver       string   `json:"Driver"`
	Count        int      `json:"Count"`
	DeviceIDs    []string `json:"DeviceIDs""`
	Capabilities []string `json:"Capabilities"`
	Options      Options  `json:"Options"`
}
type PortBindings struct {
}
type RestartPolicy struct {
	MaximumRetryCount int    `json:"MaximumRetryCount"`
	Name              string `json:"Name"`
}
type LogConfig struct {
	Type string `json:"Type"`
}
type Sysctls struct {
	NetIpv4IPForward string `json:"net.ipv4.ip_forward"`
}
type Ulimits struct {
}
type HostConfig struct {
	MaximumIOps          int                    `json:"MaximumIOps"`
	MaximumIOBps         int                    `json:"MaximumIOBps"`
	BlkioWeight          int                    `json:"BlkioWeight"`
	BlkioWeightDevice    []BlkioWeightDevice    `json:"BlkioWeightDevice"`
	BlkioDeviceReadBps   []BlkioDeviceReadBps   `json:"BlkioDeviceReadBps"`
	BlkioDeviceWriteBps  []BlkioDeviceWriteBps  `json:"BlkioDeviceWriteBps"`
	BlkioDeviceReadIOps  []BlkioDeviceReadIOps  `json:"BlkioDeviceReadIOps"`
	BlkioDeviceWriteIOps []BlkioDeviceWriteIOps `json:"BlkioDeviceWriteIOps"`
	ContainerIDFile      string                 `json:"ContainerIDFile"`
	CpusetCpus           string                 `json:"CpusetCpus"`
	CpusetMems           string                 `json:"CpusetMems"`
	CPUPercent           int                    `json:"CpuPercent"`
	CPUShares            int                    `json:"CpuShares"`
	CPUPeriod            int                    `json:"CpuPeriod"`
	CPURealtimePeriod    int                    `json:"CpuRealtimePeriod"`
	CPURealtimeRuntime   int                    `json:"CpuRealtimeRuntime"`
	Devices              []interface{}          `json:"Devices"`
	DeviceRequests       []DeviceRequests       `json:"DeviceRequests"`
	IpcMode              string                 `json:"IpcMode"`
	LxcConf              []interface{}          `json:"LxcConf"`
	Memory               int                    `json:"Memory"`
	MemorySwap           int                    `json:"MemorySwap"`
	MemoryReservation    int                    `json:"MemoryReservation"`
	KernelMemory         int                    `json:"KernelMemory"`
	OomKillDisable       bool                   `json:"OomKillDisable"`
	OomScoreAdj          int                    `json:"OomScoreAdj"`
	NetworkMode          string                 `json:"NetworkMode"`
	PidMode              string                 `json:"PidMode"`
	PortBindings         PortBindings           `json:"PortBindings"`
	Privileged           bool                   `json:"Privileged"`
	ReadonlyRootfs       bool                   `json:"ReadonlyRootfs"`
	PublishAllPorts      bool                   `json:"PublishAllPorts"`
	RestartPolicy        RestartPolicy          `json:"RestartPolicy"`
	LogConfig            LogConfig              `json:"LogConfig"`
	Sysctls              Sysctls                `json:"Sysctls"`
	Ulimits              []Ulimits              `json:"Ulimits"`
	VolumeDriver         string                 `json:"VolumeDriver"`
	ShmSize              int                    `json:"ShmSize"`
}
type Bridge struct {
	NetworkID           string `json:"NetworkID"`
	EndpointID          string `json:"EndpointID"`
	Gateway             string `json:"Gateway"`
	IPAddress           string `json:"IPAddress"`
	IPPrefixLen         int    `json:"IPPrefixLen"`
	IPv6Gateway         string `json:"IPv6Gateway"`
	GlobalIPv6Address   string `json:"GlobalIPv6Address"`
	GlobalIPv6PrefixLen int    `json:"GlobalIPv6PrefixLen"`
	MacAddress          string `json:"MacAddress"`
}
type Networks struct {
	Bridge Bridge `json:"bridge"`
}
type NetworkSettings struct {
	Bridge                 string   `json:"Bridge"`
	SandboxID              string   `json:"SandboxID"`
	HairpinMode            bool     `json:"HairpinMode"`
	LinkLocalIPv6Address   string   `json:"LinkLocalIPv6Address"`
	LinkLocalIPv6PrefixLen int      `json:"LinkLocalIPv6PrefixLen"`
	SandboxKey             string   `json:"SandboxKey"`
	EndpointID             string   `json:"EndpointID"`
	Gateway                string   `json:"Gateway"`
	GlobalIPv6Address      string   `json:"GlobalIPv6Address"`
	GlobalIPv6PrefixLen    int      `json:"GlobalIPv6PrefixLen"`
	IPAddress              string   `json:"IPAddress"`
	IPPrefixLen            int      `json:"IPPrefixLen"`
	IPv6Gateway            string   `json:"IPv6Gateway"`
	MacAddress             string   `json:"MacAddress"`
	Networks               Networks `json:"Networks"`
}
type Log struct {
	Start    time.Time `json:"Start"`
	End      time.Time `json:"End"`
	ExitCode int       `json:"ExitCode"`
	Output   string    `json:"Output"`
}
type Health struct {
	Status        string `json:"Status"`
	FailingStreak int    `json:"FailingStreak"`
	Log           []Log  `json:"Log"`
}
type State struct {
	Error      string    `json:"Error"`
	ExitCode   int       `json:"ExitCode"`
	FinishedAt time.Time `json:"FinishedAt"`
	Health     Health    `json:"Health"`
	OOMKilled  bool      `json:"OOMKilled"`
	Dead       bool      `json:"Dead"`
	Paused     bool      `json:"Paused"`
	Pid        int       `json:"Pid"`
	Restarting bool      `json:"Restarting"`
	Running    bool      `json:"Running"`
	StartedAt  time.Time `json:"StartedAt"`
	Status     string    `json:"Status"`
}
type Mounts struct {
	Name        string `json:"Name"`
	Source      string `json:"Source"`
	Destination string `json:"Destination"`
	Driver      string `json:"Driver"`
	Mode        string `json:"Mode"`
	RW          bool   `json:"RW"`
	Propagation string `json:"Propagation"`
}
