module main

go 1.18

require (
	github.com/flytam/filenamify v1.1.2
	github.com/go-chi/chi/v5 v5.0.8
	github.com/google/uuid v1.3.0
	github.com/mattn/go-sqlite3 v1.14.16
	github.com/sirupsen/logrus v1.9.0
)

require (
	github.com/andybalholm/cascadia v1.3.1 // indirect
	github.com/pquerna/cachecontrol v0.0.0-20180517163645-1555304b9b35 // indirect
	golang.org/x/net v0.0.0-20210916014120-12bc252f5db8 // indirect
)

require (
	github.com/PuerkitoBio/goquery v1.8.0
	github.com/denisbrodbeck/striphtmltags v6.6.6+incompatible
	github.com/iand/microdata v0.0.3
	github.com/jmoiron/sqlx v1.3.5
	github.com/johnjones4/keeper/core v0.0.0-00010101000000-000000000000
	github.com/piprate/json-gold v0.5.0
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
)

replace github.com/johnjones4/keeper/core => ../core
