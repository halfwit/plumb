# Plumb

## Overview

Plumb is a drop-in replacement for plan9's [plumb](https://9fans.github.io/plan9port/man/man1/plumb.html) utility.

Notably different from it, are how it handles the `type` attribute of the plumb messages.
They will be set to a proper mimetype, instead of the simple `type is text` that the traditional plumber utilized. 

This is considerably more powerful, as you no longer need to attempt to infer the content based on URIs, extensions, or directory structure.

## Rules

Your plumber rules will have to be updated to reflect this more granular message:

```
## In this example, plan9front's plumber is being used; but similar rule changes would apply for plan9 regular.

type is image/png
## No longer need to match paths
#data matches '[a-zA-Z¡-￿0-9_\-./@]+'
#data matches '([a-zA-Z¡-￿0-9_\-./@]+)\.(jpe?g|JPE?G|gif|GIF|tiff?|TIFF?|ppm|bit|png|PNG)'
arg isfile	$0
plumb to image
plumb start 9 page $file

# local html files (Your dev work, for example) can be opened in your editor
type is text/html
arg isfile $0
plumb to edit
plumb start $editor $0

# remote html files will likely be opened in your browser
type is text/html
# With the above case matching local files, we no longer need this convolution
#data matches '(https?|ftp|file|gopher|mailto|news|nntp|telnet|wais|prospero)://[a-zA-Z0-9_@\-]+([.:][a-zA-Z0-9_@\-]+)*/?[a-zA-Z0-9_?,%#~&/\-+=]+([:.][a-zA-Z0-9_?,%#~&/\-+=]+)*'
plumb to web
plumb start web $0

type is application/pdf
## No longer need to match paths
#data matches '[a-zA-Z¡-￿0-9_\-./@]+'
#data matches '([a-zA-Z¡-￿0-9_\-./@]+)\.(ps|PS|eps|EPS|pdf|PDF|dvi|DVI)'
arg isfile	$0
plumb to postscript
plumb start 9 page $file
```

In practice, setting a rule for most common mimetypes will serve you well. (For example, my last implementation of plumber used a relatively small amount of mimes, after a few years' worth of plumbing in this manner. 
See https://github.com/halfwit/Plumber/tree/master/cfg/plumber)

## About application/octet-stream mimetype
If the remote mimetype is `application/octet-stream`, which is a fallback when it cannot infer the mimetype this client will attempt to find a content-type field in any remote URL, finally setting the type to `text`.

Reference for mimetype: https://mimesniff.spec.whatwg.org/

Reasonably robust list of mimetypes: https://www.freeformatter.com/mime-types-list.html
