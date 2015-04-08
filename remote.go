package plugins

import (
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"io"
	"runtime"
)

type Encoder interface {
	Encode(e interface{}) error
}

type Decoder interface {
	Decode(e interface{}) error
}

type PluginDecl struct {
	Name       string   `json:"name"`
	Subscribes []Event  `json:"subscribes"`
	Provides   []Event  `json:"provides"`
	Formats    []Format `json:"formats"`
}

type formatResponse struct {
	Error  Error `json:"error,omitempty"`
	Format `json:"format"`
}

type Format string

var (
	JSON Format = "application/json"
	GOB  Format = "application/gob"
	XML  Format = "application/xml"
)

type RemotePlugin struct {
	PluginDecl
	enc Encoder
	dec Decoder
	r   io.ReadCloser
	w   io.WriteCloser
}

type FormatFactory struct {
	Weight    int
	Construct func(r io.ReadCloser, w io.WriteCloser) (Encoder, Decoder)
}

var SupportedFormats = map[Format]FormatFactory{
	//GOB:  FormatFactory{100, gobFactory},
	JSON: FormatFactory{50, jsonFactory},
	XML:  FormatFactory{0, xmlFactory},
}

func NewRemotePlugin(r io.ReadCloser, w io.WriteCloser) (pl *RemotePlugin, err error) {

	pl = new(RemotePlugin)
	pl.r = r
	pl.w = w

	//Pack reader in JSON decoder and decode a PluginDecl
	dec := json.NewDecoder(r)
	enc := json.NewEncoder(w)

	err = dec.Decode(&pl.PluginDecl)
	if err != nil {
		return
	}

	var best FormatFactory = FormatFactory{-1, nil}
	var bestFormat Format = ""
	for _, format := range pl.PluginDecl.Formats {
		if factory, ok := SupportedFormats[format]; ok && factory.Weight > best.Weight {
			best = factory
			bestFormat = format
		}
	}

	if bestFormat == "" {
		enc.Encode(formatResponse{Error: NoSupportedFormat})
		err = NoSupportedFormat
		return
	}

	enc.Encode(formatResponse{Format: bestFormat})

	runtime.SetFinalizer(pl, func(pl *RemotePlugin) {
		pl.r.Close()
		pl.w.Close()
	})

	pl.enc, pl.dec = best.Construct(r, w)
	return
}

func (rp RemotePlugin) Name() string {
	return rp.PluginDecl.Name
}

func (rp RemotePlugin) Subscribes() []Event {
	return rp.PluginDecl.Subscribes
}

func (rp RemotePlugin) Provides() []Event {
	return rp.PluginDecl.Provides
}

func (rp RemotePlugin) Send(e Event, args Args) error {
	return rp.enc.Encode(packet{e, args})
}

func (rp RemotePlugin) Recieve() (Event, Args, error) {
	pck := packet{}
	err := rp.dec.Decode(&pck)
	return pck.Event, pck.Args, err
}

func gobFactory(r io.ReadCloser, w io.WriteCloser) (Encoder, Decoder) {
	return gob.NewEncoder(w), gob.NewDecoder(r)
}

func jsonFactory(r io.ReadCloser, w io.WriteCloser) (Encoder, Decoder) {
	return json.NewEncoder(w), json.NewDecoder(r)
}

func xmlFactory(r io.ReadCloser, w io.WriteCloser) (Encoder, Decoder) {
	return xml.NewEncoder(w), xml.NewDecoder(r)
}
