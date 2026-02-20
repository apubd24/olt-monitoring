Final Project Structure (Clean & Scalable)
snmp-onu-monitor/
│
├── go.mod
├── go.sum
│
├── cmd/
│   └── collector/
│       └── main.go          # Application entry point
│
├── config/
│   └── db.go                # PostgreSQL connection
│
├── models/
│   └── device.go            # Device struct
│
├── repository/
│   ├── device_repo.go       # Load devices from DB
│   └── metrics_repo.go      # Insert SNMP metrics
│
├── snmp/
│   ├── oids.go              # All SNMP OID definitions
│   └── collector.go         # SNMP polling logic
│
└── README.md


=================
1. Interface + Description 
2. Distance
3 onu TX RX Power

=================


CREATE TABLE olt_onu_distancess (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    collected_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    Device_id   INT,
    if_index     TEXT,
    onu_distance  INT
);



 CREATE TABLE huawei_onu_tx_rx_powers (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    collected_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    device_id     INT,
    onu_ifindex     TEXT,
    onuRx_power      NUMERIC(6,2),
    onuTx_power      NUMERIC(6,2)
);


 CREATE TABLE olt_onu_status (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    collected_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    device_id     INT,
    onu_index     INT,
    onu_status    INT
);

// Create Interface
CREATE TABLE device_interface_descs (
    id BIGSERIAL PRIMARY KEY,
    device_id BIGINT NOT NULL,

    if_descr TEXT NOT NULL,
    if_index TEXT NOT NULL,     -- ✅ FIXED
    if_alias TEXT,

    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),

    UNIQUE (device_id, if_index)
);



==================Devices Table (PostgreSQL)===================

-- 1. Create custom ENUM types first

CREATE TYPE vendor_type AS ENUM ('BDCOM', 'VSOL', 'HUAWEI', 'ZTE', 'FIBERHOME');
CREATE TYPE device_category AS ENUM ('ROUTER', 'SWITCH', 'FIREWALL', 'OLT', 'SERVER');
CREATE TYPE pon_type AS ENUM ('EPON', 'GPON', 'XPON', 'L2', 'L3');
CREATE TYPE snmp_ver AS ENUM ('v1', 'v2', 'V3');

-- 2. Create the table
CREATE TABLE devices (
    device_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    customer_name VARCHAR(100) NOT NULL,
    device_name VARCHAR(100) NOT NULL,
    device_vendor vendor_type NOT NULL, -- Uses the ENUM created above
    device_category VARCHAR(100) NOT NULL, -- Uses the ENUM created above
    device_type pon_type NOT NULL,      -- Uses the ENUM created above
    ip_address VARCHAR(45) NOT NULL UNIQUE,
    snmp_community VARCHAR(100) NOT NULL,
    snmp_version snmp_ver NOT NULL DEFAULT 'v1',
    is_active BOOLEAN NOT NULL DEFAULT TRUE -- Uses native Boolean
);

Insert Data Example========

INSERT INTO devices (
    customer_id, 
    customer_name, 
    device_name,
    device_vendor,
    device_category,     
    device_type, 
    ip_address, 
    snmp_community, 
    snmp_version, 
    is_active
) VALUES (
    3, 
    'ORBIT Internet', 
    'VSLOL-EPON-01',
    'VSOL', 
    'OLT',
    'EPON', 
    '172.18.48.142', 
    'EarthComm', 
    'v2', 
    TRUE
);



snmp_monitoring=> select * from devices;
 device_id | customer_id |  customer_name  |    device_name    | device_vendor | device_type |  ip_address   | snmp_community | snmp_version | is_active 
-----------+-------------+-----------------+-------------------+---------------+-------------+---------------+----------------+--------------+-----------
         1 |           1 | Desh Online Ltd | BDCOM-OLT-EPON-01 | BDCOM         | EPON        | 172.16.133.14 | TestComm       | v1           | t
         2 |           2 | RaceOnline Ltd  | BDCOM-OLT-GPON-01 | BDCOM         | GPON        | q | EarthComm      | v1           | t
         3 |           3 | ORBIT Internet  | VSLOL-EPON-01     | VSOL          | EPON        | 172.18.48.142 | EarthComm      | v1           | t
(3 rows)


         snmp_monitoring=> select * from onu_status;
snmp_monitoring=> SELECT COUNT(*) FROM onu_status WHERE device_vendor = 'BDCOM' AND if_descr = 'EPON0/4:31';
 count 
-------
    28
(1 row)

snmp_monitoring=> select * from onu_status;
snmp_monitoring=> SELECT COUNT(*) FROM onu_status WHERE device_vendor = 'BDCOM' AND if_descr = 'EPON0/4:31' AND sys_name = 'switch';
 count 
-------
     0
(1 row)

snmp_monitoring=> select * from onu_status;
snmp_monitoring=> SELECT COUNT(*) FROM onu_status WHERE device_vendor = 'BDCOM' AND if_descr = 'EPON0/4:31' AND sys_name = 'Switch';
 count 