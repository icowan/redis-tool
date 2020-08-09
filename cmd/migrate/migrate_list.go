/**
 * @Time: 2020/8/8 17:29
 * @Author: solacowa@gmail.com
 * @File: migrate_list
 * @Software: GoLand
 */

package migrate

import (
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	migrateRedisListCmd = &cobra.Command{
		Use:               `list <args> [flags]`,
		Short:             "列表迁移",
		SilenceErrors:     false,
		DisableAutoGenTag: false,
		Example: `
redis-tool migrate list {key} --source-hosts 127.0.0.1:6379 --source-auth 123456 --target-redis-cluster true --target-hosts 127.0.0.1:6379,127.0.0.1:7379 --target-auth 123456
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// 关闭资源连接
			defer func() {
				log.Printf("source redis close err: %v", sourceRedis.Close())
				log.Printf("target redis close err: %v", targetRedis.Close())
			}()
			if len(args) < 1 {
				fmt.Println("至少需要一个参数")
				return errors.New("参数错误")
			}
			return migrateRedisList(args[0])
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err = prepare(); err != nil {
				fmt.Println(fmt.Sprintf("prepare error: %s", err.Error()))
				return err
			}
			return nil
		},
	}
)

func migrateRedisList(key string) error {
	begin := time.Now()

	total := sourceRedis.LLen(key)

	fmt.Println(fmt.Sprintf("Key: [%s] 总数: [%d]", key, total))

	for i := 0; i < int(total); i++ {
		if val, err := sourceRedis.RPop(key); err == nil {
			if err = targetRedis.LPush(key, val); err != nil {
				fmt.Println(fmt.Sprintf("迁移: [%s] --> failure: %s", key, err.Error()))
			}
		} else {
			fmt.Println(fmt.Sprintf("迁移: [%s] --> failure: %s", key, err.Error()))
		}
	}

	fmt.Println(fmt.Sprintf("迁移完成, 用时 [%v]", time.Since(begin)))

	return nil
}
