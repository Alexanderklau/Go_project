package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type Calc struct{}

type Args struct {
	A  float64 `json:"a"`
	B  float64 `json:"b"`
	Op string  `json:"op"`
}

type Reply struct {
	Msg  string  `json:"msg"`
	Data float64 `json:"data"`
}

// func (t *T) MethodName(argType T1, replyType *T2) error
// 一定要满足这种格式， 其中replyType要为指针。并注意大小写
func (c *Calc) Compute(args Args, reply *Reply) error {
	var (
		msg string = "ok"
	)

	switch args.Op {
	case "+":
		reply.Data = args.A + args.B
	case "-":
		reply.Data = args.A - args.B
	case "*":
		reply.Data = args.A * args.B
	case "/":
		if args.B == 0 {
			msg = "in divide op, B can't be zero"
		} else {
			reply.Data = args.A / args.B
		}
	default:
		msg = fmt.Sprintf("unsupported op:%s", args.Op)
	}
	reply.Msg = msg

	if reply.Msg == "ok" {
		return nil
	}
	return fmt.Errorf(msg)
}

// 启动server端
func Start() {
	err := rpc.Register(new(Calc))

	if err != nil {
		panic(err)
	}

	listener, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Println(err)
			continue
		}

		go jsonrpc.ServeConn(conn)
	}
}

func main() {
	server.Start()
}
