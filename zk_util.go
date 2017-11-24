package main

import (
	"fmt"
	"time"

	"github.com/samuel/go-zookeeper/zk"
	"encoding/json"
	"strconv"
)

type ZkClient struct {
	hosts 	[]string
	conn	*zk.Conn
}

func (client* ZkClient) Init()(code int){
	code = 0
	var err error
	client.conn, _, err = zk.Connect(client.hosts, time.Second*5)
	if err != nil {
		fmt.Println(err)
		code = -1
		return
	}

	return
}

func (client* ZkClient) Close()(){
	client.conn.Close()
}

func (client* ZkClient) Create(path string, data string)(code int, outpath string){
	code = 0
	var flags int32
	flags = 0
	//flags有4种取值：
	//0:永久，除非手动删除
	//zk.FlagEphemeral = 1:短暂，session断开则改节点也被删除
	//zk.FlagSequence  = 2:会自动在节点后面添加序号
	//3:Ephemeral和Sequence，即，短暂且自动添加序号
	var acls=zk.WorldACL(zk.PermAll)//控制访问权限模式

	outpath, err_create:=client.conn.Create(path, []byte(data), flags, acls)
	if err_create != nil {
		fmt.Println(err_create)
		code = -1
		return
	}
	return
}

func (client* ZkClient) SetData(path string, data string, version int32)(code int, outpath string){
	code = 0
	_, err :=client.conn.Set(path, []byte(data), version)

	if err != nil {
		fmt.Println(err)
		code = -1
		return
	}

	return
}

func (client* ZkClient) GetData(path string)(code int, stat *zk.Stat, data string){
	code = 0
	outdata, stat ,err := client.conn.Get(path)
	if err != nil {
		fmt.Println(err)
		code = -1
		return
	}

	data = string(outdata)

	return
}



func test()(){
	c, _, err := zk.Connect([]string{"127.0.0.1"}, time.Second) //*10)
	if err != nil {
		panic(err)
	}
	children, stat, ch, err := c.ChildrenW("/")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v %+v\n", children, stat)
	e := <-ch
	fmt.Printf("%+v\n", e)
}

func test_zk(){
	client := ZkClient{[]string{"127.0.0.1"}, nil}
	client.Init()

	path := "/test"
	_, stat, data := client.GetData(path)
	fmt.Println(data, stat.Version)

	body := `[{"ip":"172.16.0.98", "weight":"100"},{"ip":"172.16.0.25", "weight":"170"}]`
	client.Create("/test2", body)
	client.SetData(path, body, stat.Version)

	_,stat, data = client.GetData(path)

	fmt.Println(data, stat.Version)
	var raw_message_list []json.RawMessage
	var router_item_list []RouterItem
	err := json.Unmarshal([]byte(data), &raw_message_list)
	if err != nil {
		fmt.Println("unmarshall error:", data, err)
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
	fmt.Println(router_item_list)
	fmt.Println(data, stat.Version)
}



