package command

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/coreos/etcd/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/coreos/etcd/client"
)

func NewMemberCommand() cli.Command {
	return cli.Command{
		Name: "member",
		Subcommands: []cli.Command{
			cli.Command{
				Name:   "list",
				Usage:  "enumerate existing cluster members",
				Action: actionMemberList,
			},
			cli.Command{
				Name:   "add",
				Usage:  "add a new member to the etcd cluster",
				Action: actionMemberAdd,
			},
			cli.Command{
				Name:   "remove",
				Usage:  "remove an existing member from the etcd cluster",
				Action: actionMemberRemove,
			},
		},
	}
}

func mustNewMembersAPI(c *cli.Context) client.MembersAPI {
	peers := getPeersFlagValue(c)
	for i, p := range peers {
		if !strings.HasPrefix(p, "http") && !strings.HasPrefix(p, "https") {
			peers[i] = fmt.Sprintf("http://%s", p)
		}
	}

	hc, err := client.NewHTTPClient(&http.Transport{}, peers)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	return client.NewMembersAPI(hc, client.DefaultRequestTimeout)
}

func actionMemberList(c *cli.Context) {
	if len(c.Args()) != 0 {
		fmt.Fprintln(os.Stderr, "No arguments accepted")
		os.Exit(1)
	}
	mAPI := mustNewMembersAPI(c)
	members, err := mAPI.List()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	for _, m := range members {
		fmt.Printf("%s: name=%s peerURLs=%s clientURLs=%s\n", m.ID, m.Name, strings.Join(m.PeerURLs, ","), strings.Join(m.ClientURLs, ","))
	}
}

func actionMemberAdd(c *cli.Context) {
	args := c.Args()
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Provide a single member peerURL")
		os.Exit(1)
	}

	mAPI := mustNewMembersAPI(c)
	url := args[0]
	m, err := mAPI.Add(url)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	fmt.Printf("Added member to cluster with ID %s", m.ID)
}

func actionMemberRemove(c *cli.Context) {
	args := c.Args()
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Provide a single member ID")
		os.Exit(1)
	}

	mAPI := mustNewMembersAPI(c)
	mID := args[0]
	if err := mAPI.Remove(mID); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	fmt.Printf("Removed member %s from cluster\n", mID)
}
