## Plugins ##

This is a implementation of a "plugin" system for Go.
Go does not have dynamic loading so it is not a true plugin system, but rather a way to handle inter-process communication with intention to augment a main system.

Plugins can be connected in multiple ways.
They are defined as a interface, so as long as that is implemented, you can hook one up in what ever way you want, but the default implementation is RemotePlugin that defines a way to connect to a another process over any Reader and Writer with some protocol negotiations of format (gob, json and xml supported out-of-the-box).

This means that together with a convenience function to run all binaries from a folder, it is easy to set up a plugin api for your program.

1. Create a PluginHandler
2. Setup some events that you subscribe to.
3. Load all the plugins in the folder.
4. ???
5. Profit!

When ever something happens in your application, notify all interested plugins by using the PluginHandler.Dispatch(Event, Args) method

### Protocol ###
The protocol is divided into two stages.
The first stage is a two-way handshake where the plugin/client starts.

Client --{ PluginDeclaration }-> Server
Client <-{   FormatResponse  }-- Server

where PluginDeclaration is a JSON-encoded object with the following fields:
	name       string
	subscribes [Event]
	provides   [Event]
	formats    [Format]

	[Type] means a list of <Type>
	Event is a string representing a event, as dot-separated nodes in a tree
		Ex: "job.complete"
			"tracker.calibration.start"

	Format is a string representation of a MIME type as defined in RFC 2046.
		Ex: "text/json"
			"text/xml"
			"bin/gob"

and FormatDeclaration is a JSON encoded object with these fields:
	format Format
	error Error

where Format is as defined above and Error is a integer error-code as defined in the Errors section below.


The second stage is the communication stage, using the format negotiated in the handshake stage.

This stage use packets on the form
	Event
	Args

in that order (if important for the format chosen).
	Event is as described above
	Args is a mapping from string to any type. In JSON, this would be a simple object. In Gob, it would be a map[string]interface{}

These packets can be sent at any time by any side.
Both sides are to listen for incomming packets at all times and handle them in a timely fashion.

This stage continues until either side closes the communications channel.

### Errors ###

#### Transmission Errors ###¤
	0: Success
	1: No common format found during negotiations.

#### Server Errors ####
*To be written

#### Client Errors ####
*To be written

