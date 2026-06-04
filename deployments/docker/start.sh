#!/bin/sh
/app/server &
cloudflared tunnel --url http://localhost:8080
