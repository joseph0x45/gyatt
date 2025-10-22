#!/usr/bin/env sh
set -e

project_name=$1

if [ -z "$project_name" ]; then
  echo "Usage: $0 <project-name>"
  exit 1
fi

go mod init $1
go get github.com/a-h/templ
go get -tool github.com/a-h/templ/cmd/templ@latest

mkdir -p components
mkdir -p db
mkdir -p static
mkdir -p handlers
mkdir -p models

cat > tailwind.config.js <<EOF
/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./components/*.templ"],
  theme: {
    extend: {},
  },
  plugins: [],
};
EOF

printf 'build:\n\tgo tool templ generate && go generate && go build .\n\n' > Makefile
printf 'build-release:\n\tgo tool templ generate && go generate && go build -ldflags="-w -s"\n\n' >> Makefile
printf 'clean:\n\tgo clean\n' >> Makefile


curl https://cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js >> ./static/alpine.js

cat > static/input.css <<EOF
@import "tailwindcss";
EOF

cat > .gitignore <<EOF
$project_name
EOF

cat > components/layout.templ <<EOF
package components

templ Layout(children ...templ.Component) {
  <!DOCTYPE html>
  <html lang="en">
  <head>
    <meta charset="UTF-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1"/>
    <title>PasteBin</title>
    <link href="/static/styles.css" rel="stylesheet"/>
    <script defer src="/static/alpine.js"></script>
  </head>
  <body>
    for _, child := range children {
      @child
    }
  </body>
  </html>
}
EOF

cat > components/index.templ <<EOF
package components

templ Index() {
  @Layout(main())
}

templ main(){
  <h1 class="text-2xl bg-red-500">Welcome to GYATT</h1>
}
EOF

cat > main.go <<EOF
package main
import (
  "context"
  "$project_name/components"
  "embed"
  "flag"
  "log"
  "net/http"
)

//go:embed static/*
var static embed.FS

//go:generate tailwindcss -i static/input.css -o static/styles.css -m

func main(){
  port := flag.String("port", "8080", "The port to start $project_name on")
  flag.Parse()
  mux := http.NewServeMux()

  mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
    ctx := context.Background()
    components.Index().Render(ctx, w)
  })
  mux.HandleFunc("GET /static/", http.FileServer(http.FS(static)).ServeHTTP)

  
  server := http.Server{
    Addr: ":"+ *port,
    Handler: mux,
  }
  log.Printf("$project_name launched on port %s\n", *port)
  if err := server.ListenAndServe(); err!= nil {
    panic(err)
  }
}

EOF
