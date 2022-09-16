# extended-ceph-exporter

## Requirements

* Needs an admin user

    ```
    radosgw-admin user create --uid admin --display-name "Admin User" --caps "buckets=*;users=*;usage=read;metadata=read;zone=read"
    # Access key / "Username"
    radosgw-admin user info --uid admin | jq '.keys[0].access_key'
    # Secret key / "Password
    radosgw-admin user info --uid admin | jq '.keys[0].secret_key'
    ```

## Development

### Debugging

A VSCode debug config is available to run and debug the project.

To make the exporter talk with a Ceph RGW S3 endpoint, create a copy of the `.env.example` file and name it `.env`.
Be sure ot add your Ceph RGW S3 endpoint and credentials in it.
