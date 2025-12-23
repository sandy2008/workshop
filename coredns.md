# CoreDNS 実習：親子 DNS サーバーを構築して名前解決を理解する

## ゴール

以下の構成を **手を動かして理解**する。

- DNSサーバーの基本的な役割
- DNS フォワーディングによる名前解決の連携
- 親子関係にある2つのDNSゾーンの設定
- `dig` を使ったDNSクエリの検証方法

```
+--------------------+      1. dig @192.168.100.10 -p 10053 test.foo.sokoide.com
|   Client (Your PC) |
+--------------------+
         |
         |
         v
+--------------------------------+       +------------------------------+
| VM1 (ns-parent)                |       | VM2 (ns-child)               |
| IP: 192.168.100.10             |       | IP: 192.168.100.20           |
|                                | 2. Forward Query for foo.sokoide.com |
| CoreDNS on port 10053          |------>| CoreDNS on port 10053        |
| - Owns "sokoide.com" zone      |       | - Owns "foo.sokoide.com" zone|
| - Forwards "foo.sokoide.com"   |<------|                              |
|                                | 3. Return result                     |
+--------------------------------+       +------------------------------+
         ^
         |
         | 4. Return final result
         |
         +----------------------------------
```

- **VM1** は `sokoide.com` ゾーンを管理し、`foo.sokoide.com` への問い合わせを **VM2** へ転送（Forward）する。
- **VM2** は `foo.sokoide.com` ゾーンを管理する。
- クライアントは `dig` コマンドを使い、Port `10053` で動作するDNSサーバーに問い合わせる。

---

## 前提条件

- Ubuntu 24 がインストールされた VM 2台
  - VM1 IP: `192.168.100.10` （以下 `ns-parent`）
  - VM2 IP: `192.168.100.20` （以下 `ns-child`）
- `curl`, `tar`, `dnsutils` (`dig`コマンドのため) が利用可能であること。

確認 (両方のVMで):

```bash
sudo apt update && sudo apt install -y curl tar dnsutils
```

---

## Step 1. 両方のVMにCoreDNSをインストールする

### なぜ？

CoreDNSはGoで書かれた単一バイナリのDNSサーバーで、設定が簡単です。公式サイトからバイナリをダウンロードして展開するだけで利用できます。

**VM1 と VM2 の両方で** 以下のコマンドを実行します。

```bash
# CoreDNSの最新版をダウンロード
CORE_VERSION="1.11.1"
curl -L "https://github.com/coredns/coredns/releases/download/v${CORE_VERSION}/coredns_${CORE_VERSION}_linux_amd64.tgz" -o coredns.tgz

# 展開してパスの通った場所に配置
tar -xzvf coredns.tgz
sudo mv coredns /usr/local/bin/

# バージョンを確認
coredns -version
```

---

## Step 2. VM1で親DNSサーバー (sokoide.com) を設定する

### VM1 (`ns-parent`): Corefileの作成

`sokoide.com` ゾーンを管理し、`foo.sokoide.com` をVM2に転送する設定ファイルを作成します。

```bash
# VM1で作業
mkdir -p ~/coredns_parent
cd ~/coredns_parent

# Corefileの作成
cat <<'EOF' > Corefile
sokoide.com:10053 {
    # このサーバーがsokoide.comゾーンの権威サーバーであることを示す
    # SOAレコードなどを定義したゾーンファイルを利用する
    file db.sokoide.com

    # ログを有効化
    log
    # エラーを標準出力に表示
    errors
}

foo.sokoide.com:10053 {
    # foo.sokoide.comゾーンに関する問い合わせを
    # VM2 (192.168.100.20) の10053ポートに転送する
    forward . 192.168.100.20:10053

    # ログを有効化
    log
    # エラーを標準出力に表示
    errors
}
EOF
```

### VM1 (`ns-parent`): ゾーンファイルの作成

`sokoide.com` の具体的なレコードを定義します。

```bash
# VM1で作業 (引き続き ~/coredns_parent)

cat <<'EOF' > db.sokoide.com
$ORIGIN sokoide.com.
$TTL 3600

@   IN  SOA     ns.sokoide.com. root.sokoide.com. (
        2024010101  ; Serial
        7200        ; Refresh
        3600        ; Retry
        1209600     ; Expire
        3600 )      ; Minimum TTL

; Name servers
@   IN  NS      ns.sokoide.com.

; A records for the name server itself
ns  IN  A       192.168.100.10

; Other records
www IN  A       1.1.1.1
EOF
```

---

## Step 3. VM2で子DNSサーバー (foo.sokoide.com) を設定する

### VM2 (`ns-child`): Corefileの作成

`foo.sokoide.com` ゾーンを管理する設定ファイルを作成します。

```bash
# VM2で作業
mkdir -p ~/coredns_child
cd ~/coredns_child

# Corefileの作成
cat <<'EOF' > Corefile
foo.sokoide.com:10053 {
    # ゾーンファイルを利用する
    file db.foo.sokoide.com

    # ログを有効化
    log
    # エラーを標準出力に表示
    errors
}
EOF
```

### VM2 (`ns-child`): ゾーンファイルの作成

`foo.sokoide.com` の具体的なレコードを定義します。

```bash
# VM2で作業 (引き続き ~/coredns_child)

cat <<'EOF' > db.foo.sokoide.com
$ORIGIN foo.sokoide.com.
$TTL 3600

@   IN  SOA     ns1.foo.sokoide.com. root.foo.sokoide.com. (
        2024010101  ; Serial
        7200        ; Refresh
        3600        ; Retry
        1209600     ; Expire
        3600 )      ; Minimum TTL

; Name servers
@   IN  NS      ns1.foo.sokoide.com.

; A records
ns1  IN  A      192.168.100.20
test IN  A      2.2.2.2
EOF
```

---

## Step 4. 両方のVMでCoreDNSサーバーを起動する

### フォアグラウンドで起動

デバッグや動作確認のため、まずはフォアグラウンドで起動します。（本番環境では`systemd`などでデーモン化します）

**VM1 (`ns-parent`) で実行:**

```bash
# 新しいターミナルを開いて実行
cd ~/coredns_parent
/usr/local/bin/coredns -conf Corefile
```

**VM2 (`ns-child`) で実行:**

```bash
# 新しいターミナルを開いて実行
cd ~/coredns_child
/usr/local/bin/coredns -conf Corefile
```

---

## Step 5. 動作確認

クライアント（どちらかのVM、あるいはホストPC）から `dig` コマンドでDNSクエリを送信し、名前解決が正しく行われるか確認します。

### 確認1: `sokoide.com` ゾーンのレコードを引く

`ns-parent` (VM1) に `www.sokoide.com` を問い合わせます。

```bash
# 192.168.100.10 は ns-parent のIPアドレス
dig @192.168.100.10 -p 10053 www.sokoide.com
```

**期待される結果 (抜粋):**
`ANSWER SECTION` に `www.sokoide.com. 3600 IN A 1.1.1.1` が返ってくれば成功です。

```
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: ...
;; flags: qr aa rd; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0
...
;; ANSWER SECTION:
www.sokoide.com. 3600 IN A 1.1.1.1
```

### 確認2: `foo.sokoide.com` ゾーンのレコードを引く（フォワーディング）

`ns-parent` (VM1) に `test.foo.sokoide.com` を問い合わせます。クエリは `ns-child` (VM2) に転送されるはずです。

```bash
# 192.168.100.10 は ns-parent のIPアドレス
dig @192.168.100.10 -p 10053 test.foo.sokoide.com
```

**期待される結果 (抜粋):**
`ANSWER SECTION` に `test.foo.sokoide.com. 3600 IN A 2.2.2.2` が返ってくれば成功です。
`SERVER` フィールドが `192.168.100.10#10053` となっており、クライアントは `ns-parent` としか通信していないことがわかります。

```
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: ...
;; flags: qr rd ra; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0
...
;; ANSWER SECTION:
test.foo.sokoide.com. 3600 IN A 2.2.2.2

;; SERVER: 192.168.100.10#10053
```

CoreDNSのログを見ると、`ns-parent` が `ns-child` にクエリを転送している様子が確認できます。

---

## まとめ

- CoreDNS を使って、特定のドメイン（ゾーン）を管理する権威DNSサーバーを簡単に構築できる。
- `forward` プラグインを使うことで、特定のゾーンへの問い合わせを別のDNSサーバーに転送できる。
- これにより、複数のDNSサーバーを連携させて、階層的なドメイン空間を管理できる。
- `dig` はDNSサーバーの動作確認に不可欠なツール。

この実習はフォワーディングの例でしたが、DNSのもう一つの重要な仕組みである「委譲 (Delegation)」も CoreDNS で実現できます。委譲では、親サーバーは子のサーバーの場所（NSレコード）を教えるだけで、最終的な名前解決はクライアントが子サーバーに直接問い合わせて行います。
