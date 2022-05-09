
-- DROP TRIGGER IF EXISTS update_lastmodified_modtime ON tenant;
-- DROP FUNCTION IF EXISTS upd_lastmodified_timestamp;

             --- Function for row created timestamp ---
CREATE OR REPLACE FUNCTION row_created_timestamp()
        RETURNS TRIGGER AS '
  BEGIN
    NEW.created_datetime = NOW();
    NEW.updated_datetime = NOW();
    RETURN NEW;
  END;
' LANGUAGE 'plpgsql';

             --- Function for row updated timestamp ---
CREATE OR REPLACE FUNCTION row_updated_timestamp()
        RETURNS TRIGGER AS '
  BEGIN
    NEW.updated_datetime = NOW();
    RETURN NEW;
  END;
' LANGUAGE 'plpgsql';

             --- Create Tenant Table ---

CREATE TABLE tenant(
   tenant_id serial PRIMARY KEY,
   tenant_name VARCHAR (255) NOT NULL,
   logo bytea,
   status VARCHAR (2),
	start_date DATE NOT NULL,
	end_date DATE default '9999-12-31',
	tenant_type VARCHAR (50),
   created_datetime timestamp,
   updated_datetime timestamp
);


             --- Trigger for row created timestamp ---

CREATE OR REPLACE TRIGGER row_updated_timestamp -- FUNCTION_NAME
BEFORE UPDATE -- OR INSERT OR DELETE (OPERATIONS)
  ON TENANT  -- TABLE_NAME
  FOR EACH ROW EXECUTE PROCEDURE row_updated_timestamp(); -- FUNCTION_NAME

             --- Trigger for row updated timestamp ---

CREATE OR REPLACE TRIGGER row_created_timestamp 
BEFORE INSERT
  ON TENANT 
  FOR EACH ROW EXECUTE PROCEDURE row_created_timestamp();





CREATE TABLE business_unit(
   business_unit INT GENERATED ALWAYS AS IDENTITY,
   tenant_id INT,
   business_unit_name VARCHAR (255) NOT NULL,
   address_line1 VARCHAR (255) NOT NULL default 'Not Specified',
   address_line2 VARCHAR (255),
   city VARCHAR (100) NOT NULL default 'Not Specified',
   state VARCHAR (100) NOT NULL default 'Not Specified',
   zipcode VARCHAR (20) NOT NULL default 'Not Specified',
   head_quarters BOOLEAN NOT NULL default TRUE,
   created_datetime timestamp,
   updated_datetime timestamp,
   PRIMARY KEY(business_unit, tenant_id),
   CONSTRAINT fk_tenant
      FOREIGN KEY(tenant_id) REFERENCES tenant(tenant_id)
);


CREATE OR REPLACE TRIGGER bu_updated_timestamp -- FUNCTION_NAME
BEFORE UPDATE -- OR INSERT OR DELETE (OPERATIONS)
  ON business_unit  -- TABLE_NAME
  FOR EACH ROW EXECUTE PROCEDURE row_updated_timestamp(); -- FUNCTION_NAME

             --- Trigger for row updated timestamp ---

CREATE OR REPLACE TRIGGER bu_create_timestamp 
BEFORE INSERT
  ON business_unit 
  FOR EACH ROW EXECUTE PROCEDURE row_created_timestamp();


CREATE TABLE permission(
   tenant_id INT,
   permission_id INT,
   permission_name VARCHAR (255) NOT NULL,
   permission_active boolean NOT NULL default TRUE, 
   created_datetime timestamp,
   updated_datetime timestamp,
   PRIMARY KEY(tenant_id, permission_id),
   CONSTRAINT fk_customer
      FOREIGN KEY(tenant_id) 
	  REFERENCES tenant(tenant_id)
);


CREATE OR REPLACE TRIGGER prm_updated_timestamp -- FUNCTION_NAME
BEFORE UPDATE -- OR INSERT OR DELETE (OPERATIONS)
  ON permission  -- TABLE_NAME
  FOR EACH ROW EXECUTE PROCEDURE row_updated_timestamp(); -- FUNCTION_NAME

             --- Trigger for row updated timestamp ---

CREATE OR REPLACE TRIGGER prm_create_timestamp 
BEFORE INSERT
  ON permission 
  FOR EACH ROW EXECUTE PROCEDURE row_created_timestamp();

insert into permission (tenant_id, permission_id, permission_name) values (1,2,'BU');

CREATE TABLE permission_type(
   permission_type_id SERIAL PRIMARY KEY,
   permission_type VARCHAR (150) NOT NULL,
   permission_active boolean NOT NULL default TRUE, 
   created_datetime timestamp,
   updated_datetime timestamp
);

CREATE OR REPLACE TRIGGER prm_typ_upd_timestamp -- FUNCTION_NAME
BEFORE UPDATE -- OR INSERT OR DELETE (OPERATIONS)
  ON permission_type  -- TABLE_NAME
  FOR EACH ROW EXECUTE PROCEDURE row_updated_timestamp(); -- FUNCTION_NAME

             --- Trigger for row updated timestamp ---

CREATE OR REPLACE TRIGGER prm_typ_crt_timestamp 
BEFORE INSERT
  ON permission_type 
  FOR EACH ROW EXECUTE PROCEDURE row_created_timestamp();


CREATE TABLE tenant_projects(
   project_id serial PRIMARY KEY,
   tenant_id numeric,
   project_name VARCHAR (255) NOT NULL,
   location_id numeric,
   service_id numeric,
   created_datetime timestamp,
   updated_datetime timestamp
);

CREATE OR REPLACE TRIGGER row_updated_timestamp -- FUNCTION_NAME
BEFORE UPDATE -- OR INSERT OR DELETE (OPERATIONS)
  ON tenant_projects  -- TABLE_NAME
  FOR EACH ROW EXECUTE PROCEDURE row_updated_timestamp(); -- FUNCTION_NAME

             --- Trigger for row updated timestamp ---

CREATE OR REPLACE TRIGGER row_created_timestamp 
BEFORE INSERT
  ON tenant_projects 
  FOR EACH ROW EXECUTE PROCEDURE row_created_timestamp();

  CREATE TABLE services(
   service_id serial PRIMARY KEY,
   service_name VARCHAR (100) NOT NULL,
   status boolean,
   icon VARCHAR (10) NOT NULL,
   color VARCHAR (10) NOT NULL,
   bg_image bytea,
   route_path VARCHAR (10) NOT NULL,
   created_datetime timestamp,
   updated_datetime timestamp
);

CREATE OR REPLACE TRIGGER row_updated_timestamp -- FUNCTION_NAME
BEFORE UPDATE -- OR INSERT OR DELETE (OPERATIONS)
  ON services  -- TABLE_NAME
  FOR EACH ROW EXECUTE PROCEDURE row_updated_timestamp(); -- FUNCTION_NAME

             --- Trigger for row updated timestamp ---

CREATE OR REPLACE TRIGGER row_created_timestamp 
BEFORE INSERT
  ON services 
  FOR EACH ROW EXECUTE PROCEDURE row_created_timestamp();




  CREATE TABLE tenant_status(
   id serial PRIMARY KEY,
	code VARCHAR (25) NOT NULL,
   descr VARCHAR (255) NOT NULL, 
	start_date DATE NOT NULL,
	end_date DATE default '9999-12-31',
   created_datetime timestamp,
   updated_datetime timestamp
);

             --- Trigger for row created timestamp ---

CREATE OR REPLACE TRIGGER row_updated_timestamp -- FUNCTION_NAME
BEFORE UPDATE -- OR INSERT OR DELETE (OPERATIONS)
  ON status_master  -- TABLE_NAME
  FOR EACH ROW EXECUTE PROCEDURE row_updated_timestamp(); -- FUNCTION_NAME

             --- Trigger for row updated timestamp ---

CREATE OR REPLACE TRIGGER row_created_timestamp 
BEFORE INSERT
  ON status_master 
  FOR EACH ROW EXECUTE PROCEDURE row_created_timestamp();





DROP TABLE IF EXISTS country_master;

CREATE TABLE country_master(
   id serial,
   alpha_2_code VARCHAR (2) NOT NULL UNIQUE,
   alpha_3_code VARCHAR (3) NOT NULL UNIQUE,
   un_code VARCHAR (3) NOT NULL,
   country VARCHAR (100) NOT NULL,
   timezone VARCHAR (5) NOT NULL,
   default_currency  VARCHAR (5),
   created_datetime timestamp WITH TIME ZONE default current_timestamp,
   updated_datetime timestamp WITH TIME ZONE
);

select * from country_master;
--------------------------------------------------------------------------------
DROP TABLE IF EXISTS state_master;

CREATE TABLE state_master(
   code VARCHAR (10) primary key,
   country_code VARCHAR (3) NOT NULL REFERENCES country_master(country_code),
   state VARCHAR (100) NOT NULL,
   created_datetime timestamp WITH TIME ZONE default current_timestamp,
   updated_datetime timestamp WITH TIME ZONE
);

select * from state_master;
--------------------------------------------------------------------------------
DROP TABLE IF EXISTS district_master;

CREATE TABLE district_master(
   id serial primary key,
   state_code VARCHAR (10) NOT NULL REFERENCES state_master(code),
   district_code VARCHAR (10),
	district VARCHAR (100) NOT NULL,
   created_datetime timestamp WITH TIME ZONE default current_timestamp,
   updated_datetime timestamp WITH TIME ZONE
);

select * from district_master;
--------------------------------------------------------------------------------

DROP TABLE IF EXISTS zipcode_master;

CREATE TABLE zipcode_master(
   id serial primary key,
   district_id integer NOT NULL REFERENCES district_master(id),
   zipcode VARCHAR (15) NOT NULL,
   location VARCHAR (100) NOT NULL,
   created_datetime timestamp WITH TIME ZONE default current_timestamp,
   updated_datetime timestamp WITH TIME ZONE
);

select * from zipcode_master where zipcode = '641029';

--------------------------------------------------------------------------------
DROP TABLE IF EXISTS location;

CREATE TABLE location(
   id serial primary key,
   zipcode_id integer NOT NULL REFERENCES zipcode_master(id),
   location_code VARCHAR (15) NOT NULL,
   descr VARCHAR (100) NOT NULL,
	Address_line1 VARCHAR (100) NOT NULL,
	Address_line2 VARCHAR (100),
	land_mark VARCHAR (100),
	latitude float,
	longitude float,
   created_datetime timestamp WITH TIME ZONE default current_timestamp,
   updated_datetime timestamp WITH TIME ZONE
);

select * from location;
--------------------------------------------------------------------------------
DROP TABLE IF EXISTS client_master;

CREATE TABLE client_master(
   id  SERIAL primary key,
   client_code VARCHAR (15) NOT NULL,
   name VARCHAR (200) NOT NULL,
   location_id integer NOT NULL REFERENCES location(id),
   display_name VARCHAR (200) NOT NULL,
   descr VARCHAR (100) NOT NULL,
   status VARCHAR (5) NOT NULL,
   client_nid VARCHAR (50) NOT NULL,
	tax_ref VARCHAR (100),
   start_date date,
   end_date date default '9999-12-12'::DATE,
   timezone VARCHAR (5) NOT NULL,
   currency VARCHAR (5) NOT NULL,
   website VARCHAR (200),
   logo VARCHAR (250),
   created_datetime timestamp WITH TIME ZONE default current_timestamp,
   updated_datetime timestamp WITH TIME ZONE
);

select * from client_master;
--------------------------------------------------------------------------------
DROP TABLE IF EXISTS users;

CREATE TABLE users(
   id integer NOT NULL REFERENCES client_master(id) primary key,
   username VARCHAR (50) NOT NULL,
   password VARCHAR (50) NOT NULL,
   reset boolean,
   created_datetime timestamp WITH TIME ZONE default current_timestamp,
   updated_datetime timestamp WITH TIME ZONE
);

select * from users;

--------------------------------------------------------------------------------
DROP TABLE IF EXISTS service_master;

CREATE TABLE service_master(
   id serial primary key,
   service_code VARCHAR (5) NOT NULL,   
   service_name VARCHAR (100) NOT NULL,
   start_date date,
   end_date date default '9999-12-12'::DATE,
   created_datetime timestamp WITH TIME ZONE default current_timestamp,
   updated_datetime timestamp WITH TIME ZONE
);

select * from service_master;

--------------------------------------------------------------------------------
DROP TABLE IF EXISTS client_project;

CREATE TABLE client_project(
   id serial primary key,
   client_id integer NOT NULL REFERENCES client_master(id),
   service_id integer NOT NULL REFERENCES service_master(id),
   start_date date,
   end_date date default '9999-12-12'::DATE,
   created_datetime timestamp WITH TIME ZONE default current_timestamp,
   updated_datetime timestamp WITH TIME ZONE
);

select * from client_project;
--------------------------------------------------------------------------
DROP TABLE IF EXISTS user_profile;

CREATE TABLE user_profile(
   id serial primary key,
   first_name VARCHAR (50) NOT NULL,
   middle_name VARCHAR (50),
   last_name VARCHAR (50) NOT NULL,
   status VARCHAR (5) NOT NULL,
   email VARCHAR (50),
   phone VARCHAR (50),
   project_id integer NOT NULL REFERENCES client_project(id), 
   start_date date,
   end_date date,
   timezone VARCHAR (5) NOT NULL,
   currency VARCHAR (5) NOT NULL,
   created_datetime timestamp WITH TIME ZONE default current_timestamp,
   updated_datetime timestamp WITH TIME ZONE
);

select * from user_profile;


DROP TABLE IF EXISTS wm_devices;

CREATE TABLE wm_devices(
   id serial primary key,
   device_number VARCHAR (50) NOT NULL,
   device_descr VARCHAR (50),
   install_date date,
   status VARCHAR (50) NOT NULL,
   project_id integer NOT NULL REFERENCES client_project(id),
   location_id integer NOT NULL REFERENCES location(id),
   created_datetime timestamp WITH TIME ZONE default current_timestamp,
   updated_datetime timestamp WITH TIME ZONE
);

select * from wm_devices;


DROP TABLE IF EXISTS wm_meter_models;

CREATE TABLE wm_meter_models(
   id serial primary key,
   model_name VARCHAR (20) NOT NULL,
   model_descr VARCHAR (150) NOT NULL,
   version VARCHAR (10),
   category VARCHAR (10),
   diameter VARCHAR (10),
   consumption_type VARCHAR (10),
   status VARCHAR (50) NOT NULL,
   data_sheet_link VARCHAR (100),
   manual_link VARCHAR (100),
   manufacturer_id integer,
   created_datetime timestamp WITH TIME ZONE default current_timestamp,
   updated_datetime timestamp WITH TIME ZONE
);

select * from wm_meter_models;


DROP TABLE IF EXISTS wm_water_meters;

CREATE TABLE wm_water_meters(
   meter_id serial primary key,
   model_id integer NOT NULL REFERENCES wm_meter_models(id), 
   project_id integer NOT NULL REFERENCES client_project(id),
   install_dttm timestamp WITH TIME ZONE,
   transmitter_no integer NOT NULL REFERENCES wm_devices(id),
   status VARCHAR (10) NOT NULL,
   initial_reading numeric,
   asset_image_id varchar(200),
   warrenty_id integer,
   warrenty_notes varchar(200),
   warrenty_exp_dt date,
   created_datetime timestamp WITH TIME ZONE default current_timestamp,
   updated_datetime timestamp WITH TIME ZONE
);

select * from wm_water_meters;



DROP TABLE IF EXISTS wm_meter_group_info;

CREATE TABLE wm_meter_group_info(
   id serial primary key,
   group_code VARCHAR (10) NOT NULL,
   group_name VARCHAR (100) NOT NULL,
   descr VARCHAR (250) NOT NULL,
   status VARCHAR (5) NOT NULL,
   project_id integer NOT NULL REFERENCES client_project(id),
   created_datetime timestamp WITH TIME ZONE default current_timestamp,
   updated_datetime timestamp WITH TIME ZONE
);

select * from wm_meter_group_info;

DROP TABLE IF EXISTS wm_meter_groups;

CREATE TABLE wm_meter_groups(
   group_id integer NOT NULL REFERENCES wm_meter_group_info(id),
   meter_id integer NOT NULL REFERENCES wm_water_meters(meter_id),
   created_datetime timestamp WITH TIME ZONE default current_timestamp,
   updated_datetime timestamp WITH TIME ZONE,
	primary key (group_id, meter_id)
);

select * from wm_meter_groups;

DROP TABLE IF EXISTS wm_meter_consumption;

CREATE TABLE wm_meter_consumption(
   id serial primary key,
   meter_id integer NOT NULL REFERENCES wm_water_meters(meter_id),
   alert VARCHAR (10) NOT NULL,
   reading float NOT NULL,
   consumption float default 0,
   created_datetime timestamp WITH TIME ZONE default current_timestamp
);

select * from wm_meter_consumption;


DROP TABLE IF EXISTS wm_alerts;

CREATE TABLE wm_alerts(
   alert_id serial primary key,
   meter_id integer NOT NULL REFERENCES wm_water_meters(meter_id),
   alert_code VARCHAR (10) NOT NULL,
   status VARCHAR (10) NOT NULL,
   created_datetime timestamp WITH TIME ZONE default current_timestamp
);

select * from wm_alerts;



DROP TABLE IF EXISTS wm_notification;

CREATE TABLE wm_notification(
   id serial primary key,
   notification_type VARCHAR (10) NOT NULL,
   reference_id VARCHAR (20) NOT NULL,
   user_id integer NOT NULL REFERENCES user_profile(id),
   message VARCHAR (10) NOT NULL,
   status VARCHAR (10) NOT NULL,
   created_datetime timestamp WITH TIME ZONE default current_timestamp
);

select * from wm_notification;


DROP TABLE IF EXISTS wm_wallet;

CREATE TABLE wm_wallet(
   wallet_id VARCHAR (20) primary key,
   user_id integer NOT NULL REFERENCES user_profile(id),
   balance float NOT NULL default 0,
   status VARCHAR (10) NOT NULL,
   created_datetime timestamp WITH TIME ZONE default current_timestamp,
   updated_datetime timestamp WITH TIME ZONE
);

select * from wm_wallet;