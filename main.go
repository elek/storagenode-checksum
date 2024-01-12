package main

import (
	"context"
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/spacemonkeygo/monkit/v3"
	"go.uber.org/zap"
	"log"
	"net"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"runtime/pprof"
	"storj.io/common/storj"
	dbg "storj.io/private/debug"
	"strings"
	"syscall"
)

func main() {

	zapLog, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	if os.Getenv("STBB_PPROF") != "" {
		var output *os.File
		output, err := os.Create(os.Getenv("STBB_PPROF"))
		if err != nil {
			panic(err)
		}
		defer func() {
			output.Close()
		}()

		err = pprof.StartCPUProfile(output)
		if err != nil {
			panic(err)
		}
		defer pprof.StopCPUProfile()
	}

	if os.Getenv("STBB_MONKIT") != "" {
		filter := strings.ToLower(os.Getenv("STBB_MONKIT"))
		defer func() {
			monkit.Default.Stats(func(key monkit.SeriesKey, field string, val float64) {
				if filter == "true" || strings.Contains(strings.ToLower(key.String()), filter) {
					fmt.Println(key, field, val)
				}
			})
		}()
	}

	if os.Getenv("STBB_PPROF_ALLOCS") != "" {
		var output *os.File
		output, err := os.Create(os.Getenv("STBB_PPROF_ALLOCS"))
		if err != nil {
			panic(err)
		}
		defer func() {
			output.Close()
		}()

		defer func() {
			err = pprof.Lookup("allocs").WriteTo(output, 0)
			if err != nil {
				panic(err)
			}
		}()
	}

	if os.Getenv("STBB_DEBUG") != "" {
		fmt.Println("stating debug server")
		listener, err := net.Listen("tcp", os.Getenv("STBB_DEBUG"))
		if err != nil {
			panic(err)
		}
		dbgServer := dbg.NewServer(zapLog, listener, monkit.Default, dbg.Config{})
		go func() {
			err := dbgServer.Run(context.Background())
			if err != nil {
				fmt.Println(err)
			}
		}()
		defer dbgServer.Close()
	}

	usr1 := make(chan os.Signal, 1)
	defer close(usr1)
	signal.Notify(usr1, syscall.SIGUSR1)
	go func() {
		for {
			select {
			case _, ok := <-usr1:
				if !ok {
					return
				}
				fmt.Println(string(readStack()))
			}
		}
	}()

	var cli Checksum

	ctx := kong.Parse(&cli,
		kong.TypeMapper(reflect.TypeOf(storj.NodeURL{}), kong.MapperFunc(func(ctx *kong.DecodeContext, target reflect.Value) error {
			s := ctx.Scan.Pop().Value.(string)
			url, err := storj.ParseNodeURL(s)
			if err != nil {
				return err
			}
			target.Set(reflect.ValueOf(url))
			return nil
		})),
	)

	kong.Bind(ctx)
	err = ctx.Run(ctx)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

func readStack() []byte {
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, true)
		if n < len(buf) {
			return buf[:n]
		}
		buf = make([]byte, 2*len(buf))
	}
}
