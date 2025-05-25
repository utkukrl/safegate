# SafeGate

SafeGate is a secure and configurable reverse proxy server. It can manage HTTP and gRPC traffic, apply security measures, and modify rules in real-time.

## Features

- 🔒 HTTPS support
- 🛡️ IP whitelist filtering
- 🔑 JWT-based authentication
- ⚡ Rate limiting
- 🎮 REPL interface for real-time rule management
- 📊 Dashboard interface
- 🔄 HTTP and gRPC proxy support

## Installation

1. Clone the repository:
```bash
git clone https://github.com/utkukrl/safegate.git
cd safegate
```

2. Generate certificates:
```bash
chmod +x scripts/generate_certs.sh
./scripts/generate_certs.sh
```

3. Configure the application:
```bash
cp configs/config.yaml.example configs/config.yaml
```

4. Run the application:
```bash
cd cmd/server
go run main.go
```

## Configuration

You can configure the following settings in `configs/config.yaml`:

```yaml
proxy:
  target: "http://localhost"
  port: "8080"             

firewall:
  ip_whitelist:             
    - "127.0.0.1/32"
    - "::1/128"

rate_limit:
  enabled: true            
  rate: "100r/s"           
  burst: 200               

jwt:
  enabled: false           
  secret_key: "your-key"  

dashboard:
  enabled: true            
  port: "3000"            
          
```

## REPL Commands

You can use the following commands in the REPL interface:

- `block <path>` - Block the specified path
- `unblock <path>` - Remove block from the specified path
- `limit <METHOD> <path> rate=R burst=B` - Apply rate limiting to the specified path
- `show rules` - List active rules
- `help` - Show command list

## Security

- All communication is encrypted over HTTPS
- IP whitelist restricts access to allowed IPs only
- JWT authentication available
- Rate limiting protects against DDoS attacks


## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request
