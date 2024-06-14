package initcloud

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	chipperpb "github.com/digital-dream-labs/api/go/chipperpb"
	grpcserver "github.com/digital-dream-labs/hugh/grpc/server"
	chipperserver "github.com/kercre123/0.11-cloud/pkg/servers/chipper"
	"github.com/kercre123/0.11-cloud/pkg/stt"
	"github.com/kercre123/0.11-cloud/pkg/vars"
	wp "github.com/kercre123/0.11-cloud/pkg/voiceprocess"
	"github.com/kercre123/0.11-cloud/pkg/web"
)

var Listener net.Listener
var GRPCServer *grpcserver.Server

func MkHTTPS(port string) *http.Server {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: http.DefaultServeMux,
		TLSConfig: &tls.Config{
			GetCertificate: func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
				cert, err := tls.X509KeyPair(vars.TLSCert, vars.TLSKey)
				if err != nil {
					return nil, err
				}
				return &cert, nil
			},
		},
	}

	return srv
}

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
	Listener, err = tls.Listen("tcp", ":"+chipperPort, &tls.Config{
		Certificates: []tls.Certificate{cert},
		CipherSuites: nil,
	})
	serv, err := wp.New()
	if err != nil {
		fmt.Println("error starting wire-pod server: " + err.Error())
		os.Exit(1)
	}

	httpServ := MkHTTPS(os.Getenv(vars.WebPortEnv))

	web.AddSecretsAPI()
	web.AddSecretsWebroot()
	go httpServ.ListenAndServeTLS("", "")

	fmt.Println("serving chipper at port " + chipperPort)
	grpcServe(Listener, serv)
	for {
		time.Sleep(time.Hour * vars.PollHrs)
		certBytes, err := os.ReadFile(os.Getenv(vars.CertFileEnv))
		if err != nil {
			fmt.Println("error polling cert")
			os.Exit(1)
		}
		if string(certBytes) != string(vars.TLSCert) {
			fmt.Println("cert is different, restarting servers")
			vars.TLSCert = certBytes
			keyBytes, _ := os.ReadFile(os.Getenv(vars.KeyFileEnv))
			vars.TLSKey = keyBytes
			GRPCServer.Stop()
			cert, _ = tls.X509KeyPair(vars.TLSCert, vars.TLSKey)
			Listener, err = tls.Listen("tcp", ":"+chipperPort, &tls.Config{
				Certificates: []tls.Certificate{cert},
				CipherSuites: nil,
			})
			httpServ.Shutdown(context.TODO())
			httpServ = MkHTTPS(os.Getenv(vars.WebPortEnv))
			go httpServ.ListenAndServeTLS("", "")
			go grpcServe(Listener, serv)
		}
	}
}

func grpcServe(l net.Listener, p *wp.Server) error {
	var err error
	GRPCServer, err = grpcserver.New(
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

	chipperpb.RegisterChipperGrpcServer(GRPCServer.Transport(), s)

	return GRPCServer.Transport().Serve(l)
}
