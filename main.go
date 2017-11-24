package main

import (
	"time"
	"fmt"
)

func start()(){
	agent := RouterAgent{[]string{"127.0.0.1"}, nil}
	agent.Start()
	agent.Watch("test", "/test")

	time.Sleep(time.Second * 2000)
}

func getRouter(node string)(){
	agent := RouterAgent{[]string{"127.0.0.1"}, nil}
	agent.Start()
	_, ip := agent.GetRouter(node)
	fmt.Println(ip)
}

func main()  {
	//start()
	test_zk()
	getRouter("test")

}