package gee

import (
	"strings"
)

type node struct {
	path     string  // 待匹配路由，例如 /p/:lang
	part     string  // 路由中的一部分，例如 :lang
	children []*node // 子节点，例如 [doc, tutorial, intro]
	isMatch  bool    // 是否精确匹配，part 含有 : 或 * 时为true
}

//golang中方法接收者，使用指针和非指针的区别
// 使用指针方式接收，可以修改收到的结构的字段值, 不会随着函数的销毁而失效；
// 使用非指针方式接收，只能在函数内部修改收到的结构的字段值, 函数被销毁，这个改变也就失效了
// 第一个匹配成功的节点，用于插入
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isMatch {
			return child
		}
	}
	return nil
}

func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)

	for _, child := range n.children {
		if child.part == part || child.isMatch {
			nodes = append(nodes, child)
		}
	}

	return nodes
}

func (n *node) insert(path string, parts []string, height int) {
	//递归插入的终止条件 匹配到了第len(parts)层节点 则退出
	if len(parts) == height {
		n.path = path
		return
	}

	//对前面几层的处理
	part := parts[height]
	child := n.matchChild(part)
	//没有找到 则创建树节点 然后创建下一层
	if child == nil {
		child = &node{part: part, isMatch: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	//进行递归
	child.insert(path, parts, height+1)
}

func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.path == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}