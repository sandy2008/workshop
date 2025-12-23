# VLAN実習：macvlan + ルータコンテナで L3 分離ネットワークを作る

## ゴール

以下の構成を **手を動かして理解**する。

```
192.168.10.0/24            192.168.20.0/24

[a:192.168.10.10] ─┐
                   ├─ [router]
[b:192.168.20.20] ─┘      ├ eth0: 192.168.10.1
                          ├ eth1: 192.168.20.1
                          └ eth2: NAT → host → Internet
```

- a ⇄ b は **router 経由で通信**
- a / b から **apk add が可能**
- macvlan による **L2 分離**
- router は **NIC 3 本**

---

## 前提条件

- Ubuntu 24
- 物理 NIC 名：`eth0`（適宜`ens5`などに読み替え）

  ```bash
  $ ip link
  1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN mode DEFAULT group default qlen 1000
      link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
  2: eth0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc mq state UP mode DEFAULT group default qlen 1000
      link/ether 00:15:5d:01:09:1d brd ff:ff:ff:ff:ff:ff
  ```

- **rootful podman**
- rootlessがtrueの場合は、以下の実習を`podman`の代わりに`sudo podman`で実施してください。

  ```bash
  podman info | grep rootless
  # rootless: false
  ```

---

## Step 1. ホストに VLAN サブインタフェースを作る

### なぜ？

macvlan は **親インタフェース**にぶら下がる。
VLAN を切ることで L2 を完全分離する。

```bash
sudo ip link add link eth0 name eth0.10 type vlan id 10
sudo ip link add link eth0 name eth0.20 type vlan id 20

sudo ip link set eth0.10 up
sudo ip link set eth0.20 up
```

※ IP は付けない（L2 用途のみ）

---

## Step 2. macvlan ネットワークを 2 つ作る

### VLAN10（a 側）

```bash
sudo podman network create \
  --driver macvlan \
  --subnet 192.168.10.0/24 \
  --gateway 192.168.10.1 \
  -o parent=eth0.10 \
  net-vlan10
```

### VLAN20（b 側）

```bash
sudo podman network create \
  --driver macvlan \
  --subnet 192.168.20.0/24 \
  --gateway 192.168.20.1 \
  -o parent=eth0.20 \
  net-vlan20
```

---

## Step 3. router コンテナを作る（NIC 3 本）

```bash
sudo podman run -d --name router \
  --network net-vlan10 \
  --ip 192.168.10.1 \
  --cap-add NET_ADMIN \
  --sysctl net.ipv4.ip_forward=1 \
  alpine sleep infinity

sudo podman network connect \
  --ip 192.168.20.1 \
  net-vlan20 router

sudo podman network connect podman router
```

---

## Step 4. router に NAT を設定する

```bash
podman exec router apk add iptables

podman exec router sh -c '
iptables -t nat -A POSTROUTING -o eth2 -j MASQUERADE
iptables -A FORWARD -i eth0 -o eth2 -j ACCEPT
iptables -A FORWARD -i eth1 -o eth2 -j ACCEPT
iptables -A FORWARD -i eth2 -m state --state ESTABLISHED,RELATED -j ACCEPT
'
```

---

## Step 5. container a（192.168.10.10）

```bash
sudo podman run -d --name a \
  --network net-vlan10 \
  --ip 192.168.10.10 \
  alpine sleep infinity
```

---

## Step 6. container b（192.168.20.20）

```bash
sudo podman run -d --name b \
  --network net-vlan20 \
  --ip 192.168.20.20 \
  alpine sleep infinity
```

---

## Step 7. 動作確認

```bash
podman exec a ping -c 3 192.168.20.20
podman exec b ping -c 3 192.168.10.10
podman exec a apk add curl
```

---

## まとめ

- macvlan は L2 直結
- VLAN はブロードキャスト分離
- router コンテナで L3 + NAT を実装
