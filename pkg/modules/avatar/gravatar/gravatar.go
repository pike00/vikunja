// Vikunja is a to-do-list application to facilitate your life.
// Copyright 2018-2020 Vikunja and contributors. All rights reserved.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package gravatar

import (
	"code.vikunja.io/api/pkg/config"
	"code.vikunja.io/api/pkg/log"
	"code.vikunja.io/api/pkg/user"
	"code.vikunja.io/api/pkg/utils"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type avatar struct {
	content  []byte
	loadedAt time.Time
}

// Provider is the gravatar provider
type Provider struct {
}

// avatars is a global map which contains cached avatars of the users
var avatars map[string]*avatar

func init() {
	avatars = make(map[string]*avatar)
}

// GetAvatar implements getting the avatar for the user
func (g *Provider) GetAvatar(user *user.User, size int64) ([]byte, string, error) {
	sizeString := strconv.FormatInt(size, 10)
	cacheKey := user.Username + "_" + sizeString
	a, exists := avatars[cacheKey]
	var needsRefetch bool
	if exists {
		// elaped is alway < 0 so the next check would always succeed.
		// To have it make sense, we flip that.
		elapsed := time.Until(a.loadedAt) * -1
		needsRefetch = elapsed > time.Duration(config.AvatarGravaterExpiration.GetInt64())*time.Second
		if needsRefetch {
			log.Debugf("Refetching avatar for user %d after %v", user.ID, elapsed)
		} else {
			log.Debugf("Serving avatar for user %d from cache", user.ID)
		}
	}
	if !exists || needsRefetch {
		log.Debugf("Gravatar for user %d with size %d not cached, requesting from gravatar...", user.ID, size)
		resp, err := http.Get("https://www.gravatar.com/avatar/" + utils.Md5String(user.Email) + "?s=" + sizeString + "&d=mp")
		if err != nil {
			return nil, "", err
		}
		avatarContent, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, "", err
		}
		avatars[cacheKey] = &avatar{
			content:  avatarContent,
			loadedAt: time.Now(),
		}
	}
	return avatars[cacheKey].content, "image/jpg", nil
}
