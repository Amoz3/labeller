// Code generated by rice embed-go; DO NOT EDIT.
package main

import (
	"time"

	"github.com/GeertJohan/go.rice/embedded"
)

func init() {

	// define files
	file2 := &embedded.EmbeddedFile{
		Filename:    "index.html",
		FileModTime: time.Unix(1673703038, 0),

		Content: string("<!DOCTYPE html>\n<html lang=\"en\">\n\n<head>\n    <title>Document</title>\n</head>\n\n<body>\n    <button onclick=\"pressValid()\">Valid</button>\n    <button onclick=\"pressInvalid()\">Invalid</button>\n    <button onclick=\"pause()\">Pause</button>\n</body>\n\n</html>\n\n<script>\n    function pressValid() {\n        console.log(\"valid\")\n        fetch(\"localhost:1235/valid\")\n    }\n    \n    function pause() {\n        console.log(\"pausing\")\n        fetch(\"http://127.0.0.1:1235/pause\")\n        return\n    }\n\n    function pressInvalid() {\n        fetch(\"127.0.0.1:1235/invalid\")\n        return\n    }\n</script>"),
	}

	// define dirs
	dir1 := &embedded.EmbeddedDir{
		Filename:   "",
		DirModTime: time.Unix(1673691498, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			file2, // "index.html"

		},
	}

	// link ChildDirs
	dir1.ChildDirs = []*embedded.EmbeddedDir{}

	// register embeddedBox
	embedded.RegisterEmbeddedBox(`app`, &embedded.EmbeddedBox{
		Name: `app`,
		Time: time.Unix(1673691498, 0),
		Dirs: map[string]*embedded.EmbeddedDir{
			"": dir1,
		},
		Files: map[string]*embedded.EmbeddedFile{
			"index.html": file2,
		},
	})
}
