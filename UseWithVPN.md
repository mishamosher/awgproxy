# Getting a wireguard/amnezia/euphoria Server

This fork supports both standard WireGuard configs and AmneziaWG configs (which add
obfuscation fields like `Jc`, `Jmin`, `Jmax`, `S1`, `S2`, `H1`–`H4`).

You can get an AmneziaWG config from your VPN provider or by running your own
AmneziaWG server. Export the config from the Amnezia client app, or use
`wg genkey` / `wg pubkey` and configure the obfuscation parameters manually.

The config file you download is used as-is via `WGConfig` (see below), so no
conversion is needed.

# Simple Setup for multiple SOCKS configs for Firefox

Create a folder for your configs and startup scripts. Can be the same place as
this code. That path you will use below. For reference this text uses
`/Users/jonny/vpntabs`

For each VPN you want to run, you will download your config and name
it appropriately (e.g. `MyVPN.adblock.server.conf`) and then create two new
files from those below with similar names (e.g. `MyVPN.adblock.conf` and
`MyVPN.adblock.sh`)

You will also create a launch script, the reference below is only for macOS. The
naming should also be similar (e.g.
`/Users/jonny/Library/LaunchAgents/com.MyVPN.adblock.plist`)

## Config File

Make sure you use a unique port for every separate server.
You can optionally set proxy authentication; the same user/pass can be reused across servers.

```ini
# Link to the downloaded config (WireGuard or AmneziaWG format)
WGConfig = /Users/jonny/vpntabs/MyVPN.adblock.server.conf

# Used for Firefox containers
[Socks5]
BindAddress = 127.0.0.1:25344 # Update the port here for each new server

# Socks5 authentication parameters, specifying username and password enables
# proxy authentication.
#Username = ...
# Avoid using spaces in the password field
#Password = ...
```

## Startup Script File

This is a bash script to facilitate startup, not strictly essential, but adds
ease.
Note, you MUST update the first path to wherever you installed this binary.
Make sure you use the path for the config file above, not the one you downloaded
from your provider.

```bash
#!/bin/bash
/Users/jonny/awgproxy/awgproxy -c /Users/jonny/vpntabs/MyVPN.adblock.conf
```

## macOS LaunchAgent

To make it run every time you start your computer, you can create a launch agent
in `$HOME/Library/LaunchAgents`. Name reference above.

That file should contain the following; the label should be the same as the file
name and the paths should be set correctly:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.MyVPN.adblock</string>
    <key>Program</key>
    <string>/Users/jonny/vpntabs/MyVPN.adblock.sh</string>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
</dict>
</plist>
```

To enable it, run
`launchctl load ~/Library/LaunchAgents/com.MyVPN.adblock.plist` and
`launchctl start ~/Library/LaunchAgents/com.MyVPN.adblock.plist`

# Firefox Setup

You will need to enable the Multi Account Container Tabs extension and a proxy extension.
Sideberry works well, but Container Proxy also works.

Create a container to be dedicated to this VPN, and then add the IP, port,
username, and password from above.
