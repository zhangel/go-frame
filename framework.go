package framework

import (
	"fmt"
	"github.com/zhangel/go-frame/config"
	"log"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"syscall"
	internal_config "github.com/zhangel/go-frame/internal/config"
	"github.com/zhangel/go-frame/lifecycle"
	framework_logger "github.com/zhangel/go-frame/log"
	_ "github.com/zhangel/go-frame/profile"
	_ "github.com/zhangel/go-frame/prometheus"
	_ "github.com/zhangel/go-frame/registry"
	_ "github.com/zhangel/go-frame/retry"
	_ "github.com/zhangel/go-frame/tracing"

	_ "github.com/zhangel/go-frame/config_plugins"
	_ "github.com/zhangel/go-frame/control_plugins"

	"go.uber.org/automaxprocs/maxprocs"

	_ "github.com/zhangel/go-frame/balancer"
	_ "github.com/zhangel/go-frame/config"
	_ "github.com/zhangel/go-frame/control"
	_ "github.com/zhangel/go-frame/credentials"
	_ "github.com/zhangel/go-frame/db"
	"github.com/zhangel/go-frame/declare"



	"time"
)



const (
	frameworkPrefix               = "framework"
	flagIdleMemoryReleaseInterval = "framework.idle_mem_release_interval"
	flagForceGC                   = "framework.force_gc"
)

var (
	initialized = uint32(0)
	once        sync.Once
	presetOpt   []Option
)

func init() {
	declare.Flags(frameworkPrefix,
		declare.Flag{Name: flagIdleMemoryReleaseInterval, DefaultValue: -1, Description: "Interval in seconds of force release idle memory to os. -1 means never."},
		declare.Flag{Name: flagForceGC, DefaultValue: false, Description: "Whether do forceGC while try to release idle memory to os."},
	)

	log.SetFlags(log.Lshortfile)
	_, _ = maxprocs.Set(maxprocs.Logger(func(string, ...interface{}) {}))
}

func SetOption(opt ...Option) {
	presetOpt = append(presetOpt, opt...)
}

func Init(opt ...Option) func() {
	once.Do(func() {
		opts := &Options{}
		for _, o := range append(presetOpt, opt...) {
			if err := o(opts); err != nil {
				log.Fatal("[ERROR]", err)
			}
		}

		internal_config.PrepareConfigs(opts.beforeConfigPreparer, opts.configPreparer, opts.defaultConfigSource, opts.compactUsage, opts.flagsToShow, opts.flagsToHide)

		lifecycle.LifeCycle().Initialize(opts.finalizeTimeout)
		go func() {
			sig := make(chan os.Signal, 1)
			signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
			select {
			case <-sig:
				fmt.Println("Got SIGINT|SIGTERM signal, cleaning then terminate")
				lifecycle.Exit(1)
			}
		}()

		memReleaseInterval := config.Int(flagIdleMemoryReleaseInterval)
		forceGC := config.Bool(flagForceGC)

		if memReleaseInterval > 0 {
			if memReleaseInterval < 10 {
				memReleaseInterval = 10
			}

			go func() {
				for {
					select {
					case <-time.After(time.Duration(config.Int(flagIdleMemoryReleaseInterval)) * time.Second):
						if forceGC {
							runtime.GC()
						}
						debug.FreeOSMemory()

						framework_logger.Debugf("Framework: free idle memory to os, forceGC = %v", forceGC)
					}
				}
			}()
		}

		atomic.CompareAndSwapUint32(&initialized, 0, 1)
	})
	return Finalize
}

func SetDefaultConfig(cfg config.Config) {
	config.SetGlobalConfig(cfg)
}

func Finalize() {
	lifecycle.LifeCycle().Finalize()
}

func IsFinalizing() bool {
	return lifecycle.LifeCycle().IsFinalizing()
}



