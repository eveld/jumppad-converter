#!/bin/bash
cat << EOF >> /etc/supervisord.conf

[program:dashboard-service]
environment=PORT=9002
command=/usr/local/bin/dashboard-service
autostart=true
autorestart=true
stderr_logfile=/dev/null
stdout_logfile=/dev/null
EOF

supervisorctl reread
supervisorctl update