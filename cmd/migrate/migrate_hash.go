/**
 * @Time: 2020/8/8 16:35
 * @Author: solacowa@gmail.com
 * @File: hash
 * @Software: GoLand
 */

package migrate

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	migrateRedisHashCmd = &cobra.Command{
		Use:               `hash <args> [flags]`,
		Short:             "哈希列表",
		SilenceErrors:     false,
		DisableAutoGenTag: false,
		Example: `
redis-tool migrate hash {key} --source-redis-cluster true --source-hosts 127.0.0.1:6379,127.0.0.1:7379 --source-auth 123456 --target-redis-cluster true --target-hosts 127.0.0.1:6379,127.0.0.1:7379 --target-auth 123456
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
			return migrateRedisHGetAll(args[0])
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

func migrateRedisHGetAll(key string) error {
	begin := time.Now()

	total, err := sourceRedis.HLen(key)
	if err != nil {
		err = errors.Wrap(err, "sourceRedis.Hlen")
		return err
	}

	fmt.Println(fmt.Sprintf("Key: [%s] 总数: [%d]", key, total))

	res, err := sourceRedis.HGetAll(key)
	if err != nil {
		err = errors.Wrap(err, "sourceRedis.HGetAll")
		return err
	}

	step := math.Ceil(100 / float64(total))
	var i = 1
	for k, v := range res {
		i += int(step)
		if i >= 100-int(step) {
			i = 100
		}
		_, _ = fmt.Fprintf(os.Stdout, "%d%% [%s]\r", i, getS(i, "#")+getS(100-i, " "))
		if err = targetRedis.HSet(key, k, v); err != nil {
			continue
		}
		if i == 100 {
			fmt.Println(fmt.Sprintf("%d%% [%s]\r", i, getS(i, "#")+getS(100-i, " ")))
		}
	}

	fmt.Println(fmt.Sprintf("迁移完成, 用时 [%v]", time.Since(begin)))

	return nil
}
