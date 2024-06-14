package initcloud

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	chipperpb "github.com/digital-dream-labs/api/go/chipperpb"
	grpcserver "github.com/digital-dream-labs/hugh/grpc/server"
	chipperserver "github.com/kercre123/0.11-cloud/pkg/servers/chipper"
	"github.com/kercre123/0.11-cloud/pkg/stt"
	"github.com/kercre123/0.11-cloud/pkg/vars"
	wp "github.com/kercre123/0.11-cloud/pkg/voiceprocess"
	"github.com/kercre123/0.11-cloud/pkg/web"
)

func InitCloud() {
	chipperPort := os.Getenv(vars.ChipperPortEnv)
	if chipperPort == "" {
		fmt.Println("you must provide a chipper port (fulfill all env vars)")
		os.Exit(1)
	}
	vars.Init()
	err := stt.Init()
	if err != nil {
		fmt.Println("error initing STT: " + err.Error())
	}
	cert, err := tls.X509KeyPair(vars.TLSCert, vars.TLSKey)
	if err != nil {
		fmt.Println("error loading certs: " + err.Error())
		os.Exit(1)
	}
	tlsListener, err := tls.Listen("tcp", ":"+chipperPort, &tls.Config{
		Certificates: []tls.Certificate{cert},
		CipherSuites: nil,
	})
	serv, err := wp.New()
	if err != nil {
		fmt.Println("error starting wire-pod server: " + err.Error())
		os.Exit(1)
	}

	web.AddSecretsAPI()
	web.AddSecretsWebroot()

	go http.ListenAndServeTLS(":"+os.Getenv(vars.WebPortEnv), os.Getenv(vars.CertFileEnv), os.Getenv(vars.KeyFileEnv), nil)

	fmt.Println("serving chipper at port " + chipperPort)
	grpcServe(tlsListener, serv)
}

func grpcServe(l net.Listener, p *wp.Server) error {
	srv, err := grpcserver.New(
		grpcserver.WithViper(),
		grpcserver.WithReflectionService(),
		grpcserver.WithInsecureSkipVerify(),
	)
	if err != nil {
		log.Fatal(err)
	}

	s, _ := chipperserver.New(
		chipperserver.WithIntentProcessor(p),
		chipperserver.WithKnowledgeGraphProcessor(p),
	)

	chipperpb.RegisterChipperGrpcServer(srv.Transport(), s)

	return srv.Transport().Serve(l)
}
