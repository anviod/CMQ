package commands

import (
	"github.com/spf13/cobra"

	"fmt"
	"time"

	"github.com/micro/go-micro/util/log"

	"github.com/tian-yuan/CMQ/hub/svc"
	"github.com/tian-yuan/iot-common/util"
)

var mqttCmd = &cobra.Command{
	Use:   "mqtt",
	Short: "start mqtt server",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Start mqtt hub gateway v0.0.1 -- HEAD")

		var etcdAddr string
		cmd.Flags().StringVarP(&etcdAddr, "etcdAddress", "z", "127.0.0.1:2379", "etcd address array")
		var tracerAddr string
		cmd.Flags().StringVarP(&tracerAddr, "tracerAddress", "j", "127.0.0.1:6831", "tracer address array")
		log.Infof("start discovery client, etcd address : %s", etcdAddr)
		util.Init(
			util.WithRegistryUrls(etcdAddr),
			util.WithTracerUrl(tracerAddr),
		)
		defer util.Ctx.CloseRegisterSvc()
		defer util.Ctx.ClosePubEngineSvc()
		defer util.Ctx.CloseMessageDispatcherSvc()
		util.Ctx.InitRegisterSvc()
		util.Ctx.InitPubEngineSvc()
		util.Ctx.InitMessageDispatcherSvc()

		var redisClusterAddr string
		cmd.Flags().StringVarP(&redisClusterAddr, "redisClusterAddr", "r", "127.0.0.1:7000", "redis cluster address")
		var reqTimeout time.Duration
		cmd.Flags().DurationVarP(&reqTimeout, "redisReqTimeout", "t", 5*time.Second, "redis request timeout")
		var poolSize int32
		cmd.Flags().Int32VarP(&poolSize, "redisPoolSize", "s", 5, "redis client pool size")
		var redisSessionTimeout time.Duration
		cmd.Flags().DurationVarP(&redisSessionTimeout, "redisSessionTimeout", "o", 20*time.Minute, "redis session timeout")
		var redisSessionRefresh time.Duration
		cmd.Flags().DurationVarP(&redisSessionRefresh, "redisSessionRefresh", "f", 20*time.Minute, "redis session refresh")
		log.Infof("create redis cluster client : %s", redisClusterAddr)
		redisClient := util.GetClusterClient(redisClusterAddr, reqTimeout, int(poolSize))
		svc.Global.RedisClient = redisClient
		ss := util.NewRedisSessionStorage(redisClient)
		svc.Global.SessionStorage = ss
		svc.Global.ReqTimeOut = reqTimeout
		svc.Global.RedisSessionTimeOut = redisSessionTimeout
		svc.Global.RedisSessionRefresh = redisSessionRefresh

		conf := svc.NewMqttConf()
		cmd.Flags().StringVarP(&conf.MqttHost, "mqttHost", "m", "0.0.0.0", "mqtt hub bind host address.")
		cmd.Flags().Uint16VarP(&conf.MqttPort, "mqttPort", "p", 1883, "mqtt hub bind port.")
		mqttSvc := svc.NewMqttSvc(conf)
		mqttSvc.Start()

		httpconf := svc.NewH2cConf()
		cmd.Flags().StringVarP(&httpconf.Host, "Host", "a", "0.0.0.0", "http2 bind host address.")
		cmd.Flags().Uint16VarP(&httpconf.Port, "Port", "b", 9883, "http2 hub bind port.")
		svc.Global.SessionPrefix = httpconf.Host + ":" + fmt.Sprintf("%d", httpconf.Port)
		log.Infof("session prefix : %s", svc.Global.SessionPrefix)
		h2cSvc := svc.NewH2cSvc(httpconf)
		h2cSvc.Start(tracerAddr)
	},
}
