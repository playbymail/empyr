--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- foreign keys must be disabled to drop tables
PRAGMA foreign_keys = OFF;

-- foreign keys must be enabled with every database connection
PRAGMA foreign_keys = ON;

-- create the table for managing migrations
create table meta_migrations
(
    version    integer  not null unique,
    comment    text     not null,
    script     text     not null unique,
    created_at datetime not null default CURRENT_TIMESTAMP
);

create table orbit_codes
(
    code text not null,
    name text not null,
    primary key (code),
    unique (name)
);

insert into orbit_codes (code, name)
values ('NONE', 'no planet');
insert into orbit_codes (code, name)
values ('ASTR', 'asteroid belt');
insert into orbit_codes (code, name)
values ('GASG', 'gas giant');
insert into orbit_codes (code, name)
values ('TERR', 'terrestrial planet');

create table population_codes
(
    code          text    not null,
    name          text    not null,
    base_pay_rate real    not null check (base_pay_rate >= 0),
    sort_order    integer not null,
    primary key (code),
    unique (name)
);

insert into population_codes (code, name, base_pay_rate, sort_order)
values ('CNW', 'Construction Worker', 0.5000, 5);
insert into population_codes (code, name, base_pay_rate, sort_order)
values ('PLC', 'Police', 0.2500, 7);
insert into population_codes (code, name, base_pay_rate, sort_order)
values ('PRO', 'Professional', 0.3750, 3);
insert into population_codes (code, name, base_pay_rate, sort_order)
values ('SAG', 'Special Agents', 0.6250, 8);
insert into population_codes (code, name, base_pay_rate, sort_order)
values ('SLD', 'Soldier', 0.2500, 4);
insert into population_codes (code, name, base_pay_rate, sort_order)
values ('SPY', 'Spy', 0.6250, 6);
insert into population_codes (code, name, base_pay_rate, sort_order)
values ('TRN', 'Trainees', 0.1250, 9);
insert into population_codes (code, name, base_pay_rate, sort_order)
values ('UEM', 'Unemployable', 0.0000, 1);
insert into population_codes (code, name, base_pay_rate, sort_order)
values ('USK', 'Unskilled', 0.1250, 2);

create table sc_codes
(
    code       text    not null,
    name       text    not null,
    is_ship    integer not null check (is_ship in (0, 1)),
    is_surface integer not null check (is_surface in (0, 1)),
    primary key (code),
    unique (name)
);

insert into sc_codes (code, name, is_ship, is_surface)
values ('NONE', 'none', 0, 0);
insert into sc_codes (code, name, is_ship, is_surface)
values ('SHIP', 'Ship', 1, 0);
insert into sc_codes (code, name, is_ship, is_surface)
values ('COPN', 'Open Surface Colony', 1, 1);
insert into sc_codes (code, name, is_ship, is_surface)
values ('CENC', 'Enclosed Surface Colony', 0, 1);
insert into sc_codes (code, name, is_ship, is_surface)
values ('CORB', 'Orbital Colony', 0, 0);

create table unit_codes
(
    code           text    not null,
    name           text    not null,
    category       text    not null,
    is_operational integer not null check (is_operational in (0, 1)),
    is_consumable  integer not null check (is_consumable in (0, 1)),
    is_resource    integer not null check (is_resource in (0, 1)),
    aliases        text,
    primary key (code),
    unique (name)
);

insert into unit_codes (code, name, category, is_operational, is_consumable, is_resource)
values ('ANM', 'Anti-Missiles', 'Vehicles', 0, 1, 0);
insert into unit_codes (code, name, category, is_operational, is_consumable, is_resource)
values ('ASC', 'Assault Craft', 'Vehicles', 0, 0, 0);
insert into unit_codes (code, name, category, is_operational, is_consumable, is_resource)
values ('ASW', 'Assault Weapons', 'Vehicles', 0, 0, 0);
insert into unit_codes (code, name, category, is_operational, is_consumable, is_resource)
values ('AUT', 'Automation', 'Assembly', 1, 0, 0);
insert into unit_codes (code, name, category, is_operational, is_consumable, is_resource)
values ('CNGD', 'Consumer Goods', 'Consumables', 0, 1, 0);
insert into unit_codes (code, name, category, is_operational, is_consumable, is_resource)
values ('ESH', 'Energy Shields', 'Assembly', 1, 0, 0);
insert into unit_codes (code, name, category, is_operational, is_consumable, is_resource)
values ('EWP', 'Energy Weapons', 'Assembly', 1, 0, 0);
insert into unit_codes (code, name, category, is_operational, is_consumable, is_resource)
values ('FCT', 'Factories', 'Assembly', 1, 0, 0);
insert into unit_codes (code, name, category, is_operational, is_consumable, is_resource)
values ('FOOD', 'Food', 'Consumables', 0, 1, 1);
insert into unit_codes (code, name, category, is_operational, is_consumable, is_resource)
values ('FRM', 'Farms', 'Assembly', 1, 0, 0);
insert into unit_codes (code, name, category, is_operational, is_consumable, is_resource)
values ('FUEL', 'Fuel', 'Consumables', 0, 1, 1);
insert into unit_codes (code, name, category, is_operational, is_consumable, is_resource)
values ('GOLD', 'Gold', 'Consumables', 0, 1, 1);
insert into unit_codes (code, name, category, is_operational, is_consumable, is_resource)
values ('HEN', 'Hyper Engines', 'Assembly', 1, 0, 0);
insert into unit_codes (code, name, category, is_operational, is_consumable, is_resource)
values ('LAB', 'Laboratories', 'Assembly', 1, 0, 0);
insert into unit_codes (code, name, category, is_operational, is_consumable, is_resource)
values ('LFS', 'Life Supports', 'Assembly', 1, 0, 0);
insert into unit_codes (code, name, category, is_operational, is_consumable, is_resource)
values ('METS', 'Metals', 'Consumables', 0, 1, 1);
insert into unit_codes (code, name, category, is_operational, is_consumable, is_resource)
values ('MIN', 'Mines', 'Assembly', 1, 0, 0);
insert into unit_codes (code, name, category, is_operational, is_consumable, is_resource)
values ('MSL', 'Missile Launchers', 'Assembly', 1, 0, 0);
insert into unit_codes (code, name, category, is_operational, is_consumable, is_resource)
values ('MSS', 'Missiles', 'Vehicles', 0, 0, 0);
insert into unit_codes (code, name, category, is_operational, is_consumable, is_resource)
values ('MTBT', 'Military Robots', 'Bots', 0, 0, 0);
insert into unit_codes (code, name, category, is_operational, is_consumable, is_resource)
values ('MTSP', 'Military Supplies', 'Consumables', 0, 1, 0);
insert into unit_codes (code, name, category, is_operational, is_consumable, is_resource)
values ('NMTS', 'Non-Metals', 'Consumables', 0, 1, 1);
insert into unit_codes (code, name, category, is_operational, is_consumable, is_resource)
values ('PWP', 'Power Plants', 'Assembly', 1, 0, 0);
insert into unit_codes (code, name, category, is_operational, is_consumable, is_resource)
values ('RPV', 'Robot Probe Vehicles', 'Bots', 0, 1, 0);
insert into unit_codes (code, name, category, is_operational, is_consumable, is_resource)
values ('RSCH', 'Research', 'Consumables', 0, 1, 0);
insert into unit_codes (code, name, category, is_operational, is_consumable, is_resource)
values ('SEN', 'Sensors', 'Assembly', 1, 0, 0);
insert into unit_codes (code, name, category, is_operational, is_consumable, is_resource)
values ('SLS', 'Light Structure', 'Assembly', 1, 0, 0);
insert into unit_codes (code, name, category, is_operational, is_consumable, is_resource)
values ('SPD', 'Space Drives', 'Assembly', 1, 0, 0);
insert into unit_codes (code, name, category, is_operational, is_consumable, is_resource)
values ('STU', 'Structure', 'Assembly', 1, 0, 0);
insert into unit_codes (code, name, category, is_operational, is_consumable, is_resource)
values ('TPT', 'Transports', 'Vehicles', 0, 0, 0);

update unit_codes
set aliases = 'ASWP'
where code = 'ASW';
update unit_codes
set aliases = 'FACT, FCTU, FU'
where code = 'FCT';
update unit_codes
set aliases = 'FARM, FRMU'
where code = 'FRM';
update unit_codes
set aliases = 'MSS'
where code = 'MTSP';
update unit_codes
set aliases = 'MINU, MU'
where code = 'MIN';
update unit_codes
set aliases = 'STUN'
where code = 'STU';

create table games
(
    code           text     not null,
    name           text     not null,
    display_name   text     not null,
    current_turn   integer  not null default 0,
    home_system_id integer  not null default 0,
    home_star_id   integer  not null default 0,
    home_orbit_id  integer  not null default 0,
    created_at     datetime not null default CURRENT_TIMESTAMP,
    updated_at     datetime not null default CURRENT_TIMESTAMP
);

create table systems
(
    id           integer primary key autoincrement,
    x            integer not null check (x between 0 and 30),
    y            integer not null check (y between 0 and 30),
    z            integer not null check (z between 0 and 30),
    system_name  text    not null unique,
    nbr_of_stars integer not null default 0,
    unique (x, y, z)
);

create table stars
(
    id            integer primary key autoincrement,
    system_id     integer not null,
    sequence      text    not null check (length(sequence) = 1 and sequence between 'A' and 'Z'),
    star_name     text    not null unique,
    nbr_of_orbits integer not null default 0,
    unique (system_id, sequence),
    constraint fk_system_id foreign key (system_id) references systems (id)
);

create table orbits
(
    id              integer primary key autoincrement,
    system_id       integer not null,
    star_id         integer not null,
    orbit_no        integer not null check (orbit_no between 1 and 10),
    kind            text    not null,
    habitability    integer not null check (habitability between 0 and 25),
    nbr_of_deposits integer not null default 0,
    unique (star_id, orbit_no),
    constraint fk_system_id foreign key (system_id) references systems (id),
    constraint fk_star_id foreign key (star_id) references stars (id),
    constraint fk_orbit_kind_cd foreign key (kind) references orbit_codes (code)
);

-- deposits holds information about deposits of resources in an orbit.
create table deposits
(
    id         integer primary key autoincrement,
    orbit_id   integer not null,
    deposit_no integer not null check (deposit_no between 1 and 35),
    kind       text    not null,
    yield_pct  integer not null check (yield_pct between 0 and 100),
    unique (orbit_id, deposit_no),
    constraint fk_orbit_id foreign key (orbit_id) references orbits (id),
    constraint fk_deposit_kind_cd foreign key (kind) references unit_codes (code)
);

-- deposit_history holds information about deposits of resources in an orbit.
create table deposit_history
(
    deposit_id integer not null,
    effdt      integer not null,
    enddt      integer not null,
    qty        integer not null,
    primary key (deposit_id, effdt),
    constraint fk_deposit_id foreign key (deposit_id) references deposits (id)
);

-- deposits_summary holds the total quantity of each resource type
-- from all the deposits on a planet. It also holds the estimated
-- value of the resources (the log10 of the actual quantity).
create table deposits_summary
(
    orbit_id     integer not null,
    effdt        integer not null,
    enddt        integer not null,
    fuel_qty     integer not null,
    fuel_est_qty integer not null,
    gold_qty     integer not null,
    gold_est_qty integer not null,
    mets_qty     integer not null,
    mets_est_qty integer not null,
    nmts_qty     integer not null,
    nmts_est_qty integer not null,
    primary key (orbit_id, effdt),
    constraint fk_orbit_id foreign key (orbit_id) references orbits (id)
);

-- deposits_summary_pivot holds the quantity of each resource type
-- into a single row using deposit_id as the key. It also holds the
-- estimated value of the resources (the log10 of the actual quantity).
--
-- The data in this table is replaced every turn.
create table deposits_summary_pivot
(
    deposit_id   integer not null,
    effdt        integer not null,
    enddt        integer not null,
    fuel_qty     integer not null,
    fuel_est_qty integer not null,
    gold_qty     integer not null,
    gold_est_qty integer not null,
    mets_qty     integer not null,
    mets_est_qty integer not null,
    nmts_qty     integer not null,
    nmts_est_qty integer not null,
    primary key (deposit_id, effdt),
    constraint fk_deposit_id foreign key (deposit_id) references deposits (id)
);

create table empire
(
    id             integer primary key autoincrement,
    home_system_id integer not null,
    home_star_id   integer not null,
    home_orbit_id  integer not null,
    is_active      integer not null default 1 check (is_active in (0, 1)),
    constraint fk_system_id foreign key (home_system_id) references systems (id),
    constraint fk_star_id foreign key (home_star_id) references stars (id),
    constraint fk_orbit_id foreign key (home_orbit_id) references orbits (id)
);

create table empire_name
(
    empire_id integer not null,
    effdt     integer not null,
    enddt     integer not null,
    name      text    not null,
    primary key (empire_id, effdt),
    constraint fk_empire_id foreign key (empire_id) references empire (id)
);

create table empire_player
(
    empire_id integer not null,
    effdt     integer not null,
    enddt     integer not null,
    username  text    not null,
    email     text    not null,
    primary key (empire_id, effdt),
    constraint fk_empire_id foreign key (empire_id) references empire (id)
);

create table empire_system_name
(
    empire_id integer not null,
    system_id integer not null,
    effdt     integer not null,
    enddt     integer not null,
    name      text    not null,
    primary key (empire_id, system_id, effdt),
    constraint fk_empire_id foreign key (empire_id) references empire (id),
    constraint fk_system_id foreign key (system_id) references systems (id)
);

create table empire_star_name
(
    empire_id integer not null,
    star_id   integer not null,
    effdt     integer not null,
    enddt     integer not null,
    name      text    not null,
    primary key (empire_id, star_id, effdt),
    constraint fk_empire_id foreign key (empire_id) references empire (id),
    constraint fk_star_id foreign key (star_id) references stars (id)
);

create table scs
(
    id            integer primary key autoincrement,
    empire_id     integer not null,
    sc_cd         text    not null,
    sc_tech_level integer not null check (sc_tech_level between 1 and 10),
    constraint fk_empire_id foreign key (empire_id) references empire (id),
    constraint fk_sc_cd foreign key (sc_cd) references sc_codes (code)
);

create table sc_inventory
(
    sc_id           integer not null,
    unit_cd         text    not null,
    unit_tech_level integer not null check (unit_tech_level between 0 and 10),
    effdt           integer not null,
    enddt           integer not null,
    qty             integer not null,
    mass            real    not null,
    volume          real    not null,
    is_assembled    integer not null default 0 check (is_assembled in (0, 1)),
    is_stored       integer not null default 0 check (is_stored in (0, 1)),
    primary key (sc_id, unit_cd, unit_tech_level, effdt),
    constraint fk_sc_id foreign key (sc_id) references scs (id),
    constraint fk_unit_cd foreign key (unit_cd) references unit_codes (code)
);

create table sc_location
(
    sc_id         integer not null,
    effdt         integer not null,
    enddt         integer not null,
    orbit_id      integer not null,
    is_on_surface integer not null default 0 check (is_on_surface in (0, 1)),
    primary key (sc_id, effdt),
    constraint fk_sc_id foreign key (sc_id) references scs (id),
    constraint fk_location foreign key (orbit_id) references orbits (id)
);

create table sc_name
(
    sc_id integer not null,
    effdt integer not null,
    enddt integer not null,
    name  text    not null,
    primary key (sc_id, effdt),
    constraint fk_sc_id foreign key (sc_id) references scs (id)
);

create table sc_population
(
    sc_id         integer not null,
    population_cd text    not null,
    effdt         integer not null,
    enddt         integer not null,
    qty           integer not null check (qty >= 0),
    pay_rate      real    not null check (pay_rate >= 0),
    rebel_qty     integer not null check (rebel_qty >= 0),
    primary key (sc_id, population_cd, effdt),
    constraint fk_sc_id foreign key (sc_id) references scs (id),
    constraint fk_population_cd foreign key (population_cd) references population_codes (code)
);

create table sc_rates
(
    sc_id      integer not null,
    effdt      integer not null,
    enddt      integer not null,
    rations    real    not null check (rations >= 0),
    sol        real    not null check (sol >= 0),
    birth_rate real    not null check (birth_rate >= 0),
    death_rate real    not null check (death_rate >= 0),
    primary key (sc_id, effdt),
    constraint fk_sc_id foreign key (sc_id) references scs (id)
);

-- a group is a collection of units (factories, farms, labs, or mines) working together
-- to produce items (resources like mets and food, or units like cngd and aut-1). this
-- table is the parent to many property tables for the group.
--
-- the table is effective dated to allow the player to disassemble groups.
--
-- the group_no is the group number displayed on reports and used in orders.
--
-- for factories, the item being produced is stored in the sc_group_tooling table.
-- that table is effective-dated to allow factories to be retooled, so you will
-- need to use the as_of_date to look it up.
--
-- for mines, the deposit being mined is stored in the sc_group_deposit table.
-- that table is effective-dated to allow mines to be transferred to new deposits.
-- you will need to use the as_of_date to look it the current deposit.
--
-- farms always produce food, so there is no similar table for them.
--
-- labs aren't used yet.
create table sc_group
(
    id    integer primary key autoincrement,
    sc_id integer not null,
    kind  text    not null check (kind in ('factory', 'farm', 'lab', 'mine')),
    effdt integer not null,
    enddt integer not null,
    constraint fk_sc_id foreign key (sc_id) references scs (id)
);

-- the group number is used when assigning resources to a group. resources are
-- assigned to lower numbers first.
--
-- I don't know of a way to enforce the constraint on group_no being unique
-- for a group in Sqlite3, so the application must manage that.
create table sc_group_no
(
    group_id integer not null,
    effdt    integer not null,
    enddt    integer not null,
    group_no integer not null check (group_no between 1 and 40),
    unique (group_id, effdt),
    constraint fk_group_id foreign key (group_id) references sc_group (id)
);

-- the group deposit table stores the deposit being mined by the mine group.
-- it is effective-dated to allow mine groups to be transferred to new deposits.
create table sc_group_deposit
(
    group_id   integer not null,
    effdt      integer not null,
    enddt      integer not null,
    deposit_id integer not null,
    unique (group_id, effdt),
    constraint fk_group_id foreign key (group_id) references sc_group (id),
    constraint fk_deposit_id foreign key (deposit_id) references deposits (id)
);

-- the group tooling table stores the item being produced by the factory group.
-- it is effective-dated to allow factory groups to be retooled.
--
-- note that there is a penalty for retooling. the group will stop making new
-- items for three turns after the order is issued (meaning a new row has been
-- inserted into this table). work will resume on the fourth turn (which is
-- retooled_dt + 3).
--
-- when a new game is created, we assume that factories will deliver new items
-- immediately. this implies that they've been working at least three turns.
-- the setup is "turn 0" and effdt can't be negative, so we added a boolean
-- field, "retooled," to indicate that the record was created because of a
-- retool order. that field will be 0 (false) for records created during setup.
-- when it is 1 (true), the record was created because of a retool order and
-- items will not be produced until the fourth turn (effdt + 3).
create table sc_group_tooling
(
    group_id        integer not null,
    effdt           integer not null,
    enddt           integer not null,
    item_cd         text    not null,
    item_tech_level integer not null check (item_tech_level between 0 and 10),
    retooled        integer not null check (retooled in (0, 1)),
    unique (group_id, effdt),
    constraint fk_group_id foreign key (group_id) references sc_group (id),
    constraint fk_orders_cd foreign key (item_cd) references unit_codes (code)
);

-- the group unit table is the set of all units in a group that have the same
-- tech level. it is effective-dated to allow the player to add and remove
-- units from the group.
create table sc_group_unit
(
    group_id     integer not null,
    tech_level   integer not null,
    effdt        integer not null,
    enddt        integer not null,
    nbr_of_units integer not null check (nbr_of_units >= 0),
    primary key (group_id, tech_level, effdt),
    constraint fk_group_id foreign key (group_id) references sc_group (id)
);

-- group unit production holds the resources consumed and produced by all
-- the units in a group during a turn.
--
-- the production date on this table is the turn that items were produced.
--
-- note that the item being produced is not stored in this table.
create table sc_group_unit_production
(
    group_id      integer not null,
    tech_level    integer not null,
    production_dt integer not null,
    fuel_consumed integer not null,
    gold_consumed integer not null,
    mets_consumed integer not null,
    nmts_consumed integer not null,
    pro_consumed  integer not null,
    usk_consumed  integer not null,
    aut_consumed  integer not null,
    qty_produced  integer not null,
    primary key (group_id, tech_level, production_dt),
    constraint fk_group_id foreign key (group_id) references sc_group (id)
);

-- group unit production wip stores information about work in progress
-- in a group of factory or lab units.
--
-- note that the item being produced is not stored in this table.
create table sc_group_unit_production_wip
(
    group_id      integer not null,
    tech_level    integer not null,
    production_dt integer not null,
    wip_25pct_qty integer not null check (wip_25pct_qty >= 0),
    wip_50pct_qty integer not null check (wip_50pct_qty >= 0),
    wip_75pct_qty integer not null check (wip_75pct_qty >= 0),
    primary key (group_id, tech_level, production_dt),
    constraint fk_group_wip_id foreign key (group_id, tech_level, production_dt) references sc_group_unit_production (group_id, tech_level, production_dt)
);

-- group production summary holds the total resources consumed
-- as well as the number of units produced by an entire group each turn.
-- you can think of this as the rollup for the group unit production table.
--
-- note that the item being produced is not stored in this table.
create table sc_group_production_summary
(
    group_id      integer not null,
    production_dt integer not null,
    fuel_consumed integer not null,
    gold_consumed integer not null,
    mets_consumed integer not null,
    nmts_consumed integer not null,
    pro_consumed  integer not null,
    usk_consumed  integer not null,
    aut_consumed  integer not null,
    qty_produced  integer not null,
    primary key (group_id, production_dt),
    constraint fk_group_id foreign key (group_id) references sc_group (id)
);

-- group production wip summary stores information about work in progress
-- in a group of factory or lab units. you can think of this as the rollup
-- for the group unit production wip table.
--
-- note that the item being produced is not stored in this table.
create table sc_group_production_wip_summary
(
    group_id      integer not null,
    production_dt integer not null,
    wip_25pct_qty integer not null check (wip_25pct_qty >= 0),
    wip_50pct_qty integer not null check (wip_50pct_qty >= 0),
    wip_75pct_qty integer not null check (wip_75pct_qty >= 0),
    primary key (group_id, production_dt),
    constraint fk_group_wip_id foreign key (group_id, production_dt) references sc_group_production_wip_summary (group_id, production_dt)
);

-- mining summary holds the total resources mined on a planet in a turn
-- by a single colony.
create table sc_mining_summary
(
    sc_id         integer not null,
    production_dt integer not null,
    fuel_produced integer not null,
    gold_produced integer not null,
    mets_produced integer not null,
    nmts_produced integer not null,
    primary key (sc_id, production_dt),
    constraint fk_sc_id foreign key (sc_id) references scs (id)
);

create table sc_probe_order
(
    id        integer primary key autoincrement,
    sc_id     integer not null,
    effdt     integer not null,
    target_id integer not null,
    kind      text    not null check (kind in ('colony', 'orbit', 'ship', 'star', 'system')),
    unique (sc_id, effdt, target_id),
    constraint fk_sc_id foreign key (sc_id) references scs (id)
);

create table sc_probe_star_result
(
    probe_id      integer not null,
    effdt         integer not null,
    star_id       integer not null,
    location      text    not null,
    nbr_of_orbits integer not null,
    primary key (probe_id, effdt),
    constraint fk_probe_id foreign key (probe_id) references sc_probe_order (id),
    constraint fk_star_id foreign key (star_id) references stars (id)
);

create table sc_probe_star_orbit_result
(
    probe_id   integer not null,
    effdt      integer not null,
    star_id    integer not null,
    orbit_no   integer not null,
    orbit_kind text    not null,
    fuel_est   integer not null,
    gold_est   integer not null,
    mets_est   integer not null,
    nmts_est   integer not null,
    primary key (probe_id, effdt),
    constraint fk_probe_id foreign key (probe_id) references sc_probe_order (id),
    constraint fk_star_id foreign key (star_id) references stars (id)
);

create table sc_survey_order
(
    id        integer primary key autoincrement,
    sc_id     integer not null,
    effdt     integer not null,
    target_id integer not null,
    kind      text    not null check (kind in ('colony', 'orbit', 'ship', 'star', 'system')),
    unique (sc_id, effdt, target_id),
    constraint fk_sc_id foreign key (sc_id) references scs (id)
);

create table sc_survey_orbit_result
(
    survey_id       integer not null,
    effdt           integer not null,
    orbit_id        integer not null,
    location        text    not null,
    orbit_no        integer not null,
    habitability_no integer not null,
    farmland_in_use integer not null,
    population      integer not null,
    primary key (survey_id, effdt),
    constraint fk_survey_id foreign key (survey_id) references sc_survey_order (id),
    constraint fk_orbit_id foreign key (orbit_id) references orbits (id)
);

