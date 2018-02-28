package kong

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/franela/goreq"
	"github.com/hudl/fargo"
	"github.com/op/go-logging"
)

type Client struct {
	config Config
	logger logging.Logger
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
		logger: logging.Logger{
			Module: "kongClient",
		},
	}
}

func (c *Client) request(path string, method string, body interface{}, showDebug bool) (resp *goreq.Response) {

	res, err := goreq.Request{
		Uri:         c.config.Host + "/upstreams" + path,
		Method:      method,
		Body:        body,
		Accept:      "application/json",
		ContentType: "application/json",
	}.Do()

	if err != nil {
	}

	if err != nil {
		c.logger.Error(err)
		return nil
	}

	if res.StatusCode >= 400 {
		c.logger.Errorf("request to %s failed on status: %n", path, res.StatusCode)
		return nil
	}

	return res
}

func formatName(app *fargo.Application) string {
	reg := regexp.MustCompile("[:\\.]")
	return reg.ReplaceAllString(app.Name, "${1}-") + ".eureka.internal"
}

func (c *Client) FetchUpstream(name string) *UpstreamResObject {
	res := c.request("/"+name, http.MethodGet, nil, true)

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
		c.logger.Infof("upstream %s does not exist, now creating", upstreamName)

		createRes := c.request("", http.MethodPost, UpstreamObject{
			Name: upstreamName,
		}, false)

		if createRes == nil {
			return
		}

		var upstreamRes UpstreamResObject
		createRes.Body.FromJsonTo(&upstreamRes)
		c.RegisterTargets(upstreamRes.Name, app.Instances)
	} else {
		c.RegisterTargets(upstream.Name, app.Instances)
	}

}

func (c *Client) FetchTargetsOfUpstreams(upstreamName string) (targetList []TargetResObject) {
	req := c.request("/"+upstreamName+"/targets/active", http.MethodGet, nil, false)

	if req == nil {
		return nil
	}

	var targets TargetListResObject

	req.Body.FromJsonTo(&targets)

	return targets.Data
}

func (c *Client) RegisterTargets(upstreamName string, instances []*fargo.Instance) {
	targets := c.FetchTargetsOfUpstreams(upstreamName)

findInstance:
	for _, ins := range instances {
		targetUrl := ins.IPAddr + ":" + fmt.Sprintf("%d", ins.Port)

		if targets != nil {
			for i, target := range targets {
				if target.Target == targetUrl {
					targets = append(targets[:i], targets[i+1:]...)
					continue findInstance
				}
			}
		}

		c.logger.Infof("create new target %s for upstream %s", targetUrl, upstreamName)

		c.request("/"+upstreamName+"/targets", http.MethodPost, TargetObject{
			Target: targetUrl,
			Weight: 100,
		}, false)

	}

	if targets != nil && len(targets) > 0 {
		for _, t := range targets {
			c.logger.Infof("delete unhealthy target %s on upstream %s", t.Target, upstreamName)
			c.request("/"+upstreamName+"/targets/"+t.Id, http.MethodDelete, nil, false)
		}
	}

}
