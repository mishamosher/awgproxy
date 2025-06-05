module github.com/mishamosher/awgproxy

go 1.24.3

require (
	github.com/MakeNowJust/heredoc/v2 v2.0.1
	github.com/akamensky/argparse v1.4.0
	github.com/amnezia-vpn/amneziawg-go v0.2.12
	github.com/go-ini/ini v1.67.0
	github.com/landlock-lsm/go-landlock v0.0.0-20250303204525-1544bccde3a3
	github.com/things-go/go-socks5 v0.0.6
	golang.org/x/net v0.40.0
	suah.dev/protect v1.2.4
)

require (
	github.com/google/btree v1.1.3 // indirect
	github.com/tevino/abool/v2 v2.1.0 // indirect
	golang.org/x/crypto v0.38.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/time v0.11.0 // indirect
	golang.zx2c4.com/wintun v0.0.0-20230126152724-0fa3db229ce2 // indirect
	gvisor.dev/gvisor v0.0.0-20250605053741-4badbeb38f66 // indirect
	kernel.org/pub/linux/libs/security/libcap/psx v1.2.76 // indirect
)

replace github.com/amnezia-vpn/amneziawg-go => github.com/mishamosher/amneziawg-go v0.2.12-1
