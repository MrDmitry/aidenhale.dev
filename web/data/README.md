# Directory layout

```bash
navigation_item/    # navbar part of the url
├── article_id      # article part of the url
│   ├── assets          # (optional) location with article-specific assets
│   │   └── file.ext        # asset file to be served as-is
│   ├── extra           # (optional) nested sub-pages
│   │   └── subpage_id      # extra part of the url
│   │       └── README.md   # sub-page body
│   ├── monke.toml      # article metadata
│   └── README.md       # article body
├── monke.toml      # navbar metadata
└── README.md       # navbar page body
```
