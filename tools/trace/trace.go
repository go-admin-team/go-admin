package trace

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"

	"github.com/opentracing/opentracing-go"
	"sourcegraph.com/sourcegraph/appdash"
	appdashot "sourcegraph.com/sourcegraph/appdash/opentracing"
	"sourcegraph.com/sourcegraph/appdash/traceapp"

	"go-admin/tools"
)

var _server = &http.Server{}

// Start 启动
func Start() {
	store := appdash.NewMemoryStore()

	// Listen on any available TCP port locally.
	l, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 0})
	if err != nil {
		log.Fatal(err)
	}
	collectorPort := l.Addr().(*net.TCPAddr).Port
	collectorAdd := fmt.Sprintf(":%d", collectorPort)

	// Start an Appdash collection server that will listen for spans and
	// annotations and add them to the local collector (stored in-memory).
	cs := appdash.NewServer(l, appdash.NewLocalCollector(store))
	go cs.Start()

	// Print the URL at which the web UI will be running.
	appdashPort := 8700
	appdashURLStr := fmt.Sprintf("http://%s:%d", tools.GetLocaHonst(), appdashPort)
	appdashURL, err := url.Parse(appdashURLStr)
	if err != nil {
		log.Fatalf("Error parsing %s: %s", appdashURLStr, err)
	}

	// Start the web UI in a separate goroutine.
	tapp, err := traceapp.New(nil, appdashURL)
	if err != nil {
		log.Fatal(err)
	}
	tapp.Store = store
	tapp.Queryer = store
	_server.Addr = fmt.Sprintf(":%d", appdashPort)
	_server.Handler = tapp
	go func() {
		err := _server.ListenAndServe()
		if err != nil {
			log.Fatalf("Trace server start error: %s", err.Error())
		}
	}()
	fmt.Println(tools.Green("Trace server run at:"))
	fmt.Printf("- Local: %s/traces\n", fmt.Sprintf("http://localhost:%d", appdashPort))
	fmt.Printf("- Network: %s/traces\n", fmt.Sprintf("http://%s:%d", tools.GetLocaHonst(), appdashPort))

	tracer := appdashot.NewTracer(appdash.NewRemoteCollector(collectorAdd))
	opentracing.InitGlobalTracer(tracer)
}

// Stop 停止
func Stop(ctx context.Context) {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	fmt.Printf("%s Shutdown Server ... \r\n", tools.GetCurrentTimeStr())
	err := _server.Shutdown(ctx)
	if err != nil {
		log.Fatalf("Trace server shutdown error: %s", err.Error())
	}
	log.Println("Trace server shutdown success")
}
