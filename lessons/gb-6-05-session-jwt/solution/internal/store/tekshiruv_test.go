// VAZIFA tekshiruvi (TASK.md, 1-qadam va 🟡): RevokeToken, TokenRevoked, PurgeRevoked.
package store

import (
	"context"
	"testing"
	"time"
)

func TestRevokeToken(t *testing.T) {
	ctx := context.Background()
	s := New()
	if err := s.RevokeToken(ctx, "jti-1", time.Now().Add(time.Minute)); err != nil {
		t.Fatalf("RevokeToken: %v — TASK.md, 1-qadam: s.revoked'ga yozing va nil qaytaring", err)
	}
	if ok, err := s.TokenRevoked(ctx, "jti-1"); !ok || err != nil {
		t.Errorf("TokenRevoked(jti-1) = %v, %v; kutilgan true, nil", ok, err)
	}
	if ok, _ := s.TokenRevoked(ctx, "jti-2"); ok {
		t.Error("TokenRevoked(jti-2) = true: bu jti bekor qilinmagan edi")
	}
}

// 🟡 ixtiyoriy: muddati oʻtgan yozuvlarni tozalash.
func TestPurgeRevoked(t *testing.T) {
	ctx := context.Background()
	s := New()
	now := time.Now()
	n, err := s.PurgeRevoked(ctx, now)
	if err != nil || n != 0 {
		t.Fatalf("boʻsh denylist: PurgeRevoked = %d, %v; kutilgan 0, nil", n, err)
	}
	s.RevokeToken(ctx, "eski", now.Add(-time.Minute))
	s.RevokeToken(ctx, "yangi", now.Add(time.Minute))
	if n, err := s.PurgeRevoked(ctx, now); n != 1 || err != nil {
		t.Errorf("PurgeRevoked = %d, %v; kutilgan 1, nil (faqat exp'i oʻtgan «eski»)", n, err)
	}
	if ok, _ := s.TokenRevoked(ctx, "yangi"); !ok {
		t.Error("«yangi» denylist'dan yoʻqoldi: exp'i hali oʻtmagan yozuv qolishi kerak")
	}
}
