# Raspberry Pi - KVM Prometheus Node Exporter

Prometheus Node Exporter to collect system and hardware metrics from Raspberry Pi and KVM (Kernel-based Virtual Machine). It enables you to monitor and analyze the performance of your Raspberry Pi devices and virtual machines, allowing you to gain insights into their resource utilization, temperature, and more.


## Installation

### Prerequisites

Before installing the Prometheus Node Exporter, ensure you have the following prerequisites:

- A Raspberry Pi device with KVM installed on it.
- Raspberry Pi OS (Raspbian) or a compatible Linux distribution installed.
- Prometheus server configured (if you plan to use Prometheus for monitoring).

### Installation Steps

1. Clone the repository:

   ```bash
   git clone https://github.com/sam-rajan/raspi-kvm-exporter.git
   ```

2. Change to the project directory:

   ```bash
   cd raspi-kvm-exporter
   ```

3. Build the Node Exporter binary:

   ```bash
   make compile
   ```

4. Start the Node Exporter:

   ```bash
   ./build/raspi-kvm-exporter -port=9100 -config=/etc/collector.yml
   ```

By default, the Exporter will listen on port 9000. You can configure the port by passing `-port=<PORT>`

### Config file sample

```yaml
port: 9090
collectors:
  kvm:
    enabled: true
  raspi:
    enabled: true
    disk:
      enabled: true
      disk:
      - sda1
      - sda2
      - sda3
      - sda4
```


## Usage

Once the Node Exporter is running, you can access the collected metrics by visiting the following URL in your web browser or by configuring Prometheus to scrape the metrics:

```
http://localhost:<PORT>/metrics
```


