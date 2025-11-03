package cron_test

import (
	"context"
	"testing"
	"washwise/cron"
)

func TestGetMachineTypes(t *testing.T) {
	ctx := context.Background()
	shopId := "202401041041470000069996565184"

	resp, err := cron.GetMachineTypes(ctx, shopId)
	if err != nil {
		t.Fatalf("GetMachineTypes failed: %v", err)
	}

	if resp == nil {
		t.Fatal("Expected non-nil response")
	}

	t.Logf("GetMachineTypes success, got %d machine types", len(resp.Items))
}

func TestGetMachines(t *testing.T) {
	ctx := context.Background()
	shopId := "202401041041470000069996565184"
	machineTypeId := "c9892cb4-bd78-40f6-83c2-ba73383b090a"
	pageSize := 100
	page := 1

	resp, err := cron.GetMachines(ctx, shopId, machineTypeId, pageSize, page)
	if err != nil {
		t.Fatalf("GetMachines failed: %v", err)
	}

	if resp == nil {
		t.Fatal("Expected non-nil response")
	}

	t.Logf("GetMachines success, got %d machines", len(resp.Items))
}

func TestGetMachineDetail(t *testing.T) {
	ctx := context.Background()
	goodsId := int64(1100547706)

	resp, err := cron.GetMachineDetail(ctx, goodsId)
	if err != nil {
		t.Fatalf("GetMachineDetail failed: %v", err)
	}

	if resp == nil {
		t.Fatal("Expected non-nil response")
	}

	t.Logf("GetMachineDetail success, machine name: %s", resp.Name)
}
