package app_test

import (
	"fmt"
	"testing"

	"github.com/canonical/go-dqlite/v3/app"
	"github.com/canonical/go-dqlite/v3/client"
	"github.com/stretchr/testify/assert"
)

func TestRolesChanges_Assume(t *testing.T) {
	online := &client.NodeMetadata{}

	node := func(id uint64, role client.NodeRole) client.NodeInfo {
		return client.NodeInfo{ID: id, Address: fmt.Sprintf("127.0.0.1:900%d", id), Role: role}
	}

	cases := []struct {
		name  string
		state map[client.NodeInfo]*client.NodeMetadata
		role  client.NodeRole
	}{{
		name: "cluster too small",
		state: map[client.NodeInfo]*client.NodeMetadata{
			node(1, client.Voter): online,
			node(2, client.Spare): online,
		},
		role: -1,
	}, {
		name: "already a voter",
		state: map[client.NodeInfo]*client.NodeMetadata{
			node(1, client.Voter): online,
			node(2, client.Spare): online,
			node(3, client.Voter): online,
		},
		role: -1,
	}, {
		name: "single voter",
		state: map[client.NodeInfo]*client.NodeMetadata{
			node(1, client.Voter): online,
			node(2, client.Spare): online,
			node(3, client.Spare): online,
		},
		role: client.Voter,
	}, {
		name: "unreachable spares don't matter",
		state: map[client.NodeInfo]*client.NodeMetadata{
			node(1, client.Voter): online,
			node(2, client.Spare): nil,
			node(3, client.Spare): online,
		},
		role: client.Voter,
	}, {
		name: "missing stand-bys",
		state: map[client.NodeInfo]*client.NodeMetadata{
			node(1, client.Voter): online,
			node(2, client.Voter): online,
			node(3, client.Voter): online,
			node(4, client.Spare): online,
		},
		role: client.StandBy,
	}, {
		name: "unreachable voters",
		state: map[client.NodeInfo]*client.NodeMetadata{
			node(1, client.Voter): online,
			node(2, client.Voter): nil,
			node(3, client.Voter): nil,
			node(4, client.Spare): online,
		},
		role: -1,
	}, {
		name: "unreachable stand-by",
		state: map[client.NodeInfo]*client.NodeMetadata{
			node(1, client.Voter):   online,
			node(2, client.Voter):   online,
			node(3, client.Voter):   online,
			node(4, client.StandBy): online,
			node(5, client.StandBy): nil,
			node(6, client.Spare):   online,
		},
		role: -1,
	}}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var self uint64
			for node := range c.state {
				if node.ID > self {
					self = node.ID
				}
			}

			roles := app.RolesChanges{
				Config: app.RolesConfig{Voters: 3, StandBys: 3},
				State:  c.state,
			}
			assert.Equal(t, c.role, roles.Assume(self))
		})
	}
}
