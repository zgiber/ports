package parser

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/zgiber/ports/service"
)

const (
	lvlUp   = "[{"
	lvlDown = "]}"
)

type Parser struct {
	next chan *service.Port
	err  error
}

func New(ctx context.Context, in io.Reader) *Parser {
	p := &Parser{next: make(chan *service.Port, 1)}
	go p.parse(ctx, in)
	return p
}

func (p *Parser) Next() (*service.Port, error) {
	if p.err != nil {
		return nil, p.err
	}

	port, isParserOK := <-p.next
	if !isParserOK {
		return nil, io.EOF
	}

	return port, nil
}

func (p *Parser) parse(ctx context.Context, update io.Reader) {
	dec := json.NewDecoder(update)
	defer close(p.next)

	level := 0
	for dec.More() || level > 0 {
		if err := ctx.Err(); err != nil {
			p.err = err
			return
		}

		t, err := dec.Token()
		if err != nil {
			p.err = err
			return
		}

		// The JSON input is structured so there is only a single string
		// value on level 1, which is the key for the actual port details.
		// Therefore the decoding logic is that if the decoder is at
		// level1 and a string token is encountered, then the string
		// value is taken as an ID and the rest is decoded as details.
		if tokenValue, isString := t.(string); isString {
			if level == 1 {
				port := &service.Port{ID: tokenValue}
				err := dec.Decode(&port.Details)
				if err != nil {
					p.err = err
					return
				}
				p.next <- port
			}
		}

		if strings.Contains(lvlUp, fmt.Sprint(t)) {
			level++
		}

		if strings.Contains(lvlDown, fmt.Sprint(t)) {
			level--
		}
	}
}
