package password

// Vazifa tekshiruvi (TASK.md). Bu faylni oʻzgartirmang.

import (
	"errors"
	"testing"
)

// oldHash — videodagi parametrlar (m=19456, t=2, p=1) bilan yasalgan satr; paroli oldPassword.
// Bunday satrlar bazada allaqachon bor: parametrlar oʻzgargandan keyin ham ular tekshirilishi kerak.
const (
	oldPassword = "qovun-va-non-ikki-ming-yigirma"
	oldHash     = "$argon2id$v=19$m=19456,t=2,p=1$ZXeGATxNrjQQShHdsFh+jg$F7nIhe/iagKNWCc+GWr9n4XXOqZGhcOUwO23dJjzE2c"
)

func TestVazifaParams(t *testing.T) {
	p, err := decode(Hash(oldPassword))
	if err != nil {
		t.Fatalf("decode(Hash(…)): %v", err)
	}
	if p.memory != 47104 || p.passes != 1 || p.threads != 1 {
		t.Errorf("Hash hozir m=%d,t=%d,p=%d bilan yasayapti, kutilgan m=47104,t=1,p=1 (OWASP'ning birinchi qatori, 46 MiB). "+
			"password.go'dagi memory va passes konstantalarini oʻzgartiring (TASK.md, 1-qadam)", p.memory, p.passes, p.threads)
	}
}

func TestVazifaOldStillVerifies(t *testing.T) {
	ok, err := Verify(oldPassword, oldHash)
	if err != nil || !ok {
		t.Errorf("eski satr (m=19456,t=2) bilan Verify = %v, %v; kutilgan true, nil. "+
			"Verify parametrlarni konstantalardan emas, satrdan olishi kerak", ok, err)
	}
}

func TestVazifaNeedsRehash(t *testing.T) {
	if old, err := NeedsRehash(oldHash); err != nil || !old {
		t.Errorf("NeedsRehash(eski satr m=19456,t=2) = %v, %v; kutilgan true, nil. "+
			"decode bilan oʻqing va parametrlarni konstantalar bilan solishtiring (TASK.md, 2-qadam)", old, err)
	}
	if fresh, err := NeedsRehash(Hash(oldPassword)); err != nil || fresh {
		t.Errorf("NeedsRehash(yangi Hash) = %v, %v; kutilgan false, nil: yangi satrni qayta hash qilish shart emas", fresh, err)
	}
	if _, err := NeedsRehash("385426486641a520e865c4e6b24d379c8f8afa1eb5de2f932ba90f0622f06289"); !errors.Is(err, ErrBadHash) {
		t.Errorf("NeedsRehash(sha256 hex) xatosi %v; kutilgan ErrBadHash: buzuq satrda decode'ning xatosini qaytaring", err)
	}
}
