# Running awgproxy with rc.d

If you're on a rc.d-based distro, you'll most likely want to run AWGProxy as a systemd unit.

The provided systemd unit assumes you have the awgproxy executable installed on `/bin/awgproxy` and a configuration file stored at `/etc/awgproxy.conf`. These paths can be customized by editing the unit file.

# Setting up the unit

1. Copy the `awgproxy` file from this directory to `/usr/local/etc/rc.d`.

2. If necessary, customize the unit.
   Edit the parts with `procname`, `command`, `awgproxy_conf`  to point to the executable and the configuration file.

4. Add the following lines to `/etc/rc.conf` to enable awgproxy
   `awgproxy_enable="YES"`

5. Start awgproxy service and check status
   ```
   sudo service awgproxy start
   sudo service awgproxy status
   ```
