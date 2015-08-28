package proxy

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"

	"github.com/fsouza/go-dockerclient"
	//"github.com/k0kubun/pp"
	"github.com/mix3/illusion/config"
)

var ErrNotFound = errors.New("subdomain not found")

type Proxy struct {
	conf      config.Config
	matcher   *regexp.Regexp
	dockerCli *docker.Client
}

func NewProxy(conf config.Config) *Proxy {
	client, err := docker.NewClient(conf.DockerEndpoint)
	if err != nil {
		log.Fatal(err)
	}

	return &Proxy{
		conf:      conf,
		matcher:   regexp.MustCompile(fmt.Sprintf(`^(.+?)\.%s$`, conf.Domain)),
		dockerCli: client,
	}
}

func split(host string) (string, string) {
	ret := strings.Split(host, ":")
	if len(ret) <= 1 {
		return ret[0], ""
	}
	return ret[0], ret[1]
}

func (p *Proxy) parseHost(host string) string {
	h, _ := split(host)
	match := p.matcher.FindStringSubmatch(h)
	if 1 < len(match) {
		return match[1]
	}
	return ""
}

func (p *Proxy) searchContainer(subdomain string) (string, error) {
	for _, v := range p.conf.IgnoreSubdomain {
		if v == subdomain {
			return "", ErrNotFound
		}
	}

	containers, err := p.dockerCli.ListContainers(docker.ListContainersOptions{})
	if err != nil {
		return "", err
	}

	var container *docker.Container
	for _, v := range containers {
		if 0 < len(v.Names) && v.Names[0] == "/"+subdomain {
			container, err = p.dockerCli.InspectContainer(v.ID)
			if err != nil {
				return "", err
			}
			break
		}
	}
	if container == nil {
		return "", ErrNotFound
	}

	ip := container.NetworkSettings.IPAddress

	return ip, nil
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	subdomain := p.parseHost(r.Host)

	ip, err := p.searchContainer(subdomain)

	if err != nil {
		switch err {
		case ErrNotFound:
			http.NotFound(w, r)
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	u, err := url.Parse(fmt.Sprintf("http://%s:%d", ip, p.conf.ForwardPort))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	httputil.NewSingleHostReverseProxy(u).ServeHTTP(w, r)
}
