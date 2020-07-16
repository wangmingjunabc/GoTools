package main

import (
	"fmt"
	"hash/crc32"
	"io"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

//环的抽象
type Map struct {
	hash     Hash           //计算hash的函数
	replicas int            //这个是副本，这里影响到虚拟节点的个数
	keys     []int          //有序列表，从大到小排序的
	hashMap  map[int]string //可以理解成用来记录虚拟节点和物理节点对应关系
}

//如果有三个节点，replicas设置为3，那么就在环上有9个节点，9个元素
//hash环的初始化
func New(replicas int, fn Hash) *Map {
	return &Map{
		hash:     fn,
		replicas: replicas,
		hashMap:  make(map[int]string),
	}
}

//Hash环添加节点
/*
比如， A, B, C三个节点， replicas为3， 那么就：
节点输入：keys => [A, B, C];
用来计算hash值的输入是：i + key，也就是：[0A, 1A, 2A, 0B, 1B, 2B, 0C, 1C, 2C];
计算出来的hash序列是：m.keys = [hash(0A), hash(1A), hash(2A), hash(0B), hash(1B), hash(2B), hash(0C), hash(1C), hash(2C)]
这里认为hash函数是比较好的平衡性，那么计算出的值，应该是随机均衡打散的，我们认为是符合概率分布的；
最后会把这个hash值的序列做一个排序，做完排序之后，其实就完成了值域的打散划分;
*/
func (m *Map) Add(keys ...string) {
	//keys => [A, B, C]
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			//hash 值 = hash（i + key)
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			//虚拟节点 <-> 实际节点
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
}

//一致性hash请求
func (m *Map) Get(key string) string {
	if m.IsEmpty() {
		return ""
	}
	//根据用户输入key值，计算一个hash值
	hash := int(m.hash([]byte(key)))
	//查看值落到哪个值域范围？选择到虚节点
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	if idx == len(m.keys) {
		idx = 0
	}
	//选择到对应物理节点
	return m.hashMap[m.keys[idx]]
}

func (m *Map) IsEmpty() bool {
	if len(m.keys) == 0 {
		return true
	}
	return false
}

func hash(s []byte) uint32 {
	ieee := crc32.NewIEEE()
	_, _ = io.WriteString(ieee, string(s))
	return ieee.Sum32()
}

func main() {
	m := New(3, hash)
	m.Add("A", "B", "C")
	get := m.Get("ssssssfsdsss")
	fmt.Println(get)
}
