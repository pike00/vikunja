// Vikunja is a to-do list application to facilitate your life.
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

package v1

import (
	"code.vikunja.io/api/pkg/config"
	"code.vikunja.io/api/pkg/modules/migration/todoist"
	"code.vikunja.io/api/pkg/modules/migration/wunderlist"
	"code.vikunja.io/api/pkg/version"
	"github.com/labstack/echo/v4"
	"net/http"
)

type vikunjaInfos struct {
	Version                    string    `json:"version"`
	FrontendURL                string    `json:"frontend_url"`
	Motd                       string    `json:"motd"`
	LinkSharingEnabled         bool      `json:"link_sharing_enabled"`
	MaxFileSize                string    `json:"max_file_size"`
	RegistrationEnabled        bool      `json:"registration_enabled"`
	AvailableMigrators         []string  `json:"available_migrators"`
	TaskAttachmentsEnabled     bool      `json:"task_attachments_enabled"`
	EnabledBackgroundProviders []string  `json:"enabled_background_providers"`
	TotpEnabled                bool      `json:"totp_enabled"`
	Legal                      legalInfo `json:"legal"`
	CaldavEnabled              bool      `json:"caldav_enabled"`
}

type legalInfo struct {
	ImprintURL       string `json:"imprint_url"`
	PrivacyPolicyURL string `json:"privacy_policy_url"`
}

// Info is the handler to get infos about this vikunja instance
// @Summary Info
// @Description Returns the version, frontendurl, motd and various settings of Vikunja
// @tags service
// @Produce json
// @Success 200 {object} v1.vikunjaInfos
// @Router /info [get]
func Info(c echo.Context) error {
	info := vikunjaInfos{
		Version:                version.Version,
		FrontendURL:            config.ServiceFrontendurl.GetString(),
		Motd:                   config.ServiceMotd.GetString(),
		LinkSharingEnabled:     config.ServiceEnableLinkSharing.GetBool(),
		MaxFileSize:            config.FilesMaxSize.GetString(),
		RegistrationEnabled:    config.ServiceEnableRegistration.GetBool(),
		TaskAttachmentsEnabled: config.ServiceEnableTaskAttachments.GetBool(),
		TotpEnabled:            config.ServiceEnableTotp.GetBool(),
		CaldavEnabled:          config.ServiceEnableCaldav.GetBool(),
		Legal: legalInfo{
			ImprintURL:       config.LegalImprintURL.GetString(),
			PrivacyPolicyURL: config.LegalPrivacyURL.GetString(),
		},
	}

	// Migrators
	if config.MigrationWunderlistEnable.GetBool() {
		m := &wunderlist.Migration{}
		info.AvailableMigrators = append(info.AvailableMigrators, m.Name())
	}
	if config.MigrationTodoistEnable.GetBool() {
		m := &todoist.Migration{}
		info.AvailableMigrators = append(info.AvailableMigrators, m.Name())
	}

	if config.BackgroundsEnabled.GetBool() {
		if config.BackgroundsUploadEnabled.GetBool() {
			info.EnabledBackgroundProviders = append(info.EnabledBackgroundProviders, "upload")
		}
		if config.BackgroundsUnsplashEnabled.GetBool() {
			info.EnabledBackgroundProviders = append(info.EnabledBackgroundProviders, "unsplash")
		}
	}

	return c.JSON(http.StatusOK, info)
}
