package team

import (
	"github.com/yourlogarithm/golagno/agentic"
	"github.com/yourlogarithm/golagno/logging"
)

var logger = logging.SetupLogger("team")

type Team struct {
	initialized bool

	Members []agentic.Member
	Mode    Mode

	agentic.Options
}

func (t *Team) initialize() bool {
	t.Options.SetupID()

	if t.initialized {
		return false
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

	return true
}

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
