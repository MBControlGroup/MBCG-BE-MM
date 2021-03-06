package service

import (
    "net/http"

    "github.com/codegangsta/negroni"
    "github.com/gorilla/mux"
    "github.com/unrolled/render"
)

// NewServer configures and returns a Server.
func NewServer() *negroni.Negroni {

    formatter := render.New(render.Options{
        IndentJSON: true,
    })

    n := negroni.Classic()
    mx := mux.NewRouter()

    initRoutes(mx, formatter)

    n.UseHandler(mx)
    return n
}

func initRoutes(mx *mux.Router, formatter *render.Render) {
    //mx.HandleFunc("/hello/{id}", testHandler(formatter)).Methods("GET")
    mx.HandleFunc("/BMmanage/getAllMes", getAllBMsHandler(formatter)).Methods("POST")
    mx.HandleFunc("/BMmanage/createMes", sendBMsHandler(formatter)).Methods("POST")
    mx.HandleFunc("/BMmanage/{bm_id}", getBMHandler(formatter)).Methods("POST")
    mx.HandleFunc("/BMmanage/getAllMes", preOptionHandler(formatter)).Methods("OPTIONS")
    mx.HandleFunc("/BMmanage/createMes", preOptionHandler(formatter)).Methods("OPTIONS")
    mx.HandleFunc("/BMmanage/{bm_id}", preOptionHandler(formatter)).Methods("OPTIONS")
    //mx.HandleFunc("/service/userinfo", getUserInfoHandler(formatter))

}

func testHandler(formatter *render.Render) http.HandlerFunc {

    return func(w http.ResponseWriter, req *http.Request) {
        vars := mux.Vars(req)
        id := vars["id"]
        formatter.JSON(w, http.StatusOK, struct{ Test string }{"Hello " + id})
    }
}

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}