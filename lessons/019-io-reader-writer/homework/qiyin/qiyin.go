// Package qiyin — Qiyin: N baytdan keyin yozishni rad etadigan Writer.
// Shartnoma: n < len(p) boʻlsa, xato nil boʻlmasligi shart.
package qiyin

import (
	"errors"
	"io"
)

var ErrToldi = errors.New("chegara toʻldi")

type Chekli struct {
	W     io.Writer
	Qoldi int64
}

func (c *Chekli) Write(p []byte) (int, error) {
	if int64(len(p)) <= c.Qoldi {
		n, err := c.W.Write(p)
		c.Qoldi -= int64(n)
		return n, err
	}
	n, err := c.W.Write(p[:c.Qoldi]) // sigʻganini yozamiz
	c.Qoldi -= int64(n)
	if err == nil {
		err = ErrToldi
	}
	return n, err
}
