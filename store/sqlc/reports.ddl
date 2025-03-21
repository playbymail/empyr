--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- deposits_summary is a temporary table that is used to summarize the
-- deposits table for reporting. It holds the total quantity of each
-- resource type from all the deposits on a planet. It also holds the
-- log10 estimated value of the resources.
create table deposits_summary
(
    planet_id    integer not null,
    fuel_qty     integer not null,
    fuel_est_qty integer not null,
    gold_qty     integer not null,
    gold_est_qty integer not null,
    mets_qty     integer not null,
    mets_est_qty integer not null,
    nmts_qty     integer not null,
    nmts_est_qty integer not null
);

-- system survey report example
-- S/C 131: Survey Report for Planet #2 in System 15/15/15A
-- Habitability    =             25
-- Farmland in use =      1,950,000
-- Population      =    238,800,600
-- Deposits
--   No  Type  Yield  Qty_______  Nbr of Mines  TL
--   01  GOLD    3 %     951,003
--   02  FUEL   12 %   3,987,321
--   03  METS   47 %  12,227,283
--   04  NMTS   33 %  57,142,889

-- system survey report table.
-- this table holds the data needed to generate the system survey report.
--
create table system_survey_reports
(
    id              integer   not null primary key autoincrement,
    game_code       text      not null,
    turn_no         integer   not null,
    empire_no       integer   not null,
    sorc_no         integer   not null,
    star_name       text      not null,
    planet_no       integer   not null,
    habitability    integer   not null,
    farmland_in_use integer   not null,
    population      integer   not null,
    created_at      timestamp not null default current_timestamp,
    unique (game_code, turn_no, empire_no, sorc_no, star_name, planet_no)
);

create table system_survey_deposits
(
    report_id     integer not null,
    deposit_no    integer not null,
    resource_type text    not null,
    yield_pct     integer not null,
    qty_remaining integer not null,
    unique (report_id, deposit_no),
    foreign key (report_id) references system_survey_reports (id)
);

create table system_survey_units
(
    report_id    integer not null,
    empire_no    integer not null,
    unit_cd      text    not null,
    tech_level   integer not null,
    nbr_of_units integer not null,
    deposit_no   integer not null,
    unique (report_id, empire_no, unit_cd, tech_level),
    foreign key (report_id) references system_survey_reports (id)
);

