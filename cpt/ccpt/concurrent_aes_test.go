package ccpt

import (
	"sync"
	"testing"

	"github.com/askasoft/pango/str"
)

func TestAesCBCEncrypt(t *testing.T) {
	bits := []int{128, 192, 256}
	for _, bit := range bits {
		c := NewAesCBCCryptor("1234567890abcde", bit)

		var wg sync.WaitGroup

		tf := func(n int) {
			defer wg.Done()

			rs := str.RandString(n)
			es, err := c.EncryptString(rs)
			if err != nil {
				t.Error(n, err)
				return
			}

			// fmt.Printf("[%d] len(es) = %d\n", i, len(es))

			ds, err := c.DecryptString(es)
			if err != nil {
				t.Error(n, err)
				return
			}

			if ds != rs {
				t.Errorf("[%d] want %q, but %q", n, rs, ds)
				return
			}
		}

		for i := 64; i <= 128; i++ {
			wg.Add(1)
			go tf(i)
		}
		wg.Wait()
	}
}

func TestAesCFBEncrypt(t *testing.T) {
	bits := []int{128, 192, 256}
	for _, bit := range bits {
		c := NewAesCFBCryptor("1234567890abedefg", bit)

		var wg sync.WaitGroup

		tf := func(n int) {
			defer wg.Done()

			rs := str.RandString(n)
			es, err := c.EncryptString(rs)
			if err != nil {
				t.Error(n, err)
				return
			}

			// fmt.Printf("[%d] len(es) = %d\n", i, len(es))

			ds, err := c.DecryptString(es)
			if err != nil {
				t.Error(n, err)
				return
			}

			if ds != rs {
				t.Errorf("[%d] want %q, but %q", n, rs, ds)
				return
			}
		}

		for i := 64; i <= 128; i++ {
			wg.Add(1)
			go tf(i)
		}
		wg.Wait()
	}
}
