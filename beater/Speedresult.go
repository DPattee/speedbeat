package beater

// speedtest result
type Speedresult struct {
	Type       string
	Timestamp  string
	Ping       Ping
	Download   Download
	Upload     Upload
	PacketLoss float64
	ISP        string
	Interface  Interface
	Server     Server
	Result     Result
}

type Ping struct {
	Jitter  float64
	Latency float64
}

type Download struct {
	Bandwidth int
	Bytes     int
	Elapsed   int
}

type Upload struct {
	Bandwidth int
	Bytes     int
	Elapsed   int
}

type Interface struct {
	InternalIp string
	Name       string
	MACAddr    string
	IsVpn      bool
	ExternalIp string
}

type Server struct {
	ID       int
	Name     string
	Location string
	Country  string
	Host     string
	Port     int
	IP       string
}

type Result struct {
	ID  string
	URL string
}
