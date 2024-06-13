# Blog

This project made for learning purposes. 

As for technologies, I was interested in Golang with Templ and HTMX. 

For styles - TailwindCSS used.

For some simple frontend interactions - Alpine.JS

To bootstrap project setup - gowebly used.

## DEMO

You can look into my blog at [blog.ineersa.com](http://blog.ineersa.com/)

It serves as frontend application part for [Laravel/Filament backend](https://github.com/ineersa/blog-admin).

## Run and build

To run project locally with hot swap you need to install [gowebly](https://github.com/gowebly/gowebly)

After that you can just run
```bash
gowebly run
```
And will get your server started and watching for your local file changes with air.

You can look into Makefile and build binary or run application. 

## Web server config

I've decided not to serve static files via Golang built in web server, since I'm anyway working behind `nginx` both on local machine and on server. 

Basic nginx config:
```nginx
server {
    server_name blog.ineersa.local www.blog.ineersa.local;

    add_header X-Frame-Options "SAMEORIGIN";
    add_header X-XSS-Protection "1; mode=block";
    add_header X-Content-Type-Options "nosniff";

    charset utf-8;

    # Serve static files directly from Nginx
    location /static/ {
        root /var/www/blog;
        try_files $uri $uri/ =404;
        expires 1y;
        etag on;
        if_modified_since exact;
        add_header Cache-Control "public, no-cache";
    }

    location /errors/ {
        alias /var/www/blog/error_pages/;
        try_files $uri $uri/ =404;
    }

    error_page 404 /errors/404.html;
    error_page 500 /errors/500.html;
    error_page 501 502 503 504 505 507 /errors/50x.html;
    error_page 403 /errors/403.html;

    location / {
        proxy_pass http://localhost:7001;
	    proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
    	proxy_set_header Connection "upgrade";
    	proxy_set_header Host $host;
    	proxy_set_header X-Real-IP $remote_addr;
    	proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    	proxy_set_header X-Forwarded-Proto $scheme;
    	proxy_intercept_errors on;
    }

    location = /favicon.ico { access_log off; log_not_found off; }
    location = /robots.txt  { access_log off; log_not_found off; }

    location ~ /\.ht {
        deny all;
    }

    location ~ /\.(?!well-known).* {
        deny all;
    }

    proxy_read_timeout 60s;
    proxy_send_timeout 60s;
    proxy_connect_timeout 60s;

    error_log /var/log/nginx/blog.ineersa.local.error.log;
    access_log /var/log/nginx/blog.ineersa.local.access.log;
}
```

Also I'm serving error pages directly from html pages.

## TODO
 - Adding good semantic search
 - Adding cache for database dictionaries
 - Adding cache to HTML pages
 - Adding dependency injection 
 - Frontend refactoring to re-render smaller parts
 - Live search implementation
