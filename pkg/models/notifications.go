// Vikunja is a to-do list application to facilitate your life.
// Copyright 2018-2021 Vikunja and contributors. All rights reserved.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public Licensee as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public Licensee for more details.
//
// You should have received a copy of the GNU Affero General Public Licensee
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package models

import (
	"strconv"

	"code.vikunja.io/api/pkg/config"
	"code.vikunja.io/api/pkg/notifications"
	"code.vikunja.io/api/pkg/user"
)

// ReminderDueNotification represents a ReminderDueNotification notification
type ReminderDueNotification struct {
	User *user.User
	Task *Task
}

// ToMail returns the mail notification for ReminderDueNotification
func (n *ReminderDueNotification) ToMail() *notifications.Mail {
	return notifications.NewMail().
		To(n.User.Email).
		Subject(`Reminder for "`+n.Task.Title+`"`).
		Greeting("Hi "+n.User.GetName()+",").
		Line(`This is a friendly reminder of the task "`+n.Task.Title+`".`).
		Action("Open Task", config.ServiceFrontendurl.GetString()+"tasks/"+strconv.FormatInt(n.Task.ID, 10)).
		Line("Have a nice day!")
}

// ToDB returns the ReminderDueNotification notification in a format which can be saved in the db
func (n *ReminderDueNotification) ToDB() interface{} {
	return nil
}