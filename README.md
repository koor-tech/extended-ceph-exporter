# extended-cephmetrics-exporter

## Requirements

* Needs an admin user

    ```
    radosgw-admin user create --uid admin --display-name "Admin User" --caps "buckets=*;users=*;usage=read;metadata=read;zone=read"
    # Access key / "Username"
    radosgw-admin user info --uid admin | jq '.keys[0].access_key'
    # Secret key / "Password
    radosgw-admin user info --uid admin | jq '.keys[0].secret_key'
    ```
