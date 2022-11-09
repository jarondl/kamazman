package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
	"strings"
)

// Math is currently broken, so no sqrt or pow.
const getTimes = `
  WITH stopA as (
    SELECT * FROM arrivals WHERE route_id=@route_id and StopCode=@StopCodeA),
       stopB as (
    select * from arrivals where route_id=@route_id and StopCode=@StopCodeB)
  SELECT month,
         DayOfWeek,
         HourSourceTime,
         stopB.timeCumSum_mean - stopA.timeCumSum_mean as TimeDiff
  from stopA join stopB using (month, route_id, DayOfWeek, HourSourceTime)
  ORDER BY month, DayOfWeek, HourSourceTime;
`

const getStops = `with stopA as (SELECT DISTINCT StopSequence_Rishui, StopCode from arrivals where route_id = @route_id)
   select StopA.StopCode, stop_desc from StopA join unique_stops on (StopA.StopCode=unique_stops.stop_code) ORDER BY StopSequence_Rishui;`

const getRoutes = `select route_id, printf("%s - %s - %s", route_short_name, agency_name, route_long_name) from routes join agencies using (agency_id) where route_short_name=@route_short_name
AND route_id in applicable_routes
ORDER BY agency_name, route_id
;`

type Routes struct {
	RouteIds []string
	Names    []string
}

type Stops struct {
	StopCodes []string
	StopDescs []string
}

type Arrivals struct {
	Month          []string
	DayOfWeek      []int
	HourSourceTime []int
	TimeDiff       []float32
}

func returnJson(data any, w http.ResponseWriter, r *http.Request) {
	js, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

type Server struct {
	mux *http.ServeMux
	db  *sql.DB
}

func (s *Server) routes(w http.ResponseWriter, r *http.Request) {
	p := strings.Split(r.URL.Path, "/")
	if len(p) != 4 || p[3] == "" {
		http.NotFound(w, r)
		return
	}
	rows, err := s.db.Query(getRoutes, p[3])
	if err != nil {
		http.Error(w, "yikes", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	routes := Routes{}
	for rows.Next() {
		var routeId string
		var name string
		if err := rows.Scan(&routeId, &name); err != nil {
			log.Fatal("yiiikes")
		}
		routes.Names = append(routes.Names, name)
		routes.RouteIds = append(routes.RouteIds, routeId)
	}
	if len(routes.Names) == 0 {
		http.NotFound(w, r)
		return
	}
	returnJson(routes, w, r)
}

func (s *Server) stops(w http.ResponseWriter, r *http.Request) {
	p := strings.Split(r.URL.Path, "/")
	if len(p) != 4 || p[3] == "" {
		http.NotFound(w, r)
		return
	}
	rows, err := s.db.Query(getStops, p[3])
	if err != nil {
		http.Error(w, "yikes", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	stops := Stops{}
	for rows.Next() {
		var stopCode, stopDesc string
		if err := rows.Scan(&stopCode, &stopDesc); err != nil {
			panic("yiiikes")
		}
		stops.StopCodes = append(stops.StopCodes, stopCode)
		stops.StopDescs = append(stops.StopDescs, stopDesc)
	}
	if len(stops.StopCodes) == 0 {
		http.NotFound(w, r)
		return
	}
	returnJson(stops, w, r)
}

func (s *Server) arrivals(w http.ResponseWriter, r *http.Request) {
	p := strings.Split(r.URL.Path, "/")
	if len(p) != 6 {
		http.NotFound(w, r)
		return
	}
	rows, err := s.db.Query(getTimes,
		sql.Named("route_id", p[3]),
		sql.Named("StopCodeA", p[4]),
		sql.Named("StopCodeB", p[5]))
	if err != nil {
		log.Printf("sql error %w", err)
		http.Error(w, "yikes", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	arrivals := Arrivals{}
	for rows.Next() {
		var (
			month                     string
			dayOfWeek, hourSourceTime int
			timeDiff                  float32
		)

		if err := rows.Scan(&month, &dayOfWeek, &hourSourceTime, &timeDiff); err != nil {
			panic("yiiikes")
		}
		arrivals.Month = append(arrivals.Month, month)
		arrivals.DayOfWeek = append(arrivals.DayOfWeek, dayOfWeek)
		arrivals.HourSourceTime = append(arrivals.HourSourceTime, hourSourceTime)
		arrivals.TimeDiff = append(arrivals.TimeDiff, timeDiff)
	}

	if len(arrivals.DayOfWeek) == 0 {
		http.NotFound(w, r)
		return
	}
	returnJson(arrivals, w, r)

}

func NewServer() *Server {
	s := &Server{
		mux: http.NewServeMux(),
	}
	s.mux.HandleFunc("/api/routes/", s.routes)
	s.mux.HandleFunc("/api/stops/", s.stops)
	s.mux.HandleFunc("/api/arrivals/", s.arrivals)
	fs := http.FileServer(http.Dir("./static"))
	s.mux.Handle("/", fs)
	db, err := sql.Open("sqlite3", "file:arrivals.db?immutable=true")
	if err != nil {
		panic(err)
	}
	s.db = db
	return s
}

func main() {
	s := NewServer()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("listening %s", port)
	log.Fatal(http.ListenAndServe(":"+port, s.mux))
}
