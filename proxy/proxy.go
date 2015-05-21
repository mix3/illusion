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

	//containers, err := client.ListContainers(docker.ListContainersOptions{})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//spew.Dump(containers)

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

func (p *Proxy) searchContainer(subdomain string) (string, string, error) {
	containers, err := p.dockerCli.ListContainers(docker.ListContainersOptions{})
	if err != nil {
		return "", "", err
	}

	var container *docker.Container
	for _, v := range containers {
		if 0 < len(v.Names) && v.Names[0] == "/"+subdomain {
			container, err = p.dockerCli.InspectContainer(v.ID)
			if err != nil {
				return "", "", err
			}
			break
		}
	}
	if container == nil {
		return "", "", ErrNotFound
	}

	if len(container.NetworkSettings.Ports) != 1 {
		return "", "", ErrNotFound
	}

	var port string
	for v, _ := range container.NetworkSettings.Ports {
		if v.Proto() != "tcp" {
			return "", "", ErrNotFound
		}
		port = v.Port()
	}

	ip := container.NetworkSettings.IPAddress

	return ip, port, nil
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	subdomain := p.parseHost(r.Host)

	ip, port, err := p.searchContainer(subdomain)

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

	u, err := url.Parse(fmt.Sprintf("http://%s:%s", ip, port))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	httputil.NewSingleHostReverseProxy(u).ServeHTTP(w, r)
}
