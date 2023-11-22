<p align="center">
  <p align="center">
    <a href="https://github.com/yyewolf/gocd/releases/latest"><img alt="Release" src="https://img.shields.io/github/release/yyewolf/gocd.svg?style=flat-square"></a>
    <a href="/LICENSE"><img alt="Software License" src="https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square"></a>
    <a href="https://codeclimate.com/github/yyewolf/gocd/test_coverage"><img src="https://api.codeclimate.com/v1/badges/d9fcf617937d6026221f/test_coverage" /></a>
    <a href="https://codeclimate.com/github/yyewolf/gocd/maintainability"><img src="https://api.codeclimate.com/v1/badges/d9fcf617937d6026221f/maintainability" /></a>
    <a href="https://goreportcard.com/report/github.com/yyewolf/gocd"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/yyewolf/gocd"></a>
    <a href="https://godoc.org/github.com/yyewolf/gocd"><img src="https://godoc.org/github.com/yyewolf/gocd/backend?status.svg" alt="GoDoc"></a>
  </p>
</p>

# Go CD

**Go CD** is a utility for infrastructures using `docker compose` as a method of deployments. It allows you to update your containers using tokens.

## Usage

To get started with Go CD, follow these steps:

### 1. Deploy Go CD Container

You can deploy the Go CD container from the registry located at [ghcr.io/yyewolf/gocd](https://ghcr.io/yyewolf/gocd). Ensure that the following labels are configured in your `docker-compose.yml` or deployment manifest:

```yaml
version: '3'
services:
  gocd:
    image: ghcr.io/yyewolf/gocd
    environment:
      - "DISCORD_WEBHOOK=https://discord.com/api/webhooks/your-webhook-id/your-webhook-token"
    ports:
      - "8080:8080"
```

### 2. Deploy controlled container

You can also deploy a controlled container (it can be started before or after GoCD, it doesn't matter).

```yaml
version: '3'
services:
  container_a:
    image: your_image_a
    labels:
      - "gocd.enable=true"
      - "gocd.repo=<repo_url>" # Optional
      - "gocd.token=<token>"
```

### 3. Profit

You can now restart and update your container by doing the following :

```bash
GET /containers/<token>
``` 

Replace <token> with the token assigned to the container you want to update.

## Contributing

If you'd like to contribute to Go CD, please follow the contribution guidelines. License

This project is licensed under the MIT License - see the LICENSE file for details.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE.md) file for details.
