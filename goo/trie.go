package goo

import "strings"

type node struct {
	pattern  string
	part     string
	children []*node
	isWild   bool
}

func (n *node) matchChild(path string) *node {
	for _, child := range n.children {
		if child.part == path || child.isWild {
			return child
		}
	}

	return nil
}

func (n *node) matchChildren(path string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == path || child.isWild {
			nodes = append(nodes, child)
		}
	}

	return nodes
}

func (n *node) insert(pattern string, parts []string, height int) {
	// 最后一层才设置 pattern 用于判断匹配成功
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	if nil == child {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}

	child.insert(pattern, parts, height+1)
}

func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}

		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		node := child.search(parts, height+1)
		if nil != node {
			return node
		}
	}

	return nil
}
