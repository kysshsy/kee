package kee

import "strings"

type node struct {
	pattern string
	part    string
	child   []*node
	isWild  bool
}

// 用于插入 第一个匹配child
func (n *node) matchChild(part string) *node {
	for _, c := range n.child {
		if c.part == part || c.isWild {
			return c
		}
	}
	return nil
}

// 用于搜索  所有匹配的child
func (n *node) matchChildren(part string) []*node {
	arr := make([]*node, 0)

	for _, c := range n.child {
		if c.part == part || c.isWild {
			arr = append(arr, c)
		}
	}
	return arr
}

func (n *node) insert(pattern string, parts []string, height int) {
	if height == len(parts) {
		n.pattern = pattern
		return
	}

	child := n.matchChild(parts[height])
	part := parts[height]

	if child == nil {
		child = &node{part: part, child: make([]*node, 0, 0), isWild: part[0] == '*' || part[0] == ':'}
		n.child = append(n.child, child)
	}

	child.insert(pattern, parts, height+1)

}

func (n *node) search(parts []string, height int) *node {
	if height == len(parts) || strings.HasPrefix(parts[height], "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}
	children := n.matchChildren(parts[height])

	for _, child := range children {
		ret := child.search(parts, height+1)
		if ret != nil {
			return ret
		}
	}

	return nil
}

func (n *node) travel(nodes *[]*node) {

	if n.pattern != "" {
		*nodes = append(*nodes, n)
	}

	for _, child := range n.child {
		child.travel(nodes)
	}
}
