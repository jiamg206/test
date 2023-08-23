package service

import (
	"awesomeProject8/registry"
	"context"
	"fmt"
	"log"
	"net/http"
)

func Start(ctx context.Context, reg registry.Registration, host, port string, registerHandlersFunc func()) (context.Context, error) {
	registerHandlersFunc()
	ctx = StartService(ctx, reg.ServiceName, host, port)
	err := registry.RegisterService(reg)
	if err != nil {
		return ctx, err
	}
	return ctx, nil
}
func StartService(ctx context.Context, ServiceName registry.ServiceName, host string, port string) context.Context {
	ctx, cannel := context.WithCancel(ctx)
	var srv http.Server
	srv.Addr = ":" + port
	go func() {
		log.Println(srv.ListenAndServe())
		err := registry.ShowDownService(fmt.Sprintf("http://%s:%s", host, port))
		if err != nil {
			log.Println(err)
		}
		cannel()
	}()
	go func() {
		fmt.Printf("%v started,any key to stop \n", ServiceName)
		var s string
		fmt.Scanln(&s)
		err := registry.ShowDownService(fmt.Sprintf("http://%s:%s", host, port))
		if err != nil {
			log.Println(err)
		}
		srv.Shutdown(ctx)
		cannel()
	}()
	return ctx
}
