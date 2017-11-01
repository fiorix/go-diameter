package test

import (
	"fmt"
	"math/rand"
	"net"
	"testing"

	"github.com/fiorix/go-diameter/examples/s6a_proxy/protos"
	"github.com/fiorix/go-diameter/examples/s6a_proxy/service"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"time"
)

const TEST_LOOPS = 33

// TestS6aProxyService creates a mock S6a Diameter server, S6a S6a Proxy service
// and runs tests using GRPC client: GRPC Client <--> GRPC Server <--> S6a SCTP Diameter Server
func TestS6aProxyService(t *testing.T) {

	diamAddr := fmt.Sprintf("127.0.0.1:%d", 30000+rand.Intn(1000))
	err := StartTestS6aServer("sctp", diamAddr)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("Started S6a Server")

	lis, err := net.Listen("tcp", "")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
		return
	}
	cfg := &service.S6aProxyConfig{
		HssAddr:  diamAddr,                  // use "192.168.60.145:3868" to send diam messages to OAI HSS VM
		Protocol: "sctp",                    // tcp/sctp
		Host:     "magma-oai.openair4G.eur", // diameter host
		Realm:    "openair4G.eur",           // diameter realm
	}
	s := grpc.NewServer()
	service, err := service.NewS6aProxy(cfg)
	if err != nil {
		t.Fatalf("failed to create S6aProxy: %v", err)
		return

	}
	protos.RegisterS6AProxyServer(s, service)

	go func() {
		if err := s.Serve(lis); err != nil {
			t.Fatalf("failed to serve: %v", err)
			return
		}
	}()
	addr := lis.Addr()
	t.Logf("Started S6a GRPC Proxy on %s", addr.String())

	// Set up a connection to the server.
	conn, err := grpc.Dial(addr.String(), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("GRPC connect error: %v", err)
		return
	}
	defer conn.Close()

	c := protos.NewS6AProxyClient(conn)
	req := &protos.AuthenticationInformationRequest{
		UserName:                   TEST_IMSI,
		VisitedPlmn:                []byte(TEST_PLMN_ID),
		NumRequestedEutranVectors:  3,
		ImmediateResponsePreferred: true,
	}
	complChan := make(chan error, TEST_LOOPS+1)
	testLoopF := func() {
		// AIR
		r, err := c.AuthenticationInformation(context.Background(), req)
		if err != nil {
			t.Fatalf("GRPC AIR Error: %v", err)
			complChan <- err
			return
		}
		t.Logf("GRPC AIA: %#+v", *r)
		if r.ErrorCode != protos.ErrorCode_UNDEFINED {
			t.Errorf("Unexpected AIA Error Code: %d", r.ErrorCode)
		}
		if len(r.EutranVectors) != 3 {
			t.Errorf("Unexpected Number of EutranVectors: %d, Expected: 3", len(r.EutranVectors))
		}
		ulReq := &protos.UpdateLocationRequest{
			UserName:           TEST_IMSI,
			VisitedPlmn:        []byte(TEST_PLMN_ID),
			SkipSubscriberData: false,
			InitialAttach:      true,
		}
		// ULR
		ulResp, err := c.UpdateLocation(context.Background(), ulReq)
		if err != nil {
			t.Fatalf("GRPC ULR Error: %v", err)
			complChan <- err
			return
		}
		t.Logf("GRPC ULA: %#+v", *ulResp)
		if r.ErrorCode != protos.ErrorCode_UNDEFINED {
			t.Errorf("Unexpected AIA Error Code: %d", r.ErrorCode)
		}
		complChan <- nil
	}
	go testLoopF()

	select {
	case testErr := <-complChan:
		if testErr != nil {
			t.Fatal(err)
			return
		}
	case <-time.After(time.Second):
		t.Fatal("TestS6aProxyService Timed out")
	}

	// return

	for round := 0; round < TEST_LOOPS; round++ {
		go testLoopF()
	}
	for round := 0; round < TEST_LOOPS; round++ {
		testErr := <-complChan
		if testErr != nil {
			t.Fatal(err)
			return
		}
	}
}
