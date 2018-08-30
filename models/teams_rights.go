package models

// TeamRight defines the rights teams can have for lists/namespaces
type TeamRight int

// define unknown team right
const (
	TeamRightUnknown = -1
)

// Enumerate all the team rights
const (
	// Can read lists in a Team
	TeamRightRead TeamRight = iota
	// Can write tasks in a Team like lists and todo tasks. Cannot create new lists.
	TeamRightWrite
	// Can manage a list/namespace, can do everything
	TeamRightAdmin
)

func (r TeamRight) isValid() error {
	if r != TeamRightAdmin && r != TeamRightRead && r != TeamRightWrite {
		return ErrInvalidTeamRight{r}
	}

	return nil
}

// CanCreate checks if the user can create a new team
func (t *Team) CanCreate(user *User) bool {
	// This is currently a dummy function, later on we could imagine global limits etc.
	return true
}

// CanUpdate checks if the user can update a team
func (t *Team) CanUpdate(user *User) bool {

	// Check if the current user is in the team and has admin rights in it
	exists, _ := x.Where("team_id = ?", t.ID).
		And("user_id = ?", user.ID).
		And("admin = ?", true).
		Get(&TeamMember{})

	return exists
}

// CanDelete checks if a user can delete a team
func (t *Team) CanDelete(user *User) bool {
	//t.ID = id
	return t.IsAdmin(user)
}

// IsAdmin returns true when the user is admin of a team
func (t *Team) IsAdmin(user *User) bool {
	exists, _ := x.Where("team_id = ?", t.ID).
		And("user_id = ?", user.ID).
		And("admin = ?", true).
		Get(&TeamMember{})
	return exists
}

// CanRead returns true if the user has read access to the team
func (t *Team) CanRead(user *User) bool {
	// Check if the user is in the team
	exists, _ := x.Where("team_id = ?", t.ID).
		And("user_id = ?", user.ID).
		Get(&TeamMember{})
	return exists
}
