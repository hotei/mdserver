// doc.go (c) David Rook - all rights reserved

// Program mdserver serves markdown documents via web interface
//
//	usage: mdserver & (usually run in background but not always if testing)
//
//	Reads directory tree and picks out all files that are markdown documents.
//	These are then presented as filenames to user via browser.
//	When clicked, the file is run through the markdown process and shown to user.
//
// The list is not refreshed until the server is restarted. An enhancement would
// be to present a form at the top to allow refresh.
//
//	Change Log
//	2014-03-25 working
//	2014-03-25 started
package main

