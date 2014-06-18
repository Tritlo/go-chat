Go chat
=======

[![Build Status](https://travis-ci.org/Tritlo/go-chat.svg?branch=master)](https://travis-ci.org/Tritlo/go-chat)

A very simple chat program using go-socket.io (from googollee) and martini
(from codegansta).

Uses either a port or a socket.

After installing dependencies with:
    go get     github.com/codegangsta/martini
    go get github.com/codegangsta/martini-contrib/render
    go get github.com/googollee/go-socket.io


Run with:
    go run chat.go

Reads PORT from environment, but uses a socket if SOCKET environment variable
exists. Default port is 3000.

Example nginx.conf (a little changed from a live config):

    upstream chat_server {
        # server 0.0.0.0:3000 fail_timeout=0;
        server unix:/webapps/chat/run/chat.sock fail_timeout=0;
    }

    server {
        listen 80;

        location /socket.io {
            proxy_pass http://chat_server;
            proxy_redirect off;
            proxy_buffering off;

            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;

            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection "Upgrade";
        }


        location / {
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto https;

            proxy_set_header Host $http_host;
            proxy_redirect off;
            proxy_pass http://chat_server;
        }
    }
