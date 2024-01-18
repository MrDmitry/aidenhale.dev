# Directory layout

```bash
category_id     # category part of the url
├── article_id      # article part of the url
│   ├── assets          # (optional) location with article-specific assets
│   │   └── file.ext        # asset file to be served as-is
│   ├── extra           # (optional) nested sub-pages
│   │   └── subpage_id      # extra part of the url
│   │       └── README.md   # sub-page body
│   ├── monke.toml      # article metadata
│   └── README.md       # article body
└── monke.toml  # category metadata
```

## Developer commentary

Layout above is the skeleton for processing and routing of the articles. Initially I wanted each category to also have
some description, but gave up on the idea because I could not find a use for it. It's cool to have a `README.md` on
GitHub so you could take a peek, but it just did not make sense for a website.

Apart from being a part of the URL, categories can define tags to be inherited by all articles under them. But the main
purpose is to have another dimension for the URLs indicating the "vibe" of the article.
