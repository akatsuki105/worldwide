# Worldwide
[![Build Status](https://travis-ci.com/Akatsuki-py/Worldwide.svg?branch=master)](https://travis-ci.com/Akatsuki-py/Worldwide)
[![Go Report Card](https://goreportcard.com/badge/github.com/Akatsuki-py/Worldwide)](https://goreportcard.com/report/github.com/Akatsuki-py/Worldwide)

Go言語で書かれたゲームボーイカラーエミュレータです。  

ほぼ全てのROMが問題なく動作し、サウンド機能やセーブ機能、一部通信機能など幅広い機能を備えたエミュレータです。

<img src="https://imgur.com/UnmQnVE.gif" width="320" height="288">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<img src="https://imgur.com/cFugCTA.gif" width="320" height="288">


<img src="https://imgur.com/8YR987D.png" width="320" height="288">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<img src="https://imgur.com/2zwsb84.png" width="320" height="288">

## このエミュレータの特徴 & 今後実装予定の機能
- [x] 60fpsで動作
- [x] [cpu_instrs](https://github.com/retrio/gb-test-roms/tree/master/cpu_instrs) というテストROMをクリアしています
- [x] 少ないCPU使用率
- [x] サウンドの実装
- [x] ゲームボーイカラーのソフトに対応
- [x] WindowsやLinuxなど様々なプラットフォームに対応
- [x] いくつかのゲームパッドをサポート
- [x] [WebAssemblyでエミュレータの一部機能ををWebアプリとして実装](https://akatsuki-py.github.io/Worldwide/wasm.html)
- [x] MBC1に対応
- [x] MBC2に対応
- [x] MBC3に対応
- [x] MBC5に対応
- [x] RTCも実装しています
- [x] セーブ機能をサポート(得られたsavファイルは実機やBGBなどの一般的なエミュレータで利用できます)
- [x] クイックセーブのサポート
- [x] ウィンドウの縮小拡大が可能
- [x] ゲームボーイモードでパレットカラーの変更をサポート
- [x] ポケモンクリスタルの日本語ROMをサポート
- [x] ローカルネットワーク内のゲームボーイの通信機能をサポート(未対応のROMもあります テトリス、ポケモン赤などが動作します)
- [ ] ローカルネットワーク内のゲームボーイカラーの通信機能をサポート
- [ ] ネットワークをまたいだ通信機能のサポート
- [ ] GUIの操作メニュー
- [ ] wasm版のサウンドのサポート
- [ ] ラズパイ対応
- [ ] スーパーゲームボーイのエミュレーション機能

## 実行方法

- go 1.13 以上で動作します。

```
go run ./cmd/main.go "xxx.gb"
```

## ダウンロード

[ここ](https://github.com/Akatsuki-py/Worldwide/releases)からダウンロードできます。最新版をダウンロードすることをお勧めします。

## コマンド

| キー入力             | コマンド      |
| -------------------- | ------------- |
| <kbd>&larr;</kbd>    | &larr; ボタン |
| <kbd>&uarr;</kbd>    | &uarr; ボタン |
| <kbd>&darr;</kbd>    | &darr; ボタン |
| <kbd>&rarr;</kbd>    | &rarr; ボタン |
| <kbd>X</kbd>         | A ボタン      |
| <kbd>Z</kbd>         | B ボタン      |
| <kbd>Enter</kbd>     | Start ボタン  |
| <kbd>Right shift</kbd> | Select ボタン |
| <kbd>E</kbd>         | ウィンドウの拡大  |
| <kbd>R</kbd>         | ウィンドウの縮小 |
| <kbd>D + S</kbd>     | クイックセーブ  |
| <kbd>L</kbd>         | クイックロード |

## Webバージョンの紹介

Goのwasmビルド機能を利用して作成した[Webアプリ版](https://akatsuki-py.github.io/Worldwide/)もあります。


<img src="https://imgur.com/7ZJxQIu.png">
