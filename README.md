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

## TODO
 - Adding good semantic search
 - Adding cache for database dictionaries
 - Adding cache to HTML pages
 - Adding dependency injection 
 - Frontend refactoring to re-render smaller parts
 - Live search implementation