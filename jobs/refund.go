package jobs

import (
	"fmt"
	"github.com/go-redsync/redsync"
	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
	"github.com/tietang/go-utils"
	"imooc.com/resk/core/envelopes"
	"imooc.com/resk/infra"
	"time"
)

type RefundExpiredJobStarter struct {
	infra.BaseStarter
	ticker *time.Ticker
	mutex  *redsync.Mutex
}

func (r *RefundExpiredJobStarter) Init(ctx infra.StarterContext) {
	d := ctx.Props().GetDurationDefault("jobs.refund.interval", time.Minute)
	r.ticker = time.NewTicker(d)
	maxIdle := ctx.Props().GetIntDefault("redis.maxIdle", 2)
	maxActive := ctx.Props().GetIntDefault("redis.maxActive", 5)
	timeout := ctx.Props().GetDurationDefault("redis.timeout", 20*time.Second)
	addr := ctx.Props().GetDefault("redis.addr", "127.0.0.1:6379")
	pools := make([]redsync.Pool, 0)
	pool := &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: timeout,
		Dial: func() (conn redis.Conn, e error) {
			return redis.Dial("tcp", addr)
		},
	}
	pools = append(pools, pool)
	rsync := redsync.New(pools)
	ip, err := utils.GetExternalIP()
	if err != nil {
		ip = "127.0.0.1"
	}
	r.mutex = rsync.NewMutex("lock:RefundExpired",
		redsync.SetExpiry(50*time.Second),
		redsync.SetRetryDelay(3),
		redsync.SetGenValueFunc(func() (s string, e error) {
			now := time.Now()
			log.Infof("节点%s正在执行过期红包的退款任务", ip)
			return fmt.Sprintf("%d:%s", now.Unix(), ip), nil
		}),
	)
}

func (r *RefundExpiredJobStarter) Start(ctx infra.StarterContext) {
	go func() {
		for {
			c := <-r.ticker.C

			err := r.mutex.Lock()
			if err == nil {
				log.Debug("过期红包退款开始...", c)
				//红包过期退款的业务逻辑代码
				domain := envelopes.ExpiredEnvelopeDomain{}
				domain.Expired()
			} else {
				log.Info("已经有节点在运行该任务了")
			}
			r.mutex.Unlock()

		}
	}()

}

func (r *RefundExpiredJobStarter) Stop(ctx infra.StarterContext) {
	r.ticker.Stop()
}

//这里把测试和debug留给同学们
