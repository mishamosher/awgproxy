module github.com/mishamosher/awgproxy

go 1.24.6

require (
	github.com/MakeNowJust/heredoc/v2 v2.0.1
	github.com/akamensky/argparse v1.4.0
	github.com/amnezia-vpn/amneziawg-go v0.2.13
	github.com/go-ini/ini v1.67.0
	github.com/landlock-lsm/go-landlock v0.0.0-20250303204525-1544bccde3a3
	github.com/things-go/go-socks5 v0.0.6
	golang.org/x/net v0.42.0
	suah.dev/protect v1.2.4
)

require (
	github.com/google/btree v1.1.3 // indirect
	github.com/tevino/abool v1.2.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	golang.org/x/crypto v0.40.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/time v0.12.0 // indirect
	golang.zx2c4.com/wintun v0.0.0-20230126152724-0fa3db229ce2 // indirect
	gvisor.dev/gvisor v0.0.0-20250718015824-35000683b6d7 // indirect
	kernel.org/pub/linux/libs/security/libcap/psx v1.2.76 // indirect
)

replace github.com/things-go/go-socks5 => github.com/mishamosher/go-socks5 v0.0.6-1
