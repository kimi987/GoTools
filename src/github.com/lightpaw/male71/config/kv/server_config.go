package kv

import (
	"fmt"
	"github.com/lightpaw/logrus"
	"net"
	"strconv"
	"strings"
	"github.com/lightpaw/male7/util/timeutil"
	"time"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/gen/pb/util"
	"github.com/lightpaw/male7/constants"
	"github.com/lightpaw/male7/util/check"
)

const file = "server.yaml"

func NewIndividualServerConfig() *IndividualServerConfig {
	reader := newConfigReader()
	result := &IndividualServerConfig{}

	var err error
	result.DBSN = reader.DefString("dbsn", "root:my-secret-pw@tcp(192.168.1.5:3306)/male7_common")

	if result.DBMaxOpenConns, err = reader.DefInt("db_max_open_conns", 800); err != nil {
		logrus.WithError(err).Panic("从配置获取db_max_open_conns失败")
	}

	if result.DBMaxIdleConns, err = reader.DefInt("db_max_idle_conns", 30); err != nil {
		logrus.WithError(err).Panic("从配置获取db_max_idle_conns失败")
	}

	if result.DBConnMaxLifetime, err = time.ParseDuration(reader.DefString("db_conn_max_lifetime", "30s")); err != nil {
		logrus.WithError(err).Panic("从配置获取db_conn_max_lifetime失败")
	}

	// datastore
	if result.EnableMysql, err = reader.DefBool("enable_mysql", false); err != nil {
		logrus.WithError(err).Panic("从配置获取enable_datastore失败")
	}
	result.DatastoreHost = reader.DefString("datastore_host", "192.168.1.5:8432")
	result.ProjectID = reader.DefString("datastore_project_id", "male7")

	if result.IsCheckConfig, err = reader.DefBool("check_config", false); err != nil {
		logrus.WithError(err).Panic("从配置获取check_config失败")
	}

	if result.IsDebug, err = reader.DefBool("debug", false); err != nil {
		logrus.WithError(err).Panic("从配置获取debug失败")
	}

	if result.IsDebugYuanbao, err = reader.DefBool("debug_yuanbao", false); err != nil {
		logrus.WithError(err).Panic("从配置获取debug_yuanbao失败")
	}

	if result.IsDebugFight, err = reader.DefBool("debug_fight", false); err != nil {
		logrus.WithError(err).Panic("从配置获取debug_fight失败")
	}

	if result.IsAllowRobot, err = reader.DefBool("allow_robot", false); err != nil {
		logrus.WithError(err).Panic("从配置获取allow_robot失败")
	}

	if result.IgnoreHeartBeat, err = reader.DefBool("ignore_heart_beat", false); err != nil {
		logrus.WithError(err).Panic("从配置获取ignore_heart_beat失败")
	}

	if result.DontEncrypt, err = reader.DefBool("dont_encrypt", false); err != nil {
		logrus.WithError(err).Panic("从配置获取dont_encrypt失败")
	}

	if result.SkipHeader, err = reader.DefBool("skip_header", false); err != nil {
		logrus.WithError(err).Panic("从配置获取skip_header失败")
	}

	if result.HeroOfflineCacheCount, err = reader.DefInt("hero_offline_cache_count", 5000); err != nil {
		logrus.WithError(err).Panic("从配置获取hero_offline_cache_count失败")
	}

	if result.HeroEvictNotAccessedDuration, err = time.ParseDuration(reader.DefString("hero_evict_not_accessed_duration", "10m")); err != nil {
		logrus.WithError(err).Panic("从配置获取hero_evict_not_accessed_duration失败")
	}

	if result.httpPort, err = reader.DefInt("http_port", 18080); err != nil {
		logrus.WithError(err).Panic("从配置获取http_port失败")
	}

	result.replayPrefix = reader.DefString("replay_prefix", "replay/")

	if result.ServerID, err = reader.DefInt("server_id", 1); err != nil {
		logrus.WithError(err).Panic("从配置获取server_id失败")
	}

	if result.PlatformID, err = reader.DefInt("platform_id", 1); err != nil {
		logrus.WithError(err).Panic("从配置获取platform_id失败")
	}

	if result.Port, err = reader.DefInt("port", 8080); err != nil {
		logrus.WithError(err).Panic("从配置获取port失败")
	}

	if result.RpcPort, err = reader.DefInt("rpc_port", 0); err != nil {
		logrus.WithError(err).Panic("从配置获取rpc_port失败")
	}

	if result.StartTime, err = timeutil.ParseDayLayout(reader.DefString("start_time", "2018-01-01")); err != nil {
		logrus.WithError(err).Panic("从配置获取start_time失败")
	}

	constants.PlayerNamePrefix = reader.DefString("player_name", "君主_")

	result.ServerName = reader.DefString("server_name", "测试服")

	result.ConnectionServerAddr = reader.DefString("conn_addr", "116.247.86.50:7778")

	result.LocalAddrStr = reader.DefString("local_addr", "")
	result.LocalAddr = parseLocalAddr(result.LocalAddrStr)

	b := result.LocalAddr // 内网地址
	internalHttpAddr := fmt.Sprintf("%d.%d.%d.%d:%d", b[0], b[1], b[2], b[3], result.httpPort)
	result.MetricsAddr = internalHttpAddr

	result.httpAddr = reader.DefString("http_addr", "")
	if len(result.httpAddr) <= 0 {
		result.httpAddr = "http://" + internalHttpAddr
	}

	if rpcAddr := reader.DefString("rpc_addr", ""); len(rpcAddr) > 0 {
		result.RpcAddr = rpcAddr
	} else {
		result.RpcAddr = fmt.Sprintf("%d.%d.%d.%d", result.LocalAddr[0], result.LocalAddr[1], result.LocalAddr[2], result.LocalAddr[3])
	}

	result.ZkAddr = parseZkAddr(reader.DefString("zk", "111.230.185.163:2181;193.112.127.115:2181;118.126.107.11:2181"))

	if addr := reader.DefString("etcd", ""); len(addr) > 0 {
		result.EtcdAddr = strings.Split(addr, ";")
	}

	if util.IsDebug, err = reader.DefBool("enable_debug_msg", false); err != nil {
		logrus.WithError(err).Panic("从配置获取enable_debug_msg失败")
	}

	if util.DiffByte, err = reader.DefInt("debug_msg_diff_byte", 0); err != nil {
		logrus.WithError(err).Panic("从配置获取debug_msg_diff_byte失败")
	}

	if result.IsCheckConfig, err = reader.DefBool("check_config", false); err != nil {
		logrus.WithError(err).Panic("从配置获取check_config失败")
	}

	if result.EnableOnlineCountMetrics, err = reader.DefBool("enable_online_count_metrics", true); err != nil {
		logrus.WithError(err).Panic("从配置获取enable_online_count_metrics失败")
	}

	if result.EnableRegisterMetrics, err = reader.DefBool("enable_register_count_metrics", true); err != nil {
		logrus.WithError(err).Panic("从配置获取enable_register_count_metrics失败")
	}

	if result.EnableDBMetrics, err = reader.DefBool("enable_db_metrics", true); err != nil {
		logrus.WithError(err).Panic("从配置获取enable_db_metrics失败")
	}

	if result.EnableMsgMetrics, err = reader.DefBool("enable_msg_metrics", true); err != nil {
		logrus.WithError(err).Panic("从配置获取enable_msg_metrics失败")
	}

	if result.EnablePanicMetrics, err = reader.DefBool("enable_panic_metrics", true); err != nil {
		logrus.WithError(err).Panic("从配置获取enable_panic_metrics失败")
	}

	if result.EnableCosUploader, err = reader.DefBool("enable_cos_uploader", false); err != nil {
		logrus.WithError(err).Panic("从配置获取enable_cos_uploader失败")
	}
	result.CosUrl = reader.DefString("cos_url", "https://male7replay-1256076575.cos-website.ap-guangzhou.myqcloud.com")

	result.LightpawKey = []byte(reader.DefString("lightpaw_key", "suibianqige"))

	result.OrderKey = reader.DefString("order_key", "luanquyige")

	result.LoginClusterAddr = reader.DefString("login_cluster_addr", "login.lightpaw.com:7890")

	result.CombatClusterAddr = reader.DefString("combat_cluster_addr", "")

	result.TssClusterAddr = reader.DefString("tss_cluster_addr", "")

	if result.DisablePush, err = reader.DefBool("disable_push", false); err != nil {
		logrus.WithError(err).Panic("从配置获取disable_push失败")
	}

	// ==== tlog ====
	result.GameAppID = reader.DefString("game_app_id", "7777777")

	if result.TlogStart, err = reader.DefBool("tlog_start", false); err != nil {
		logrus.WithError(err).Panic("从配置获取 tlog_start 失败")
	}

	if result.TlogWriteBufSize, err = reader.DefInt("tlog_write_buf_size", 4096); err != nil {
		logrus.WithError(err).Panic("从配置获取 tlog_write_buf_size 失败")
	}

	if result.TlogRotateSize, err = reader.DefInt("tlog_max_size", 1024*1024*1024); err != nil {
		logrus.WithError(err).Panic("从配置获取 tlog_max_size 失败")
	}

	if result.TlogCacheSize, err = reader.DefInt("tlog_cache_size", 65536); err != nil {
		logrus.WithError(err).Panic("从配置获取 tlog_max_size 失败")
	}

	result.TlogCurrentBaseDir = reader.DefString("tlog_current_base_dir", "tlogs/current/")

	result.TlogArchiveBaseDir = reader.DefString("tlog_archive_base_dir", "tlogs/archive/")

	if result.TlogRotateDuration, err = time.ParseDuration(reader.DefString("tlog_rotate_duration", "1h")); err != nil {
		logrus.WithError(err).Panic("从配置获取 tlog_max_size 失败")
	}

	// kafka fail log
	if result.TlogKafkaFailWriteBufSize, err = reader.DefInt("tlog_kafka_fail_write_buf_size", 4096); err != nil {
		logrus.WithError(err).Panic("从配置获取 tlog_write_buf_size 失败")
	}

	if result.TlogKafkaFailRotateSize, err = reader.DefInt("tlog_kafka_fail_max_size", 1024*1024*1024); err != nil {
		logrus.WithError(err).Panic("从配置获取 tlog_max_size 失败")
	}

	if result.TlogKafkaFailCacheSize, err = reader.DefInt("tlog_kafka_fail_cache_size", 65536); err != nil {
		logrus.WithError(err).Panic("从配置获取 tlog_max_size 失败")
	}

	result.TlogKafkaFailCurrentBaseDir = reader.DefString("tlog_kafka_fail_current_base_dir", "tlogs/kafka_current/")

	result.TlogKafkaFailArchiveBaseDir = reader.DefString("tlog_kafka_fail_archive_base_dir", "tlogs/kafka_archive/")

	if result.TlogKafkaFailRotateDuration, err = time.ParseDuration(reader.DefString("tlog_kafka_fail_rotate_duration", "1h")); err != nil {
		logrus.WithError(err).Panic("从配置获取 tlog_kafka_fail_rotate_duration 失败")
	}

	// ==== kafka ====
	if result.KafkaStart, err = reader.DefBool("kafka_start", false); err != nil {
		logrus.WithError(err).Panic("从配置获取 kafka_start 失败")
	}

	kafkaBrokerAddr := reader.DefString("kafka_broker_addr", "")
	if result.KafkaStart {
		check.PanicNotTrue(len(kafkaBrokerAddr) > 0, "如果 kafka_start = true, 则 kafka_broker_addr 必须配置。格式: ip1:port1;ip2:port2")

		result.KafkaBrokerAddr = strings.Split(kafkaBrokerAddr, ";")
		for _, addr := range result.KafkaBrokerAddr {
			if _, _, err := net.SplitHostPort(addr); err != nil {
				logrus.WithField("kafka broker", addr).WithError(err).Panic("kafka broker not a valid ip address")
			}
		}
	}

	result.TlogTopic = reader.DefString("tlog_topic", "tlog")

	if result.P8RechargeDebug, err = reader.DefBool("p8_recharge_debug", false); err != nil {
		logrus.WithError(err).Panic("从配置获取 p8_recharge_debug 失败")
	}

	result.LuaConfAddr = reader.DefString("lua_conf_addr", "http://706.lightpaw.com:7759")

	// aws配置
	result.AwsRegion = reader.DefString("aws_region", "")
	result.AwsAccessKeyID = reader.DefString("aws_access_key_id", "")
	result.AwsSecretAccessKey = reader.DefString("aws_secret_access_key", "")
	result.AwsSessionToken = reader.DefString("aws_session_token", "")

	// 数据后台日志上报
	result.AwsFirehoseDeliveryStreamName = reader.DefString("aws_firehose_delivery_stream_name", "")

	// 初始化
	replayUrlPrefix := removeLastSlash(result.httpAddr) + "/" + removeLastSlash(result.replayPrefix)

	result.serverInfoProto = &shared_proto.HeroServerInfoProto{
		Sid:       int32(result.ServerID),
		StartTime: timeutil.Marshal32(result.StartTime),
		LocalUrl:  replayUrlPrefix,
		CosUrl:    removeLastSlash(result.CosUrl),
	}

	return result
}

func removeLastSlash(s string) string {
	if strings.HasSuffix(s, "/") {
		return s[:len(s)-1]
	}
	return s
}

// 解析zk地址
func parseZkAddr(zkAddr string) []string {
	// zk address
	zkAddrArray := strings.Split(zkAddr, ";")
	for _, zk := range zkAddrArray {
		if _, _, err := net.SplitHostPort(zk); err != nil {
			logrus.WithField("zk", zk).WithError(err).Panic("zkAddr not a valid ip address")
		}
	}

	return zkAddrArray
}

// 解析本地地址
func parseLocalAddr(localAddr string) []byte {
	var err error
	if localAddr == "" {
		localAddr, err = getLocalAddr()
		if err != nil {
			logrus.WithError(err).Panic("无法自动获得本机的内网地址, 请通过配置server.yaml中的local_addr字段手动设置")
		}
	}

	localAddrArray := make([]byte, 4)
	localAddrSplit := strings.Split(localAddr, ".")

	if len(localAddrSplit) != 4 {
		logrus.WithField("addr", localAddr).Panic("本地的机器地址local_addr必须是ip的形式. eg 192.168.1.10")
	}

	for i := 0; i < 4; i++ {
		num, err := strconv.Atoi(localAddrSplit[i])
		if err != nil || num < 0 || num >= 256 {
			logrus.WithField("addr", localAddr).Panic("本地的机器地址local_addr必须是个合法的ip. eg 192.168.1.10")
		}

		localAddrArray[i] = uint8(num)
	}
	return localAddrArray
}

//gogen:iface
type IndividualServerConfig struct {
	IsCheckConfig   bool // 检查配置模式
	IsDebug         bool
	IsDebugYuanbao  bool
	IsDebugFight    bool
	IsAllowRobot    bool
	IgnoreHeartBeat bool
	DontEncrypt     bool
	SkipHeader      bool // 跳过开头

	HeroOfflineCacheCount int

	HeroEvictNotAccessedDuration time.Duration

	httpAddr string
	httpPort int // http服务器监听端口号

	replayPrefix string // 回放录像的前缀，默认 replay/

	MetricsAddr              string // 监控url
	EnableOnlineCountMetrics bool   // 开启在线用户数监控，默认开启
	EnableRegisterMetrics    bool   // 开启Register监控，默认开启
	EnableDBMetrics          bool   // 开启db错误监控，默认开启
	EnableMsgMetrics         bool   // 开启消息耗时监控，默认开启
	EnablePanicMetrics       bool   // 开启Panic监控，默认开启

	// mysql
	EnableMysql bool
	DBSN        string // db 连接信息

	DBMaxOpenConns    int
	DBMaxIdleConns    int
	DBConnMaxLifetime time.Duration

	// datastore (默认使用datastore)
	DatastoreHost string
	ProjectID     string

	ServerID             int    // 本服id
	PlatformID           int    // 平台id, 暂时只有国服 1
	ServerName           string // 本服的名字, 只注册时用
	ConnectionServerAddr string // 连接服的地址, 只注册时用
	LocalAddr            []byte // 本机的ip地址
	LocalAddrStr         string // 本机的ip地址str
	Port                 int    // 监听的端口
	RpcAddr              string // rpc访问地址，默认使用本机内网地址，实际可以根据配置更改
	RpcPort              int    // rpc监听的端口
	ZkAddr               []string
	EtcdAddr             []string

	// uploader
	EnableCosUploader bool
	CosUrl            string

	// 登录服集群
	LoginClusterAddr string

	// 战斗服集群
	CombatClusterAddr string

	// tss集群
	TssClusterAddr string

	// 开服时间
	StartTime time.Time

	// 管理Key
	LightpawKey []byte

	// OrderKey
	OrderKey string

	// 禁止推送消息
	DisablePush bool

	serverInfoProto *shared_proto.HeroServerInfoProto

	KafkaStart      bool
	KafkaBrokerAddr []string
	TlogTopic       string

	GameAppID string // 腾讯分配的appId

	P8RechargeDebug bool

	/* tlog */
	TlogStart          bool          // 是否开启腾讯日志，默认关闭
	TlogWriteBufSize   int           // File Writer 缓存
	TlogCacheSize      int           // tlog 缓存队列大小
	TlogRotateDuration time.Duration // 新tlog文件创建间隔时间
	TlogRotateSize     int           // 单个tlog文件最大大小
	TlogCurrentBaseDir string        // tlog 临时文件地址
	TlogArchiveBaseDir string        // tlog 存档地址

	TlogKafkaFailWriteBufSize   int           // File Writer 缓存
	TlogKafkaFailCacheSize      int           // tlog 缓存队列大小
	TlogKafkaFailRotateDuration time.Duration // 新tlog文件创建间隔时间
	TlogKafkaFailRotateSize     int           // 单个tlog文件最大大小
	TlogKafkaFailCurrentBaseDir string        // tlog kafka 失败时临时文件地址
	TlogKafkaFailArchiveBaseDir string        // tlog kafka 失败时存档地址

	// 获取客户端配置地址
	LuaConfAddr string

	// aws配置
	AwsRegion          string
	AwsAccessKeyID     string
	AwsSecretAccessKey string
	AwsSessionToken    string

	// 数据后台日志上报
	AwsFirehoseDeliveryStreamName string
}

func (c *IndividualServerConfig) GetGameAppID() string {
	return c.GameAppID
}

func (c *IndividualServerConfig) GetLocalAddStr() string {
	return c.LocalAddrStr
}

func (c *IndividualServerConfig) GetTlogStart() bool {
	return c.TlogStart
}

func (c *IndividualServerConfig) GetTlogTopic() string {
	return c.TlogTopic
}

func (c *IndividualServerConfig) GetKafkaStart() bool {
	return c.KafkaStart
}

func (c *IndividualServerConfig) GetKafkaBrokerAddr() []string {
	return c.KafkaBrokerAddr
}

func (c *IndividualServerConfig) GetPort() int {
	return c.Port
}

func (c *IndividualServerConfig) GetHttpPort() int {
	return c.httpPort
}

func (c *IndividualServerConfig) GetReplayPrefix() string {
	return c.replayPrefix
}

func (c *IndividualServerConfig) GetIsDebug() bool {
	return c.IsDebug
}

func (c *IndividualServerConfig) GetIsAllowRobot() bool {
	return c.IsAllowRobot
}

func (c *IndividualServerConfig) GetIgnoreHeartBeat() bool {
	return c.IgnoreHeartBeat
}

func (c *IndividualServerConfig) GetDontEncrypt() bool {
	return c.DontEncrypt
}

func (c *IndividualServerConfig) GetServerID() int {
	return c.ServerID
}

func (c *IndividualServerConfig) GetZoneAreaID() int {
	// todo 大区 ID
	return c.ServerID
}

func (c *IndividualServerConfig) GetPlatformID() int {
	return c.PlatformID
}

func (c *IndividualServerConfig) GetServerStartTime() time.Time {
	return c.StartTime
}

func (c *IndividualServerConfig) GetServerInfo() *shared_proto.HeroServerInfoProto {
	return c.serverInfoProto
}

func (c *IndividualServerConfig) GetSkipHeader() bool {
	return c.SkipHeader
}

func (c *IndividualServerConfig) GetDisablePush() bool {
	return c.DisablePush
}
