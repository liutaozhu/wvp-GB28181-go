package redis

import (
	"context"
	"fmt"

	"wvp-pro-go/internal/config"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var (
	Client *redis.Client
	Ctx    = context.Background()
)

// Redis key prefixes matching Java VideoManagerConstants
const (
	WVPPrefix = "wvp:"

	// Device keys
	DevicePrefix     = WVPPrefix + "device:"
	DeviceOnline     = WVPPrefix + "device:online"
	DeviceRegisterTS = WVPPrefix + "device:register:"
	DeviceKeepalive  = WVPPrefix + "device:keepalive:"

	// Stream/session keys
	InvitePrefix       = WVPPrefix + "invite:"
	InviteSession      = WVPPrefix + "stream:session:"
	StreamInfoPrefix   = WVPPrefix + "stream:info:"
	SSRCPrefix         = WVPPrefix + "ssrc:"
	SSRCPool           = WVPPrefix + "ssrc:pool"
	SSRCTransaction    = WVPPrefix + "ssrc:tx:"

	// SIP subscribe keys
	SubscribePrefix = WVPPrefix + "subscribe:"

	// Platform keys
	PlatformPrefix = WVPPrefix + "platform:"

	// Channel cache
	ChannelCachePrefix = WVPPrefix + "channel:cache:"

	// Alarm cache
	AlarmCachePrefix = WVPPrefix + "alarm:cache:"

	// RTP server
	RTPServerPrefix = WVPPrefix + "rtp:server:"

	// Media server
	MediaServerPrefix = WVPPrefix + "media:server:"

	// Redis pub/sub channels
	StreamChangeChannel = WVPPrefix + "stream:change"
	AlarmChannel        = WVPPrefix + "alarm"
	DeviceOnlineChannel = WVPPrefix + "device:online"
	DeviceOfflineChannel = WVPPrefix + "device:offline"
	GPSChannel          = WVPPrefix + "gps"
	PlatformPlayChannel = WVPPrefix + "platform:play"

	// Server
	WVPServerPrefix = WVPPrefix + "server:"

	// CSEQ
	CSEQKey = WVPPrefix + "cseq"
)

func Init(cfg config.RedisConfig, log *zap.Logger) error {
	rdb := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr(),
		Password:     cfg.Password,
		DB:           cfg.DB,
		ReadTimeout:  0,
		WriteTimeout: 0,
	})

	if err := rdb.Ping(Ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect redis: %w", err)
	}

	Client = rdb
	log.Info("redis initialized", zap.String("addr", cfg.Addr()))
	return nil
}
