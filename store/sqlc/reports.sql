-- CreateReport creates a new report for the given Sorc.
--
-- name: CreateReport :one
INSERT INTO reports (sorc_id, turn_no)
VALUES (:sorc_id, :turn_no)
RETURNING id;

-- CreateReportProductionInput adds a production input to the given report.
--
-- name: CreateReportProductionInput :one
INSERT INTO report_production_inputs (report_id, category, fuel, gold, metals, non_metals)
VALUES (:report_id, :category, :fuel, :gold, :metals, :non_metals)
RETURNING id;

-- CreateReportProductionOutput adds a production output to the given report.
--
-- name: CreateReportProductionOutput :one
INSERT INTO report_production_outputs (report_id, category, unit_cd, tech_level, farmed, mined, manufactured)
VALUES (:report_id, :category, :unit_cd, :tech_level, :farmed, :mined, :manufactured)
RETURNING id;

-- CreateReportSurvey adds a survey to the given report.
--
-- name: CreateReportSurvey :one
INSERT INTO report_surveys (report_id, orbit_id)
VALUES (:report_id, :orbit_id)
RETURNING id;

-- CreateReportSurveyDeposit adds a survey of a deposit to the given report.
--
-- name: CreateReportSurveyDeposit :exec
INSERT INTO report_survey_deposits (report_id, deposit_no, deposit_qty, deposit_kind, deposit_yield_pct)
VALUES (:report_id, :deposit_no, :deposit_qty, :deposit_kind, :deposit_yield_pct);

-- DeleteAllTurnReports deletes all reports for a turn in the game.
--
-- name: DeleteAllTurnReports :exec
DELETE
FROM reports
WHERE sorc_id in (SELECT sorcs.id
                  FROM empires,
                       sorcs
                  WHERE empires.game_id = :game_id
                    AND empires.id = sorcs.empire_id)
  AND turn_no = :turn_no;

-- ReadEmpireReports returns a list of reports for an empire for a turn.
--
-- name: ReadEmpireReports :many
SELECT DISTINCT sorcs.id AS sorc_id, reports.id AS report_id
FROM empires,
     sorcs,
     reports
WHERE empires.id = :empire_id
  AND sorcs.empire_id = empires.id
  AND reports.sorc_id = sorcs.id
  AND reports.turn_no = :turn_no
ORDER BY sorcs.id, reports.id;

-- ReadReport returns the report id for a sorc and turn.
--
-- name: ReadReport :one
SELECT id
FROM reports
WHERE sorc_id = :sorc_id
  AND turn_no = :turn_no;

-- ReadReportProductionInputs returns a list of production inputs for a sorc and turn.
--
-- name: ReadReportProductionInputs :many
SELECT category, fuel, gold, metals, non_metals
FROM reports,
     report_production_inputs
WHERE reports.sorc_id = :sorc_id
  AND reports.turn_no = :turn_no
  AND report_id = reports.id
ORDER BY category;

-- ReadReportProductionOutputs returns a list of production outputs for a sorc and turn.
--
-- name: ReadReportProductionOutputs :many
SELECT category, farmed, mined, manufactured
FROM reports,
     report_production_outputs
WHERE reports.sorc_id = :sorc_id
  AND reports.turn_no = :turn_no
  AND report_id = reports.id
ORDER BY unit_cd, tech_level;

-- ReadSystemSurveyReports returns a list of survey reports for a sorc in a given turn.
--
-- name: ReadSystemSurveyReports :many
SELECT report_id AS report_id,
       id        AS system_survey_id,
       orbit_id  AS orbit_id
FROM report_surveys
WHERE report_id = :report_id
ORDER BY id, orbit_id;


-- ReadSystemSurveyDeposits returns a list of deposits for a survey report.
--
-- name: ReadSystemSurveyDeposits :many
SELECT deposit_no, deposit_qty, deposit_kind, deposit_yield_pct
FROM report_survey_deposits
WHERE report_id = :report_id
ORDER BY deposit_no;

-- SelectDepositsSummary returns the summary of deposits for all planets.
--
-- name: SelectDepositsSummary :many
select planet_id,
       case when kind = 'FUEL' then sum(remaining_qty) else 0 end as fuel_qty,
       case when kind = 'GOLD' then sum(remaining_qty) else 0 end as gold_qty,
       case when kind = 'METS' then sum(remaining_qty) else 0 end as mets_qty,
       case when kind = 'NMTS' then sum(remaining_qty) else 0 end as nmts_qty
from deposits
group by planet_id
order by planet_id;