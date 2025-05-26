# Consul DNS Resolution Setup Guide

This guide explains how to set up DNS resolution for Consul services in your system.

## Prerequisites
- Consul server running locally
- Systemd-resolved service (default on most modern Linux distributions)

## DNS Resolution Setup

### 1. Verify Consul DNS Service
First, verify that Consul's DNS service is running and responding:
```bash
dig @127.0.0.1 -p 8600 config-service.service.consul
```

### 2. Configure Systemd-resolved
Create a configuration file for Consul DNS:

```bash
# Create the configuration directory
sudo mkdir -p /etc/systemd/resolved.conf.d

# Create and edit the Consul configuration file
sudo nano /etc/systemd/resolved.conf.d/consul.conf
```

Add the following configuration:
```ini
[Resolve]
DNS=127.0.0.1:8600
Domains=~consul
```

### 3. Apply Configuration
Restart the systemd-resolved service to apply the changes:
```bash
sudo systemctl restart systemd-resolved
```

### 4. Verify Configuration
Test the DNS resolution:
```bash
# Check if the service is resolvable
getent hosts config-service.service.consul

# Test connection to the service
telnet config-service.service.consul 50051
```

## How It Works
- Consul runs a DNS server on port 8600
- The configuration tells systemd-resolved to:
  - Use Consul's DNS server (127.0.0.1:8600) for DNS queries
  - Handle all domains ending in `.consul`
- This allows you to use service names like `service-name.service.consul` for service discovery

## Troubleshooting
If DNS resolution is not working:
1. Ensure Consul is running and healthy
2. Verify the DNS configuration is correct
3. Check systemd-resolved service status: `systemctl status systemd-resolved`
4. Check DNS resolution logs: `journalctl -u systemd-resolved`
