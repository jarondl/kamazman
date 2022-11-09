
CREATE TABLE IF NOT EXISTS stops(
  stop_id integer not null primary key,
  stop_code text,
  stop_desc text,
  parent_station integer);

CREATE TABLE IF NOT EXISTS agencies(
  agency_id integer not null primary key,
  agency_name text);

CREATE TABLE IF NOT EXISTS routes(
  route_id text not null primary key,
  agency_id integer not null,
  route_short_name text,
  route_long_name text,
  route_desc text,
  FOREIGN KEY(agency_id) references agencies(agency_id));

CREATE TABLE IF NOT EXISTS arrivals_import(
  month text not null,
  route_id text not null,
  DayOfWeek integer not null,
  HourSourceTime integer not null,
  StopSequence_Rishui integer not null,
  StopCode text not null,
  count_common integer not null,
  timeCumSum_mean real,
  timeCumSum_std real,
  distCumSum_mean real);
