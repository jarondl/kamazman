.timer on

CREATE INDEX IF NOT EXISTS arrivals_stops on arrivals(StopCode);

CREATE INDEX IF NOT EXISTS route_stops_idx on arrivals(route_id, StopCode, StopSequence_Rishui);

CREATE VIEW IF NOT EXISTS route_stops as select distinct route_id,  StopCode, StopSequence_Rishui from arrivals order by route_id, StopCode, StopSequence_Rishui;

CREATE TABLE applicable_routes as select distinct route_id from arrivals order by route_id;
