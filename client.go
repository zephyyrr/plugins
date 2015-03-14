package plugins

import (
	"encoding/json"
	"io"
)

type Client struct {
	Muxer
	enc Encoder
	dec Decoder
}

//Creates a new plugin client from a declaration, a reader and the corresponding writer.
//The Format field of the declaration can be omitted (nil) and will in that case be set
// to a list of formats supported by the system (All keys in SupportedFormats).
func NewClient(decl PluginDecl, r io.ReadCloser, w io.WriteCloser) (c *Client, err error) {
	c = new(Client)
	c.Muxer = make(mapMuxr)

	//Begin protocol handshake
	enc := json.NewEncoder(w)
	dec := json.NewDecoder(r)

	if decl.Formats == nil {
		//Perhaps this can be cached.
		//On the other hand, I expect this to be done once per process anyways, so no point unless hardcoded.
		//But that would defeat the purpose of custom encoder/decoders.
		for format := range SupportedFormats {
			decl.Formats = append(decl.Formats, format)
		}
	}

	if decl.Subscribes == nil {
		//Perhaps this can be cached.
		for format := range SupportedFormats {
			decl.Formats = append(decl.Formats, format)
		}
	}

	//Send declaration
	enc.Encode(decl)

	//Read response
	var resp formatResponse
	dec.Decode(&resp)
	if resp.Error != Success {
		err = resp.Error
		return
	}

	c.enc, c.dec = SupportedFormats[resp.Format].Construct(r, w)
	return
}

func (c *Client) Run() (err error) {
	for {
		//Read packets
		var pck packet
		err = c.dec.Decode(&pck)
		if err != nil {
			return err
		}

		//Send to muxr
		go c.Muxer.HandleEvent(pck.Event, pck.Args)
	}
}

func (c *Client) Dispatch(event Event, args Args) error {
	return c.enc.Encode(packet{event, args})
}
