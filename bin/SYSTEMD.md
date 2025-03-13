# Systemd Configuration

## Users

```bash
adduser empyr --shell /usr/sbin/nologin --home /home/empyr
gpasswd -d empyr users
chmod 755 /home/empyr/
cd /home/empyr/
rm ~empyr/.bash* ~empyr/.cloud-locale-test.skip ~empyr/.profile
```

## Files

```bash
root@epimethian:/etc/systemd/system# ll /etc/systemd/system/empyr-*.service
-rw-r--r-- 1 root root 421 Sep 30 22:06 empyr-a01.service
-rw-r--r-- 1 root root 421 Sep 30 22:06 empyr-a02.service

root@epimethian:/var/www/a01/bin# ll /var/www/a01/bin
drwxrwxr-x 2 empyr empyr     4096 Oct 13 22:24 ./
drwxr-xr-x 6 empyr empyr     4096 Oct 13 21:57 ../
-rwxr-xr-x 1 empyr empyr 18013459 Oct 13 21:57 empyr

root@epimethian:/var/www/a01/bin# ll /var/www/a02/bin
drwxrwxr-x 2 empyr empyr     4096 Oct 13 22:24 ./
drwxr-xr-x 6 empyr empyr     4096 Oct 13 21:57 ../
-rwxr-xr-x 1 empyr empyr 18013459 Oct 13 21:57 empyr
```

## Install

```bash
root@epimethian:/etc/systemd/system# systemctl daemon-reload

root@epimethian:/etc/systemd/system# systemctl status empyr-a01.service
○ empyr-a01.service - Empyr A01 Web Service
     Loaded: loaded (/etc/systemd/system/empyr-a01.service; disabled; preset: enabled)
     Active: inactive (dead)

root@epimethian:/etc/systemd/system# systemctl enable empyr-a01.service
Created symlink '/etc/systemd/system/multi-user.target.wants/empyr-a01.service' → '/etc/systemd/system/empyr-a01.service'.

root@epimethian:/etc/systemd/system# systemctl status empyr-a01.service
○ empyr-a01.service - Empyr A01 Web Service
     Loaded: loaded (/etc/systemd/system/empyr-a01.service; enabled; preset: enabled)
     Active: inactive (dead)

root@epimethian:/etc/systemd/system# systemctl start empyr-a01.service

root@epimethian:/etc/systemd/system# systemctl status empyr-a01.service
● empyr-a01.service - Empyr A01 Web Service
     Loaded: loaded (/etc/systemd/system/empyr-a01.service; enabled; preset: enabled)
     Active: active (running) since Thu 2025-03-13 17:39:31 UTC; 3s ago
 Invocation: 07cb01b33aa74daea8ce7b3f95479665
   Main PID: 9889 (empyr)
      Tasks: 6 (limit: 1109)
     Memory: 1.5M (peak: 1.6M)
        CPU: 18ms
     CGroup: /system.slice/empyr-a01.service
             └─9889 /var/www/a01/bin/empyr start server

Mar 13 17:39:31 epimethian systemd[1]: Started empyr-a01.service - Empyr A01 Web Service.
Mar 13 17:39:31 epimethian empyr[9889]: load.go:92: env: loaded ".env"
Mar 13 17:39:31 epimethian empyr[9889]: open.go:104: store: open: /var/www/a01/data/a01.db
Mar 13 17:39:31 epimethian empyr[9889]: routes.go:32: assets: /var/www/a01/public
Mar 13 17:39:31 epimethian empyr[9889]: server.go:76: server: listening on http://localhost:8881

root@epimethian:/etc/systemd/system# systemctl stop empyr-a01.service
root@epimethian:/etc/systemd/system# systemctl status empyr-a01.service
○ empyr-a01.service - Empyr A01 Web Service
     Loaded: loaded (/etc/systemd/system/empyr-a01.service; enabled; preset: enabled)
     Active: inactive (dead) since Thu 2025-03-13 17:41:00 UTC; 3s ago
   Duration: 1min 29.064s
 Invocation: 07cb01b33aa74daea8ce7b3f95479665
    Process: 9889 ExecStart=/var/www/a01/bin/empyr start server (code=exited, status=0/SUCCESS)
   Main PID: 9889 (code=exited, status=0/SUCCESS)
   Mem peak: 1.7M
        CPU: 20ms

Mar 13 17:40:18 epimethian empyr[9889]: views.go:33: responders: parsing [/var/www/a01/templates/logout.gohtml]
Mar 13 17:41:00 epimethian systemd[1]: Stopping empyr-a01.service - Empyr A01 Web Service...
Mar 13 17:41:00 epimethian empyr[9889]: server.go:85: server: signal terminated: received after 1m29.013996123s
Mar 13 17:41:00 epimethian empyr[9889]: server.go:90: server: timeout 5s: creating context (137ns)
Mar 13 17:41:00 epimethian empyr[9889]: server.go:95: server: canceling idle connections (43.253µs)
Mar 13 17:41:00 epimethian empyr[9889]: server.go:98: server: shutting down the server (52.534µs)
Mar 13 17:41:00 epimethian empyr[9889]: server.go:103: server: ¡stopped gracefully! (82.689µs)
Mar 13 17:41:00 epimethian empyr[9889]: start.go:72: server: shut down after 1m29.016373338s
Mar 13 17:41:00 epimethian systemd[1]: empyr-a01.service: Deactivated successfully.
Mar 13 17:41:00 epimethian systemd[1]: Stopped empyr-a01.service - Empyr A01 Web Service.
```

## Monitor

```bash
root@epimethian:/etc/systemd/system# journalctl -f -u empyr-a01.service
```

