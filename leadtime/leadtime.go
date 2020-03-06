package leadtime

import (
	"fmt"
	"time"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 github.com/concourse/concourse/go-concourse/concourse.Client

func LeadTime(repository Repository) time.Duration {
	node := findStart(repository)
	Print(repository, node, map[int]bool{}, 0)
	return time.Since(node.Finish)
}

func findStart(repository Repository) Node {
	for _, node := range repository.Nodes("concourse", "publish-binaries") {
		if node.Succeeded {
			return node
		}
	}
	panic("couldn't find a start")
}

func Print(repository Repository, node Node, visited map[int]bool, depth int) {
	fmt.Println(fmt.Sprintf("%*s", depth, ""), node.Finish, node)
	visited[node.ID] = true
	if node.Job == "unit" {
		return
	}
	for _, input := range repository.InputsForNode(node) {
		link, err := repository.LinkForInput(node, input)
		if err != nil {
			// fmt.Println(err)
		} else {
			for _, n := range repository.NodesForLink(link) {
				fresh := !visited[n.ID]
				success := n.Succeeded
				precedes := n.Finish.Before(node.Finish)
				if fresh && success && precedes {
					Print(repository, n, visited, depth+1)
				}
			}
		}
	}
}
