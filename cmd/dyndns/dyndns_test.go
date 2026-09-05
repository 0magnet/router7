package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/0magnet/router7/internal/dyndns"
	"github.com/libdns/libdns"
)

func TestLogic(t *testing.T) {
	cfg := DynDNSRecord{
		Zone:             "zekjur.net",
		RecordName:       "dyndns.zekjur.net",
		RecordTTLSeconds: 300, // 5 minutes
	}
	update = func(ctx context.Context, zone string, record libdns.Record, _ dyndns.RecordGetterSetter) error {
		if got, want := zone, cfg.Zone; got != want {
			return fmt.Errorf("update(): unexpected zone: got %q, want %q", got, want)
		}
		if got, want := record.RR().Name, cfg.RecordName; got != want {
			return fmt.Errorf("update(): unexpected record name: got %q, want %q", got, want)
		}
		return nil
	}
	if err := logic("lo", []DynDNSRecord{cfg}); err != nil {
		t.Fatalf("logic: %v", err)
	}

}
