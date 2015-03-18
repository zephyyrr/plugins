package main

import (
	"github.com/zephyyrr/plugins"
	"log"
	"net"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	ph := plugins.NewManager()

	if pls, err := plugins.LoadAll("./plugins"); err == nil {
		for _, pl := range pls {
			ph.Handle(pl)
		}
	} else {
		log.Println("Error loading local plugins:", err)
	}

	go openTCPEntrance(ph, ":5436")

	ph.Muxer().AddHandler("ping", plugins.HandlerFunc(func(event plugins.Event, args plugins.Args) {
		ph.Dispatch("pong", args)
	}))

	ph.Muxer().AddHandler("log", plugins.HandlerFunc(func(event plugins.Event, args plugins.Args) {
		log.Printf("(%s): %v", event, args)
	}))

	ph.Muxer().AddHandler("square", plugins.HandlerFunc(func(event plugins.Event, args plugins.Args) {
		if x, ok := args["Argument"].(int); ok {
			ph.Dispatch("square.answer", plugins.Args{
				"argument": x,
				"answer":   x * x,
			})
		}
	}))

	ph.ListenAndServe()
}

func openTCPEntrance(man *plugins.Manager, addr string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			break
		}

		log.Println("Recieved connection from %s", conn.RemoteAddr().String())

		go func(conn net.Conn) {
			plugin, err := plugins.NewRemotePlugin(conn, conn)
			if err != nil {
				return
			}

			man.Handle(plugin)
		}(conn)
	}
}
