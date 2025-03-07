# 
- [Install Sqlc](https://sqlc.dev)
- [Install Templ](https://templ.guide)
- [Install Goose](https://github.com/pressly/goose)

# 
```bash
(cd sqlc && sqlc generate)
goose up
templ generate
go run dstfrsh/main.go serve
```