# PortPecker

PortPecker is a network checking tool designed to verify connectivity between hosts based on a set of predefined rules. It's built in Go and runs in a Docker environment, making it easy to deploy and test network configurations.

## Features

-   Checks network connectivity based on a JSON configuration file
-   Supports both TCP and UDP protocols
-   Can specify rules for specific hosts or rules that apply to all hosts
-   Provides colorized and structured output for easy readability
-   Runs in a Docker container for easy deployment and testing

## Prerequisites

-   Docker
-   Docker Compose

## Quick Start

1. Clone the repository:

```
git clone https://github.com/yourusername/portpecker.git
cd portpecker
```

2. Create a `config.json` file in the project root directory. Example:

```json
{
    "rules": [
        {
            "source": "172.20.0.2",
            "destination": "172.20.0.3",
            "ports": [
                { "port": "80", "protocol": "TCP" },
                { "port": "443", "protocol": "TCP" }
            ],
            "note": "Web server access"
        },
        {
            "source": "*",
            "destination": "172.20.0.4",
            "ports": [{ "port": "5432", "protocol": "TCP" }],
            "note": "Database access (all hosts)"
        }
    ]
}
```

3. Build and run the Docker environment:
   `docker-compose up --build -d`

4. Check the logs to see the results:
   `docker-compose logs portpecker`

## Configuration

The config.json file defines the rules for network checking. Each rule has the following structure:

-   source: The source IP address. Use "\*" to apply the rule to all hosts.
-   destination: The destination IP address to check connectivity to.
-   ports: An array of port objects, each containing:
    -   port: The port number to check.
    -   protocol: The protocol to use (TCP or UDP).
-   note: A description of the rule.

## Docker Environment

The project uses Docker Compose to set up a testing environment with three containers:

1. portpecker: Runs the PortPecker tool.
2. webserver: A simple container listening on ports 80 and 443.
3. database: A simple container listening on port 5432.

These containers are connected to a custom Docker network with predefined IP addresses.

## Development

To modify the PortPecker tool:

1. Update the Go code in main.go.
2. Modify the Dockerfile if you need to change the build process.
3. Update docker-compose.yml if you need to change the Docker environment.

After making changes, rebuild and run the Docker environment:
`docker-compose up --build -d`

## Troubleshooting

If you encounter "connection refused" errors, ensure that the destination containers are running and listening on the specified ports.
If the tool can't find the config file, check that the config.json file is in the correct location and properly mounted in the Docker container.
