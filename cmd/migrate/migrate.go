/**
 * @Time: 2020/8/8 12:43
 * @Author: solacowa@gmail.com
 * @File: main
 * @Software: GoLand
 */

package migrate

import (
	"flag"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/icowan/redis-tool/redis"
)

var (
	sourceHosts, targetHosts, sourceAuth, targetAuth, sourcePrefix, targetPrefix string
	err                                                                          error
	sourceRedisCluster, targetRedisCluster                                       bool

	sourceDatabase, targetDatabase int

	sourceRedis redis.RedisInterface
	targetRedis redis.RedisInterface

	rootCmd = &cobra.Command{
		Use:               "redis-tool",
		Short:             "",
		SilenceErrors:     true,
		DisableAutoGenTag: true,
		Long: `# redis 迁移工具
可用的配置类型：
[migrate]
有关本系统的相关概述，请参阅 https://github.com/icowan/redis-tool
`,
	}

	migrateRedisCmd = &cobra.Command{
		Use:               `migrate command <args> [flags]`,
		Short:             "数据迁移命令",
		SilenceErrors:     false,
		DisableAutoGenTag: false,
		Example: `
支持命令:
[hash, set, sorted-set]
`,
	}
)

func init() {

	migrateRedisCmd.PersistentFlags().StringVar(&sourceHosts, "source-hosts", "127.0.0.1:6379", "源redis地址, 多个ip用','隔开")
	migrateRedisCmd.PersistentFlags().StringVar(&targetHosts, "target-hosts", "127.0.0.1:6379", "目标redis地址, 多个ip用','隔开")
	migrateRedisCmd.PersistentFlags().IntVar(&sourceDatabase, "source-database", 0, "源database")
	migrateRedisCmd.PersistentFlags().IntVar(&targetDatabase, "target-database", 0, "目标database")
	migrateRedisCmd.PersistentFlags().StringVar(&sourceAuth, "source-auth", "", "源密码")
	migrateRedisCmd.PersistentFlags().StringVar(&targetAuth, "target-auth", "", "目标密码")
	migrateRedisCmd.PersistentFlags().BoolVar(&sourceRedisCluster, "source-redis-cluster", false, "源redis是否是集群")
	migrateRedisCmd.PersistentFlags().BoolVar(&targetRedisCluster, "target-redis-cluster", false, "目标redis是否是集群")
	migrateRedisCmd.PersistentFlags().StringVar(&sourcePrefix, "source-prefix", "", "源redis前缀")
	migrateRedisCmd.PersistentFlags().StringVar(&targetPrefix, "target-prefix", "", "目标redis前缀")

	migrateRedisCmd.AddCommand(migrateRedisHashCmd, migrateRedisSortedSetCmd, migrateRedisSetCmd, migrateRedisAllCmd)
	addFlags(rootCmd)
	rootCmd.AddCommand(migrateRedisCmd)
}

func addFlags(rootCmd *cobra.Command) {
	flag.CommandLine.VisitAll(func(gf *flag.Flag) {
		rootCmd.PersistentFlags().AddGoFlag(gf)
	})
}

func Run() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

func prepare() error {
	var sourceDrive, targetDrive = redis.RedisSingle, redis.RedisSingle
	if sourceRedisCluster {
		sourceDrive = redis.RedisCluster
	}

	sourceRedis, err = redis.NewRedisClient(sourceDrive, sourceHosts, sourceAuth, sourcePrefix, sourceDatabase)
	if err != nil {
		return errors.Wrap(err, "源redis连接失败!")
	}

	if targetRedisCluster {
		targetDrive = redis.RedisCluster
	}

	targetRedis, err = redis.NewRedisClient(targetDrive, targetHosts, targetAuth, targetPrefix, targetDatabase)
	if err != nil {
		return errors.Wrap(err, "目标redis连接失败!")
	}
	return nil
}

func getS(n int, char string) (s string) {
	for i := 1; i <= n; i++ {
		s += char
	}
	return
}
