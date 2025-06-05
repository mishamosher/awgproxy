# Running awgproxy with systemd

If you're on a systemd-based distro, you'll most likely want to run awgproxy as a systemd unit.

The provided systemd unit assumes you have the awgproxy executable installed on `/opt/awgproxy/awgproxy` and a configuration file stored at `/etc/awgproxy.conf`. These paths can be customized by editing the unit file.

# Setting up the unit

1. Copy the `awgproxy.service` file from this directory to `/etc/systemd/system/`, or use the following cURL command to download it:
   ```bash
   curl https://raw.githubusercontent.com/pufferffish/awgproxy/master/systemd/awgproxy.service | sudo tee /etc/systemd/system/awgproxy.service
   ```

2. If necessary, customize the unit.

   Edit the parts with `LoadCredential`, `ExecStartPre=` and `ExecStart=` to point to the executable and the configuration file. For example, if awgproxy is installed on `/usr/bin` and the configuration file is located in `/opt/myfiles/awgproxy.conf` do the following change:
   ```service
   LoadCredential=conf:/opt/myfiles/awgproxy.conf
   ExecStartPre=/usr/bin/awgproxy -n -c ${CREDENTIALS_DIRECTORY}/conf
   ExecStart=/usr/bin/awgproxy -c ${CREDENTIALS_DIRECTORY}/conf
   ```

4. Reload systemd and enable the unit.
   ```bash
   sudo systemctl daemon-reload
   sudo systemctl enable --now awgproxy.service
   ```

5. Make sure it's working correctly.

   Finally, check out the unit status to confirm `awgproxy.service` has started without problems. You can use commands like `systemctl status awgproxy.service` and/or `sudo journalctl -u awgproxy.service`.

# Additional notes

If you want to disable the extensive logging that's done by awgproxy, simply add `-s` parameter to `ExecStart=`. This will enable the silent mode that was implemented with [pull/67](https://github.com/pufferffish/wireproxy/pull/67).
