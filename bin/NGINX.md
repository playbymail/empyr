# Directories

```bash
root@epimethian:/var/www/a01# ll /var/www/a01
drwxr-xr-x 2 root  root  4096 Mar 11 02:32 bin/
drwxr-xr-x 2 root  root  4096 Mar 13 08:46 build/
drwxr-xr-x 2 empyr empyr 4096 Mar 11 02:32 data/
drwxr-xr-x 4 root  root  4096 Mar 11 02:32 public/

root@epimethian:/var/www/a01# ll /var/www/a02
total 24
drwxr-xr-x 2 root  root  4096 Mar 11 02:32 bin/
drwxr-xr-x 2 root  root  4096 Mar 11 02:32 build/
drwxr-xr-x 2 empyr empyr 4096 Mar 11 02:32 data/
drwxr-xr-x 4 root  root  4096 Mar 11 02:32 public/
```

```nginx
upstream empyr {
    server localhost:8881;
}

server {
    server_name a01.epimethean.dev;
    root        /var/www/a01/public;
    access_log  /var/log/nginx/a01.access.log;
    error_log   /var/log/nginx/a01.error.log  crit;

    index index.html;

    location / {
        proxy_pass        http://localhost:8881;
        proxy_set_header  Host $host;
        proxy_set_header  X-Real-IP $remote_addr;
        proxy_set_header  X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header  X-Forwarded-Proto $scheme;

        # Serve maintenance page if the service is unavailable
        error_page 502 503 504 /maintenance.html;
        location = /maintenance.html {
            root /var/www/a01/public;
        }
    }
}
```