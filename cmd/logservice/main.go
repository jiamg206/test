package main

import (
	"awesomeProject8/log"
	"awesomeProject8/registry"
	"awesomeProject8/service"
	"context"
	"fmt"
	stlog "log"
)

func main() {
	log.Run("./dis.log")
	host, port := "localhost", "4055"
	ServiceAddr := fmt.Sprintf("http://%s:%s", host, port)
	r := registry.Registration{ServiceName: registry.LogService, ServiceURL: ServiceAddr, RequiredServices: []registry.ServiceName{},
		ServiceUpdateURL: ServiceAddr + "/services"}
	ctx, err := service.Start(context.Background(), r, host, port, log.RegisterHandlers)
	if err != nil {
		stlog.Fatalln(err)
	}
	<-ctx.Done()
	fmt.Println("Service shutdown")
}
