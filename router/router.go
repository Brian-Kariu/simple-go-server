package router

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
)

func NewRouter(lg *log.Logger, initTasks bool) http.Handler {
  mux := http.NewServeMux()
  th := NewTasksHandler(lg, initTasks)
  mux.HandleFunc("/", th.taskListPage())
  mux.HandleFunc("/new", th.taskCreatePage())
  mux.Handle("/view/", http.StripPrefix("/view", http.HandlerFunc(th.taskViewPage())))
  mux.Handle("/edit/", http.StripPrefix("/edit", http.HandlerFunc(th.taskEditPage())))

  mux.Handle("/api/tasks/", http.StripPrefix("/api/tasks", th))
  mux.HandleFunc("/api/comments", th.addComment)
  mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))

  return logMiddleware(lg, recoverMiddleware(lg, mux)) 
}

var notFoundPage = func() http.HandlerFunc {
  files := tmplLayout("./web/templates/_notFound.html")
  tmpl := template.Must(template.New("index").Funcs(defaultFuncs).ParseFiles(files...))
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    var buf bytes.Buffer
    if err := tmpl.ExecuteTemplate(&buf, "base", nil); err != nil {
      fmt.Printf("ERR: %v\n", err)
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
   }
    w.WriteHeader(http.StatusOK)
    io.Copy(w, &buf)
  })
}()