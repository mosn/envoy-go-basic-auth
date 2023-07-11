package test

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

const ContainerName = "e2e_envoy_test"

type User struct {
	Username string
	Password string
}

func generateUsersConfig(users []User) string {
	config := `users:`
	for _, user := range users {
		config += fmt.Sprintf(`
                             - username: "%s"
                               password: "%s"`,
			user.Username,
			user.Password)
	}

	return config
}

func generateEnvoyConfig(users []User) {
	usersConfig := generateUsersConfig(users)
	config := fmt.Sprintf(`
static_resources:

  listeners:
    - name: listener_0
      address:
        socket_address:
          address: 0.0.0.0
          port_value: 10000
      filter_chains:
        - filters:
            - name: envoy.filters.network.http_connection_manager
              typed_config:
                "@type": type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager
                stat_prefix: ingress_http
                access_log:
                  - name: envoy.access_loggers.stdout
                    typed_config:
                      "@type": type.googleapis.com/envoy.extensions.access_loggers.stream.v3.StdoutAccessLog
                http_filters:
                  - name: envoy.filters.http.golang
                    typed_config:
                      "@type": type.googleapis.com/envoy.extensions.filters.http.golang.v3alpha.Config
                      library_id: example
                      library_path: /etc/envoy/libgolang.so
                      plugin_name: basic-auth
                      plugin_config:
                        "@type": type.googleapis.com/xds.type.v3.TypedStruct
                        value:
                            %s

                  - name: envoy.filters.http.router
                    typed_config:
                      "@type": type.googleapis.com/envoy.extensions.filters.http.router.v3.Router
                route_config:
                  name: local_route
                  virtual_hosts:
                    - name: local_service
                      domains: ["*"]
                      routes:
                        - match:
                            prefix: "/"
                          route:
                            host_rewrite_literal: mosn.io
                            cluster: service_mosn_io

  clusters:
    - name: service_mosn_io
      type: LOGICAL_DNS
      # Comment out the following line to test on v6 networks
      dns_lookup_family: V4_ONLY
      load_assignment:
        cluster_name: service_mosn_io
        endpoints:
          - lb_endpoints:
              - endpoint:
                  address:
                    socket_address:
                      address: mosn.io
                      port_value: 443
      transport_socket:
        name: envoy.transport_sockets.tls
        typed_config:
          "@type": type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.UpstreamTlsContext
          sni: mosn.io
`, usersConfig)
	// Write the configuration to the specified file
	err := os.WriteFile("envoy.yaml", []byte(config), 0644)
	if err != nil {
		panic(fmt.Sprintf("failed to write Envoy configuration to file: %v", err))
	}

}

func startEnvoyContainer(users []User) {
	generateEnvoyConfig(users)
	cmd := exec.Command("docker", "run", "--rm",
		"--name", ContainerName,
		"-v", fmt.Sprintf("%s/envoy.yaml:/etc/envoy/envoy.yaml", getPwd()),
		"-v", fmt.Sprintf("%s/../libgolang.so:/etc/envoy/libgolang.so", getPwd()),
		"-p", "10000:10000", "-e", "GODEBUG=cgocheck=0",
		"envoyproxy/envoy:contrib-v1.26-latest", "envoy", "-c", "/etc/envoy/envoy.yaml")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func getPwd() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func stopEnvoyContainer() {
	cmd := exec.Command("docker", "stop", ContainerName)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
