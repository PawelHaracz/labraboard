package routers

import "github.com/gin-gonic/gin"

func UnitedSetup() gin.HandlerFunc {
	//s3c, _ := newS3Client(cfg.KeyArn)
	//rc := newRedisClient(cfg.RedisConn)

	return func(c *gin.Context) {
		//c.Set("s3c", s3c)
		//c.Set("rc", rc)
		//c.Set("prefix", cfg.BucketPrefix)
	}
}
