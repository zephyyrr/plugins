package main

import (
	"github.com/zephyyrr/plugins"
	"github.com/zephyyrr/plugins/structhandler"
	"log"
	"net"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags)
}

func main() {
	ph := plugins.NewManager()

	if pls, err := plugins.LoadAll("./plugins"); err == nil {
		for _, pl := range pls {
			log.Printf("%s, Provides: %v, Subscribes: %v", pl.Name(), pl.Provides(), pl.Subscribes())
			ph.Handle(pl)
		}
	} else {
		log.Println("Error loading local plugins:", err)
	}

	go openTCPEntrance(ph, ":5436")

	ph.Muxer().AddHandler("ping", plugins.HandlerFunc(func(event plugins.Event, args plugins.Args) {
		Log(LogEntry{"DEBUG", "main", "Recieved ping", args})
		ph.Dispatch("pong", args)
	}))

	ph.Muxer().AddHandler("log", structhandler.New(Log))

	ph.Muxer().AddHandler("square", plugins.HandlerFunc(func(event plugins.Event, args plugins.Args) {
		if x, ok := args["Argument"].(int); ok {
			ph.Dispatch("square.answer", plugins.Args{
				"argument": x,
				"answer":   x * x,
			})
		}
	}))

	go func() {
		for {
			time.Sleep(10 * time.Second)
			ph.Dispatch("ping", nil)
		}
	}()

	ph.ListenAndServe()
}

type LogEntry struct {
	Level, Origin, Message string
	Additional             map[string]interface{}
}

func Log(le LogEntry) {
	log.Printf("[%s] (%s) %s : %v", le.Level, le.Origin, le.Message, le.Additional)
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
