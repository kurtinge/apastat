package collector

type Collector interface {
	GetStats() (*ServerStatus, error)
}

type ServerStatus struct {
	Uptime      string
	RequestSec  int
	ServerSlots []Slot
}

type Slot struct {
	ServerSlot          string
	Pid                 int
	Mode                ServerMode
	Cpu                 float64
	SecondsSinceRequest int
	Client              string
	Protocol            string
	Vhost               string
	Request             string
}

type ServerMode string

const (
	ServerModeWaiting             ServerMode = "_"
	ServerModeStartingUp          ServerMode = "S"
	ServerModeReadingRequest      ServerMode = "R"
	ServerModeSendingReply        ServerMode = "W"
	ServerModeKeepalive           ServerMode = "K"
	ServerModeDNSLookup           ServerMode = "D"
	ServerModeClosingConnection   ServerMode = "C"
	ServerModeLogging             ServerMode = "L"
	ServerModeGracefullyFinishing ServerMode = "G"
	ServerModeIdle                ServerMode = "I"
	ServerModeOpenSlot            ServerMode = "."
)
