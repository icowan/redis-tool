/**
 * @Time: 2020/8/8 16:40
 * @Author: solacowa@gmail.com
 * @File: sorted-set
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
	migrateRedisSortedSetCmd = &cobra.Command{
		Use:               `sorted-set <args> [flags]`,
		Short:             "有序集合迁移",
		SilenceErrors:     false,
		DisableAutoGenTag: false,
		Example: `
redis-tool migrate sorted-set {key} --source-hosts 127.0.0.1:6379 --source-auth 123456 --target-redis-cluster true --target-hosts 127.0.0.1:6379,127.0.0.1:7379 --target-auth 123456
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
			//return readFile()
			return migrateRedisSortedSet(args[0])
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

func migrateRedisSortedSet(key string) error {
	begin := time.Now()

	total, err := sourceRedis.ZCard(key)
	if err != nil {
		err = errors.Wrap(err, "sourceRedis.ZCard")
		return err
	}

	fmt.Println(fmt.Sprintf("Key: [%s] 总数: [%d]", key, total))
	var base float64 = 50000

	step := math.Ceil(float64(total) / base)

	for n := 0; n < int(step); n++ {
		res, err := sourceRedis.ZRangeWithScores(key, int64(n)*50000, int64(n+1)*50000)
		if err != nil {
			err = errors.Wrap(err, "sourceRedis.ZRangeWithScores")
			return err
		}

		s := 100 / base
		var i float64 = 0
		for _, v := range res {
			i += s
			if i >= 100-s {
				i = 100
			}
			_, _ = fmt.Fprintf(os.Stdout, "[%d/%d] %.2f%% \r", n+1, int(step), i)
			if err = targetRedis.ZAdd(key, v.Score, v.Member); err != nil {
				continue
			}
		}
		fmt.Println(fmt.Sprintf("迁移第: [%d/%d], 用时 [%v]", n+1, int(step), time.Since(begin)))
	}

	fmt.Println(fmt.Sprintf("迁移完成, 用时 [%v]", time.Since(begin)))

	return nil
}
