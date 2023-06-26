package cmd

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/cnk3x/metaman/name"
	"github.com/fsnotify/fsnotify"
	"github.com/mattn/go-runewidth"
	"github.com/spf13/cobra"
)

func init() {
	AddToRoot(commandLinkAll())
}

func commandLinkAll() *cobra.Command {
	var (
		srcRoot string
		dstRoot string
		minSize int
		watch   bool
	)

	c := &cobra.Command{
		Use:   "link",
		Short: "对目录下的视频文件硬链接",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if srcRoot == "" || dstRoot == "" {
				return fmt.Errorf("缺少参数")
			}

			if srcRoot, err = filepath.Abs(srcRoot); err != nil {
				return
			}

			if dstRoot, err = filepath.Abs(dstRoot); err != nil {
				return
			}

			if srcRoot == dstRoot {
				return fmt.Errorf("源目录和目标目录不能相同")
			}

			if err = linkDir(cmd.Context(), dstRoot, srcRoot, minSize); err != nil {
				return
			}

			if watch {
				err = name.Watch(cmd.Context(), srcRoot, func(fsw *fsnotify.Watcher, evts []fsnotify.Event) {
					for _, evt := range evts {
						log.Printf("[wat] %s", evt)
						dstPath, err := name.Link(dstRoot, evt.Name)
						if err != nil {
							log.Printf("*%s -> %s", evt.Name, err)
						} else {
							log.Printf("%s <- %s", runewidth.FillRight(dstPath, 60), evt.Name)
						}
					}
				}, time.Second*3)
			}

			return
		},
	}

	flags := c.Flags()
	flags.StringVarP(&srcRoot, "src", "s", srcRoot, "*源目录")
	flags.StringVarP(&dstRoot, "dst", "d", dstRoot, "*目标目录")
	flags.IntVarP(&minSize, "min-size", "m", minSize, "最小文件大小,单位M")
	flags.BoolVarP(&watch, "watch", "w", watch, "监听源目录，新文件自动硬链接")
	return c
}

func linkDir(ctx context.Context, dstRoot, srcRoot string, minSize int) (err error) {
	return name.Walk(srcRoot, func(path string, info fs.FileInfo) error {
		if err = ctx.Err(); err != nil {
			return fs.SkipAll
		}
		dstPath, err := name.Link(dstRoot, path)
		if err != nil {
			if os.IsExist(err) {
				err = nil
			}
			return err
		}
		log.Printf("%s <- %s", runewidth.FillRight(dstPath, 60), path)
		return nil
	}, minSize)
}
