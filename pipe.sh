#!/usr/bin/bash

set -euo pipefail

sqlite3 arrivals_import.db "PRAGMA journal_mode=WAL"

sqlite3 arrivals_import.db ".read import_to_sqlite.sql"

csvcut -c "stop_id,stop_code,stop_desc,parent_station" stops.txt | \
  sqlite3 arrivals_import.db ".import --csv --skip 1 '|cat -' stops"

csvcut -c "agency_id,agency_name" agency.txt | \
  sqlite3 arrivals_import.db ".import --csv --skip 1 '|cat -' agencies"

csvcut -c "route_id,agency_id,route_short_name,route_long_name,route_desc" routes.txt | \
  sqlite3 arrivals_import.db ".import --csv --skip 1 '|cat -' routes"


sqlite3 arrivals_import.db -cmd "PRAGMA synchronous=OFF" \
  -cmd '.separator "|"' -cmd '.timer on' '.import  --skip 1 arrivaltostationdayandhours.csv arrivals_import'

sqlite3 arrivals_import.db ".read create_arrivals.sql"

sqlite3 arrivals_import.db ".read unique_stops.sql"
sqlite3 arrivals_import.db ".read create_index.sql"

sqlite3 arrivals_import.db 'VACUUM INTO "arrivals.db"'
sqlite3 arrivals.db "PRAGMA journal_mode=delete"
