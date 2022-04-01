
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
   address_line1 VARCHAR (255) NOT NULL default 'Not Specified',
   address_line2 VARCHAR (255),
   city VARCHAR (100) NOT NULL default 'Not Specified',
   state VARCHAR (100) NOT NULL default 'Not Specified',
   zipcode VARCHAR (20) NOT NULL default 'Not Specified',
   head_quarters BOOLEAN NOT NULL default TRUE,
   created_datetime timestamp,
   updated_datetime timestamp
);

INSERT INTO tenant(tenant_name)
VALUES ('YardHome Residency'),
       ('WorkRoom Builders'),
       ('DoorCell HBD'),
       ('ComHouse'),
       ('Hallobby'),
       ('HomeFit'),
       ('ArchiFix'),
       ('Benchmate'); 

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