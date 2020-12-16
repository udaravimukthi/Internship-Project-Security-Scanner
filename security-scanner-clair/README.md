 psql postgres://postgres:password@postgres-svc:5432/postgres?sslmode=disable

psql \c
You are now connected to database "postgres" as user "postgres".

psql \l
                                List of databases
   Name    |  Owner   | Encoding |  Collate   |   Ctype    |   Access privileges   
-----------+----------+----------+------------+------------+-----------------------
 postgres  | postgres | UTF8     | en_US.utf8 | en_US.utf8 | 

psql \du
                                   List of roles
 Role name |                         Attributes                         | Member of 
-----------+------------------------------------------------------------+-----------
 postgres  | Superuser, Create role, Create DB, Replication, Bypass RLS | {}


 psql \dt
                        List of relations
 Schema |                 Name                 | Type  |  Owner   
--------+--------------------------------------+-------+----------
 public | feature                              | table | postgres
 public | featureversion                       | table | postgres
 public | keyvalue                             | table | postgres
 public | layer                                | table | postgres
 public | layer_diff_featureversion            | table | postgres
 public | lock                                 | table | postgres
 public | namespace                            | table | postgres
 public | schema_migrations                    | table | postgres
 public | vulnerability                        | table | postgres
 public | vulnerability_affects_featureversion | table | postgres
 public | vulnerability_fixedin_feature        | table | postgres
 public | vulnerability_notification           | table | postgres

 psql \d+ vulnerability
                                                            Table "public.vulnerability"
    Column    |           Type           | Collation | Nullable |                  Default                  | Storage  | Stats target | Description 
--------------+--------------------------+-----------+----------+-------------------------------------------+----------+--------------+-------------
 id           | integer                  |           | not null | nextval('vulnerability_id_seq'::regclass) | plain    |              | 
 namespace_id | integer                  |           | not null |                                           | plain    |              | 
 name         | character varying(128)   |           | not null |                                           | extended |              | 

  select * from vulnerability