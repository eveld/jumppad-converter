#!/bin/bash
cat <<EOF >> /etc/supervisord.conf

[program:counting-service]
environment=PORT=9003
command=/usr/local/bin/counting-service
autostart=true
autorestart=true
stderr_logfile=/dev/null
stdout_logfile=/dev/null
EOF

supervisorctl reread
supervisorctl update