package team

import (
	"github.com/yourlogarithm/l337/agentic"
	"github.com/yourlogarithm/l337/internal/logging"
)

var logger = logging.SetupLogger("team")

type Team struct {
	initialized bool

	// Members of the team, either `agent.Agent`, `Team`, or custom `agentic.Member` implementations.
	Members []agentic.Member
	// One of "collaborate", "coordinate", or "route".
	Mode Mode
	agentic.Options
}

func (t *Team) initialize() error {
	if err := t.Options.Initialize(); err != nil {
		return err
	}

	if t.initialized {
		return nil
	}

	switch t.Mode {
	case ModeCollaborate:
		t.Tools.AddTool(t.generateCollaborateTool())
	case ModeCoordinate:
		t.Tools.AddTool(t.generateCoordinateTool())
	case ModeRoute:
		t.Tools.AddTool(t.generateRouteTool())
	}

	t.initialized = true

	return nil
}

// Returns the `agentic.MemberTypeTeam` constant.
func (t *Team) Type() agentic.MemberType {
	return agentic.MemberTypeTeam
}

func (t *Team) GetOptions() *agentic.Options {
	return &t.Options
}

func (t *Team) findMemberByID(memberID string) *agentic.Member {
	for _, member := range t.Members {
		opts := member.GetOptions()
		if opts.ID == memberID {
			return &member
		}
	}
	return nil
}
