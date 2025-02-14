#!/bin/bash
############################################################################
#
ddlScript="../pkg/store/sqlc/schema.sql"
[ -f "${ddlScript}" ] || {
  echo "error: missing ddl script"
  echo "       ddlScript == '${ddlScript}'"
  exit 2
}

############################################################################
# rebuild the reference database used by datagrip and the ide
reference_db="reference.db"
sqlite3 "${reference_db}" ".read ../pkg/store/sqlc/schema.sql" || exit 2
echo " info: rebuilt '${reference_db}'"

############################################################################
# rebuild the server's database
for server_db in alpha01.db; do
  rm -f "${server_db}"
  sqlite3 "${server_db}" ".read ../pkg/store/sqlc/schema.sql" || exit 2
  echo " info: rebuilt '${server_db}'"
done

exit 0
