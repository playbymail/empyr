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