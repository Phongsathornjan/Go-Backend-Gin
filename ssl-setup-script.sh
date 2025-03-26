#!/bin/bash

# Domain name 
DOMAIN="test.com"

# Request initial certificate
docker-compose run --rm certbot certonly \
    --webroot \
    --webroot-path /var/www/certbot/ \
    -d $DOMAIN \
    -d www.$DOMAIN

# Create a renewal script
cat << EOF > renew-ssl.sh
#!/bin/bash
docker-compose run --rm certbot renew
docker-compose kill -s HUP nginx
EOF

chmod +x renew-ssl.sh

# Set up crontab for automatic renewal
(crontab -l 2>/dev/null; echo "0 0,12 * * * /path/to/renew-ssl.sh") | crontab -