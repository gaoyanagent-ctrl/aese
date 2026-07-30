package main

import "testing"

func TestLocalDevelopmentAuthenticationRequiresLoopbackListen(t *testing.T) {
	for _, address := range []string{"127.0.0.1:8090", "localhost:8090", "[::1]:8090"} {
		if !isLoopbackListen(address) {
			t.Fatalf("expected loopback address: %s", address)
		}
	}
	for _, address := range []string{":8090", "0.0.0.0:8090", "192.168.50.222:8090"} {
		if isLoopbackListen(address) {
			t.Fatalf("unexpected public local-dev address: %s", address)
		}
	}
}
