package constant

const (
	EnvName = "B0B_ENV"
	// EtcdConfig /conf/json/dev
	EtcdConfig = "/conf/%s/%s"
	// EtcdServers /servers/[env]/[Type]/[serversName]/[nodeId]
	EtcdServers    = "/servers/%s/%s/%s/%s"
	EtcdServersPre = "/servers/%s/"

	ProtoHeader = "Proto"
)

type ENV string

const (
	Dev  ENV = "dev"
	Test ENV = "test"
	Prod ENV = "prod"
)

type OtelType string

const (
	HttpOtelType OtelType = "http"
	RpcOtelType  OtelType = "rpc"
)
