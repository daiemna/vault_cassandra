# Read-Write permission on 'kv/data/*' and 'sys/policies/acl/*' path
path "kv/data/*" {
  capabilities = [ "read", "create", "update", "list" ]
}

path "sys/policies/acl/*" {
   capabilities = [ "read", "create", "update", "list" ]
}

path "database/*" {
  capabilities = [ "read", "create", "update", "list" ]
}

path "database/config/*" {
  capabilities = [ "read", "create", "update", "list" ]
}

path "database/static-roles/postgres-client" {
  capabilities = [ "read", "create", "update", "list" ]
}


path "database/static-roles/vernemq_db-static-role-staging" {
  capabilities = [ "read", "create", "update", "list" ]
}

path "database/static-roles/vernemq_db-static-role-development" {
  capabilities = [ "read", "create", "update", "list" ]
}
