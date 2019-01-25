# Plumb

## Overview

plumb is a drop-in replacement for plan9's plumb command, that understands mime-types of local files and URLs. 

## Rules

plumb will write plumb messages as the usual `type is text` for non-files, and non-URLs. When the content is one of those, however, it will set the type to that mime.

```
# Examples of types sent to plumber by this implementation of plumb
type is image/png
type is text/html
type is text/plain
type is application/pdf

```

In practice, setting a rule for most common mimetypes will serve you well. For example, my last implementation of plumber used a relatively small amount of mimes, after a few years' worth of plumbing in this manner. 
See https://github.com/halfwit/Plumber/tree/master/cfg/plumber

## About application/octet-stream mimetype
If the remote mimetype is `application/octet-stream`, which is a fallback when it cannot infer the mimetype this client will attempt to find a content-type field in any remote URL, finally setting the type to `text`.

Reference for mimetype: https://mimesniff.spec.whatwg.org/

Reasonably robust list of mimetypes: https://www.freeformatter.com/mime-types-list.html
