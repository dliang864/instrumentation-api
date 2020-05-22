-- Users and Roles for HHD Instrumentation Webapp

-- User instrumentation_user
-- Note: Substitute real password for 'password'
CREATE USER instrumentation_user WITH ENCRYPTED PASSWORD 'password';

-- Role instrumentation_reader
-- Tables specific to instrumentation app
CREATE ROLE instrumentation_reader;
GRANT SELECT ON
    instrument,
    instrument_group,
    instrument_group_instruments,
    instrument_type,
    parameter,
    timeseries,
    timeseries_measurement,
    unit
TO instrumentation_reader;

-- Role instrumentation_writer
-- Tables specific to instrumentation app
CREATE ROLE instrumentation_writer;
GRANT INSERT,UPDATE,DELETE ON 
    instrument,
    instrument_group,
    instrument_group_instruments,
    instrument_type,
    timeseries,
    timeseries_measurement
TO instrumentation_writer;

-- Role postgis_reader
CREATE ROLE postgis_reader;
GRANT SELECT ON geometry_columns TO postgis_reader;
GRANT SELECT ON geography_columns TO postgis_reader;
GRANT SELECT ON spatial_ref_sys TO postgis_reader;

-- Grant Permissions to instrument_user
GRANT postgis_reader TO instrumentation_user;
GRANT instrumentation_reader TO instrumentation_user;
GRANT instrumentation_writer TO instrumentation_user;