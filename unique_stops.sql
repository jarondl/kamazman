CREATE  TABLE unique_stops AS select * from stops where 
  stop_id in (select distinct coalesce(nullif(parent_station, ""), stop_id) from stops) order by stop_code asc;
