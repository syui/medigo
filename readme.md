medium client

```bash
$ mkdir -p ~/.config/medigo
$ cp medium.json ~/.config/medigo/.medium.json 
$ go run main.go
$ go build
$ ./medigo
```

https://github.com/Medium/medium-sdk-go/blob/master/medium.go

```go
ContentFormatHTML     contentFormat = "html"
ContentFormatMarkdown               = "markdown"

PublishStatusDraft    publishStatus = "draft"
PublishStatusUnlisted               = "unlisted"
PublishStatusPublic                 = "public"

LicenseAllRightsReserved license = "all-rights-reserved"
LicenseCC40By                    = "cc-40-by"
LicenseCC40BySA                  = "cc-40-by-sa"
LicenseCC40ByND                  = "cc-40-by-nd"
LicenseCC40ByNC                  = "cc-40-by-nc"
LicenseCC40ByNCND                = "cc-40-by-nc-nd"
LicenseCC40ByNCSA                = "cc-40-by-nc-sa"
LicenseCC40Zero                  = "cc-40-zero"
LicensePublicDomain              = "public-domain"

formatJSON = "json"
formatForm = "form"
formatFile = "file"

UserID        string        `json:"-"`
Title         string        `json:"title"`
Content       string        `json:"content"`
ContentFormat contentFormat `json:"contentFormat"`
Tags          []string      `json:"tags,omitempty"`
CanonicalURL  string        `json:"canonicalUrl,omitempty"`
PublishStatus publishStatus `json:"publishStatus,omitempty"`
License       license       `json:"license,omitempty"`
```

client-id, client-secret-id

```bash
$ go run -ldflags "-X main.cid=xxx -X main.secret=xxx" main.go
$ go build -ldflags "-X main.cid=xxx -X main.secret=xxx"
```

https://github.com/Medium/medium-api-docs

```
## publications
scope +listPublications	
GET https://api.medium.com/v1/users/{{userId}}/publications

## contributors
GET https://api.medium.com/v1/publications/{{publicationId}}/contributors
```

