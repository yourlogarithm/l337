package team

type Mode string

const (
	// In Collaborate Mode, all team members respond to the user query at once. This gives the team coordinator to review whether the team has reached a consensus on a particular topic and then synthesize the responses from all team members into a single response.​
	ModeCollaborate Mode = "collaborate"
	// In Coordinate Mode, the Team Leader delegates tasks to team members and synthesizes their outputs into a cohesive response.
	ModeCoordinate Mode = "coordinate"
	// In Route Mode, the Team Leader directs user queries to the most appropriate team member based on the content of the request.
	// The Team Leader acts as a smart router, analyzing the query and selecting the best-suited agent to handle it. The member’s response is then returned directly to the user.
	ModeRoute Mode = "route"
)
