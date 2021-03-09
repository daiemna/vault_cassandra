#!/bin/bash
# export VAULT_TOKEN="thisisarealroottoken"

vault_bin_path=environment/bin
vault_ver=1.5.0
vault_zip_name=vault_${vault_ver}_linux_amd64.zip
vault_zip=$vault_bin_path/$vault_zip_name
vault_exec=$vault_bin_path/vault

if test -f "$vault_exec"; then
    echo "$vault_exec exist"
else
    mkdir $vault_bin_path
    if test -f "$vault_zip"; then
        echo "$vault_zip exist"
    else
        wget --directory-prefix=$vault_bin_path/ https://releases.hashicorp.com/vault/${vault_ver}/$vault_zip_name
    fi
    unzip $vault_zip -d $vault_bin_path
    chmod +x $vault_bin_path/vault
fi
export PATH="$(pwd)/$vault_bin_path":$PATH
export VAULT_ADDR=http://127.0.0.1:58200
export VAULT_TOKEN="testtoken"

vault secrets enable kv && \
vault kv enable-versioning kv/ && \
cat environment/read-write-kv-policy.hcl | vault policy write go-fidelio-read-write - && \
vault auth enable userpass && \
vault secrets enable database && \
vault write auth/userpass/users/john \
   password=doe \
   policies=go-fidelio-read-write && \
vault write database/config/scylla-prod \
      plugin_name="cassandra-database-plugin" \
      allowed_roles="*" \
      hosts="int-test-scylla-c,127.0.0.1" \
      protocol_version=4 \
      username="vaultadmin" \
      password="vaultpass" && \
vault write database/roles/my-role \
    db_name=scylla-prod \
    creation_statements="CREATE USER '{{username}}' WITH PASSWORD '{{password}}' NOSUPERUSER; \
          GRANT SELECT ON ALL KEYSPACES TO {{username}};" \
    default_ttl="60s" \
    max_ttl="60s"


