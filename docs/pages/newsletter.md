---
date: 2025-01-01T06:00:00Z
title: "Newsletter - Markdown Ninja"
type: "page"
tags: ["docs"]
authors: ["Markdown Ninja"]
url: "/docs/newsletter"
---

Your posts can be sent by email to your subscribed contacts.

Simply set the `newsletter` field to `true` and `type` to `post` to send a post as a newsletter.

```yaml
---
date: 2025-01-01T06:00:00Z
title: "Newsletter example"
type: "post"
authors: ["Markdown Ninja"]
url: "/newsletter-example"
newsletter: true
---
```

The newsletter will be sent on the value of the `date` field, here `2025-01-01T06:00:00Z` or immediately if `date` is in the past.
