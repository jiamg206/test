package main

import (
	"awesomeProject8/registry"
	"context"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.Handle("/services", &registry.RegistryService{})
	ctx, cannel := context.WithCancel(context.Background())
	defer cannel()
	var srv http.Server
	srv.Addr = registry.ServerPort
	go func() {
		log.Println(srv.ListenAndServe())
		cannel()
	}()
	go func() {
		fmt.Println("registry service started,any key to stop ")
		var s string
		fmt.Scanln(&s)
		srv.Shutdown(ctx)
		cannel()
	}()
	<-ctx.Done()
	fmt.Println("registry service stopped")
}
