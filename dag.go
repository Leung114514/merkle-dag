package merkledag

import (
	"hash"
)

func Add(store KVStore, node Node, h hash.Hash) []byte {
	// 将分片写入到KVStore中
	Store(store, node)

	// 计算Merkle Root
	root := calMkR(node, h)

	return root
}

func Store(store KVStore, node Node) {
	switch n := node.(type) {
	case File:
		_ = store.Put([]byte("file"), n.Bytes())
	case Dir:
		iter := n.It()
		for iter.Next() {
			childNode := iter.Node()
			Store(store, childNode)
		}
	}
}

func calMkR(node Node, h hash.Hash) []byte {
	switch n := node.(type) {
	case File:
		h.Write(n.Bytes())
	case Dir:
		iter := n.It()
		for iter.Next() {
			childNode := iter.Node()
			childHash := calMkR(childNode, h)
			h.Write(childHash)
		}
	}
	return h.Sum(nil)
}
