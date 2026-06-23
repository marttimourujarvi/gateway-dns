# HTTPRoute DNS Reconciler Architecture

> Flow of the HTTPRoute DNS reconciler showing how Kubernetes HTTPRoute events are translated into in-memory DNS A records served by a UDP DNS server.

```mermaid
flowchart TD
    %% Purpose: Shows how the HTTPRoute reconciler watches K8s resources, extracts Gateway IPs, and populates a shared in-memory DNS store consumed by a UDP DNS server.

    subgraph K8s["Kubernetes Cluster"]
        HR[HTTPRoute resource]
        GW[Gateway resource]
    end

    subgraph Controller["Controller Runtime"]
        W[Manager Watcher]
        R[HTTPRouteReconciler]
    end

    subgraph DNS["DNS Server"]
        S[UDP DNS Server :5353]
        H[Request Handler]
    end

    subgraph Store["Shared State"]
        DS[("DNSStore<br/>map[hostname] → ipv4")]
    end

    HR -->|watch events| W
    W -->|Reconcile Request| R

    R --> D{"DeletionTimestamp<br/>IsZero?"}

    D -->|"No<br/>(Deleting)"| F1{"Has Finalizer<br/>home.dns/finalizer?"}
    F1 -->|Yes| Del["Delete A records<br/>from DNSStore"]
    Del --> RF["Remove Finalizer<br/>Update HTTPRoute"]
    F1 -->|No| Done1["Return"]
    RF --> Done1

    D -->|"Yes<br/>(Create/Update)"| F2{"Has Finalizer?"}
    F2 -->|No| AF["Add Finalizer<br/>Update HTTPRoute"]
    F2 -->|Yes| GL["Lookup Parent Gateway(s)"]
    AF --> GL

    GL --> GA["Read Gateway.Status.Addresses"]
    GA --> IPV4{"IPv4 Address?"}
    IPV4 -->|No| Skip["Skip non-IPv4"]
    IPV4 -->|Yes| SA["Store A record<br/>DNSStore.Set(hostname., ip)"]
    Skip --> Done2["Return"]
    SA --> Done2

    R -.->|writes| DS
    Del -.->|deletes| DS

    Client["DNS Client"] -->|A query| S
    S --> H
    H --> L["DNSStore.Lookup(hostname)"]
    L -.->|reads| DS
    L --> Found{"Record Found?"}
    Found -->|No| NX["RcodeNameError"]
    Found -->|Yes| RA["Build A RR<br/>TTL 60s"]
    NX --> Resp["Write Response"]
    RA --> Resp
    Resp --> Client
```

---
*Generated on 2026-09-19 · diagram type: `flowchart`*
