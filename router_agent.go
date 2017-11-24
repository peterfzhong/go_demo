package main

import (
	"github.com/samuel/go-zookeeper/zk"
	"fmt"
	"math/rand"
	"encoding/json"
	"sync"
	"strconv"
)

type RouterItem struct {
	ip  	string `json:"ip"`
	weight   int 	`json:"weight,string"`
}


var routerMap map[string][]RouterItem = make(map[string][]RouterItem)

var mlock sync.RWMutex

type RouterAgent struct {
	hosts	[]string
	client  *ZkClient
}

func (agent *RouterAgent) Start()(code int){
	code = 0

	agent.client = &ZkClient{agent.hosts, nil}
	code = agent.client.Init()
	if (0 != code){
		return
	}
	//zk.WithEventCallback(callback)

	return
}

func (agent* RouterAgent) UpdateRouterMap(node string, value string ){
	var raw_message_list []json.RawMessage
	var router_item_list []RouterItem
	err := json.Unmarshal([]byte(value), &raw_message_list)
	if err != nil {
		fmt.Println("unmarshall error:", value, err)
		return
	}
	for _, item := range raw_message_list{
		fmt.Println(string(item))
		router_item := RouterItem{}
		var result map[string]interface{}
		err := json.Unmarshal(item, &result)
		if err != nil {
			fmt.Println("unmarshall error:", string(item), err)
			return
		}
		router_item.ip, _ = result["ip"].(string)
		data, _ := result["weight"].(string)
		router_item.weight, _ = strconv.Atoi(data)
		fmt.Println(result)
		fmt.Println(router_item)
		router_item_list = append(router_item_list, router_item)
	}

	mlock.Lock()
	defer mlock.Unlock()

	routerMap[node] = router_item_list

	fmt.Println(routerMap)

}

func (agent* RouterAgent) Watch( node string, path string) () {
	go func() {
		for {
			//agent.client.conn.GetW()
			snapshot, _, events, err := agent.client.conn.GetW(path)
			if err != nil {
				fmt.Println(err)
				return
			}
			evt := <-events
			if evt.Err != nil {
				fmt.Println(evt.Err)
				return
			}

			if zk.EventNodeDataChanged == evt.Type{
				code, _, data := agent.client.GetData(path)
				if code != 0{
					continue
				}

				agent.UpdateRouterMap(node, data)
			}
			fmt.Println(string(snapshot))
			fmt.Println(evt.Path, evt.Server, evt.State, evt.State, evt.Type)
		}
	}()

}

func (agent *RouterAgent) GetRouter(node string)(code int, ip string){
	code = 0

	router_item_list , ok := routerMap[node]
	if ok == false{
		fmt.Println("get node not exists: ", node)
		code = -1
		return
	}

	total_value := 0
	for _, router_item := range router_item_list{
		total_value += router_item.weight
	}

	ran := rand.Intn(total_value)

	pre_value := 0
	cur_value := 0
	for _, router_item := range router_item_list{
		pre_value = cur_value
		cur_value += router_item.weight

		if pre_value < ran && cur_value < ran{
			ip = router_item.ip
			return
		}
	}

	return
}





