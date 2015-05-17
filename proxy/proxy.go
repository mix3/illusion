package proxy

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/fsouza/go-dockerclient"
	"github.com/mix3/illusion/config"
)

type Proxy struct {
	conf      config.Config
	matcher   *regexp.Regexp
	dockerCli *docker.Client
}

func NewProxy(conf config.Config) *Proxy {
	var client *docker.Client
	var err error
	if path := os.Getenv("DOCKER_CERT_PATH"); path == "" {
		client, err = docker.NewClient(conf.DockerEndpoint)
	} else {
		// for boot2docker
		ca := fmt.Sprintf("%s/ca.pem", path)
		cert := fmt.Sprintf("%s/cert.pem", path)
		key := fmt.Sprintf("%s/key.pem", path)
		client, err = docker.NewTLSClient(conf.DockerEndpoint, cert, key, ca)
	}
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
		return "", "", fmt.Errorf("container not found, name(%s)", subdomain)
	}

	var ports []docker.PortBinding
	if len(container.NetworkSettings.Ports) != 1 {
		return "", "", fmt.Errorf("PortBinding is invalid, name(%s)", subdomain)
	}
	for _, ports = range container.NetworkSettings.Ports {
		if len(ports) != 1 {
			return "", "", fmt.Errorf("PortBinding is invalid, name(%s)", subdomain)
		}
	}

	ip := container.NetworkSettings.IPAddress
	if path := os.Getenv("DOCKER_HOST"); path != "" {
		// for boot2docker
		u, err := url.Parse(path)
		if err != nil {
			return "", "", fmt.Errorf("DOCKER_HOST(%s) is invalid", path)
		}
		ip, _ = split(u.Host)
	}

	return ip, ports[0].HostPort, nil
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	subdomain := p.parseHost(r.Host)
	ip, port, err := p.searchContainer(subdomain)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u, err := url.Parse(fmt.Sprintf("http://%s:%s", ip, port))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	httputil.NewSingleHostReverseProxy(u).ServeHTTP(w, r)
}
