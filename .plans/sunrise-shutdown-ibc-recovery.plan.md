---
name: Sunrise shutdown IBC recovery
overview: 高さ H のハンドラが回収前に全モジュールの genesis JSON をファイルへ書く。読みやすい残高台帳はそのファイルから後で作る。アップグレードは IBC バウチャーの集権口座への回収と、未着パケットを触ってもチェーンが止まらないガードを行う。
todos:
  - id: pre-state-file
    content: 高さ H で全ノードの DAEMON_HOME/shutdown に genesis JSON を書く。ディスク不足でも回収は続行し、ログに失敗を残す
    status: pending
  - id: ledger-cli
    content: 保存した genesis JSON から bank.json / positions.json / lockup.json / in_flight.json を後から作るコマンドを追加する
    status: pending
  - id: inflight-guard
    content: 未着パケットの ack と timeout は再送も返金もせずエラーで戻し、チェーンを止めない
    status: pending
  - id: sweep
    content: transfer エスクローを除き ibc/ を sunrise1xxgjt7yqkmn63m2d0nrf0vt5uuc2hr6l45xaa9 へ送る
    status: pending
  - id: guards
    content: 新規入金 IBC を拒否し、集権口座からの戻し MsgTransfer は通す
    status: pending
  - id: tests
    content: 台帳の分離、ポジション数量、回収範囲のテストを追加する
    status: pending
isProject: false
---

# Sunrise 廃止前の IBC 回収と保有者台帳

## 記録の出どころ

バイナリの更新、Proposal、投票は状態を保存しない。高さ H でハンドラが `ibc/` を動かすと、その前の銀行残高は稼働中のチェーンから参照できなくなる。通常の pruning では過去高さの状態も残らない。

正本は、高さ H のハンドラが回収より先に書くファイルである。実装に含める。別ノードで `sunrised export` を実行する必要はない。ハンドラはプロセス内の `ExportGenesisForModules` を使う。手動で export する場合は、PATH の `/root/go/bin/sunrised` ではなく、cosmovisor が起動している `/root/.sunrise/cosmovisor/upgrades/v1.2.0/bin/sunrised` を使う。前者は `libwasmvm.x86_64.so` が無く起動しない。

全モジュール指定の `sunrised export` は `module ibchooks does not exist` で失敗する。`ibchooks` は [app/app_config.go](app/app_config.go) の InitGenesis 順にあるが、AppModule として登録されていない。ハンドラの書き出しもこの名前を除く。稼働中ノードに対する export はデータベースのロックで失敗するので、CLI で取るときは `systemctl stop cosmovisor` のあと、同じバイナリで実行し、終わったら `systemctl start cosmovisor` する。

2026-09-28 に consensus-a でこの手順を確認した。ファイルは `/root/.sunrise/shutdown/pre-upgrade-state.json`（342MB、sha256 `8eec7e295eb840264d8a7f2900aed59336d3df9c1b79f869f356660480bc1ccd`、`initial_height` 6394813）。これは高さ H の正本ではなく、手順の確認用である。`ibchooks` は含まれない。本番のハンドラは同じパスへ、高さ H の処理直前の状態で書き直す。consensus-a の確認用ファイルはそこで置き換わる。残す場合は高さ H の前に別名でコピーする。他のノードにはこの確認用ファイルは無く、高さ H で新規に作られる。

別ノードの `sunrised export` は、高さ H より前の状態になる。その後のブロックで残高が動くので、処理直前の正本にはしない。

高さ H を実行する全ノード（コンセンサスもフルノードも）が、それぞれの `DAEMON_HOME/shutdown/pre-upgrade-state.json` に同じ内容を書く。consensus-a では `DAEMON_HOME` が `/root/.sunrise` で、このマウント上にできる。バリデータだけではない。

ディスクが足りないノードはファイルを書けない。その失敗でハンドラをエラーにすると、書けたノードと書けないノードでブロックの成否が割れ、チェーンが分岐する。だから書き込み失敗はログに残し、回収は全ノードで続ける。確認用ファイルは 342MB だったので、高さ H の前に `DAEMON_HOME` の空きを 2GB 以上見ておく。書けなかったノードは、書けたノードのファイルを sha256 で照合してコピーする。全ノードが書けなかった場合、回収後のチェーンからは処理前状態を復元できない。

この JSON が処理前のアプリケーション状態である。銀行残高、ポジション、未請求報酬、lockup、未着パケット、wasm、IBC のクライアントとチャネル、ステーキングがモジュールごとの `app_state` に入る。ブロック履歴と IAVL の旧バージョンは入らない。

ブロック履歴まで残す場合だけ、高さ H の前にノードを1台止め、その `data` ディレクトリをコピーする。これは `sunrised export` とは別のオペレータ作業であり、ハンドラの中ではコピーしない。

読みやすい台帳は、保存した JSON を入力に後から作る。チェーンが動いていなくても、シャットダウン後でも再生成できる。

- `shutdown-ledger/bank.json`
- `shutdown-ledger/positions.json`
- `shutdown-ledger/lockup.json`
- `shutdown-ledger/in_flight.json`

export の genesis JSON では、銀行残高とポジションは最初から別モジュールに入る。

- `app_state.bank.balances` — アドレスごとの銀行残高。プールアドレスの残高もここにある
- `app_state.liquiditypool.positions` — 保有者、pool、tick、liquidity。トークン数量は入っていない
- `app_state.liquiditypool.accumulator_positions` — 未請求報酬
- lockup と swap の in-flight も、それぞれの `app_state` にある

台帳ファイルはオペレータの作業用であり、リポジトリにはコミットしない。正本の `pre-upgrade-state.json` もコミットしない。

- `shutdown-ledger/bank.json` — 銀行残高をそのまま。ポジションの分解額は混ぜない
- `shutdown-ledger/positions.json` — ポジションごとに、全額引き出しと同じ向きで計算した base / quote と未請求報酬
- `shutdown-ledger/lockup.json` — lockup 派生アドレスの銀行残高を owner に対応づけた一覧。`bank.json` の行は書き換えない
- `shutdown-ledger/in_flight.json` — 未着パケット。回収額に入れない

`positions.json` はプールアドレスが持つ銀行残高の内訳である。`bank.json` のプール残高と足して供給量にしない。数量の計算はアップグレード前の export にあるプール価格を使う。アップグレードはポジション状態を書き換えない。

### positions.json の1行

数量は `DecreaseLiquidity` と同じ計算にする。liquidity を負にして `Pool.CalcActualAmounts` を呼び、返った `TruncateInt` の絶対値を使う。ポジションクエリ（正の liquidity で切り上げ）とは、最小単位で1ずれることがある。台帳は引き出し額に合わせる。

各ポジションは、その liquidity とプールの現在価格だけで決まる。他のポジションを先に引き出しても現在価格は変わらないので、計算順は数量を変えない。端数でポジション合計とプールの銀行残高がずれた分は、`bank.json` のプールアドレスに残る。

入力は、そのポジションの lower tick、upper tick、liquidity と、同じ export のプールの `current_tick` と `current_sqrt_price` である。

- `position_id`、`owner`（`Position.address`）、`pool_id`
- `lower_tick`、`upper_tick`、`liquidity`（保存されている値のまま）
- `denom_base`、`denom_quote`（プールの denom）
- `amount_base`、`amount_quote`
  - 現在価格がレンジ内なら両方
  - 現在価格が下限より下なら base のみ、quote は 0
  - 現在価格が上限以上なら quote のみ、base は 0
- `unclaimed_rewards`（そのポジションの accumulator。denom と数量）

プールアドレスの銀行残高は `bank.json` にだけ置く。ポジション合計との差はプールアドレスの残高として残り、`positions.json` には書かない。

## アップグレードの実行

認識どおり、ガバナンス Proposal の `MsgSoftwareUpgrade` で名前 `v1.3.0` と高さ H を指定する。可決後、高さ H のブロックでハンドラが走る。バリデータはその前に、ハンドラ入りのバイナリを cosmovisor へ入れておく。入っていないノードは高さ H で停止する。

台帳コマンドの入力は、ハンドラが回収前に書いた `pre-upgrade-state.json` である。高さ H の後に `sunrised export` したものは、回収後の状態なので使わない。

## トークンの扱い

- 回収先は `sunrise1xxgjt7yqkmn63m2d0nrf0vt5uuc2hr6l45xaa9`
- `urise` は送金可能な RISE。回収しない。出金もしない。lockup 分は `lockup.json` に owner 付きで残す
- `uvrise` は銀行送金が停止されたステーク側。委任と shareclass を保有者に換算し、RISE の記録として `bank.json` とは別の換算行を `lockup.json` に隣接する `staking.json` に書く
- `uusdrise` は `bank.json` 上の denom として残す。出金依頼は USDN として扱う。オンチェーン変換はしない
- `ibc/` は `bank.json` に保有者を残したあと、アップグレードで集権口座へ集める。コントラクトは銀行残高のアドレスを保有者とし、内部台帳は分解しない
- 未着の実トークンは transfer エスクローに残す

## アップグレード v1.3.0 がやること

1. 回収の前に全モジュールの genesis JSON を `DAEMON_HOME/shutdown/pre-upgrade-state.json` へ書く。ディスク不足でも回収は止めない
2. `ibc/` を集権口座へ `SendCoins` する。transfer エスクローと集権口座自身は除く
3. 未着パケットの ack と timeout が届いたら、`IbcKeeperFn` を呼ばずにエラーを返す。再送も返金もしない。通常の `MsgTransfer`（未着記録がないパケット）はそのまま通す
4. その後の新規入金 IBC は拒否する。集権口座からの戻し `MsgTransfer` は通す

[PR #311](https://github.com/sunriselayer/sunrise/pull/311) の再送、timeout 既定値、返金順の変更は入れない。戻し送金は未着の再送を使わず、集権口座の新しい `MsgTransfer` で行う。

元チェーンへの返済はアップグレード後の手動送金。ガス用の `uusdrise` は集権口座に事前に残す。Hermes は戻し送金の間も、未着の `excluded_sequences` を外さない。外すと timeout 処理が届く。ガードがその処理を失敗させるのでチェーンは止まらないが、未着のトークンはエスクローに残ったままになる。

手順で入れ替えられないのは次だけである。Proposal の指定高さ H で、先に全状態をファイルへ書き、同じブロックで回収する。その後に集権口座から元チェーンへ戻し、戻しが終わってからシャットダウンする。台帳 JSON の生成は、保存したファイルがあればアップグレード後でもシャットダウン後でもよい。ポジション数量の計算順は数量を変えない。
