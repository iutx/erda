package manager

import (
	"database/sql/driver"
	"encoding/json"
)

type ClusterMetaData struct {
	ClusterName    string        // 集群名称
	DisplayName    string        // 集群显示名称
	ClusterType    string        // 集群类型 'dcos','edas','k8s','localdocker','swarm'
	ManageType     string        // 集群管理方式 'inet', 'kubeconfig', 'sa', 'proxy'
	ManageConfig   *ManageConfig // 管理配置信息, address ca client-cert client-key token
	WildcardDomain string        // 泛域名
	CloudVendor    string        // 云服务供应商，'erda', 'alicloud-ecs', 'alicoud-cs', 'alicoud-cs-managed'
	Options        *Options
}

func (c *ClusterMetaData) TableName() string {
	return "co_clusters"
}

type Options struct {
	CpuNumQuota              int
	EnableTag                bool
	EnableOrg                bool
	EnableWorkspace          bool
	SparkVersion             float64
	DevCpuSubscribeRatio     int64
	TestCpuSubscribeRatio    int64
	StagingCpuSubscribeRatio int64
	EdasConfig               *EdasConfig
}

func (o *Options) Value() (driver.Value, error) {
	b, err := json.Marshal(o)
	return string(b), err
}

func (o *Options) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), o)
}

type EdasConfig struct {
	UserName string
	Password string
}

//func (e *EdasConfig) Value() (driver.Value, error) {
//	b, err := json.Marshal(o)
//	return string(b), err
//}
//
//func (e *EdasConfig) Scan(input interface{}) error {
//	return json.Unmarshal(input.([]byte), o)
//}

type ManageConfig struct {
	Address  string `json:"address"`
	Insecure bool   `json:"insecure"`
	CaData   string `json:"caData"`
	CertData string `json:"certData"`
	KeyData  string `json:"keyData"`
	Token    string `json:"token"`
}

func (m *ManageConfig) Value() (driver.Value, error) {
	b, err := json.Marshal(m)
	return string(b), err
}

func (m *ManageConfig) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), m)
}

type ExecutorConfigs map[string]ExecutorConfig

type ExecutorConfig struct {
	Kind    string            // 执行器类型
	Options map[string]string // 不同环境超买比，Spark Version，以及一些开关
}

func (e *ExecutorConfigs) Value() (driver.Value, error) {
	b, err := json.Marshal(e)
	return string(b), err
}

func (e *ExecutorConfigs) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), e)
}
