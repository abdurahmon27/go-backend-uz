//go:build linux

// Uy vazifasi javoblari (README.md, "Uy vazifasi").
//
//	go run ./uyvazifa guruh      O'rta: har yozuvga fsync va 10 ta yozuvga bitta fsync
//	go run ./uyvazifa unix       Qiyin: localhost round trip, TCP va Unix socket
//	go run ./uyvazifa parallel   Guru: SSD'dan 4 KiB, 1 / 4 / 32 goroutine bir vaqtda (O_DIRECT, Linux)
//
// parallel shu papkadagi sinov.dat ni o'qiydi: uni avval `go run ./olcha disk` yaratadi.
package main

import (
	"fmt"
	"io"
	"math/rand/v2"
	"net"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"syscall"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "ishlatish: go run ./uyvazifa guruh|unix|parallel")
		os.Exit(2)
	}
	switch os.Args[1] {
	case "guruh":
		guruh()
	case "unix":
		unixSocket()
	case "parallel":
		parallel()
	}
}

func tekshir(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "xato:", err)
		os.Exit(1)
	}
}

func median(t []time.Duration) time.Duration {
	s := slices.Clone(t)
	slices.Sort(s)
	return s[len(s)/2]
}

func p(t []time.Duration, q float64) time.Duration {
	s := slices.Clone(t)
	slices.Sort(s)
	return s[int(q*float64(len(s)-1)+0.5)]
}

// guruh: ma'lumotlar bazalaridagi "group commit" g'oyasi: bir nechta yozuv, bitta fsync.
func guruh() {
	f, err := os.Create("guruh.dat")
	tekshir(err)
	defer os.Remove("guruh.dat")
	defer f.Close()
	b := make([]byte, 4096)
	const n = 300
	for _, k := range []int{1, 10} { // k ta yozuvdan keyin bitta fsync
		t0 := time.Now()
		for i := range n {
			_, err := f.WriteAt(b, int64(i%256)*4096)
			tekshir(err)
			if (i+1)%k == 0 {
				tekshir(f.Sync())
			}
		}
		d := time.Since(t0)
		fmt.Printf("%2d ta yozuvga 1 fsync: bitta yozuvga %v\n", k, (d / n).Round(time.Microsecond))
	}
}

// unixSocket: bitta baytli round trip, TCP (127.0.0.1) va Unix domain socket.
func unixSocket() {
	yol := filepath.Join(os.TempDir(), fmt.Sprintf("uyvazifa-%d.sock", os.Getpid()))
	defer os.Remove(yol)
	for _, t := range [][2]string{{"tcp", "127.0.0.1:0"}, {"unix", yol}} {
		ln, err := net.Listen(t[0], t[1])
		tekshir(err)
		go func() {
			c, err := ln.Accept()
			if err == nil {
				io.Copy(c, c)
			}
		}()
		c, err := net.Dial(t[0], ln.Addr().String())
		tekshir(err)
		b := []byte{1}
		var v []time.Duration
		for range 20000 {
			t0 := time.Now()
			_, err := c.Write(b)
			tekshir(err)
			_, err = io.ReadFull(c, b)
			tekshir(err)
			v = append(v, time.Since(t0))
		}
		c.Close()
		ln.Close()
		fmt.Printf("%-4s round trip: mediana %v, p90 %v\n", t[0], median(v), p(v, 0.9))
	}
}

// parallel: bir nechta goroutine bir vaqtda SSD'dan o'qiydi. Har bir o'qish kechikishi (latency)
// va soniyasiga nechta o'qish (IOPS) qanday o'zgaradi?
func parallel() {
	const (
		kb4  = 4096
		hajm = 1 << 30
		n    = 8000
	)
	for _, g := range []int{1, 4, 32} {
		f, err := os.OpenFile("sinov.dat", os.O_RDONLY|syscall.O_DIRECT, 0)
		tekshir(err)
		var mu sync.Mutex
		var v []time.Duration
		var wg sync.WaitGroup
		t0 := time.Now()
		for j := range g {
			wg.Go(func() {
				buf, err := syscall.Mmap(-1, 0, kb4, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_ANON|syscall.MAP_PRIVATE)
				tekshir(err)
				r := rand.New(rand.NewPCG(uint64(j), 7))
				var my []time.Duration
				for range n / g {
					t := time.Now()
					_, err := f.ReadAt(buf, int64(r.IntN(hajm/kb4))*kb4)
					tekshir(err)
					my = append(my, time.Since(t))
				}
				mu.Lock()
				v = append(v, my...)
				mu.Unlock()
			})
		}
		wg.Wait()
		d := time.Since(t0)
		f.Close()
		fmt.Printf("%2d goroutine: mediana %6v, p90 %6v, %6.0f o'qish/s\n", g,
			median(v).Round(time.Microsecond), p(v, 0.9).Round(time.Microsecond), float64(len(v))/d.Seconds())
	}
}
