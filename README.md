# Docker Hub Service with Aquila FCM

## Overview

`AquilaFCM` is a lightweight service designed to send notifications via Firebase Cloud Messaging (FCM) to specific tokens. It streamlines push notification delivery for mobile and web apps, ensuring efficient token management and message dispatching. The service is containerized and hosted on Docker Hub, allowing for easy deployment and integration into microservice architectures.

---

## Service Description

- **Service Name**: `AquilaFCM`
- **Docker Image**: `vkhangstack/aquilafcm`
- **Functionality**: Send notifications via FCM to specific tokens.

---

## Steps to Run the Service with Mount Volume

### 1. Prepare the Host Directory

Create a directory on your host machine to store configuration files or logs:

```bash
mkdir -p ~/aquilafcm-config

touch serviceAccount.json
```

Add a configuration file for the service:

```bash
echo "{
    "type": "service_account",
    "project_id": "",
    "private_key_id": "",
    "private_key": "",
    "client_email": "",
    "client_id": "",
    "auth_uri": "",
    "token_uri": "",
    "auth_provider_x509_cert_url": "",
    "client_x509_cert_url": "",
    "universe_domain": "googleapis.com"
}" > ~/aquilafcm-config/serviceAccount.json
```

### 2. Run the Docker Container

Use the following command to start the AquilaFCM container with the host directory mounted:

```bash
docker run -d \
    -v ~/aquilafcm-config/:/app/config \
    --name aquilafcm-service \
    vkhangstack:aquilafcm
```

#### Explanation of Flags:

- `-d`: Run the container in detached mode (in the background).
- `-v ~/aquilafcm-config:/app/config`: Mount the `~/aquilafcm-config` directory on the host to `/app/config` inside the container.
- `--name aquilafcm-service`: Assign the name `aquilafcm-service` to the container.

### 3. Access the Service

Once the container is running, you can use the service to send FCM notifications. For example, you might make an API request to the service endpoint (details depend on the service's implementation).

---

## Additional Notes

- **Dynamic Updates**: Changes to files in the `~/aquilafcm-config` directory (e.g., `config.json`) will be reflected immediately in the container.
- **Cleaning Up**: To stop and remove the container:
  ```bash
  docker stop aquilafcm-service
  docker rm aquilafcm-service
  ```

---

## Use Cases

- Sending push notifications to mobile apps, and web apps.
- Manage FCM tokens and deliver notifications efficiently.
- Integrating FCM notification functionality into microservices.

---

## How to use

### Send message single token

```bash
curl --location --request POST 'http://localhost:8080/send' \
--header 'Content-Type: application/json' \
--data '{
    "token":"token",
    "title": "",
    "body": "",
    "imageUrl": "",
    "data": {
      "key": "value",
    }
}'

```

### Send message single multiple tokens

```bash
curl --location --request PUT 'http://localhost:8080/send' \
--header 'Content-Type: application/json' \
--data '{
    "token":["token", "token2],
    "title": "",
    "body": "",
    "imageUrl": "",
    "data": {
      "key": "value",
    }
}'



For further customization or issues, refer to the official documentation of the `aquilafcm` service.
```
