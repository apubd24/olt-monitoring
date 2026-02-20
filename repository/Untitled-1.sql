
Table name: olt_onu_distance

id  |         collected_at          | device_id | onu_index | if_descr | onu_distance 
------+-------------------------------+-----------+-----------+----------+--------------
    1 | 2026-02-01 04:08:10.113224+00 |         2 |       161 | N/A      |         1331
    2 | 2026-02-01 04:08:10.12432+00  |         1 |        34 | N/A      |         2716
    3 | 2026-02-01 04:08:10.12617+00  |         2 |       164 | N/A      |         5278
    4 | 2026-02-01 04:08:10.136327+00 |         2

Table name:  device_interfaces;

    id   | device_id |     if_descr     |         created_at         |         updated_at         | if_index 
-------+-----------+------------------+----------------------------+----------------------------+----------
    71 |         2 | GPON0/5          | 2026-02-01 03:20:02.956786 | 2026-02-02 06:51:15.218185 |      114
    24 |         2 | GigaEthernet0/0  | 2026-02-01 03:19:59.908034 | 2026-02-02 06:51:15.187841 |      109
     8 |         2 | GigaEthernet0/5  | 2026-02-01 03:19:59.856168 | 2026-02-02 06:51:15.137935 |      101
    21 |         1 | EPON0/1          | 2026-02-01 03:19:59.900437 | 2026-02-02 06:51:15.180398 |       15
    69 |         2 | GPON0/3          | 2026-02-01 03:20:02.943057 | 2026-02-02 06:51:15.20513  |      112



Table name:  devices;
 device_id | customer_id |  customer_name  |    device_name    | device_vendor | device_type |  ip_address   | snmp_community | snmp_version | is_active 
-----------+-------------+-----------------+-------------------+---------------+-------------+---------------+----------------+--------------+-----------
         1 |           1 | Desh Online Ltd | BDCOM-OLT-EPON-01 | BDCOM         | EPON        | 172.16.133.14 | TestComm       | v1           | t
         2 |           2 | RaceOnline Ltd  | BDCOM-OLT-GPON-01 | BDCOM         | GPON        | 172.16.4.214  | EarthComm      | v1           | t
         3 |           3 | ORBIT Internet  | VSLOL-EPON-01     | VSOL          | EPON        | 172.18.48.142 | EarthComm      | v1           | t




Using this three tabe output I want to creste teable like below for onu_distance. 

From Table: olt_onu_distance for the onu_distance, match with collam onu_index, device_id  with table device_interfaces collam if_index, device_id and display like below:

id  |   collected_at          | device_id |  customer_name  |    device_name    | device_vendor | device_type |   device_ip   |  sys_name   | onu_index | if_descr | onu_distance |  if_descr | 







✅ PostgreSQL SELECT Query


SELECT
    dsn.id,
    dsn.collected_at,
    dsn.device_id,

    dev.customer_name,
    dev.device_name,
    dev.device_vendor,
    dev.device_type,
    dev.ip_address AS device_ip,

    dsn.sys_name,
    dsn.onu_index,

    di.if_descr,
    dsn.onu_distance

FROM olt_onu_distance dsn

JOIN device_interfaces di
    ON dsn.device_id = di.device_id
   AND dsn.onu_index = di.if_index

JOIN devices dev
    ON dsn.device_id = dev.device_id

ORDER BY dsn.collected_at DESC;








🚀 (Recommended) Create a VIEW

So you don’t have to write this query again:

CREATE VIEW v_onu_distance_details AS
SELECT
    dsn.id,
    dsn.collected_at,
    dsn.device_id,

    dev.customer_name,
    dev.device_name,
    dev.device_vendor,
    dev.device_type,
    dev.ip_address AS device_ip,

    dsn.sys_name,
    dsn.onu_index,

    di.if_descr,
    dsn.onu_distance

FROM olt_onu_distance dsn
JOIN device_interfaces di
    ON dsn.device_id = di.device_id
   AND dsn.onu_index = di.if_index
JOIN devices dev
    ON dsn.device_id = dev.device_id;
========================================




