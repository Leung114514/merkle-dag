package main

import (
	"fmt"
	"math/big"
	"math/rand"
)

const (
	K             = 3  // 每个桶中节点的数量
	IDLengthBytes = 20 // 节点 ID 的字节长度
	B             = 8  // 桶的数量
)

// Node 结构体表示 Kademlia 网络中的节点
type Node struct {
	ID string
	// 其他节点信息
}

// Bucket 结构体表示 Kademlia 网络中的节点桶
type Bucket struct {
	MinID  *big.Int
	MaxID  *big.Int
	Nodes  []Node
	Length int
}

// K_Bucket 结构体用于管理 Kademlia 网络中的节点桶
type K_Bucket struct {
	Buckets []*Bucket
}

// NewBucket 创建一个新的桶
func NewBucket(minID, maxID *big.Int) *Bucket {
	return &Bucket{
		MinID:  minID,
		MaxID:  maxID,
		Length: 0,
		Nodes:  make([]Node, 0, K),
	}
}

// NewK_Bucket 创建一个新的 K_Bucket 实例
func NewK_Bucket() *K_Bucket {
	buckets := make([]*Bucket, 0, B)
	for i := 0; i < B; i++ {
		minID := new(big.Int).Lsh(big.NewInt(1), uint(IDLengthBytes*i))
		maxID := new(big.Int).Lsh(big.NewInt(1), uint(IDLengthBytes*(i+1)))
		buckets = append(buckets, NewBucket(minID, maxID))
	}
	return &K_Bucket{
		Buckets: buckets,
	}
}

func (k *K_Bucket) insertNode(nodeId string) {
	newNode := Node{ID: nodeId}
	idInt := new(big.Int)
	idInt.SetString(nodeId, 16)

	for i, bucket := range k.Buckets {
		if idInt.Cmp(bucket.MinID) >= 0 && idInt.Cmp(bucket.MaxID) < 0 {
			if bucket.Length < K {
				bucket.Nodes = append(bucket.Nodes, newNode)
				bucket.Length++
				return
			} else if bucket.Length == K {
				// Bucket full, need split
				fmt.Println("Bucket is full, need to split")
				leftID := new(big.Int).Set(bucket.MinID)
				rightID := new(big.Int).Set(bucket.MaxID)

				// Split bucket into two equal parts
				midID := new(big.Int).Add(leftID, rightID)
				midID.Div(midID, big.NewInt(2))

				// Create new buckets
				leftBucket := NewBucket(leftID, midID)
				rightBucket := NewBucket(new(big.Int).Add(midID, big.NewInt(1)), rightID)

				// Reallocate nodes to new buckets
				for _, node := range bucket.Nodes {
					nodeIDInt := new(big.Int)
					nodeIDInt.SetString(node.ID, 16)
					if nodeIDInt.Cmp(midID) < 0 {
						leftBucket.Nodes = append(leftBucket.Nodes, node)
						leftBucket.Length++
					} else {
						rightBucket.Nodes = append(rightBucket.Nodes, node)
						rightBucket.Length++
					}
				}

				// Replace the full bucket with the two new buckets
				k.Buckets[i] = leftBucket
				k.Buckets = append(k.Buckets, rightBucket)
				return
			}
		}
	}
}

// printBucketContents 方法打印每个桶中存在的 NodeID
func (k *K_Bucket) printBucketContents() {
	for i, bucket := range k.Buckets {
		fmt.Printf("Bucket %d: [", i)
		for _, node := range bucket.Nodes {
			fmt.Printf("%s, ", node.ID)
		}
		fmt.Println("]")
	}
}
func (k *K_Bucket) findNode(nodeId string) []string {
	idInt := new(big.Int)
	idInt.SetString(nodeId, 16)

	for _, bucket := range k.Buckets {
		if idInt.Cmp(bucket.MinID) >= 0 && idInt.Cmp(bucket.MaxID) < 0 {
			for _, node := range bucket.Nodes {
				if node.ID == nodeId {
					// 如果找到指定节点，直接返回节点的ID
					return []string{node.ID}
				}
			}

			// 没有找到指定节点，从对应桶中随机返回两个节点的ID
			ids := make([]string, 0)
			for i := 0; i < 2; i++ {
				if len(bucket.Nodes) > 0 {
					randIdx := rand.Intn(len(bucket.Nodes))
					ids = append(ids, bucket.Nodes[randIdx].ID)
				}
			}
			return ids
		}
	}
	return nil
}
func test() {
	kBucket := NewK_Bucket()

	// 插入节点
	kBucket.insertNode("00112233445566778899")
	kBucket.insertNode("aabbccddeeff00112233")
	kBucket.insertNode("11223344556677889900")
	kBucket.insertNode("22334455667788990011")
	kBucket.insertNode("33445566778899001122")
	kBucket.insertNode("33445566778899001123")
	kBucket.insertNode("33445566778899001124")
	kBucket.insertNode("33445566778899001125")

	// 打印每个桶中的节点内容
	kBucket.printBucketContents()

	// 查询节点
	fmt.Println(kBucket.findNode("00112233445566778899"))
	fmt.Println(kBucket.findNode("11223344556677889900"))
	fmt.Println(kBucket.findNode("22334455667788990011"))
	fmt.Println(kBucket.findNode("33445566778899001122"))
	fmt.Println(kBucket.findNode("33445566778899001123"))
	fmt.Println(kBucket.findNode("33445566778899001124"))
	fmt.Println(kBucket.findNode("33445566778899001125"))
}
func main() {
	test()
}
