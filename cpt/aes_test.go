package cpt

import (
	"fmt"
	"sync"
	"testing"

	"github.com/askasoft/pango/gog"
	"github.com/askasoft/pango/ran"
)

func TestAesGCMEncrypt(t *testing.T) {
	bits := []int{128, 192, 256}
	for _, bit := range bits {
		c := gog.Must(NewAesGCMCryptor("1234567890abcde", bit))
		for i := 64; i <= 128; i++ {
			rs := ran.RandString(i)
			es, err := c.EncryptString(rs)
			if err != nil {
				t.Fatal(i, err)
			}

			fmt.Printf("%d [%d] len(es) = %d\n", bit, i, len(es))

			ds, err := c.DecryptString(es)
			if err != nil {
				t.Fatal(i, err)
			}

			if ds != rs {
				t.Fatalf("[%d] want %q, but %q", i, rs, ds)
			}
		}
	}
}

func TestAesGCMConcurrentEncrypt(t *testing.T) {
	bits := []int{128, 192, 256}
	for _, bit := range bits {
		c := gog.Must(NewAesGCMCryptor("1234567890abcde", bit))

		var wg sync.WaitGroup

		tf := func(n int) {
			defer wg.Done()

			for range 1000 {
				rs := ran.RandString(n)
				es, err := c.EncryptString(rs)
				if err != nil {
					t.Error(n, err)
					return
				}

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
		}

		for i := 64; i <= 128; i++ {
			wg.Add(1)
			go tf(i)
		}
		wg.Wait()
	}
}

func TestAesCBCEncrypt(t *testing.T) {
	bits := []int{128, 192, 256}
	for _, bit := range bits {
		c := gog.Must(NewAesCBCCryptor("1234567890abcde", bit))
		for i := 64; i <= 128; i++ {
			rs := ran.RandString(i)
			es, err := c.EncryptString(rs)
			if err != nil {
				t.Fatal(i, err)
			}

			fmt.Printf("%d [%d] len(es) = %d\n", bit, i, len(es))

			ds, err := c.DecryptString(es)
			if err != nil {
				t.Fatal(i, err)
			}

			if ds != rs {
				t.Fatalf("[%d] want %q, but %q", i, rs, ds)
			}
		}
	}
}

func TestAesCBCConcurrentEncrypt(t *testing.T) {
	bits := []int{128, 192, 256}
	for _, bit := range bits {
		c := gog.Must(NewAesCBCCryptor("1234567890abcde", bit))

		var wg sync.WaitGroup

		tf := func(n int) {
			defer wg.Done()

			for range 1000 {
				rs := ran.RandString(n)
				es, err := c.EncryptString(rs)
				if err != nil {
					t.Error(n, err)
					return
				}

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
		}

		for i := 64; i <= 128; i++ {
			wg.Add(1)
			go tf(i)
		}
		wg.Wait()
	}
}
