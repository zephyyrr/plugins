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
