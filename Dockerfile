FROM shyim/shopware-platform-nginx:php74

COPY . /var/www/html/public
COPY --from=composer /usr/bin/composer /usr/bin/composer

RUN composer install -d /var/www/html/public
COPY nginx.conf.sigil /var/www/html
