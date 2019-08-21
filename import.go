// Package wfimport provides library functions for parsing and importing
// writefreely posts from various files.
//
// Current support is for files, directories and zip archives.
// Support is planned for exported data from Medium, Ghost, Wordpress and
// writefreely exports in zip or json.

package wfimport

const (
	// DraftsKey is the key for all parsed draft posts in a ZipCollections map
	DraftsKey = "drafts"
)
