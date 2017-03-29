package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"syscall"
	"io/ioutil"
	"encoding/json"
	"path/filepath"
	"golang.org/x/crypto/ssh/terminal"
	"github.com/urfave/cli"
	"github.com/medium/medium-sdk-go"
	"github.com/skratchdot/open-golang/open"
	//cregex "github.com/mingrammer/commonregex"
)

// Oauth medium 
type Oauth struct {
	Cid	string `json:"client_id"`
	Secret  string `json:"client_secret"`
	Self	string `json:"self_token"`
	Token	string `json:"access_token"`
}

// PostConfig date
type PostConfig struct {
	Title		string		`json:"title"`
	Tags		[]string	`json:"tags"`
	Content		string		`json:"content"`
	CanonicalURL	string		`json:"url"`
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

	var o Oauth
	var b PostConfig
	callurl := "https://syui.github.io/medigo/callback/medium"

	app := cli.NewApp()
	app.Version = "0.03"
	cli.HelpFlag = cli.BoolFlag{
		Name: "help",
		Usage: "subcommands : $ medium-go p d # post draft",

	}
	dir := filepath.Join(os.Getenv("HOME"), ".config", "medium-go")
	dirPost := filepath.Join(dir, "posts")
	dirFile := filepath.Join(dir, "body.json")
	dirConf := filepath.Join(dir, ".medium.json")

	if err := os.MkdirAll(dirPost, os.ModePerm); err != nil {
		panic(err)
	}

	_, err := os.Stat(dirConf)
	if err != nil {
		fmt.Print("Self-Issued Access Tokens(https://medium.com/me/settings -> Integration tokens): ")
		password, e := terminal.ReadPassword(int(syscall.Stdin))
		if e != nil {
			log.Fatal(e)
		} else {
			fmt.Printf("\nYour token is %v\n", string(password))
		}
		d1 := make(map[string]interface{})
		d1["self_token"] = string(password)
		r, _ := json.Marshal(d1);
		ioutil.WriteFile(dirConf, r, os.ModePerm)
	}
	_, err = os.Stat(dirFile)
	if err != nil {
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

	m := medium.NewClient(cid, secret)

	//if len(o.Token) != 0 {
	//	fmt.Printf("your token is access %s\n")
	//}
	token := o.Self

	m2 := medium.NewClientWithAccessToken(token)
	m.AccessToken = m2.AccessToken

	u, e := m2.GetUser()
	if e != nil {
		fmt.Printf("rm %s\n", dirConf)
		os.Remove(dirConf)
		log.Fatal(e)
	}

	//log.Println(at, u, p)

	app.Commands = []cli.Command{
		{
			Name:    "post",
			Aliases: []string{"p"},
			Usage:   "create post",
			Action:  func(c *cli.Context) error {

				fileinfos, _ :=ioutil.ReadDir(dirPost)
				for _,fileinfo := range fileinfos {
					fmt.Println(fileinfo.Name())
				}

				return nil
			},
			Subcommands: cli.Commands{
				cli.Command{
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
						fmt.Println(p)
						if e != nil {
							log.Fatal(e)
						}
						return nil
					},
				},
				cli.Command{
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
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "\n\tsub : id(i), secret-id(s), self-token(t)",
			Action:  func(c *cli.Context) error {
				fmt.Println(o.Cid, o.Secret, o.Token)
				return nil
			},
			Subcommands: cli.Commands{
				cli.Command{
					Name:   "id",
					Usage:   "id",
					Aliases: []string{"i"},
					Action:  func(c *cli.Context) error {
						fmt.Println(o.Cid)
						return nil
					},
				},
				cli.Command{
					Name:   "secret",
					Usage:   "client_secret",
					Aliases: []string{"s"},
					Action:  func(c *cli.Context) error {
						fmt.Println(o.Secret)
						return nil
					},
				},
				cli.Command{
					Name:   "token",
					Usage:   "self_token",
					Aliases: []string{"t"},
					Action:  func(c *cli.Context) error {
						fmt.Println(o.Token)
						return nil
					},
				},
			},
		},
		{
			Name:    "user",
			Usage:   "user",
			Aliases: []string{"u"},
			Action:  func(c *cli.Context) error {
				fmt.Println(u)
				return nil
			},
		},
		{
			Name:    "oauth",
			Usage:   "oauth",
			Aliases: []string{"o"},
			Action:  func(c *cli.Context) error {
				// Build the URL where you can send the user to obtain an authorization code.
				url := m.GetAuthorizationURL("secretstate", callurl,
				medium.ScopeBasicProfile, medium.ScopePublishPost)
				println(url)
				time.Sleep(1 * time.Second)
				open.Run(url)
				time.Sleep(1 * time.Second)
				//log.Printf("Authentication URL: %s\n", url)
				//https://example.com/callback/medium?state=secretstate&code=XXXXXX
				fmt.Print("authorization code : ")
				code, e := terminal.ReadPassword(int(syscall.Stdin))
				if e != nil {
					log.Fatal(e)
				}
				at, err := m.ExchangeAuthorizationCode(string(code), callurl)
				if err != nil {
					log.Fatal(err)
				}
				outputJSON, err := json.Marshal(&at)
				if err != nil {
					panic(err)
				}
				fmt.Printf("\nYour token is %s\n", outputJSON)
				ioutil.WriteFile(dirConf, outputJSON, os.ModePerm)
				return nil
			},
		},
	}
	app.Run(os.Args)
}

