![logo](./logo.png)

# 🌏 worldwide
![Go](https://github.com/pokemium/worldwide/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/pokemium/worldwide)](https://goreportcard.com/report/github.com/pokemium/worldwide)
[![GitHub stars](https://img.shields.io/github/stars/pokemium/worldwide)](https://github.com/pokemium/worldwide/stargazers)
[![GitHub license](https://img.shields.io/github/license/pokemium/worldwide)](https://github.com/pokemium/worldwide/blob/master/LICENSE)

Go言語で書かれたゲームボーイカラーエミュレータです。  

多くのROMが問題なく動作し、サウンド機能やセーブ機能など幅広い機能を備えたエミュレータです。

<img src="https://imgur.com/ZlrXAW9.png" width="320px"> <img src="https://imgur.com/xVqjkrk.png" width="320px"><br/>
<img src="https://imgur.com/E7oob9c.png" width="320px"> <img src="https://imgur.com/nYpkH95.png" width="320px">

## 🚩 このエミュレータの特徴 & 今後実装予定の機能
- [x] 60fpsで動作
- [x] [cpu_instrs](https://github.com/retrio/gb-test-roms/tree/master/cpu_instrs) と [instr_timing](https://github.com/retrio/gb-test-roms/tree/master/instr_timing)というテストROMをクリアしています
- [x] 少ないCPU使用率
- [x] サウンドの実装
- [x] ゲームボーイカラーのソフトに対応
- [x] WindowsやLinuxなど様々なプラットフォームに対応
- [x] MBC1, MBC2, MBC3, MBC5に対応
- [x] RTCの実装
- [x] セーブ機能をサポート(得られたsavファイルは実機やBGBなどの一般的なエミュレータで利用できます)
- [x] ウィンドウの縮小拡大が可能
- [x] HTTPサーバーAPI
- [ ] プラグイン機能
- [ ] [Libretro](https://docs.libretro.com/) APIのサポート
- [ ] ローカルネットワーク内での通信プレイ
- [ ] グローバルネットワーク内での通信プレイ
- [ ] SGBのサポート
- [ ] シェーダのサポート

## 🎮 使い方

[ここ](https://github.com/pokemium/worldwide/releases)から実行ファイルをダウンロードした後、次のように起動します。

```sh
./worldwide "***.gb" # もしくは `***.gbc`
```

## 🐛 HTTPサーバー

`worldwide`はHTTPサーバーを内包しており、ユーザーはHTTPリクエストを通じて `worldwide`にさまざまな指示を出すことが可能です。

[サーバードキュメント](./server/README.md)を参照してください。

## 🔨 ビルド

ソースコードからビルドしたい方向けです。

requirements
- Go 1.16
- make

```sh
make build                              # Windowsなら `make build-windows`
./build/darwin-amd64/worldwide "***.gb" # Windowsなら `./build/windows-amd64/worldwide.exe "***.gb"`
```

## 📄 コマンド

| キー入力             | コマンド      |
| -------------------- | ------------- |
| <kbd>&larr;</kbd>    | &larr; ボタン |
| <kbd>&uarr;</kbd>    | &uarr; ボタン |
| <kbd>&darr;</kbd>    | &darr; ボタン |
| <kbd>&rarr;</kbd>    | &rarr; ボタン |
| <kbd>X</kbd>         | A ボタン      |
| <kbd>Z</kbd>         | B ボタン      |
| <kbd>Enter</kbd>     | Start ボタン  |
| <kbd>Backspace</kbd> | Select ボタン |
