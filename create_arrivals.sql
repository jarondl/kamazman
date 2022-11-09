.timer on
create table arrivals as select route_id, StopCode, StopSequence_Rishui, month, DayOfWeek, HourSourceTime, timeCumSum_mean from arrivals_import order by route_id, StopSequence_Rishui, month, DayOfWeek, HourSourceTime;

DROP TABLE arrivals_import;
