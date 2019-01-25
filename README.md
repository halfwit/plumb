# Plumb

## Overview

plumb is a drop-in replacement for plan9's plumb command, that understands mime-types of local files and URLs. 

## Rules

plumb will write plumb messages as the usual `type is text` for non-files, and non-URLs. When the content is one of those, however, it will set the type to that mime.

```

type is image/png
type is text/html
type is text/plain
type is application/pdf

```

If the remote mimetype is `application/octet-stream`, which is a fallback when it cannot infer the mimetype this client will attempt to find a content-type field in any remote URL, finally setting the type to `text`.
