package merkledag

// Hash to file
func Hash2File(store KVStore, hash []byte, path string, hp HashPool) []byte {
	// 先从 KVStore 中获取 hash 对应的节点
	nodeBytes, err := store.Get(hash)
	if err != nil {
		// 错误处理：无法找到对应的 hash
		return nil
	}
	// 根据 nodeBytes 构建节点
	node, err := constructNodeFromBytes(nodeBytes, hp)
	if err != nil {
		// 错误处理：无法构建节点
		return nil
	}
	// 确保根节点是一个目录
	rootDir, ok := node.(Dir)
	if !ok {
		// 错误处理：根节点不是目录
		return nil
	}

	// 根据路径寻找文件内容
	fileBytes, found := findFileInDir(rootDir, path)
	if !found {
		// 错误处理：无法在目录中找到对应路径的文件
		return nil
	}

	return fileBytes
}
func constructNodeFromBytes(nodeBytes []byte, hp HashPool) (Node, error) {
	// 这里根据读取到的字节和具体实现进行节点的构建
	// 注意：此处的实现与具体数据结构有关，需要根据具体情况进行解析
	// 这里省略具体地构建逻辑
	return nil, nil
}

func findFileInDir(dir Dir, path string) ([]byte, bool) {
	iter := dir.It()
	for iter.Next() {
		childNode := iter.Node()
		switch child := childNode.(type) {
		case File:
			// 判断文件路径是否匹配
			if path == "file.txt" { // 假设这里使用 "file.txt" 作为文件路径
				return child.Bytes(), true
			}
		case Dir:
			// 如果是目录，则继续在子目录中查找
			fileBytes, found := findFileInDir(child, path)
			if found {
				return fileBytes, true
			}
		}
	}

	return nil, false
}
