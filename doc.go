// Copyright © 2019-2020 A Bunch Tell LLC. and contributors.
//
// This is free software: you can redistribute it and/or modify
// it under the terms of the Mozilla Public License, included
// in the LICENSE file in this source code package.

/*
Package wfimport provides library functions for parsing and importing
writefreely posts from various files. https://github.com/writeas/writefreely

Status

Current support is for files, directories and zip archives.
Support is planned for exported data from Medium, Ghost, Wordpress and
writefreely exports in zip or json.

About Posts

In the context of this package a post is actually referring to a PostParams from
https://github.com/writeas/go-writeas.

	// PostParams holds values for creating or updating a post.
	PostParams struct {
		// Parameters only for updating
		ID    string `json:"-"`
		Token string `json:"token,omitempty"`

		// Parameters for creating or updating
		Slug     string     `json:"slug"`
		Created  *time.Time `json:"created,omitempty"`
		Updated  *time.Time `json:"updated,omitempty"`
		Title    string     `json:"title,omitempty"`
		Content  string     `json:"body,omitempty"`
		Font     string     `json:"font,omitempty"`
		IsRTL    *bool      `json:"rtl,omitempty"`
		Language *string    `json:"lang,omitempty"`

		// Parameters only for creating
		Crosspost []map[string]string `json:"crosspost,omitempty"`

		// Parameters for collection posts
		Collection string `json:"-"`
	}

*/
package wfimport
