package kong

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/franela/goreq"
	"github.com/hudl/fargo"
)

type Client struct {
	config Config
}

type Config struct {
	Host string
}

type UpstreamObject struct {
	Name string `json:"name"`
}

type UpstreamResObject struct {
	UpstreamObject
	Id string `json:"id"`
}

type UpstreamListResObject struct {
	Data []UpstreamResObject `json:"data"`
}

type TargetObject struct {
	Target string `json:"target"`
	Weight int    `json:"weight"`
}

type TargetResObject struct {
	TargetObject
	Id         string `json:"id"`
	UpstreamId string `json:"upstream_id"`
}

type TargetListResObject struct {
	Data []TargetResObject `json:"data"`
}

type ResObject struct {
	StatusCode int
	Message    string
}

func NewClient(c Config) (client Client) {
	return Client{
		config: c,
	}
}

func (c *Client) request(path string, method string, body interface{}) (resp *goreq.Response) {

	res, err := goreq.Request{
		Uri:         c.config.Host + "/upstreams" + path,
		Method:      method,
		Body:        body,
		Accept:      "application/json",
		ContentType: "application/json",
	}.Do()

	if err != nil {
		return nil
	}

	if res.StatusCode >= 400 {
		return nil
	}

	return res
}

func formatName(app *fargo.Application) string {
	reg := regexp.MustCompile("[:\\.]")
	return reg.ReplaceAllString(app.Name, "${1}-") + ".eureka.internal"
}

func (c *Client) FetchUpstream(name string) *UpstreamResObject {
	res := c.request("/"+name, http.MethodGet, nil)

	if res == nil {
		return nil
	}

	var upstream UpstreamResObject
	res.Body.FromJsonTo(&upstream)

	return &upstream
}

func (c *Client) RegisterUpstream(app *fargo.Application) {
	upstreamName := formatName(app)

	upstream := c.FetchUpstream(upstreamName)

	if upstream == nil {
		createRes := c.request("", http.MethodPost, UpstreamObject{
			Name: upstreamName,
		})

		if createRes == nil {
			return
		}

		var upstreamRes UpstreamResObject
		createRes.Body.FromJsonTo(&upstreamRes)
		c.RegisterTargets(upstreamRes.Id, app.Instances)
	} else {
		c.RegisterTargets(upstream.Id, app.Instances)
	}

}

func (c *Client) FetchTargetsOfUpstreams(upstreamId string) (targetList []TargetResObject) {
	req := c.request("/"+upstreamId+"/targets", http.MethodGet, nil)

	if req == nil {
		return nil
	}

	var targets TargetListResObject

	req.Body.FromJsonTo(&targets)

	return targets.Data
}

func (c *Client) RegisterTargets(upstreamId string, instances []*fargo.Instance) {
	targets := c.FetchTargetsOfUpstreams(upstreamId)

	for _, ins := range instances {
		targetUrl := ins.IPAddr + ":" + fmt.Sprintf("%d", ins.Port)

		if targets != nil {
			for i, target := range targets {
				if target.Target == targetUrl {
					targets = append(targets[:i], targets[i+1:]...)
					continue
				}
			}
		}

		c.request("/"+upstreamId+"/targets", http.MethodPost, TargetObject{
			Target: targetUrl,
			Weight: 100,
		})

	}

	if targets != nil && len(targets) > 0 {
		for _, t := range targets {
			c.request("/"+upstreamId+"/targets/"+t.Id, http.MethodDelete, nil)
		}
	}

}
