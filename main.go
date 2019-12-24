package main

import (
    "fmt"
    "log"
    "os"
    "time"
    "syscall"
    "io/ioutil"
    "net/http"
    "strings"
    "encoding/json"
    "path/filepath"
    "golang.org/x/crypto/ssh/terminal"
    "github.com/urfave/cli/v2"
    //"github.com/medium/medium-sdk-go"
    "github.com/syui/medium-sdk-go"
    "github.com/skratchdot/open-golang/open"
    "github.com/hokaccha/go-prettyjson"
    //cregex "github.com/mingrammer/commonregex"
)

// Oauth Access Token
type Oauth struct {
    AccessToken  string   `json:"access_token"`
    ExpiresAt    int64    `json:"expires_at"`
    RefreshToken string   `json:"refresh_token"`
    Scope        []string `json:"scope"`
    TokenType    string   `json:"token_type"`
}

//Self Access Token
type Self struct {
    SelfToken    string   `json:"self_token"`
}

// PostConfig date
type PostConfig struct {
    Title		string		`json:"title"`
    Tags		[]string	`json:"tags"`
    Content		string		`json:"content"`
    CanonicalURL	string		`json:"url"`
}

// UserJSON date
type UserJSON struct {
    ID       string `json:"id"`
    ImageURL string `json:"imageUrl"`
    Name     string `json:"name"`
    URL      string `json:"url"`
    Username string `json:"username"`
}

//// Post date
//type Post struct {
//	CanonicalURL	string `json:"canonicalUrl"`
//	Content		string `json:"content"`
//	ID		string `json:"id"`
//	MediumURL	string `json:"mediumUrl"`
//}

// Post date
type Post struct {
    Content string   `json:"content"`
    Lisense string   `json:"lisense"`
    Tags    []string `json:"tags"`
    Title   string   `json:"title"`
}

var cid string
var secret string
var token string

func check(e error) {
    if e != nil {
	panic(e)
    }
}

func main() {
    // Go to https://medium.com/me/applications to get your applicationId and applicationSecret.
    // Self-Issued Access Tokens : https://github.com/Medium/medium-api-docs#22-self-issued-access-tokens

    //var s Self
    var o Oauth
    var b PostConfig
    var userjson UserJSON
    var bodyjson Post

    callurl := "https://syui.github.io/medigo/callback/medium"

    app := cli.NewApp()
    app.Version = "0.2"
    dir := filepath.Join(os.Getenv("HOME"), ".config", "medigo")
    dirPost := filepath.Join(dir, "posts")
    dirFile := filepath.Join(dir, "body.json")
    dirConf := filepath.Join(dir, "medium.json")
    dirSelf := filepath.Join(dir, ".self.json")
    dirUser := filepath.Join(dir, "user.json")
    dirArti := filepath.Join(dir, "article.json")
    dirPubl := filepath.Join(dir, "publication.json")

    if e := os.MkdirAll(dirPost, os.ModePerm); e != nil {
	panic(e)
    }

    m := medium.NewClient(cid, secret)
    _, e := os.Stat(dirConf)
    if e != nil {
	url := m.GetAuthorizationURL("secretstate", callurl,
	medium.ScopeBasicProfile, medium.ScopePublishPost, medium.ScopeListPublications)
	println(url)
	time.Sleep(1 * time.Second)
	open.Run(url)
	time.Sleep(1 * time.Second)
	fmt.Print("authorization code : ")
	code, e := terminal.ReadPassword(int(syscall.Stdin))
	if e != nil {
	    log.Fatal(e)
	}
	fmt.Print("\ninput code : ", string(code))
	at, e := m.ExchangeAuthorizationCode(string(code), callurl)
	if e != nil {
	    log.Fatal(e)
	}
	outputJSON, e := json.Marshal(&at)
	if e != nil {
	    panic(e)
	}
	jat, _ := prettyjson.Marshal(at)

	fmt.Printf("\nYour token is %s\n", jat)
	ioutil.WriteFile(dirConf, outputJSON, os.ModePerm)

    }

    _, e = os.Stat(dirFile)
    if e != nil {
	var bodytmp = []byte(`
	{
	    "title":"test",
	    "tags":["test"],
	    "lisense":"LicenseCC40Zero",
	    "content":"body"
	}
	`)
	ioutil.WriteFile(dirFile, bodytmp, os.ModePerm)
    }

    file,e := ioutil.ReadFile(dirConf)
    if e != nil {
	fmt.Printf("File eor: %v\n", e)
	os.Exit(1)
    }
    json.Unmarshal(file, &o)

    body,e := ioutil.ReadFile(dirFile)
    if e != nil {
	fmt.Printf("File eor: %v\n", e)
	os.Exit(1)
    }
    json.Unmarshal(body, &b)

    m2 := medium.NewClientWithAccessToken(o.AccessToken)
    m.AccessToken = m2.AccessToken

    u, e := m2.GetUser()
    //if e != nil {
    //	fmt.Printf("rm %s\n", dirConf)
    //	os.Remove(dirConf)
    //	log.Fatal(e)
    //}

    // Refresh Token : dirConf(ctime -49)
    rt, e := m.ExchangeRefreshToken(o.RefreshToken)
    if e != nil {
	log.Fatal(e)
    }
    outputRF, e := json.Marshal(&rt)
    if e != nil {
	panic(e)
    }
    ioutil.WriteFile(dirConf, outputRF, os.ModePerm)

    app.Commands = []*cli.Command{
	{
	    Name:    "post",
	    Aliases: []string{"p"},
	    Usage:   "carte post\n\t\tsub-command : draft(d), public(p)",
	    Action:  func(c *cli.Context) error {
		fmt.Println(dirFile)
		fileinfos, _ := ioutil.ReadDir(dirPost)
		for _,fileinfo := range fileinfos {
		    fmt.Println(fileinfo.Name())

		}
		bodyfile, _ := ioutil.ReadFile(dirFile)
		json.Unmarshal(bodyfile, &bodyjson)
		fmt.Println(string(bodyfile))
		return nil
	    },
	    Subcommands: []*cli.Command{
		{
		    Name:   "draft",
		    Usage:   "draft",
		    Aliases: []string{"d"},
		    Action:  func(c *cli.Context) error {
			p,e := m.CreatePost(medium.CreatePostOptions{
			    UserID:		u.ID,
			    Title:		b.Title,
			    Content:	b.Content,
			    Tags:		b.Tags,
			    //CanonicalURL:	b.CanonicalURL,
			    ContentFormat: medium.ContentFormatMarkdown,
			    PublishStatus: medium.PublishStatusDraft,
			})
			jp, _ := prettyjson.Marshal(p)
			fmt.Println(string(jp))
			if e != nil {
			    log.Fatal(e)
			}
			return nil
		    },
		},
		{
		    Name:   "public",
		    Usage:   "public",
		    Aliases: []string{"p"},
		    Action:  func(c *cli.Context) error {
			p,e := m.CreatePost(medium.CreatePostOptions{
			    UserID:		u.ID,
			    Title:		b.Title,
			    Content:	b.Content,
			    Tags:		b.Tags,
			    //CanonicalURL:		b.CanonicalURL,
			    ContentFormat: medium.ContentFormatMarkdown,
			    PublishStatus: medium.PublishStatusPublic,
			})
			jp, _ := prettyjson.Marshal(p)
			fmt.Println(string(jp))
			fmt.Println(p)
			if e != nil {
			    log.Fatal(e)
			}
			return nil
		    },
		},
	    },
	},
	{
	    Name:    "key",
	    Aliases: []string{"k"},
	    Action:  func(c *cli.Context) error {
		jo, _ := prettyjson.Marshal(o)
		fmt.Println(string(jo))
		return nil
	    },
	},
	{
	    Name:    "user",
	    Usage:   "user info",
	    Aliases: []string{"u"},
	    Action:  func(c *cli.Context) error {
		outputJSON, e := json.Marshal(&u)
		if e != nil {
		    panic(e)
		}
		ioutil.WriteFile(dirUser, outputJSON, os.ModePerm)

		ju, _ := prettyjson.Marshal(u)
		fmt.Println(string(ju))
		return nil
	    },
	},
	{
	    Name:    "oauth",
	    Usage:   "get oauth-access-token",
	    Aliases: []string{"o"},
	    Action:  func(c *cli.Context) error {
		url := m.GetAuthorizationURL("secretstate", callurl,
		medium.ScopeBasicProfile, medium.ScopePublishPost, medium.ScopeListPublications)
		println(url)
		time.Sleep(1 * time.Second)
		open.Run(url)
		time.Sleep(1 * time.Second)
		fmt.Print("authorization code : ")
		code, e := terminal.ReadPassword(int(syscall.Stdin))
		if e != nil {
		    log.Fatal(e)
		}
		fmt.Print("\ninput code : ", string(code))
		at, e := m.ExchangeAuthorizationCode(string(code), callurl)
		if e != nil {
		    log.Fatal(e)
		}
		jat, _ := prettyjson.Marshal(at)

		fmt.Printf("\nYour token is %s\n", jat)
		ioutil.WriteFile(dirConf, jat, os.ModePerm)
		return nil
	    },
	},
	{
	    Name:    "self",
	    Usage:   "get self-access-token",
	    Aliases: []string{"s"},
	    Action:  func(c *cli.Context) error {
		fmt.Print("Self-Issued Access Tokens(https://medium.com/me/settings -> Integration tokens): ")
		password, e := terminal.ReadPassword(int(syscall.Stdin))
		if e != nil {
		    log.Fatal(e)
		} else {
		    fmt.Printf("\nYour token is %v\n", string(password))
		}
		d1 := make(map[string]interface{})
		d1["self_token"] = string(password)
		r, _ := prettyjson.Marshal(d1)
		ioutil.WriteFile(dirSelf, r, os.ModePerm)
		return nil
	    },
	},
	{
	    Name:    "refresh",
	    Usage:   "refresh access-token",
	    Aliases: []string{"r"},
	    Action:  func(c *cli.Context) error {
		rt, e := m.ExchangeRefreshToken(o.RefreshToken)
		if e != nil {
		    log.Fatal(e)
		}
		outputRF, e := json.Marshal(&rt)
		if e != nil {
		    panic(e)
		}
		jrt, _ := prettyjson.Marshal(rt)

		fmt.Printf("\nRefresh token done %s\n", jrt)
		ioutil.WriteFile(dirConf, outputRF, os.ModePerm)
		return nil
	    },
	},
	{
	    Name:    "publication",
	    Usage:   "get publications",
	    Aliases: []string{"pub"},
	    Action:  func(c *cli.Context) error {
		pu, e := m.GetPublications(u.ID)
		if e != nil {
		    log.Fatal(e)
		}
		outputPub, e := json.Marshal(&pu)
		if e != nil {
		    panic(e)
		}
		jpu, _ := prettyjson.Marshal(pu)

		fmt.Printf("%s", jpu)
		ioutil.WriteFile(dirPubl, outputPub, os.ModePerm)
		return nil
	    },
	},
	{
	    Name:    "article",
	    Usage:   "get post article",
	    Aliases: []string{"a"},
	    Action:  func(c *cli.Context) error {
		userbody,e := ioutil.ReadFile(dirUser)
		if e != nil {
		    fmt.Printf("File eor: %v\n", e)
		    os.Exit(1)
		}
		json.Unmarshal(userbody, &userjson)
		userurl := userjson.URL
		url := fmt.Sprintln(userurl, "/latest")
		url = strings.TrimSpace(url)
		req, _ := http.NewRequest("GET",url, nil)
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", m.AccessToken))
		client := new(http.Client)
		resp, e := client.Do(req)
		if e != nil {
		    panic(e)
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		body = body[16:] //"payload":
		fmt.Printf("%s", body)
		ioutil.WriteFile(dirArti, body, os.ModePerm)

		return nil
	    },
	},
    }
    app.Run(os.Args)
}

