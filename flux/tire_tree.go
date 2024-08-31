package flux

import "strings"

type node struct {
	urlPattern string  // 待匹配路由，例如 /p/:lang
	part       string  // 路由中的一部分，例如 :lang
	children   []*node // 子节点，例如 [doc, tutorial, intro]
	isExact    bool    // 是否精确匹配，part 含有 : 或 * 时为false
}

// 第一个匹配成功的节点，用于插入
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || !child.isExact {
			return child
		}
	}
	return nil
}

// 所有匹配成功的节点，用于查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || !child.isExact {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func (n *node) insert(pattern string, parts []string, deep int) {
	//匹配到了最后一个位置,添加pattern
	if len(parts) == deep {
		n.urlPattern = pattern
		return
	}
	//获取当前深度的part
	part := parts[deep]
	//查找是否有匹配的子节点
	child := n.matchChild(part)
	//如果没有匹配的子节点,则创建一个新的子节点
	if child == nil {
		child = &node{part: part, isExact: part[0] != ':' && part[0] != '*'}
		n.children = append(n.children, child)
	}
	//递归插入
	child.insert(pattern, parts, deep+1)
}

func (n *node) search(parts []string, deep int) *node {
	if len(parts) == deep || strings.HasPrefix(n.part, "*") {
		if n.urlPattern == "" {
			return nil
		}
		return n
	}

	part := parts[deep]
	children := n.matchChildren(part)
	for _, child := range children {
		result := child.search(parts, deep+1)
		if result != nil {
			return result
		}
	}

	return nil
}

// 解析urlPattern
func parseUrlPattern(pattern string) []string {
	rawParts := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, part := range rawParts {
		if part == "" {
			continue
		}
		parts = append(parts, part)
		if part[0] == '*' {
			break
		}
	}
	return parts
}
