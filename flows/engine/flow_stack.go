package engine

import (
	"github.com/nyaruka/goflow/flows"
)

type flowFrame struct {
	flow         flows.Flow
	visitedNodes map[flows.NodeUUID]bool
}

type flowStack struct {
	stack []*flowFrame
}

// creates a new empty flow stack
func newFlowStack() *flowStack {
	return &flowStack{stack: make([]*flowFrame, 0)}
}

// creates a new flow stack based on the ancestry of the given run
func flowStackFromRun(run flows.FlowRun) *flowStack {
	s := newFlowStack()
	ancestors := run.Ancestors()
	for a := len(ancestors) - 1; a >= 0; a-- {
		s.push(ancestors[a].Flow())
	}
	s.push(run.Flow())
	return s
}

// creates a new frame for the given flow and pushes it onto the stack
func (s *flowStack) push(flow flows.Flow) {
	s.stack = append(s.stack, &flowFrame{flow: flow, visitedNodes: make(map[flows.NodeUUID]bool)})
}

// pops the current frame off the stack
func (s *flowStack) pop() {
	s.stack = s.stack[:len(s.stack)-1]
}

// records the given node as visited in the current frame
func (s *flowStack) visit(nodeUUID flows.NodeUUID) {
	s.stack[len(s.stack)-1].visitedNodes[nodeUUID] = true
}

// checks whether the given node has already been visited in the current frame
func (s *flowStack) hasVisited(nodeUUID flows.NodeUUID) bool {
	return s.stack[len(s.stack)-1].visitedNodes[nodeUUID]
}

func (s *flowStack) hasFlow(flowUUID flows.FlowUUID) bool {
	for _, f := range s.stack {
		if f.flow.UUID() == flowUUID {
			return true
		}
	}
	return false
}

func (s *flowStack) depth() int { return len(s.stack) }
