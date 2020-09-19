package aws

import (
	"github.com/lightpaw/male7/config/kv"
	"github.com/lightpaw/eventlog"
	"github.com/aws/aws-sdk-go/service/firehose"
	"time"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/lightpaw/logrus"
)

func NewAwsService(serverConfig *kv.IndividualServerConfig) *AwsService {
	s := &AwsService{
		serverConfig: serverConfig,
	}

	s.init()

	return s
}

//gogen:iface
type AwsService struct {
	serverConfig *kv.IndividualServerConfig

	session *session.Session
}

func (s *AwsService) init() {

	region := s.serverConfig.AwsRegion
	if len(region) <= 0 {
		return
	}

	creds := credentials.NewChainCredentials([]credentials.Provider{
		&credentials.StaticProvider{Value: credentials.Value{
			AccessKeyID:     s.serverConfig.AwsAccessKeyID,
			SecretAccessKey: s.serverConfig.AwsSecretAccessKey,
			SessionToken:    s.serverConfig.AwsSessionToken,
		}},
		&credentials.EnvProvider{},
		&credentials.SharedCredentialsProvider{},
	})

	// 初始化相关的aws服务
	session, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: creds,
	})
	if err != nil {
		logrus.WithError(err).Panicf("初始化AwsService失败")
	}

	s.session = session
}

func (s *AwsService) InitFirehoseEventLog() bool {

	if s.session == nil {
		return false
	}

	if len(s.serverConfig.AwsFirehoseDeliveryStreamName) <= 0 {
		return false
	}

	firehose := firehose.New(s.session)
	eventlog.Start(eventlog.NewFirehoseDestination(s.serverConfig.AwsFirehoseDeliveryStreamName, firehose), 1*time.Second)
	return true
}
